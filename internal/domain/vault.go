package domain

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/zalando/go-keyring"
	"golang.org/x/crypto/bcrypt"
	"io"
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
	exists, err := checkIfExists(env)
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

	return vault, nil
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

func BCryptHash(passphrase string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(passphrase), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("failed to generate hash: %w", err)
	}
	return string(hash), nil
}

func GenerateSalt() (string, error) {
	n := 16
	salt := make([]byte, n)
	if _, err := io.ReadFull(rand.Reader, salt); err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(salt), nil
}

// VaultPath returns the default path for a vault
func VaultPath(env string) string {
	return filepath.Join(".envsecrets", fmt.Sprintf("%s.vault", env))
}

// check if secrets repo exists
func checkIfExists(env string) (bool, error) {
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
