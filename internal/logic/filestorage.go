package logic

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

const DefaultDirMode = 0o700
const DefaultFileMode uint32 = 0o600

func CreateVault(vault *Vault) error {
	if vault == nil {
		return fmt.Errorf("vault cannot be nil")
	}

	// Create .envsecrets directory first
	if err := os.MkdirAll(".envsecrets", os.FileMode(DefaultDirMode)); err != nil {
		return fmt.Errorf("failed to create .envsecrets directory: %w", err)
	}

	// Check if vault exists using env name (not path)
	exists, err := CheckIfExists(vault.Meta.Env)
	if err != nil {
		return fmt.Errorf("failed to check if vault exists: %w", err)
	}
	if exists {
		return fmt.Errorf("vault already exists for environment %q", vault.Meta.Env)
	}

	vault.path = createPath(vault.Meta.Env)

	data, err := json.MarshalIndent(vault, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal vault: %w", err)
	}

	if err := write(vault.path, data); err != nil {
		return fmt.Errorf("failed to write vault: %w", err)
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

func SaveVault(vault *Vault) error {
	if vault == nil {
		return fmt.Errorf("vault cannot be nil")
	}
	if vault.path == "" {
		return fmt.Errorf("vault path is not set")
	}

	data, err := json.MarshalIndent(vault, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal vault: %w", err)
	}

	if err := write(vault.path, data); err != nil {
		return fmt.Errorf("failed to write vault: %w", err)
	}

	return nil
}

func createPath(env string) string {
	return filepath.Join(".envsecrets", fmt.Sprintf("%s.vault", env))
}

func write(path string, data []byte) error {
	return os.WriteFile(path, data, os.FileMode(DefaultFileMode))
}

func readFile(path string) ([]byte, error) {
	return os.ReadFile(path)
}
