package auth

import (
	"github.com/99designs/keyring"
)

const (
	KeyringService = "notte-cli"
	KeyringKey     = "api_key"
)

// getKeyring returns configured keyring
func getKeyring() (keyring.Keyring, error) {
	return keyring.Open(keyring.Config{
		ServiceName: KeyringService,
	})
}

// GetKeyringAPIKey retrieves API key from OS keychain
func GetKeyringAPIKey() (string, error) {
	kr, err := getKeyring()
	if err != nil {
		return "", err
	}

	item, err := kr.Get(KeyringKey)
	if err != nil {
		return "", err
	}

	return string(item.Data), nil
}

// SetKeyringAPIKey stores API key in OS keychain
func SetKeyringAPIKey(apiKey string) error {
	kr, err := getKeyring()
	if err != nil {
		return err
	}

	return kr.Set(keyring.Item{
		Key:  KeyringKey,
		Data: []byte(apiKey),
	})
}

// DeleteKeyringAPIKey removes API key from OS keychain
func DeleteKeyringAPIKey() error {
	kr, err := getKeyring()
	if err != nil {
		return err
	}

	return kr.Remove(KeyringKey)
}
