package git

import (
	"os/exec"
	"regexp"
	"strings"
)

var ticketRegex = regexp.MustCompile(`([A-Z][A-Z0-9]+-[0-9]+)`)

// ExtractTicket extracts a Jira ticket ID from the current git branch name.
func ExtractTicket() (string, error) {
	out, err := exec.Command("git", "rev-parse", "--abbrev-ref", "HEAD").Output()
	if err != nil {
		return "", err
	}

	branch := strings.TrimSpace(string(out))
	return ExtractTicketFromString(branch)
}

// ExtractTicketFromString extracts a Jira ticket ID from any string.
func ExtractTicketFromString(s string) (string, error) {
	match := ticketRegex.FindString(s)
	if match == "" {
		return "", &NoTicketError{Input: s}
	}
	return match, nil
}

// IsGitRepo returns true if the current directory is inside a git repository.
func IsGitRepo() bool {
	err := exec.Command("git", "rev-parse", "HEAD").Run()
	return err == nil
}

type NoTicketError struct {
	Input string
}

func (e *NoTicketError) Error() string {
	return "no Jira ticket found in: " + e.Input
}
