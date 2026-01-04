package logic

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"io"
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
