package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

// KeyEntry represents a single entry in bansos.csv with id and key
type KeyEntry struct {
	ID  string
	Key string
}

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

func readKeysList(bansosPath string) ([]KeyEntry, error) {
	file, err := os.Open(bansosPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open bansos.csv: %w", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)

	// Read and skip header row
	header, err := reader.Read()
	if err != nil {
		return nil, fmt.Errorf("failed to read bansos.csv header: %w", err)
	}

	// Validate header
	if len(header) < 2 || strings.TrimSpace(strings.ToLower(header[0])) != "id" || strings.TrimSpace(strings.ToLower(header[1])) != "key" {
		return nil, fmt.Errorf("bansos.csv header harus berformat: id,key")
	}

	var entries []KeyEntry
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("failed to read bansos.csv: %w", err)
		}

		if len(record) < 2 {
			continue
		}

		id := strings.TrimSpace(record[0])
		key := strings.TrimSpace(record[1])

		if id == "" || key == "" {
			continue
		}

		entries = append(entries, KeyEntry{ID: id, Key: key})
	}

	if len(entries) == 0 {
		return nil, fmt.Errorf("bansos.csv is empty")
	}

	return entries, nil
}

func getNextEntry(currentKey string, entries []KeyEntry) (currentEntry *KeyEntry, nextEntry KeyEntry) {
	for i, entry := range entries {
		if entry.Key == currentKey {
			current := &entries[i]
			if i == len(entries)-1 {
				return current, entries[0]
			}
			return current, entries[i+1]
		}
	}
	// If current key is not found in the list, use the first entry
	return nil, entries[0]
}

func maskKey(key string) string {
	if len(key) <= 12 {
		return key
	}
	return key[:8] + "..." + key[len(key)-4:]
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
		path := filepath.Join(cwd, "bansos.csv")
		if _, err := os.Stat(path); err == nil {
			return path, nil
		}
	}

	// Try executable directory
	exe, err := os.Executable()
	if err == nil {
		exeDir := filepath.Dir(exe)
		path := filepath.Join(exeDir, "bansos.csv")
		if _, err := os.Stat(path); err == nil {
			return path, nil
		}
	}

	return "", fmt.Errorf("bansos.csv not found in current directory or executable directory")
}

func printFileNotFoundHelp() {
	fmt.Fprintln(os.Stderr, "")
	fmt.Fprintln(os.Stderr, "  ERROR: File bansos.csv tidak ditemukan!")
	fmt.Fprintln(os.Stderr, "")
	fmt.Fprintln(os.Stderr, "  Untuk menggunakan program ini, Anda perlu membuat file bansos.csv")
	fmt.Fprintln(os.Stderr, "  di direktori yang sama dengan program ini.")
	fmt.Fprintln(os.Stderr, "")
	fmt.Fprintln(os.Stderr, "  Cara membuat:")
	fmt.Fprintln(os.Stderr, "    1. Buat file bernama bansos.csv di direktori saat ini")
	fmt.Fprintln(os.Stderr, "    2. Baris pertama harus header: id,key")
	fmt.Fprintln(os.Stderr, "    3. Baris berikutnya berisi ID dan API key dipisahkan koma")
	fmt.Fprintln(os.Stderr, "")
	fmt.Fprintln(os.Stderr, "  Format penulisan (CSV dengan header id,key):")
	fmt.Fprintln(os.Stderr, "  ┌─────────────────────────────────────────────────────┐")
	fmt.Fprintln(os.Stderr, "  │ id,key                                              │")
	fmt.Fprintln(os.Stderr, "  │ anu,sk-ant-api03-xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx   │")
	fmt.Fprintln(os.Stderr, "  │ awg,sk-ant-api03-yyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyy   │")
	fmt.Fprintln(os.Stderr, "  │ bukped,sk-ant-api03-zzzzzzzzzzzzzzzzzzzzzzzzzzzzzz  │")
	fmt.Fprintln(os.Stderr, "  └─────────────────────────────────────────────────────┘")
	fmt.Fprintln(os.Stderr, "")
	fmt.Fprintln(os.Stderr, "  Note: ID bisa berupa angka atau huruf (bebas).")
	fmt.Fprintln(os.Stderr, "")
	fmt.Fprintln(os.Stderr, "  Contoh perintah untuk membuat file:")
	fmt.Fprintln(os.Stderr, "    Linux/macOS : nano bansos.csv")
	fmt.Fprintln(os.Stderr, "    Windows     : notepad bansos.csv")
	fmt.Fprintln(os.Stderr, "")
}

