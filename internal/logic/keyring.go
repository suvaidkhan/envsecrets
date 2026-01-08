package logic

import (
	"fmt"

	"github.com/zalando/go-keyring"
)

func Get(key string) (string, error) {
	return keyring.Get("envsecrets", key)
}

// Clear removes the cached passphrase for an environment from the keyring
func Clear(env string) error {
	ring := fmt.Sprintf("env:%s", env)
	return keyring.Delete("envsecrets", ring)
}

func Set(key, value string) error {
	return keyring.Set("envsecrets", key, value)
}
