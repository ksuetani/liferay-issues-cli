package git

import "testing"

func TestExtractTicketFromString(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    string
		wantErr bool
	}{
		{
			name:  "simple ticket",
			input: "LPS-12345",
			want:  "LPS-12345",
		},
		{
			name:  "ticket in branch name",
			input: "LPS-12345-fix-login-bug",
			want:  "LPS-12345",
		},
		{
			name:  "ticket with feature prefix",
			input: "feature/LPS-99999-add-new-api",
			want:  "LPS-99999",
		},
		{
			name:  "ticket at end of string",
			input: "some-work-on-GROW-42",
			want:  "GROW-42",
		},
		{
			name:  "multi-letter project key",
			input: "LRQA-100",
			want:  "LRQA-100",
		},
		{
			name:  "alphanumeric project key",
			input: "LPS3-500",
			want:  "LPS3-500",
		},
		{
			name:    "no ticket in string",
			input:   "main",
			wantErr: true,
		},
		{
			name:    "empty string",
			input:   "",
			wantErr: true,
		},
		{
			name:    "lowercase letters only",
			input:   "fix-login-bug",
			wantErr: true,
		},
		{
			name:    "numbers but no project key",
			input:   "12345",
			wantErr: true,
		},
		{
			name:    "single letter project key (invalid)",
			input:   "A-123",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ExtractTicketFromString(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ExtractTicketFromString(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ExtractTicketFromString(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

func TestNoTicketError(t *testing.T) {
	err := &NoTicketError{Input: "main"}
	if err.Error() != "no Jira ticket found in: main" {
		t.Errorf("unexpected error message: %s", err.Error())
	}
}
