package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

var Version = "dev"
var rootCmd = &cobra.Command{
	Use:   "envsecrets",
	Short: "envsecrets securely manages your .env files and secrets",
	Long: `envsecrets is a lightweight CLI tool for securely managing environment variables

envsecrets encrypts your environment variables using AES-GCM encryption with Argon2 key derivation.
Your secrets are protected with a passphrase that can be stored securely in your system's keyring.`,
	Version: Version,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Fprintln(os.Stderr, "Welcome to envsecrets! Use --help to see available commands.")
	},
}

func Execute() error {
	return rootCmd.Execute()
}
