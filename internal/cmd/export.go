package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/suvaidkhan/envsecrets/internal/logic"
)

var exportCmd = &cobra.Command{
	Use:   "export",
	Short: "Export decrypted vault entries",
	Long:  `Decrypts and exports all vault entries to stdout in dotenv or JSON format.`,
	Example: `  envsecrets export --env prod > .env
  envsecrets export --env staging --format json > env.json`,
	RunE: runExport,
}

var (
	exportEnvFlag    string
	exportFormatFlag string
)

func init() {
	exportCmd.Flags().StringVarP(&exportEnvFlag, "env", "e", "", "environment name (required)")
	exportCmd.Flags().StringVar(&exportFormatFlag, "format", "dotenv", "output format (dotenv or json)")
	exportCmd.MarkFlagRequired("env")
	rootCmd.AddCommand(exportCmd)
}

func runExport(cmd *cobra.Command, args []string) error {
	// Validate format
	if exportFormatFlag != "dotenv" && exportFormatFlag != "json" {
		return fmt.Errorf("invalid format %q, must be dotenv or json", exportFormatFlag)
	}

	// Open vault
	vault, err := logic.OpenVault(exportEnvFlag)
	if err != nil {
		return fmt.Errorf("failed to open vault: %w", err)
	}

	// Decrypt all entries
	decrypted := make(map[string]string)
	for key, entry := range vault.Entries {
		plaintext, err := logic.Decrypt(entry.Value, vault.Meta.Salt, vault.Passphrase())
		if err != nil {
			return fmt.Errorf("failed to decrypt entry %q: %w", key, err)
		}
		decrypted[key] = string(plaintext)
	}

	// Output in requested format
	switch exportFormatFlag {
	case "dotenv":
		for key, value := range decrypted {
			fmt.Printf("%s=%s\n", key, value)
		}
	case "json":
		data, err := json.MarshalIndent(decrypted, "", "  ")
		if err != nil {
			return fmt.Errorf("failed to marshal JSON: %w", err)
		}
		fmt.Println(string(data))
	}

	return nil
}
