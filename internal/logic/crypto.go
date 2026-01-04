package logic

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"golang.org/x/crypto/argon2"
	"golang.org/x/crypto/bcrypt"
	"io"
	"runtime"
)

const (
	DefaultArgonTime     uint32 = 3
	DefaultArgonMemoryKB uint32 = 64
	DefaultArgonThreads  uint8  = 4
	DefaultKeyLength     uint32 = 32
	DefaultNonceSize     int    = 12
)

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

func Verify(hashedPassphrase, passphrase string) error {
	if hashedPassphrase == "" {
		return fmt.Errorf("hashed passphrase cannot be empty")
	}
	if err := bcrypt.CompareHashAndPassword([]byte(hashedPassphrase), []byte(passphrase)); err != nil {
		return fmt.Errorf("invalid passphrase: %w", err)
	}
	return nil
}

func getAEAD(encodedSalt, passphrase string) (cipher.AEAD, error) {
	if encodedSalt == "" {
		return nil, fmt.Errorf("salt cannot be empty")
	}
	if passphrase == "" {
		return nil, fmt.Errorf("passphrase cannot be empty")
	}

	salt, err := base64.StdEncoding.DecodeString(encodedSalt)
	if err != nil {
		return nil, fmt.Errorf("invalid salt encoding: %w", err)
	}
	if len(salt) == 0 {
		return nil, fmt.Errorf("salt cannot be empty")
	}

	key := deriveKey([]byte(passphrase), salt)

	block, err := aes.NewCipher(key)
	if err != nil {
		clearBytes(key, salt)
		return nil, fmt.Errorf("failed to create cipher: %w", err)
	}

	aead, err := cipher.NewGCM(block)
	if err != nil {
		clearBytes(key, salt)
		return nil, fmt.Errorf("failed to create GCM: %w", err)
	}

	return aead, nil
}

// Encrypt encrypts plaintext and returns base64-encoded ciphertext
func Encrypt(plaintext []byte, encodedSalt, passphrase string) (string, error) {
	if plaintext == nil {
		return "", fmt.Errorf("plaintext cannot be nil")
	}

	nonce := make([]byte, DefaultNonceSize)
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", fmt.Errorf("failed to generate nonce: %w", err)
	}

	aead, err := getAEAD(encodedSalt, passphrase)
	if err != nil {
		return "", err
	}

	ciphertext := aead.Seal(nil, nonce, plaintext, nil)
	result := make([]byte, 0, len(nonce)+len(ciphertext))
	result = append(result, nonce...)
	result = append(result, ciphertext...)

	encoded := base64.StdEncoding.EncodeToString(result)

	clearBytes(nonce, ciphertext)

	return encoded, nil
}

// Decrypt decrypts base64-encoded ciphertext and returns plaintext
func Decrypt(ciphertext, encodedSalt, passphrase string) ([]byte, error) {
	if ciphertext == "" {
		return nil, fmt.Errorf("ciphertext cannot be empty")
	}

	raw, err := base64.StdEncoding.DecodeString(ciphertext)
	if err != nil {
		return nil, fmt.Errorf("invalid ciphertext encoding: %w", err)
	}

	aead, err := getAEAD(encodedSalt, passphrase)
	if err != nil {
		return nil, err
	}

	if err := validateCiphertextLength(raw, aead.Overhead()); err != nil {
		return nil, err
	}

	// Extract nonce and ciphertext
	nonce := raw[:DefaultNonceSize]
	ciphertextBytes := raw[DefaultNonceSize:]
	plaintext, err := aead.Open(nil, nonce, ciphertextBytes, nil)
	if err != nil {
		clearBytes(nonce, ciphertextBytes)
		return nil, fmt.Errorf("decryption failed: %w", err)
	}

	clearBytes(nonce, ciphertextBytes)

	if plaintext == nil {
		return []byte{}, nil
	}

	return plaintext, nil
}

// validateCiphertextLength checks if the ciphertext meets the minimum length requirement
// The minimum length is nonce size + AEAD overhead (authentication tag)
func validateCiphertextLength(ciphertext []byte, overhead int) error {
	minLen := DefaultNonceSize + overhead
	if len(ciphertext) < minLen {
		return fmt.Errorf(
			"ciphertext too short: expected at least %d bytes, got %d",
			minLen,
			len(ciphertext),
		)
	}
	return nil
}

// deriveKey derives a key from a passphrase using Argon2id
func deriveKey(passphrase, salt []byte) []byte {
	return argon2.IDKey(
		passphrase,
		salt,
		DefaultArgonTime,
		DefaultArgonMemoryKB,
		DefaultArgonThreads,
		DefaultKeyLength,
	)
}

// clearBytes clears sensitive data from memory by overwrites a byte slice with zeros
func clearBytes(b ...[]byte) {
	if b == nil {
		return
	}
	for _, eachB := range b {
		for i := range eachB {
			eachB[i] = 0
		}
	}
	runtime.KeepAlive(b)
}
