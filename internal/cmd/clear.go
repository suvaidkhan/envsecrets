package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/suvaidkhan/envsecret/internal/logic"
)

var clearCmd = &cobra.Command{
	Use:   "clear",
	Short: "Clear cached passphrase from keyring",
	Long:  `Removes the cached passphrase for an environment from the system keyring.`,
	Example: `  envsecrets clear --env prod
  envsecrets clear -e staging`,
	RunE: runClear,
}

var clearEnvFlag string

func init() {
	clearCmd.Flags().StringVarP(&clearEnvFlag, "env", "e", "", "environment name (required)")
	clearCmd.MarkFlagRequired("env")
	rootCmd.AddCommand(clearCmd)
}

func runClear(cmd *cobra.Command, args []string) error {
	if err := logic.Clear(clearEnvFlag); err != nil {
		// Keyring delete returns error if key doesn't exist, which is fine
		// Just report that cache was cleared
	}

	fmt.Printf("✓ Cleared cached passphrase for %s environment\n", clearEnvFlag)
	return nil
}
