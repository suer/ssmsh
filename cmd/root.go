package cmd

import (
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
	"github.com/spf13/cobra"

	"github.com/suer/ssmsh/internal/awsconfig"
	"github.com/suer/ssmsh/internal/instance"
	"github.com/suer/ssmsh/internal/session"
)

var (
	flagProfile string
	flagRegion  string
	flagCommand string
)

func newRootCmd() *cobra.Command {
	root := &cobra.Command{
		Use:          "ssmsh <name-or-instance-id>",
		Short:        "Log in to an EC2 instance via AWS SSM Session Manager",
		Args:         cobra.ExactArgs(1),
		SilenceUsage: true,
		RunE:         runLogin,
	}
	root.Flags().StringVar(&flagProfile, "profile", "", "AWS profile name")
	root.Flags().StringVar(&flagRegion, "region", "", "AWS region")
	root.Flags().StringVarP(&flagCommand, "command", "c", "", "Command to run instead of an interactive shell")
	return root
}

func Execute() error {
	return newRootCmd().Execute()
}

func runLogin(cmd *cobra.Command, args []string) error {
	ctx := cmd.Context()
	target := args[0]

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
		Profile:    flagProfile,
		Region:     cfg.Region,
		Endpoint:   fmt.Sprintf("https://ssm.%s.amazonaws.com", cfg.Region),
	})
}
