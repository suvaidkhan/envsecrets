package logic

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

type Entry struct {
	Value     string `json:"value"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

type Meta struct {
	Env         string `json:"env"`
	Salt        string `json:"salt"`
	FingerPrint string `json:"fingerprint"`
}

type Vault struct {
	Meta       Meta             `json:"meta"`
	Entries    map[string]Entry `json:"entries"`
	path       string
	passphrase string
}

func NewVault(env, fingerprint, salt string) (*Vault, error) {
	if env == "" {
		return nil, errors.New("environment cannot be empty")
	}
	if fingerprint == "" {
		return nil, errors.New("fingerprint cannot be empty")
	}
	if salt == "" {
		return nil, errors.New("salt cannot be empty")
	}

	vault := &Vault{
		Meta: Meta{
			Env:         env,
			Salt:        salt,
			FingerPrint: fingerprint,
		},
		Entries: make(map[string]Entry),
	}

	return vault, nil
}

// Create new vault
func Create(env string) (*Vault, error) {
	exists, err := CheckIfExists(env)
	if err != nil {
		return nil, fmt.Errorf("failed to check vault existence: %w", err)
	}
	if exists {
		return nil, fmt.Errorf("vault %s already exists", env)
	}

	passPhrase := NewPassphrase("")
	pass, err := passPhrase.Get(env)
	if err != nil {
		return nil, fmt.Errorf("failed to get passphrase: %w", err)
	}

	fingerprint, err := BCryptHash(pass)
	if err != nil {
		return nil, fmt.Errorf("failed to hash passphrase: %w", err)
	}
	salt, err := GenerateSalt()
	if err != nil {
		return nil, fmt.Errorf("failed to generate salt: %w", err)
	}
	vault, err := NewVault(env, fingerprint, salt)
	if err != nil {
		return nil, fmt.Errorf("failed to create vault: %w", err)
	}
	err = CreateVault(vault)
	if err != nil {
		return nil, fmt.Errorf("failed to create vault: %w", err)
	}

	return vault, nil
}

func OpenVault(env string) (*Vault, error) {
	if exists, err := CheckIfExists(env); err != nil || !exists {
		return nil, fmt.Errorf("env %s does not exist", env)
	}
	passPhrase := NewPassphrase("")
	pass, err := passPhrase.Get(env)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve passphrase: %w", err)
	}
	vault, err := LoadVault(env)
	if err != nil || vault == nil {
		return nil, fmt.Errorf("failed to load vault: %w", err)
	}
	err = Verify(vault.Meta.FingerPrint, pass)
	if err != nil {
		Clear(env)
		return nil, fmt.Errorf("invalid credentials: %w", err)
	}
	vault.passphrase = pass
	return vault, err
}

func (v *Vault) SetEntry(key, encryptedValue string) error {
	if key == "" {
		return errors.New("key cannot be empty")
	}
	if encryptedValue == "" {
		return errors.New("encrypted value cannot be empty")
	}

	now := time.Now().UTC().Format(time.RFC3339)
	entry, exists := v.Entries[key]

	if exists {
		entry.Value = encryptedValue
		entry.UpdatedAt = now
	} else {
		entry = Entry{
			Value:     encryptedValue,
			CreatedAt: now,
			UpdatedAt: now,
		}
	}

	if v.Entries == nil {
		v.Entries = make(map[string]Entry)
	}
	v.Entries[key] = entry
	return nil
}

// VaultPath returns the default path for a vault
func VaultPath(env string) string {
	return filepath.Join(".envsecrets", fmt.Sprintf("%s.vault", env))
}

// check if secrets repo exists
func CheckIfExists(env string) (bool, error) {
	if env == "" {
		return false, fmt.Errorf("environment cannot be empty")
	}

	vaultPath := VaultPath(env)
	_, err := os.Stat(vaultPath)
	if err != nil {
		return false, nil
	}

	return true, nil
}
