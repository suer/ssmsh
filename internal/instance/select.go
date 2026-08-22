package instance

import (
	"fmt"
	"os"
	"strings"

	"github.com/manifoldco/promptui"
	"golang.org/x/term"
)

type AmbiguousError struct {
	Name       string
	Candidates []Candidate
}

func (e *AmbiguousError) Error() string {
	var b strings.Builder
	fmt.Fprintf(&b, "multiple instances matched %q:\n", e.Name)
	for _, c := range e.Candidates {
		fmt.Fprintf(&b, "  %s (%s)\n", c.InstanceID, c.Name)
	}
	return b.String()
}

var isInteractiveTerminal = func() bool {
	return term.IsTerminal(int(os.Stdin.Fd())) && term.IsTerminal(int(os.Stdout.Fd()))
}

func InteractiveSelector() Selector {
	return func(candidates []Candidate) (string, error) {
		if !isInteractiveTerminal() {
			return "", &AmbiguousError{Candidates: candidates}
		}

		searcher := func(input string, index int) bool {
			c := candidates[index]
			target := strings.ToLower(c.Name + " " + c.InstanceID)
			return strings.Contains(target, strings.ToLower(input))
		}

		prompt := promptui.Select{
			Label: fmt.Sprintf("Multiple instances matched (%d)", len(candidates)),
			Items: candidates,
			Templates: &promptui.SelectTemplates{
				Label:    "{{ . }}",
				Active:   "\U0001F449 {{ .Name | cyan }} ({{ .InstanceID }})",
				Inactive: "  {{ .Name }} ({{ .InstanceID }})",
				Selected: "Selected: {{ .Name }} ({{ .InstanceID }})",
			},
			Searcher: searcher,
			Size:     10,
		}

		idx, _, err := prompt.Run()
		if err != nil {
			return "", fmt.Errorf("selection cancelled: %w", err)
		}
		return candidates[idx].InstanceID, nil
	}
}
