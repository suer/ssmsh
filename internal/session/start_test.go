package session

import (
	"encoding/json"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
)

func TestBuildStartSessionInput_Login(t *testing.T) {
	input := buildStartSessionInput("i-05239d460d79ad389", "", 0, 0)

	if aws.ToString(input.Target) != "i-05239d460d79ad389" {
		t.Fatalf("got Target %q, want %q", aws.ToString(input.Target), "i-05239d460d79ad389")
	}
	if input.DocumentName != nil {
		t.Fatalf("got DocumentName %q, want nil", aws.ToString(input.DocumentName))
	}
	if input.Parameters != nil {
		t.Fatalf("got Parameters %v, want nil", input.Parameters)
	}
}

func TestBuildStartSessionInput_Command(t *testing.T) {
	input := buildStartSessionInput("i-05239d460d79ad389", "echo hi", 0, 0)

	if aws.ToString(input.DocumentName) != "AWS-StartNonInteractiveCommand" {
		t.Fatalf("got DocumentName %q, want %q", aws.ToString(input.DocumentName), "AWS-StartNonInteractiveCommand")
	}
	got := input.Parameters["command"]
	if len(got) != 1 || got[0] != "echo hi" {
		t.Fatalf("got Parameters[command] %v, want [echo hi]", got)
	}
}

func TestBuildStartSessionInput_PortForward(t *testing.T) {
	input := buildStartSessionInput("i-05239d460d79ad389", "", 10080, 80)

	if aws.ToString(input.DocumentName) != "AWS-StartPortForwardingSession" {
		t.Fatalf("got DocumentName %q, want %q", aws.ToString(input.DocumentName), "AWS-StartPortForwardingSession")
	}
	if got := input.Parameters["localPortNumber"]; len(got) != 1 || got[0] != "10080" {
		t.Fatalf("got Parameters[localPortNumber] %v, want [10080]", got)
	}
	if got := input.Parameters["portNumber"]; len(got) != 1 || got[0] != "80" {
		t.Fatalf("got Parameters[portNumber] %v, want [80]", got)
	}
}

func TestBuildPluginArgs(t *testing.T) {
	input := buildStartSessionInput("i-05239d460d79ad389", "echo hi", 0, 0)
	output := &ssm.StartSessionOutput{
		SessionId:  aws.String("session-123"),
		TokenValue: aws.String("token-abc"),
		StreamUrl:  aws.String("wss://example.com/stream"),
	}

	args, err := buildPluginArgs(output, input, "ap-northeast-1", "suer", "https://ssm.ap-northeast-1.amazonaws.com")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(args) != 6 {
		t.Fatalf("got %d args, want 6: %v", len(args), args)
	}

	var resp struct {
		SessionId  string
		TokenValue string
		StreamUrl  string
	}
	if err := json.Unmarshal([]byte(args[0]), &resp); err != nil {
		t.Fatalf("failed to unmarshal response json: %v", err)
	}
	if resp.SessionId != "session-123" || resp.TokenValue != "token-abc" || resp.StreamUrl != "wss://example.com/stream" {
		t.Fatalf("unexpected response json: %+v", resp)
	}

	if args[1] != "ap-northeast-1" {
		t.Fatalf("got region arg %q, want %q", args[1], "ap-northeast-1")
	}
	if args[2] != "StartSession" {
		t.Fatalf("got operation arg %q, want %q", args[2], "StartSession")
	}
	if args[3] != "suer" {
		t.Fatalf("got profile arg %q, want %q", args[3], "suer")
	}

	var params struct {
		Target       string
		DocumentName string
		Parameters   map[string][]string
	}
	if err := json.Unmarshal([]byte(args[4]), &params); err != nil {
		t.Fatalf("failed to unmarshal params json: %v", err)
	}
	if params.Target != "i-05239d460d79ad389" {
		t.Fatalf("got Target %q, want %q", params.Target, "i-05239d460d79ad389")
	}
	if params.DocumentName != "AWS-StartNonInteractiveCommand" {
		t.Fatalf("got DocumentName %q, want %q", params.DocumentName, "AWS-StartNonInteractiveCommand")
	}
	if len(params.Parameters["command"]) != 1 || params.Parameters["command"][0] != "echo hi" {
		t.Fatalf("got Parameters %v, want command=[echo hi]", params.Parameters)
	}

	if args[5] != "https://ssm.ap-northeast-1.amazonaws.com" {
		t.Fatalf("got endpoint arg %q, want %q", args[5], "https://ssm.ap-northeast-1.amazonaws.com")
	}
}

func TestBuildPluginArgs_EmptyProfile(t *testing.T) {
	input := buildStartSessionInput("i-05239d460d79ad389", "", 0, 0)
	output := &ssm.StartSessionOutput{SessionId: aws.String("s"), TokenValue: aws.String("t"), StreamUrl: aws.String("u")}

	args, err := buildPluginArgs(output, input, "ap-northeast-1", "", "https://ssm.ap-northeast-1.amazonaws.com")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if args[3] != "" {
		t.Fatalf("got profile arg %q, want empty string", args[3])
	}
}
