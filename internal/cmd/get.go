package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/suvaidkhan/envsecret/internal/logic"
	"os"
)

var getCmd = &cobra.Command{
	Use:     "get",
	Short:   "Get secrets from your vault based on env and key provided",
	Long:    `Gets the decrypted secrets from your vault based on env and key provided the output is printed on stdout`,
	Example: `  envsecrets get --env prod --key API_KEY`,
	RunE:    runGet,
}

var (
	envGetFlag string
	keyGetFlag string
)

func init() {
	getCmd.Flags().StringVarP(&envGetFlag, "env", "e", "", "The environment you want to get (required)")
	getCmd.Flags().StringVarP(&keyGetFlag, "key", "k", "", "The secret key you want to get (required)")
	_ = getCmd.MarkFlagRequired("env")
	_ = getCmd.MarkFlagRequired("key")
	rootCmd.AddCommand(getCmd)
}

func runGet(cmd *cobra.Command, args []string) error {
	env := envGetFlag
	key := keyGetFlag

	vault, err := logic.OpenVault(env)
	if err != nil {
		return fmt.Errorf("Vault cannot be opened: %w", err)
	}

	entry, err := vault.GetEntry(key)
	if err != nil {
		return fmt.Errorf("Vault cannot be retrieved: %w", err)
	}

	fmt.Fprintln(os.Stdout, entry.Value)
	return nil
}
