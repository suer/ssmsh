package session

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
)

const nonInteractiveCommandDocument = "AWS-StartNonInteractiveCommand"

type StartOptions struct {
	InstanceID string
	Command    string
	Profile    string
	Region     string
	Endpoint   string
}

func Start(ctx context.Context, client SSMAPI, opts StartOptions) error {
	input := buildStartSessionInput(opts.InstanceID, opts.Command)

	output, err := client.StartSession(ctx, input)
	if err != nil {
		return fmt.Errorf("start session for %s: %w", opts.InstanceID, err)
	}

	args, err := buildPluginArgs(output, input, opts.Region, opts.Profile, opts.Endpoint)
	if err != nil {
		return err
	}

	cmd := exec.CommandContext(ctx, pluginBinary, args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		var exitErr *exec.ExitError
		if errors.As(err, &exitErr) {
			os.Exit(exitErr.ExitCode())
		}
		return fmt.Errorf("run %s: %w", pluginBinary, err)
	}
	return nil
}

func buildStartSessionInput(instanceID, command string) *ssm.StartSessionInput {
	input := &ssm.StartSessionInput{Target: aws.String(instanceID)}
	if command != "" {
		input.DocumentName = aws.String(nonInteractiveCommandDocument)
		input.Parameters = map[string][]string{"command": {command}}
	}
	return input
}

func buildPluginArgs(output *ssm.StartSessionOutput, input *ssm.StartSessionInput, region, profile, endpoint string) ([]string, error) {
	respJSON, err := json.Marshal(struct {
		SessionId  string `json:"SessionId"`
		TokenValue string `json:"TokenValue"`
		StreamUrl  string `json:"StreamUrl"`
	}{
		SessionId:  aws.ToString(output.SessionId),
		TokenValue: aws.ToString(output.TokenValue),
		StreamUrl:  aws.ToString(output.StreamUrl),
	})
	if err != nil {
		return nil, fmt.Errorf("marshal start session response: %w", err)
	}

	params := struct {
		Target       string              `json:"Target"`
		DocumentName string              `json:"DocumentName,omitempty"`
		Parameters   map[string][]string `json:"Parameters,omitempty"`
	}{
		Target:       aws.ToString(input.Target),
		DocumentName: aws.ToString(input.DocumentName),
		Parameters:   input.Parameters,
	}
	paramsJSON, err := json.Marshal(params)
	if err != nil {
		return nil, fmt.Errorf("marshal start session params: %w", err)
	}

	return []string{
		string(respJSON),
		region,
		"StartSession",
		profile,
		string(paramsJSON),
		endpoint,
	}, nil
}
