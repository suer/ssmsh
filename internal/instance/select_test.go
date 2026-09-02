package instance

import (
	"errors"
	"testing"
)

func TestInteractiveSelector_NonTTY_ReturnsAmbiguousError(t *testing.T) {
	original := isInteractiveTerminal
	isInteractiveTerminal = func() bool { return false }
	defer func() { isInteractiveTerminal = original }()

	candidates := []Candidate{
		{InstanceID: "i-aaa", Name: "dup"},
		{InstanceID: "i-bbb", Name: "dup"},
	}

	_, err := InteractiveSelector()("dup", candidates)

	var ambiguous *AmbiguousError
	if !errors.As(err, &ambiguous) {
		t.Fatalf("got %v, want AmbiguousError", err)
	}
	if ambiguous.Name != "dup" {
		t.Fatalf("got Name %q, want %q", ambiguous.Name, "dup")
	}
	if len(ambiguous.Candidates) != 2 {
		t.Fatalf("got %d candidates, want 2", len(ambiguous.Candidates))
	}
}