func printFileEmptyHelp(bansosPath string) {
	fmt.Fprintln(os.Stderr, "")
	fmt.Fprintf(os.Stderr, "  ERROR: File bansos.csv kosong! (%s)\n", bansosPath)
	fmt.Fprintln(os.Stderr, "")
	fmt.Fprintln(os.Stderr, "  File bansos.csv ditemukan, tetapi tidak berisi API key apapun.")
	fmt.Fprintln(os.Stderr, "  Silakan isi dengan format CSV: id,key (baris pertama adalah header).")
	fmt.Fprintln(os.Stderr, "")
	fmt.Fprintln(os.Stderr, "  Format penulisan (CSV dengan header id,key):")
	fmt.Fprintln(os.Stderr, "  ┌─────────────────────────────────────────────────────┐")
	fmt.Fprintln(os.Stderr, "  │ id,key                                              │")
	fmt.Fprintln(os.Stderr, "  │ anu,sk-ant-api03-xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx   │")
	fmt.Fprintln(os.Stderr, "  │ awg,sk-ant-api03-yyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyy   │")
	fmt.Fprintln(os.Stderr, "  │ bukped,sk-ant-api03-zzzzzzzzzzzzzzzzzzzzzzzzzzzzzz  │")
	fmt.Fprintln(os.Stderr, "  └─────────────────────────────────────────────────────┘")
	fmt.Fprintln(os.Stderr, "")
	fmt.Fprintln(os.Stderr, "  Tips:")
	fmt.Fprintln(os.Stderr, "    - Baris pertama harus header: id,key")
	fmt.Fprintln(os.Stderr, "    - ID bisa angka atau huruf (bebas)")
	fmt.Fprintln(os.Stderr, "    - Setiap baris berikutnya: id,api_key")
	fmt.Fprintln(os.Stderr, "    - Jangan ada spasi di awal/akhir key")
	fmt.Fprintln(os.Stderr, "    - Baris kosong akan diabaikan")
	fmt.Fprintln(os.Stderr, "")
}

func main() {
	// Get home directory
	homeDir, err := getHomeDir()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	settingsPath := filepath.Join(homeDir, ".claude", "settings.json")

	// Find bansos.csv
	bansosPath, err := findBansosFile()
	if err != nil {
		printFileNotFoundHelp()
		os.Exit(1)
	}

	fmt.Printf("Using bansos.csv: %s\n", bansosPath)

	// Read the list of keys from bansos.csv
	entries, err := readKeysList(bansosPath)
	if err != nil {
		if strings.Contains(err.Error(), "bansos.csv is empty") {
			printFileEmptyHelp(bansosPath)
		} else {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		}
		os.Exit(1)
	}

	// Read current key from existing settings.json
	currentKey, err := readCurrentKey(settingsPath)
	if err != nil {
		fmt.Printf("Could not read current key from settings.json: %v\n", err)
		fmt.Println("Using first key from bansos.csv")
		currentKey = ""
	}

	// Get next entry
	currentEntry, nextEntry := getNextEntry(currentKey, entries)

	if currentEntry != nil {
		fmt.Printf("Current key: ID %s (%s)\n", currentEntry.ID, maskKey(currentEntry.Key))
	} else {
		fmt.Printf("Current key: tidak ditemukan di bansos.csv\n")
	}

	fmt.Printf("Switched to: ID %s (%s)\n", nextEntry.ID, maskKey(nextEntry.Key))

	// Generate and write new settings
	settings := generateSettings(nextEntry.Key)
	if err := writeSettings(settingsPath, settings); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Successfully wrote settings to: %s\n", settingsPath)
}
