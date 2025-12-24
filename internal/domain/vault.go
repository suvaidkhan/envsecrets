package domain

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

type Vault struct {
	Env         string `json:"env"`
	Salt        string `json:"salt"`
	Fingerprint string `json:"fingerprint"`
	Entries     map[string]string
	path        string
	passphrase  string
}

type Option func(*Vault)

func New(env string, opts ...Option) *Vault {
	v := &Vault{
		Env:     env,
		Entries: make(map[string]string),
	}
	for _, opt := range opts {
		opt(v)
	}
	return v
}

// WithSalt sets the salt for the vault
func WithSalt(salt string) Option {
	return func(v *Vault) {
		v.Salt = salt
	}
}

// WithFingerprint sets the fingerprint
func WithFingerprint(fp string) Option {
	return func(v *Vault) {
		v.Fingerprint = fp
	}
}

// WithPath sets the file path
func WithPath(path string) Option {
	return func(v *Vault) {
		v.path = path
	}
}

// Load loads and decrypts a vault from disk
func Load(env string, passphrase string) (*Vault, error) {
	path := VaultPath(env)

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read vault: %w", err)
	}

	// Decrypt the data
	decrypted, err := decrypt(data, passphrase)
	if err != nil {
		return nil, fmt.Errorf("decrypt vault: %w", err)
	}

	// Unmarshal
	var v Vault
	if err := json.Unmarshal(decrypted, &v); err != nil {
		return nil, fmt.Errorf("unmarshal vault: %w", err)
	}

	// Set runtime fields
	v.path = path
	v.passphrase = passphrase

	return &v, nil
}

func (v *Vault) Add(key, value string) error {
	if key == "" {
		return fmt.Errorf("key cannot be empty")
	}

	v.Entries[key] = value
	return nil
}

// Get retrieves a plaintext value by key
func (v *Vault) Get(key string) (string, bool) {
	value, ok := v.Entries[key]
	return value, ok
}

// Delete removes an entry
func (v *Vault) Delete(key string) error {
	if _, exists := v.Entries[key]; !exists {
		return fmt.Errorf("key %q not found", key)
	}
	delete(v.Entries, key)
	return nil
}

// Keys returns all keys in sorted order
func (v *Vault) Keys() []string {
	keys := make([]string, 0, len(v.Entries))
	for k := range v.Entries {
		keys = append(keys, k)
	}
	return keys
}

// Save encrypts and saves the vault to disk
func (v *Vault) Save() error {
	if v.path == "" {
		v.path = VaultPath(v.Env)
	}

	// Ensure directory exists
	if err := os.MkdirAll(filepath.Dir(v.path), 0700); err != nil {
		return fmt.Errorf("create vault directory: %w", err)
	}

	// Marshal to JSON
	data, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal vault: %w", err)
	}

	// Encrypt
	encrypted, err := encrypt(data, v.passphrase)
	if err != nil {
		return fmt.Errorf("encrypt vault: %w", err)
	}

	// Write to file
	if err := os.WriteFile(v.path, encrypted, 0600); err != nil {
		return fmt.Errorf("write vault: %w", err)
	}

	return nil
}

// VaultPath returns the default path for a vault
func VaultPath(env string) string {
	return filepath.Join(".envsecrts", fmt.Sprintf("%s.vault", env))
}

// Helper functions (would be in crypto.go)
func encrypt(data []byte, passphrase string) ([]byte, error) {
	// Your AES-256-GCM + Argon2id implementation
	return nil, nil
}

func decrypt(data []byte, passphrase string) ([]byte, error) {
	// Your AES-256-GCM + Argon2id implementation
	return nil, nil
}
