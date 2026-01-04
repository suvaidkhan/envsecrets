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

func LoadVault(env string) (*Vault, error) {
	if env == "" {
		return nil, fmt.Errorf("environment cannot be empty")
	}

	path := createPath(env)
	data, err := readFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("vault not found for environment %q: %w", env, err)
		}
		return nil, fmt.Errorf("failed to read vault file: %w", err)
	}
	var vault Vault
	if err := json.Unmarshal(data, &vault); err != nil {
		return nil, fmt.Errorf("failed to unmarshal vault: %w", err)
	}
	vault.path = path

	if vault.Meta.Env != env {
		return nil, fmt.Errorf(
			"vault environment mismatch: expected %q, got %q",
			env,
			vault.Meta.Env,
		)
	}

	return &vault, nil

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

func readFile(path string) ([]byte, error) {
	return os.ReadFile(path)
}
