package ui

import "testing"

func TestExtractText_String(t *testing.T) {
	got := extractText("plain text")
	if got != "plain text" {
		t.Errorf("extractText(string) = %q, want 'plain text'", got)
	}
}

func TestExtractText_Nil(t *testing.T) {
	got := extractText(nil)
	if got != "" {
		t.Errorf("extractText(nil) = %q, want ''", got)
	}
}

func TestExtractText_ADF_Paragraph(t *testing.T) {
	adf := map[string]interface{}{
		"version": 1,
		"type":    "doc",
		"content": []interface{}{
			map[string]interface{}{
				"type": "paragraph",
				"content": []interface{}{
					map[string]interface{}{
						"type": "text",
						"text": "Hello world",
					},
				},
			},
		},
	}

	got := extractText(adf)
	if got != "Hello world" {
		t.Errorf("extractText(ADF paragraph) = %q, want 'Hello world'", got)
	}
}

func TestExtractText_ADF_MultipleParagraphs(t *testing.T) {
	adf := map[string]interface{}{
		"version": 1,
		"type":    "doc",
		"content": []interface{}{
			map[string]interface{}{
				"type": "paragraph",
				"content": []interface{}{
					map[string]interface{}{"type": "text", "text": "First paragraph"},
				},
			},
			map[string]interface{}{
				"type": "paragraph",
				"content": []interface{}{
					map[string]interface{}{"type": "text", "text": "Second paragraph"},
				},
			},
		},
	}

	got := extractText(adf)
	want := "First paragraph\nSecond paragraph"
	if got != want {
		t.Errorf("extractText(ADF multi-paragraph) = %q, want %q", got, want)
	}
}

func TestExtractText_ADF_BulletList(t *testing.T) {
	adf := map[string]interface{}{
		"version": 1,
		"type":    "doc",
		"content": []interface{}{
			map[string]interface{}{
				"type": "bulletList",
				"content": []interface{}{
					map[string]interface{}{
						"type": "listItem",
						"content": []interface{}{
							map[string]interface{}{
								"type": "paragraph",
								"content": []interface{}{
									map[string]interface{}{"type": "text", "text": "Item one"},
								},
							},
						},
					},
					map[string]interface{}{
						"type": "listItem",
						"content": []interface{}{
							map[string]interface{}{
								"type": "paragraph",
								"content": []interface{}{
									map[string]interface{}{"type": "text", "text": "Item two"},
								},
							},
						},
					},
				},
			},
		},
	}

	got := extractText(adf)
	if got == "" {
		t.Error("extractText(ADF bulletList) should not be empty")
	}
}

func TestExtractText_ADF_CodeBlock(t *testing.T) {
	adf := map[string]interface{}{
		"version": 1,
		"type":    "doc",
		"content": []interface{}{
			map[string]interface{}{
				"type": "codeBlock",
				"content": []interface{}{
					map[string]interface{}{"type": "text", "text": "fmt.Println(\"hi\")"},
				},
			},
		},
	}

	got := extractText(adf)
	want := "```\nfmt.Println(\"hi\")\n```"
	if got != want {
		t.Errorf("extractText(ADF codeBlock) = %q, want %q", got, want)
	}
}

func TestExtractText_ADF_InlineCard(t *testing.T) {
	adf := map[string]interface{}{
		"version": 1,
		"type":    "doc",
		"content": []interface{}{
			map[string]interface{}{
				"type": "paragraph",
				"content": []interface{}{
					map[string]interface{}{
						"type": "inlineCard",
						"attrs": map[string]interface{}{
							"url": "https://example.com",
						},
					},
				},
			},
		},
	}

	got := extractText(adf)
	if got != "https://example.com" {
		t.Errorf("extractText(ADF inlineCard) = %q, want 'https://example.com'", got)
	}
}

func TestExtractText_ADF_Mention(t *testing.T) {
	adf := map[string]interface{}{
		"version": 1,
		"type":    "doc",
		"content": []interface{}{
			map[string]interface{}{
				"type": "paragraph",
				"content": []interface{}{
					map[string]interface{}{
						"type": "mention",
						"attrs": map[string]interface{}{
							"text": "@johndoe",
						},
					},
				},
			},
		},
	}

	got := extractText(adf)
	if got != "@johndoe" {
		t.Errorf("extractText(ADF mention) = %q, want '@johndoe'", got)
	}
}
