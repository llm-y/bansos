package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func TestGetNextKey(t *testing.T) {
	keys := []string{"key1", "key2", "key3", "key4"}

	tests := []struct {
		name       string
		currentKey string
		expected   string
	}{
		{"first key returns second", "key1", "key2"},
		{"middle key returns next", "key2", "key3"},
		{"last key wraps to first", "key4", "key1"},
		{"unknown key returns first", "unknown", "key1"},
		{"empty key returns first", "", "key1"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := getNextKey(tt.currentKey, keys)
			if result != tt.expected {
				t.Errorf("getNextKey(%q, keys) = %q, want %q", tt.currentKey, result, tt.expected)
			}
		})
	}
}

func TestReadKeysList(t *testing.T) {
	// Create a temporary bansos.txt
	tmpDir := t.TempDir()
	bansosPath := filepath.Join(tmpDir, "bansos.txt")

	content := "key1\nkey2\n\nkey3\n  key4  \n"
	if err := os.WriteFile(bansosPath, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	keys, err := readKeysList(bansosPath)
	if err != nil {
		t.Fatalf("readKeysList() error: %v", err)
	}

	expected := []string{"key1", "key2", "key3", "key4"}
	if len(keys) != len(expected) {
		t.Fatalf("got %d keys, want %d", len(keys), len(expected))
	}

	for i, key := range keys {
		if key != expected[i] {
			t.Errorf("keys[%d] = %q, want %q", i, key, expected[i])
		}
	}
}

func TestReadCurrentKey(t *testing.T) {
	tmpDir := t.TempDir()
	settingsPath := filepath.Join(tmpDir, "settings.json")

	settings := Settings{
		Env: map[string]string{
			"ANTHROPIC_API_KEY":  "test-api-key-123",
			"ANTHROPIC_BASE_URL": "https://cc.freemodel.dev",
		},
		Permissions: Permissions{
			Allow: []string{"Bash", "Write", "Edit"},
			Deny:  []string{},
		},
		APIKeyHelper: "echo 'test-api-key-123'",
	}

	data, _ := json.MarshalIndent(settings, "", "    ")
	if err := os.WriteFile(settingsPath, data, 0644); err != nil {
		t.Fatal(err)
	}

	key, err := readCurrentKey(settingsPath)
	if err != nil {
		t.Fatalf("readCurrentKey() error: %v", err)
	}

	if key != "test-api-key-123" {
		t.Errorf("got %q, want %q", key, "test-api-key-123")
	}
}

func TestGenerateSettings(t *testing.T) {
	apiKey := "sk-ant-test123"
	settings := generateSettings(apiKey)

	if settings.Env["ANTHROPIC_API_KEY"] != apiKey {
		t.Errorf("ANTHROPIC_API_KEY = %q, want %q", settings.Env["ANTHROPIC_API_KEY"], apiKey)
	}

	if settings.Env["ANTHROPIC_BASE_URL"] != "https://cc.freemodel.dev" {
		t.Errorf("ANTHROPIC_BASE_URL = %q, want %q", settings.Env["ANTHROPIC_BASE_URL"], "https://cc.freemodel.dev")
	}

	if settings.Env["CLAUDE_CODE_DISABLE_NONESSENTIAL_TRAFFIC"] != "1" {
		t.Errorf("CLAUDE_CODE_DISABLE_NONESSENTIAL_TRAFFIC = %q, want %q", settings.Env["CLAUDE_CODE_DISABLE_NONESSENTIAL_TRAFFIC"], "1")
	}

	expectedHelper := "echo 'sk-ant-test123'"
	if settings.APIKeyHelper != expectedHelper {
		t.Errorf("APIKeyHelper = %q, want %q", settings.APIKeyHelper, expectedHelper)
	}

	if len(settings.Permissions.Allow) != 3 {
		t.Errorf("Allow has %d items, want 3", len(settings.Permissions.Allow))
	}

	if len(settings.Permissions.Deny) != 0 {
		t.Errorf("Deny has %d items, want 0", len(settings.Permissions.Deny))
	}
}

func TestWriteSettings(t *testing.T) {
	tmpDir := t.TempDir()
	settingsPath := filepath.Join(tmpDir, ".claude", "settings.json")

	settings := generateSettings("test-key")
	if err := writeSettings(settingsPath, settings); err != nil {
		t.Fatalf("writeSettings() error: %v", err)
	}

	// Verify file was written
	data, err := os.ReadFile(settingsPath)
	if err != nil {
		t.Fatalf("failed to read written file: %v", err)
	}

	var readBack Settings
	if err := json.Unmarshal(data, &readBack); err != nil {
		t.Fatalf("failed to unmarshal written file: %v", err)
	}

	if readBack.Env["ANTHROPIC_API_KEY"] != "test-key" {
		t.Errorf("ANTHROPIC_API_KEY = %q, want %q", readBack.Env["ANTHROPIC_API_KEY"], "test-key")
	}
}

func TestEndToEnd(t *testing.T) {
	tmpDir := t.TempDir()

	// Create bansos.txt
	bansosPath := filepath.Join(tmpDir, "bansos.txt")
	bansosContent := "key-alpha\nkey-beta\nkey-gamma\n"
	if err := os.WriteFile(bansosPath, []byte(bansosContent), 0644); err != nil {
		t.Fatal(err)
	}

	// Create initial settings.json with key-beta
	claudeDir := filepath.Join(tmpDir, ".claude")
	if err := os.MkdirAll(claudeDir, 0755); err != nil {
		t.Fatal(err)
	}

	settingsPath := filepath.Join(claudeDir, "settings.json")
	initialSettings := generateSettings("key-beta")
	data, _ := json.MarshalIndent(initialSettings, "", "    ")
	if err := os.WriteFile(settingsPath, data, 0644); err != nil {
		t.Fatal(err)
	}

	// Read current key
	currentKey, err := readCurrentKey(settingsPath)
	if err != nil {
		t.Fatalf("readCurrentKey() error: %v", err)
	}

	if currentKey != "key-beta" {
		t.Errorf("currentKey = %q, want %q", currentKey, "key-beta")
	}

	// Read keys list
	keys, err := readKeysList(bansosPath)
	if err != nil {
		t.Fatalf("readKeysList() error: %v", err)
	}

	// Get next key (should be key-gamma since key-beta is current)
	nextKey := getNextKey(currentKey, keys)
	if nextKey != "key-gamma" {
		t.Errorf("nextKey = %q, want %q", nextKey, "key-gamma")
	}

	// Write new settings
	newSettings := generateSettings(nextKey)
	if err := writeSettings(settingsPath, newSettings); err != nil {
		t.Fatalf("writeSettings() error: %v", err)
	}

	// Verify
	finalKey, err := readCurrentKey(settingsPath)
	if err != nil {
		t.Fatalf("readCurrentKey() after write error: %v", err)
	}

	if finalKey != "key-gamma" {
		t.Errorf("finalKey = %q, want %q", finalKey, "key-gamma")
	}
}
