package logic

import (
	"fmt"
	"os"

	"github.com/AlecAivazis/survey/v2"
)

type Passphrase struct {
	envVar string
}

func NewPassphrase(envVar string) *Passphrase {
	if envVar == "" {
		envVar = "ENVSECRET_PASSPHRASE"
	}
	return &Passphrase{envVar}
}

func (s *Passphrase) Get(env string) (string, error) {
	if env == "" {
		return "", fmt.Errorf("environment cannot be empty")
	}
	if passphrase := os.Getenv(s.envVar); passphrase != "" {
		return passphrase, nil
	}

	ring := fmt.Sprintf("env:%s", env)
	passphrase, err := Get(ring)
	if err == nil && passphrase != "" {
		return passphrase, nil
	}

	return s.getFromUser(env)
}

func (s *Passphrase) getFromUser(env string) (string, error) {
	passphrase := ""
	prompt := &survey.Password{
		Message: fmt.Sprintf("Enter passphrase for environment %q:", env),
	}
	if err := survey.AskOne(prompt, &passphrase); err != nil {
		return "", fmt.Errorf("failed to get passphrase: %w", err)
	}
	if passphrase == "" {
		return "", fmt.Errorf("passphrase cannot be empty")
	}

	// Cache to keyring for future use
	ring := fmt.Sprintf("env:%s", env)
	if err := Set(ring, passphrase); err != nil {
		// Non-fatal: just warn if keyring is unavailable
		fmt.Fprintf(os.Stderr, "Warning: failed to cache passphrase to keyring: %v\n", err)
	}

	return passphrase, nil
}
