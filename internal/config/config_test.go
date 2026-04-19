package config

import (
	"net/http"
	"os"
	"path/filepath"
	"testing"
)

func TestResolveAuth_EnvVars(t *testing.T) {
	t.Setenv("JIRA_API_TOKEN", "test-token")
	t.Setenv("JIRA_USER", "test@example.com")

	cfg := &Config{
		Jira: JiraConfig{Instance: "test.atlassian.net"},
	}

	header, err := ResolveAuth(cfg)
	if err != nil {
		t.Fatalf("ResolveAuth() error = %v", err)
	}
	if header.Get("Authorization") == "" {
		t.Error("expected Authorization header")
	}
}

func TestResolveAuth_ConfigFile(t *testing.T) {
	// Ensure env vars are not set
	t.Setenv("JIRA_API_TOKEN", "")
	t.Setenv("JIRA_USER", "")

	cfg := &Config{
		Jira: JiraConfig{Instance: "test.atlassian.net"},
		Auth: AuthConfig{
			Email: "user@example.com",
			Token: "config-token",
		},
	}

	header, err := ResolveAuth(cfg)
	if err != nil {
		t.Fatalf("ResolveAuth() error = %v", err)
	}
	if header.Get("Authorization") == "" {
		t.Error("expected Authorization header")
	}
}

func TestResolveAuth_Netrc(t *testing.T) {
	t.Setenv("JIRA_API_TOKEN", "")
	t.Setenv("JIRA_USER", "")

	// Create a temporary .netrc file
	tmpDir := t.TempDir()
	t.Setenv("HOME", tmpDir)

	netrcContent := "machine test.atlassian.net\nlogin user@example.com\npassword netrc-token\n"
	if err := os.WriteFile(filepath.Join(tmpDir, ".netrc"), []byte(netrcContent), 0600); err != nil {
		t.Fatalf("writing .netrc: %v", err)
	}

	cfg := &Config{
		Jira: JiraConfig{Instance: "test.atlassian.net"},
	}

	header, err := ResolveAuth(cfg)
	if err != nil {
		t.Fatalf("ResolveAuth() error = %v", err)
	}
	if header.Get("Authorization") == "" {
		t.Error("expected Authorization header from .netrc")
	}
}

func TestResolveAuth_NoCredentials(t *testing.T) {
	t.Setenv("JIRA_API_TOKEN", "")
	t.Setenv("JIRA_USER", "")

	// Point HOME to empty dir so .netrc won't be found
	tmpDir := t.TempDir()
	t.Setenv("HOME", tmpDir)

	cfg := &Config{
		Jira: JiraConfig{Instance: "test.atlassian.net"},
	}

	_, err := ResolveAuth(cfg)
	if err == nil {
		t.Fatal("expected error when no credentials found")
	}
}

func TestResolveAuth_EnvVarsPriority(t *testing.T) {
	// Both env vars and config are set; env vars should win
	t.Setenv("JIRA_API_TOKEN", "env-token")
	t.Setenv("JIRA_USER", "env@example.com")

	cfg := &Config{
		Jira: JiraConfig{Instance: "test.atlassian.net"},
		Auth: AuthConfig{
			Email: "config@example.com",
			Token: "config-token",
		},
	}

	header, err := ResolveAuth(cfg)
	if err != nil {
		t.Fatalf("ResolveAuth() error = %v", err)
	}

	// The header should use env var credentials (env@example.com:env-token)
	expected := basicAuthHeader("env@example.com", "env-token")
	if header.Get("Authorization") != expected.Get("Authorization") {
		t.Error("env vars should take priority over config file")
	}
}

func TestBasicAuthHeader(t *testing.T) {
	header := basicAuthHeader("user@example.com", "token123")

	auth := header.Get("Authorization")
	if auth == "" {
		t.Fatal("expected Authorization header")
	}
	if len(auth) < 6 || auth[:6] != "Basic " {
		t.Errorf("Authorization should start with 'Basic ', got %q", auth)
	}
}

func TestConfigFilePath(t *testing.T) {
	path := ConfigFilePath()
	if path == "" {
		t.Error("ConfigFilePath() should not return empty string")
	}
	if filepath.Base(path) != "config.yaml" {
		t.Errorf("config file should be named config.yaml, got %q", filepath.Base(path))
	}
}

func TestConfigFilePath_XDG(t *testing.T) {
	tmpDir := t.TempDir()
	t.Setenv("XDG_CONFIG_HOME", tmpDir)

	path := ConfigFilePath()
	expected := filepath.Join(tmpDir, "issues", "config.yaml")
	if path != expected {
		t.Errorf("ConfigFilePath() = %q, want %q", path, expected)
	}
}

func TestEnsureConfigDir(t *testing.T) {
	tmpDir := t.TempDir()
	t.Setenv("XDG_CONFIG_HOME", tmpDir)

	if err := EnsureConfigDir(); err != nil {
		t.Fatalf("EnsureConfigDir() error = %v", err)
	}

	dir := filepath.Join(tmpDir, "issues")
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		t.Error("config directory was not created")
	}
}

// Verify the header format is compatible with what Jira expects
func TestBasicAuthHeader_Format(t *testing.T) {
	h := basicAuthHeader("test@example.com", "api-token")

	// Should produce a valid HTTP header
	req := &http.Request{Header: h}
	user, pass, ok := req.BasicAuth()
	if !ok {
		t.Fatal("header should be valid basic auth")
	}
	if user != "test@example.com" {
		t.Errorf("user = %q, want 'test@example.com'", user)
	}
	if pass != "api-token" {
		t.Errorf("pass = %q, want 'api-token'", pass)
	}
}
