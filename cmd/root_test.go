package cmd

import "testing"

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
