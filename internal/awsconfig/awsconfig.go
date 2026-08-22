package awsconfig

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
)

func Load(ctx context.Context, profile, region string) (aws.Config, error) {
	var optFns []func(*config.LoadOptions) error
	if profile != "" {
		optFns = append(optFns, config.WithSharedConfigProfile(profile))
	}
	if region != "" {
		optFns = append(optFns, config.WithRegion(region))
	}
	return config.LoadDefaultConfig(ctx, optFns...)
}
