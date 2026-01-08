package cmd

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/suvaidkhan/envsecret/internal/logic"
)

var importCmd = &cobra.Command{
	Use:   "import [file]",
	Short: "Import entries from file into vault",
	Long:  `Imports plaintext entries from a dotenv or JSON file into an encrypted vault.`,
	Example: `  envsecrets import .env --env prod --format dotenv
  envsecrets import config.json --env staging --format json
  cat .env | envsecrets import --env local --format dotenv
  envsecrets import .env --env prod --format dotenv --overwrite`,
	RunE: runImport,
}

var (
	importEnvFlag       string
	importFormatFlag    string
	importOverwriteFlag bool
)

func init() {
	importCmd.Flags().StringVarP(&importEnvFlag, "env", "e", "", "environment name (required)")
	importCmd.Flags().StringVar(&importFormatFlag, "format", "", "input format: dotenv or json (required)")
	importCmd.Flags().BoolVar(&importOverwriteFlag, "overwrite", false, "overwrite existing keys")
	importCmd.MarkFlagRequired("env")
	importCmd.MarkFlagRequired("format")
	rootCmd.AddCommand(importCmd)
}

func runImport(cmd *cobra.Command, args []string) error {
	// Validate format
	if importFormatFlag != "dotenv" && importFormatFlag != "json" {
		return fmt.Errorf("invalid format %q, must be dotenv or json", importFormatFlag)
	}

	// Open input (file or stdin)
	var reader io.Reader
	if len(args) > 0 {
		file, err := os.Open(args[0])
		if err != nil {
			return fmt.Errorf("failed to open file: %w", err)
		}
		defer file.Close()
		reader = file
	} else {
		reader = os.Stdin
	}

	// Parse input
	var entries map[string]string
	var err error
	switch importFormatFlag {
	case "dotenv":
		entries, err = parseDotEnv(reader)
	case "json":
		entries, err = parseJSON(reader)
	}
	if err != nil {
		return fmt.Errorf("failed to parse input: %w", err)
	}

	if len(entries) == 0 {
		return fmt.Errorf("no entries found in input")
	}

	// Open vault
	vault, err := logic.OpenVault(importEnvFlag)
	if err != nil {
		return fmt.Errorf("failed to open vault: %w", err)
	}

	// Import entries
	imported := 0
	skipped := 0
	for key, value := range entries {
		// Check if key exists
		_, err := vault.GetEntry(key)
		if err == nil && !importOverwriteFlag {
			skipped++
			continue
		}

		// Encrypt value
		encryptedValue, err := logic.Encrypt([]byte(value), vault.Meta.Salt, vault.Passphrase())
		if err != nil {
			return fmt.Errorf("failed to encrypt entry %q: %w", key, err)
		}

		// Add to vault
		if err := vault.SetEntry(key, encryptedValue); err != nil {
			return fmt.Errorf("failed to set entry %q: %w", key, err)
		}
		imported++
	}

	// Save vault if any entries were imported
	if imported > 0 {
		if err := logic.SaveVault(vault); err != nil {
			return fmt.Errorf("failed to save vault: %w", err)
		}
	}

	fmt.Printf("✓ Imported %d entry(s), skipped %d entry(s)\n", imported, skipped)
	return nil
}

// parseDotEnv parses dotenv format (KEY=value)
func parseDotEnv(r io.Reader) (map[string]string, error) {
	entries := make(map[string]string)
	scanner := bufio.NewScanner(r)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		// Skip empty lines and comments
		if line == "" || strings.HasPrefix(line, "#") || strings.HasPrefix(line, "//") {
			continue
		}

		// Split on first '='
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])

		if key == "" {
			continue
		}

		// Remove surrounding quotes if present
		if len(value) >= 2 {
			if (value[0] == '"' && value[len(value)-1] == '"') ||
				(value[0] == '\'' && value[len(value)-1] == '\'') {
				value = value[1 : len(value)-1]
			}
		}

		entries[key] = value
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return entries, nil
}

// parseJSON parses JSON format
func parseJSON(r io.Reader) (map[string]string, error) {
	var entries map[string]string
	decoder := json.NewDecoder(r)
	if err := decoder.Decode(&entries); err != nil {
		return nil, err
	}
	return entries, nil
}
