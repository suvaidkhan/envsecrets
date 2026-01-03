package logic

import (
	"encoding/json"
	"fmt"
	"os"
)

const DefaultDirMode = 0o700
const DefaultFileMode uint32 = 0o600

func CreateVault(vault *Vault) error {
	if vault == nil {
		return fmt.Errorf("vault cannot be nil")
	}

	vault.path = createPath(vault.Meta.Env)
	if err := mkDir(vault.path); err != nil {
		return fmt.Errorf("failed to create vault directory: %w", err)
	}

	exists, err := CheckIfExists(vault.path)
	if err != nil {
		return fmt.Errorf("failed to check if path exists: %w", err)
	}
	if exists {
		return fmt.Errorf("vault already exists for environment %q", vault.Meta.Env)
	}

	data, err := json.MarshalIndent(vault, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal vault: %w", err)
	}

	if err := write(vault.path, data); err != nil {
		return fmt.Errorf("failed to write to vault: %w", err)
	}

	return nil
}

func createPath(env string) string {
	return env + ".vault.enc"
}

func mkDir(path string) error {
	return os.MkdirAll(path, os.FileMode(DefaultDirMode))
}

func write(path string, data []byte) error {
	return os.WriteFile(path, data, os.FileMode(DefaultFileMode))
}
