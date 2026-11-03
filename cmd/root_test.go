package cmd

import (
	"bytes"
	"strings"
	"testing"
)

func TestVersionFlag(t *testing.T) {
	for _, flag := range []string{"-v", "--version"} {
		cmd := newRootCmd()
		var out bytes.Buffer
		cmd.SetOut(&out)
		cmd.SetArgs([]string{flag})

		if err := cmd.Execute(); err != nil {
			t.Fatalf("%s: unexpected error: %v", flag, err)
		}
		if !strings.Contains(out.String(), version()) {
			t.Fatalf("%s: got output %q, want it to contain %q", flag, out.String(), version())
		}
	}
}

func TestParseLocalForward(t *testing.T) {
	local, remote, err := parseLocalForward("10080:80")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if local != 10080 || remote != 80 {
		t.Fatalf("got local=%d remote=%d, want local=10080 remote=80", local, remote)
	}
}

func TestParseLocalForward_Invalid(t *testing.T) {
	cases := []string{"", "10080", "10080:", ":80", "abc:80", "10080:abc", "0:80", "10080:0", "70000:80"}
	for _, spec := range cases {
		if _, _, err := parseLocalForward(spec); err == nil {
			t.Fatalf("parseLocalForward(%q): expected error, got nil", spec)
		}
	}
}

func TestUsageShownOnInvalidInput(t *testing.T) {
	cases := map[string][]string{
		"no args":        {},
		"unknown flag":   {"--unknown", "web"},
		"conflict flags": {"web", "-c", "ls", "-L", "1:2"},
		"bad forward":    {"web", "-L", "abc"},
	}
	for name, args := range cases {
		t.Run(name, func(t *testing.T) {
			cmd := newRootCmd()
			var out bytes.Buffer
			cmd.SetOut(&out)
			cmd.SetErr(&out)
			cmd.SetArgs(args)

			if err := cmd.Execute(); err == nil {
				t.Fatal("expected error, got nil")
			}
			if !strings.Contains(out.String(), "Usage:") {
				t.Fatalf("expected usage in output, got %q", out.String())
			}
		})
	}
}
