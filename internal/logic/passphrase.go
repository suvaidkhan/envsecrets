package logic

import (
	"fmt"
	"github.com/AlecAivazis/survey/v2"
	"github.com/zalando/go-keyring"
	"os"
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

	keyring := fmt.Sprintf("env:%s", env)
	passphrase, err := Get(keyring)
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
	//TODO keyring cache logic
	return passphrase, nil
}
