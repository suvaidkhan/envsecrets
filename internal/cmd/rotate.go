package cmd

import (
	"fmt"
	"os"

	"github.com/AlecAivazis/survey/v2"
	"github.com/spf13/cobra"
	"github.com/suvaidkhan/envsecret/internal/logic"
)

var rotateCmd = &cobra.Command{
	Use:   "rotate",
	Short: "Rotate the passphrase for a vault",
	Long:  `Changes the passphrase for a vault by re-encrypting all entries with a new passphrase and salt.`,
	Example: `  envsecrets rotate --env prod`,
	RunE:    runRotate,
}

var rotateEnvFlag string

func init() {
	rotateCmd.Flags().StringVarP(&rotateEnvFlag, "env", "e", "", "environment name (required)")
	rotateCmd.MarkFlagRequired("env")
	rootCmd.AddCommand(rotateCmd)
}

func runRotate(cmd *cobra.Command, args []string) error {
	// Open vault with current passphrase
	vault, err := logic.OpenVault(rotateEnvFlag)
	if err != nil {
		return fmt.Errorf("failed to open vault: %w", err)
	}

	oldPassphrase := vault.Passphrase()
	oldSalt := vault.Meta.Salt

	// Prompt for new passphrase
	var newPassphrase string
	prompt := &survey.Password{Message: "Enter new passphrase:"}
	if err := survey.AskOne(prompt, &newPassphrase); err != nil {
		return fmt.Errorf("failed to get new passphrase: %w", err)
	}
	if newPassphrase == "" {
		return fmt.Errorf("new passphrase cannot be empty")
	}

	// Confirm new passphrase
	var confirmPassphrase string
	confirmPrompt := &survey.Password{Message: "Confirm new passphrase:"}
	if err := survey.AskOne(confirmPrompt, &confirmPassphrase); err != nil {
		return fmt.Errorf("failed to confirm passphrase: %w", err)
	}
	if newPassphrase != confirmPassphrase {
		return fmt.Errorf("passphrases do not match")
	}

	// Generate new salt and fingerprint
	newSalt, err := logic.GenerateSalt()
	if err != nil {
		return fmt.Errorf("failed to generate salt: %w", err)
	}

	newFingerprint, err := logic.BCryptHash(newPassphrase)
	if err != nil {
		return fmt.Errorf("failed to hash new passphrase: %w", err)
	}

	// Re-encrypt all entries
	for key, entry := range vault.Entries {
		// Decrypt with old passphrase
		plaintext, err := logic.Decrypt(entry.Value, oldSalt, oldPassphrase)
		if err != nil {
			return fmt.Errorf("failed to decrypt entry %q: %w", key, err)
		}

		// Encrypt with new passphrase
		encryptedValue, err := logic.Encrypt(plaintext, newSalt, newPassphrase)
		if err != nil {
			return fmt.Errorf("failed to encrypt entry %q: %w", key, err)
		}

		// Update entry (preserves timestamps)
		entry.Value = encryptedValue
		vault.Entries[key] = entry
	}

	// Update vault metadata
	vault.Meta.Salt = newSalt
	vault.Meta.FingerPrint = newFingerprint

	// Update keyring cache with new passphrase
	ring := fmt.Sprintf("env:%s", rotateEnvFlag)
	if err := logic.Set(ring, newPassphrase); err != nil {
		// Non-fatal warning
		fmt.Fprintf(os.Stderr, "Warning: failed to update passphrase in keyring: %v\n", err)
	}

	// Save vault
	if err := logic.SaveVault(vault); err != nil {
		return fmt.Errorf("failed to save vault: %w", err)
	}

	fmt.Printf("✓ Passphrase rotated successfully for %s vault\n", rotateEnvFlag)
	fmt.Printf("  %d entries re-encrypted\n", len(vault.Entries))
	return nil
}
