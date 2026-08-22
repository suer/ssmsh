package awsconfig

import (
	"context"
	"testing"
)

func TestLoad_RegionOverride(t *testing.T) {
	cfg, err := Load(context.Background(), "", "ap-northeast-1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.Region != "ap-northeast-1" {
		t.Fatalf("got region %q, want %q", cfg.Region, "ap-northeast-1")
	}
}
