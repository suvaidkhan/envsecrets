package cmd

import (
	"fmt"

	"github.com/AlecAivazis/survey/v2"
	"github.com/spf13/cobra"
	"github.com/suvaidkhan/envsecret/internal/logic"
)

var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete secret",
	Long:  `Deletes secret from the vault`,
	Example: `  envsecrets delete --env prod --key API_KEY
  envsecrets delete --env dev --secretKey`,
	RunE: runDel,
}

var (
	deleteEnvFlag string
	deleteKeyFlag string
)

func init() {
	deleteCmd.Flags().StringVarP(&deleteEnvFlag, "env", "e", "", "Environment variable name (required)")
	deleteCmd.Flags().StringVarP(&deleteKeyFlag, "key", "k", "", "Secret Key (required)")
	deleteCmd.MarkFlagRequired("env")
	deleteCmd.MarkFlagRequired("key")
	rootCmd.AddCommand(deleteCmd)
}

func runDel(cmd *cobra.Command, args []string) error {
	env := deleteEnvFlag
	key := deleteKeyFlag

	if key == "" {
		prompt := &survey.Input{Message: "Enter key:"}
		survey.AskOne(prompt, &key)
	}
	if key == "" {
		return fmt.Errorf("key cannot be empty")
	}

	if env == "" {
		prompt := &survey.Input{Message: "Enter env:"}
		survey.AskOne(prompt, &env)
	}
	if env == "" {
		return fmt.Errorf("env cannot be empty")
	}

	vault, err := logic.OpenVault(env)
	if err != nil {
		return fmt.Errorf("error opening vault: %w", err)
	}

	err = vault.DeleteEntry(key)
	if err != nil {
		return fmt.Errorf("error deleting secret: %w", err)
	}
	return nil
}
