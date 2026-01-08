package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/suvaidkhan/envsecrets/internal/logic"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize a new vault for an environment",
	Long: `Creates a new encrypted vault for storing environment secrets.

The vault will be created at .envsecrets/{env}.vault with 0o600 permissions.
You'll be prompted for a passphrase if not provided via environment variable or keyring.`,
	Example: `  # Initialize a production vault
  envsecrets init --env prod

  # Initialize with short flag
  envsecrets init -e staging`,
	RunE: runInit,
}

var envFlag string

func init() {
	initCmd.Flags().StringVarP(&envFlag, "env", "e", "", "environment name (required)")
	initCmd.MarkFlagRequired("env")
	rootCmd.AddCommand(initCmd)
}

func runInit(cmd *cobra.Command, args []string) error {
	_, err := logic.Create(envFlag)
	if err != nil {
		return fmt.Errorf("failed to create vault: %w", err)
	}

	vaultPath := logic.VaultPath(envFlag)
	fmt.Printf("✓ Vault created successfully for environment '%s'\n", envFlag)
	fmt.Printf("  Location: %s\n", vaultPath)

	return nil
}
