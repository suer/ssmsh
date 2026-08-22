package cmd

import (
	"fmt"
	"runtime/debug"
	"strconv"
	"strings"

	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
	"github.com/spf13/cobra"

	"github.com/suer/ssmsh/internal/awsconfig"
	"github.com/suer/ssmsh/internal/instance"
	"github.com/suer/ssmsh/internal/session"
)

var (
	flagProfile      string
	flagRegion       string
	flagCommand      string
	flagLocalForward string
)

func newRootCmd() *cobra.Command {
	root := &cobra.Command{
		Use:          "ssmsh <name-or-instance-id>",
		Short:        "Log in to an EC2 instance via AWS SSM Session Manager",
		Args:         cobra.ExactArgs(1),
		SilenceUsage: true,
		RunE:         runLogin,
		Version:      version(),
	}
	root.Flags().StringVar(&flagProfile, "profile", "", "AWS profile name")
	root.Flags().StringVar(&flagRegion, "region", "", "AWS region")
	root.Flags().StringVarP(&flagCommand, "command", "c", "", "Command to run instead of an interactive shell")
	root.Flags().StringVarP(&flagLocalForward, "local-forward", "L", "", "Forward a local port to a remote port on the instance (local:remote)")
	return root
}

func Execute() error {
	return newRootCmd().Execute()
}

func version() string {
	info, ok := debug.ReadBuildInfo()
	if !ok || info.Main.Version == "" {
		return "(devel)"
	}
	return info.Main.Version
}

func runLogin(cmd *cobra.Command, args []string) error {
	ctx := cmd.Context()
	target := args[0]

	if flagCommand != "" && flagLocalForward != "" {
		return fmt.Errorf("--command and --local-forward cannot be used together")
	}

	var localPort, remotePort int
	if flagLocalForward != "" {
		var err error
		localPort, remotePort, err = parseLocalForward(flagLocalForward)
		if err != nil {
			return err
		}
	}

	if err := session.CheckPluginInstalled(); err != nil {
		return err
	}

	cfg, err := awsconfig.Load(ctx, flagProfile, flagRegion)
	if err != nil {
		return fmt.Errorf("load AWS config: %w", err)
	}
	if cfg.Region == "" {
		return fmt.Errorf("AWS region is not resolved, specify --region or set AWS_REGION")
	}

	instanceID, err := instance.Resolve(ctx, ec2.NewFromConfig(cfg), target, instance.InteractiveSelector())
	if err != nil {
		return err
	}

	return session.Start(ctx, ssm.NewFromConfig(cfg), session.StartOptions{
		InstanceID: instanceID,
		Command:    flagCommand,
		LocalPort:  localPort,
		RemotePort: remotePort,
		Profile:    flagProfile,
		Region:     cfg.Region,
		Endpoint:   fmt.Sprintf("https://ssm.%s.amazonaws.com", cfg.Region),
	})
}

func parseLocalForward(spec string) (local, remote int, err error) {
	localStr, remoteStr, ok := strings.Cut(spec, ":")
	if !ok {
		return 0, 0, fmt.Errorf("invalid --local-forward %q, want local:remote", spec)
	}

	local, err = strconv.Atoi(localStr)
	if err != nil {
		return 0, 0, fmt.Errorf("invalid local port in --local-forward %q: %w", spec, err)
	}
	remote, err = strconv.Atoi(remoteStr)
	if err != nil {
		return 0, 0, fmt.Errorf("invalid remote port in --local-forward %q: %w", spec, err)
	}
	if local < 1 || local > 65535 || remote < 1 || remote > 65535 {
		return 0, 0, fmt.Errorf("ports in --local-forward %q must be between 1 and 65535", spec)
	}
	return local, remote, nil
}
