package logic

import (
	"github.com/zalando/go-keyring"
)

func Get(key string) (string, error) {
	return keyring.Get("envsecrets", key)
}

// TODO Implemeny cache cleaning logic
func Clear(env string) {

}

func Set(key, value string) error {
	return keyring.Set("envsecrets", key, value)
}
