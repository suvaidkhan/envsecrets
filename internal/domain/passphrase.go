package domain

import (
	"fmt"
	"os"
)

type PassphraseService struct {
	envVar string
}

func NewPassphraseService(envVar string) *PassphraseService {
	if envVar == "" {
		envVar = "ENVSECRET_PASSPHRASE"
	}
	return &PassphraseService{envVar}
}

func (s *PassphraseService) Get(env string) (string, error) {
	if env == "" {
		return "", fmt.Errorf("environment cannot be empty")
	}
	if passphrase := os.Getenv(s.envVar); passphrase != "" {
		return passphrase, nil
	}
}
