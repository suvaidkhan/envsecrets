package cmd

import (
	"fmt"

	"github.com/AlecAivazis/survey/v2"
	"github.com/spf13/cobra"
	"github.com/suvaidkhan/envsecret/internal/logic"
)

var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Add or update an entry in the vault",
	Long:  `Adds a new secret or updates an existing one in the encrypted vault.`,
	Example: `  envsecrets add --env prod --key API_KEY --value secret123
  envsecrets add --env dev --secret`,
	RunE: runAdd,
}

var (
	addEnvFlag    string
	addKeyFlag    string
	addValueFlag  string
	addSecretFlag bool
)

func init() {
	addCmd.Flags().StringVarP(&addEnvFlag, "env", "e", "", "environment name (required)")
	addCmd.Flags().StringVarP(&addKeyFlag, "key", "k", "", "entry key")
	addCmd.Flags().StringVarP(&addValueFlag, "value", "v", "", "entry value")
	addCmd.Flags().BoolVarP(&addSecretFlag, "secret", "s", false, "hide value input")
	addCmd.MarkFlagRequired("env")
	rootCmd.AddCommand(addCmd)
}

func runAdd(*cobra.Command, []string) error {
	// Get key and value (from flags or prompt)
	key := addKeyFlag
	value := addValueFlag

	if key == "" {
		prompt := &survey.Input{Message: "Enter key:"}
		survey.AskOne(prompt, &key)
	}
	if key == "" {
		return fmt.Errorf("key cannot be empty")
	}

	if value == "" {
		if addSecretFlag {
			prompt := &survey.Password{Message: "Enter secret:"}
			survey.AskOne(prompt, &value)
		} else {
			prompt := &survey.Input{Message: "Enter value:"}
			survey.AskOne(prompt, &value)
		}
	}
	if value == "" {
		return fmt.Errorf("value cannot be empty")
	}

	// Open vault (loads from disk and verifies passphrase)
	vault, err := logic.OpenVault(addEnvFlag)
	if err != nil {
		return fmt.Errorf("failed to open vault: %w", err)
	}

	// Encrypt the value
	encryptedValue, err := logic.Encrypt([]byte(value), vault.Meta.Salt, vault.Passphrase())
	if err != nil {
		return fmt.Errorf("failed to encrypt value: %w", err)
	}

	// Add entry to vault (in memory)
	if err := vault.SetEntry(key, encryptedValue); err != nil {
		return fmt.Errorf("failed to set entry: %w", err)
	}

	// Save vault back to disk
	if err := logic.SaveVault(vault); err != nil {
		return fmt.Errorf("failed to save vault: %w", err)
	}

	fmt.Printf("✓ Entry '%s' added successfully to %s vault\n", key, addEnvFlag)
	return nil
}
