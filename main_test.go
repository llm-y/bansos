package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func TestGetNextEntry(t *testing.T) {
	entries := []KeyEntry{
		{ID: "anu", Key: "key1"},
		{ID: "awg", Key: "key2"},
		{ID: "bukped", Key: "key3"},
		{ID: "4", Key: "key4"},
	}

	tests := []struct {
		name            string
		currentKey      string
		expectedNextKey string
		expectedNextID  string
		expectCurrent   bool
	}{
		{"first key returns second", "key1", "key2", "awg", true},
		{"middle key returns next", "key2", "key3", "bukped", true},
		{"last key wraps to first", "key4", "key1", "anu", true},
		{"unknown key returns first", "unknown", "key1", "anu", false},
		{"empty key returns first", "", "key1", "anu", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			currentEntry, nextEntry := getNextEntry(tt.currentKey, entries)
			if nextEntry.Key != tt.expectedNextKey {
				t.Errorf("getNextEntry(%q, entries).nextKey = %q, want %q", tt.currentKey, nextEntry.Key, tt.expectedNextKey)
			}
			if nextEntry.ID != tt.expectedNextID {
				t.Errorf("getNextEntry(%q, entries).nextID = %q, want %q", tt.currentKey, nextEntry.ID, tt.expectedNextID)
			}
			if tt.expectCurrent && currentEntry == nil {
				t.Errorf("getNextEntry(%q, entries).currentEntry = nil, want non-nil", tt.currentKey)
			}
			if !tt.expectCurrent && currentEntry != nil {
				t.Errorf("getNextEntry(%q, entries).currentEntry = %v, want nil", tt.currentKey, currentEntry)
			}
		})
	}
}

func TestReadKeysList(t *testing.T) {
	// Create a temporary bansos.csv
	tmpDir := t.TempDir()
	bansosPath := filepath.Join(tmpDir, "bansos.csv")

	content := "id,key\nanu,key1\nawg,key2\nbukped,key3\n4,  key4  \n"
	if err := os.WriteFile(bansosPath, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	entries, err := readKeysList(bansosPath)
	if err != nil {
		t.Fatalf("readKeysList() error: %v", err)
	}

	expected := []KeyEntry{
		{ID: "anu", Key: "key1"},
		{ID: "awg", Key: "key2"},
		{ID: "bukped", Key: "key3"},
		{ID: "4", Key: "key4"},
	}
	if len(entries) != len(expected) {
		t.Fatalf("got %d entries, want %d", len(entries), len(expected))
	}

	for i, entry := range entries {
		if entry.ID != expected[i].ID {
			t.Errorf("entries[%d].ID = %q, want %q", i, entry.ID, expected[i].ID)
		}
		if entry.Key != expected[i].Key {
			t.Errorf("entries[%d].Key = %q, want %q", i, entry.Key, expected[i].Key)
		}
	}
}

func TestReadKeysListInvalidHeader(t *testing.T) {
	tmpDir := t.TempDir()
	bansosPath := filepath.Join(tmpDir, "bansos.csv")

	content := "nama,value\n1,key1\n"
	if err := os.WriteFile(bansosPath, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	_, err := readKeysList(bansosPath)
	if err == nil {
		t.Fatal("readKeysList() expected error for invalid header, got nil")
	}

	if !contains(err.Error(), "header harus berformat") {
		t.Errorf("unexpected error message: %v", err)
	}
}

func TestReadKeysListEmpty(t *testing.T) {
	tmpDir := t.TempDir()
	bansosPath := filepath.Join(tmpDir, "bansos.csv")

	content := "id,key\n"
	if err := os.WriteFile(bansosPath, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	_, err := readKeysList(bansosPath)
	if err == nil {
		t.Fatal("readKeysList() expected error for empty file, got nil")
	}

	if !contains(err.Error(), "bansos.csv is empty") {
		t.Errorf("unexpected error message: %v", err)
	}
}

func TestMaskKey(t *testing.T) {
	tests := []struct {
		name     string
		key      string
		expected string
	}{
		{"long key", "sk-ant-api03-key1xxxxxxxx", "sk-ant-a...xxxx"},
		{"short key", "shortkey", "shortkey"},
		{"exactly 12 chars", "123456789012", "123456789012"},
		{"13 chars", "1234567890123", "12345678...0123"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := maskKey(tt.key)
			if result != tt.expected {
				t.Errorf("maskKey(%q) = %q, want %q", tt.key, result, tt.expected)
			}
		})
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

	// Create bansos.csv
	bansosPath := filepath.Join(tmpDir, "bansos.csv")
	bansosContent := "id,key\nanu,key-alpha\nawg,key-beta\nbukped,key-gamma\n"
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
	entries, err := readKeysList(bansosPath)
	if err != nil {
		t.Fatalf("readKeysList() error: %v", err)
	}

	// Get next entry (should be key-gamma since key-beta is current)
	currentEntry, nextEntry := getNextEntry(currentKey, entries)
	if currentEntry == nil {
		t.Fatal("currentEntry should not be nil")
	}
	if currentEntry.ID != "awg" {
		t.Errorf("currentEntry.ID = %q, want %q", currentEntry.ID, "awg")
	}
	if nextEntry.Key != "key-gamma" {
		t.Errorf("nextEntry.Key = %q, want %q", nextEntry.Key, "key-gamma")
	}
	if nextEntry.ID != "bukped" {
		t.Errorf("nextEntry.ID = %q, want %q", nextEntry.ID, "bukped")
	}

	// Write new settings
	newSettings := generateSettings(nextEntry.Key)
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

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > 0 && containsSubstring(s, substr))
}

func containsSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
