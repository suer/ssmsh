package session

import (
	"errors"
	"testing"
)

func TestCheckPluginInstalled_Found(t *testing.T) {
	original := lookPath
	lookPath = func(file string) (string, error) { return "/usr/local/bin/" + file, nil }
	defer func() { lookPath = original }()

	if err := CheckPluginInstalled(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestCheckPluginInstalled_NotFound(t *testing.T) {
	original := lookPath
	lookPath = func(file string) (string, error) { return "", errors.New("not found") }
	defer func() { lookPath = original }()

	err := CheckPluginInstalled()
	if err == nil {
		t.Fatal("expected error")
	}
}
