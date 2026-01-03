package domain

import (
	"github.com/zalando/go-keyring"
)

func Get(key string) (string, error) {
	return keyring.Get("envsecrets", key)
}
