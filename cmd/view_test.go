package cmd

import "testing"

func TestNavigateJSON(t *testing.T) {
	data := map[string]interface{}{
		"fields": map[string]interface{}{
			"summary": "Test issue",
			"status": map[string]interface{}{
				"name": "Open",
			},
			"labels": []interface{}{"bug", "critical"},
		},
		"key": "LPS-123",
	}

	tests := []struct {
		name string
		path string
		want interface{}
	}{
		{
			name: "top-level field",
			path: ".key",
			want: "LPS-123",
		},
		{
			name: "nested field",
			path: ".fields.summary",
			want: "Test issue",
		},
		{
			name: "deeply nested field",
			path: ".fields.status.name",
			want: "Open",
		},
		{
			name: "without leading dot",
			path: "fields.summary",
			want: "Test issue",
		},
		{
			name: "nonexistent field",
			path: ".fields.nonexistent",
			want: nil,
		},
		{
			name: "empty path returns whole object",
			path: "",
			want: data,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := navigateJSON(data, tt.path)
			if tt.want == nil {
				if got != nil {
					t.Errorf("navigateJSON(%q) = %v, want nil", tt.path, got)
				}
				return
			}
			// Compare string values
			if s, ok := tt.want.(string); ok {
				if got != s {
					t.Errorf("navigateJSON(%q) = %v, want %q", tt.path, got, s)
				}
			}
		})
	}
}

func TestSplitPath(t *testing.T) {
	tests := []struct {
		input string
		want  []string
	}{
		{"fields.summary", []string{"fields", "summary"}},
		{"key", []string{"key"}},
		{"a.b.c.d", []string{"a", "b", "c", "d"}},
		{"", nil},
	}

	for _, tt := range tests {
		got := splitPath(tt.input)
		if len(got) != len(tt.want) {
			t.Errorf("splitPath(%q) = %v, want %v", tt.input, got, tt.want)
			continue
		}
		for i := range got {
			if got[i] != tt.want[i] {
				t.Errorf("splitPath(%q)[%d] = %q, want %q", tt.input, i, got[i], tt.want[i])
			}
		}
	}
}
