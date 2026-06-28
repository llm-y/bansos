package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

// Settings represents the structure of ~/.claude/settings.json
type Settings struct {
	Env          map[string]string `json:"env"`
	Permissions  Permissions       `json:"permissions"`
	APIKeyHelper string            `json:"apiKeyHelper"`
}

// Permissions represents the permissions block in settings.json
type Permissions struct {
	Allow []string `json:"allow"`
	Deny  []string `json:"deny"`
}

func getHomeDir() (string, error) {
	if runtime.GOOS == "windows" {
		home := os.Getenv("USERPROFILE")
		if home == "" {
			home = os.Getenv("HOMEDRIVE") + os.Getenv("HOMEPATH")
		}
		if home == "" {
			return "", fmt.Errorf("cannot determine home directory on Windows")
		}
		return home, nil
	}
	home := os.Getenv("HOME")
	if home == "" {
		return "", fmt.Errorf("cannot determine home directory on Linux/macOS")
	}
	return home, nil
}

func readCurrentKey(settingsPath string) (string, error) {
	data, err := os.ReadFile(settingsPath)
	if err != nil {
		return "", err
	}

	var settings Settings
	if err := json.Unmarshal(data, &settings); err != nil {
		return "", fmt.Errorf("failed to parse settings.json: %w", err)
	}

	key, ok := settings.Env["ANTHROPIC_API_KEY"]
	if !ok || key == "" {
		return "", fmt.Errorf("ANTHROPIC_API_KEY not found in settings.json")
	}

	return key, nil
}

func readKeysList(bansosPath string) ([]string, error) {
	file, err := os.Open(bansosPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open bansos.txt: %w", err)
	}
	defer file.Close()

	var keys []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line != "" {
			keys = append(keys, line)
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("failed to read bansos.txt: %w", err)
	}

	if len(keys) == 0 {
		return nil, fmt.Errorf("bansos.txt is empty")
	}

	return keys, nil
}

func getNextKey(currentKey string, keys []string) string {
	for i, key := range keys {
		if key == currentKey {
			// If it's the last key, wrap around to the first
			if i == len(keys)-1 {
				return keys[0]
			}
			return keys[i+1]
		}
	}
	// If current key is not found in the list, use the first key
	return keys[0]
}

func generateSettings(apiKey string) Settings {
	return Settings{
		Env: map[string]string{
			"ANTHROPIC_API_KEY":                        apiKey,
			"ANTHROPIC_BASE_URL":                       "https://cc.freemodel.dev",
			"CLAUDE_CODE_DISABLE_NONESSENTIAL_TRAFFIC": "1",
		},
		Permissions: Permissions{
			Allow: []string{"Bash", "Write", "Edit"},
			Deny:  []string{},
		},
		APIKeyHelper: fmt.Sprintf("echo '%s'", apiKey),
	}
}

func writeSettings(settingsPath string, settings Settings) error {
	// Ensure the directory exists
	dir := filepath.Dir(settingsPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create directory %s: %w", dir, err)
	}

	data, err := json.MarshalIndent(settings, "", "    ")
	if err != nil {
		return fmt.Errorf("failed to marshal settings: %w", err)
	}

	if err := os.WriteFile(settingsPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write settings.json: %w", err)
	}

	return nil
}

func findBansosFile() (string, error) {
	// Try current working directory first
	cwd, err := os.Getwd()
	if err == nil {
		path := filepath.Join(cwd, "bansos.txt")
		if _, err := os.Stat(path); err == nil {
			return path, nil
		}
	}

	// Try executable directory
	exe, err := os.Executable()
	if err == nil {
		exeDir := filepath.Dir(exe)
		path := filepath.Join(exeDir, "bansos.txt")
		if _, err := os.Stat(path); err == nil {
			return path, nil
		}
	}

	return "", fmt.Errorf("bansos.txt not found in current directory or executable directory")
}

func main() {
	// Get home directory
	homeDir, err := getHomeDir()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	settingsPath := filepath.Join(homeDir, ".claude", "settings.json")

	// Find bansos.txt
	bansosPath, err := findBansosFile()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Using bansos.txt: %s\n", bansosPath)

	// Read the list of keys from bansos.txt
	keys, err := readKeysList(bansosPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	// Read current key from existing settings.json
	currentKey, err := readCurrentKey(settingsPath)
	if err != nil {
		fmt.Printf("Could not read current key from settings.json: %v\n", err)
		fmt.Println("Using first key from bansos.txt")
		currentKey = ""
	} else {
		fmt.Printf("Current key: %s...%s\n", currentKey[:min(8, len(currentKey))], currentKey[max(0, len(currentKey)-4):])
	}

	// Get next key
	nextKey := getNextKey(currentKey, keys)
	fmt.Printf("Next key: %s...%s\n", nextKey[:min(8, len(nextKey))], nextKey[max(0, len(nextKey)-4):])

	// Generate and write new settings
	settings := generateSettings(nextKey)
	if err := writeSettings(settingsPath, settings); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Successfully wrote settings to: %s\n", settingsPath)
}
