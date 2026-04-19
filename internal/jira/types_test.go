package jira

import "testing"

func TestParseJiraTime(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
		wantY   int
		wantM   int
		wantD   int
	}{
		{
			name:  "standard Jira timestamp",
			input: "2024-03-15T10:30:00.000+0000",
			wantY: 2024, wantM: 3, wantD: 15,
		},
		{
			name:  "with timezone offset",
			input: "2023-12-25T08:00:00.000-0500",
			wantY: 2023, wantM: 12, wantD: 25,
		},
		{
			name:    "invalid format",
			input:   "not-a-date",
			wantErr: true,
		},
		{
			name:    "empty string",
			input:   "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseJiraTime(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseJiraTime(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if got.Year() != tt.wantY || int(got.Month()) != tt.wantM || got.Day() != tt.wantD {
					t.Errorf("ParseJiraTime(%q) = %v, want %d-%d-%d", tt.input, got, tt.wantY, tt.wantM, tt.wantD)
				}
			}
		})
	}
}

func TestErrorResponse_String(t *testing.T) {
	tests := []struct {
		name string
		resp ErrorResponse
		want string
	}{
		{
			name: "single error message",
			resp: ErrorResponse{
				ErrorMessages: []string{"Issue does not exist"},
			},
			want: "Issue does not exist",
		},
		{
			name: "field errors only",
			resp: ErrorResponse{
				Errors: map[string]string{
					"summary": "Summary is required",
				},
			},
			want: "summary: Summary is required",
		},
		{
			name: "both error messages and field errors",
			resp: ErrorResponse{
				ErrorMessages: []string{"Validation failed"},
				Errors: map[string]string{
					"summary": "cannot be empty",
				},
			},
			want: "Validation failed\nsummary: cannot be empty",
		},
		{
			name: "empty error response",
			resp: ErrorResponse{},
			want: "unknown error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.resp.String()
			if got != tt.want {
				t.Errorf("ErrorResponse.String() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestMakeADFBody(t *testing.T) {
	result := MakeADFBody("Hello world")

	doc, ok := result.(map[string]interface{})
	if !ok {
		t.Fatal("MakeADFBody should return a map")
	}

	if doc["version"] != 1 {
		t.Errorf("version = %v, want 1", doc["version"])
	}
	if doc["type"] != "doc" {
		t.Errorf("type = %v, want doc", doc["type"])
	}

	content, ok := doc["content"].([]interface{})
	if !ok || len(content) != 1 {
		t.Fatal("content should have 1 element")
	}

	para, ok := content[0].(map[string]interface{})
	if !ok {
		t.Fatal("first content element should be a map")
	}
	if para["type"] != "paragraph" {
		t.Errorf("paragraph type = %v, want paragraph", para["type"])
	}

	inline, ok := para["content"].([]interface{})
	if !ok || len(inline) != 1 {
		t.Fatal("paragraph content should have 1 element")
	}

	textNode, ok := inline[0].(map[string]interface{})
	if !ok {
		t.Fatal("inline element should be a map")
	}
	if textNode["type"] != "text" {
		t.Errorf("text node type = %v, want text", textNode["type"])
	}
	if textNode["text"] != "Hello world" {
		t.Errorf("text = %v, want Hello world", textNode["text"])
	}
}
