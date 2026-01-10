package cmd

import (
	"fmt"
	"os"

	"github.com/AlecAivazis/survey/v2"
	"github.com/spf13/cobra"
	"github.com/suvaidkhan/envsecrets/internal/logic"
)

var destroyCmd = &cobra.Command{
	Use:   "destroy",
	Short: "Destroy an entire vault",
	Long:  `Permanently deletes a vault and all its secrets. This action cannot be undone.`,
	Example: `  envsecrets destroy --env prod
  envsecrets destroy -e staging`,
	RunE: runDestroy,
}

var destroyEnvFlag string

func init() {
	destroyCmd.Flags().StringVarP(&destroyEnvFlag, "env", "e", "", "environment name (required)")
	destroyCmd.MarkFlagRequired("env")
	rootCmd.AddCommand(destroyCmd)
}

func runDestroy(cmd *cobra.Command, args []string) error {
	env := destroyEnvFlag

	// 1. Check if vault exists
	exists, err := logic.CheckIfExists(env)
	if err != nil {
		return fmt.Errorf("failed to check vault existence: %w", err)
	}
	if !exists {
		return fmt.Errorf("vault not found for environment %q", env)
	}

	// 2. Clear cache first to force fresh password prompt
	logic.Clear(env)

	// 3. Load vault to get fingerprint for verification
	vault, err := logic.LoadVault(env)
	if err != nil {
		return fmt.Errorf("failed to load vault: %w", err)
	}

	// 4. Prompt for password (cache already cleared, so will prompt)
	var passphrase string
	prompt := &survey.Password{
		Message: fmt.Sprintf("Enter passphrase to destroy vault for environment %q:", env),
	}
	if err := survey.AskOne(prompt, &passphrase); err != nil {
		return fmt.Errorf("failed to get passphrase: %w", err)
	}
	if passphrase == "" {
		return fmt.Errorf("passphrase cannot be empty")
	}

	// 5. Verify password against fingerprint
	if err := logic.Verify(vault.Meta.FingerPrint, passphrase); err != nil {
		return fmt.Errorf("invalid passphrase: cannot destroy vault")
	}

	// 6. Ask for confirmation
	confirm := false
	confirmPrompt := &survey.Confirm{
		Message: fmt.Sprintf("This will permanently delete the %s vault and all its secrets. Are you sure?", env),
		Default: false,
	}
	if err := survey.AskOne(confirmPrompt, &confirm); err != nil {
		return fmt.Errorf("confirmation prompt failed: %w", err)
	}
	if !confirm {
		fmt.Println("Vault destruction cancelled")
		return nil
	}

	// 7. Delete vault file
	vaultPath := logic.VaultPath(env)
	if err := os.Remove(vaultPath); err != nil {
		return fmt.Errorf("failed to delete vault file: %w", err)
	}

	// 8. Success message
	fmt.Printf("✓ Vault %s destroyed successfully\n", env)
	return nil
}
