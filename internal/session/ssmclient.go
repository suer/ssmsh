package session

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/ssm"
)

type SSMAPI interface {
	StartSession(ctx context.Context, params *ssm.StartSessionInput, optFns ...func(*ssm.Options)) (*ssm.StartSessionOutput, error)
}
