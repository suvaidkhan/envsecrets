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
