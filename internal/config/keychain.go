package config

import (
	"fmt"
	"strings"

	"github.com/99designs/keyring"
)

// Keychain provides secure storage for sensitive data.
type Keychain struct {
	ring keyring.Keyring
}

// NewKeychain creates a new keychain for the given service name.
func NewKeychain(service string) (*Keychain, error) {
	ring, err := keyring.Open(keyring.Config{
		ServiceName: service,
	})
	if err != nil {
		return nil, fmt.Errorf("open keyring: %w", err)
	}
	return &Keychain{ring: ring}, nil
}

// Set stores a secret in the keychain.
func (k *Keychain) Set(key, value string) error {
	item := keyring.Item{
		Key:         key,
		Data:        []byte(value),
		Label:       fmt.Sprintf("PairAdmin: %s", key),
		Description: "API key or other sensitive data",
	}
	return k.ring.Set(item)
}

// Get retrieves a secret from the keychain.
func (k *Keychain) Get(key string) (string, error) {
	item, err := k.ring.Get(key)
	if err != nil {
		if strings.Contains(err.Error(), "not found") || strings.Contains(err.Error(), "no such") {
			return "", fmt.Errorf("key %q not found in keychain", key)
		}
		return "", fmt.Errorf("get from keyring: %w", err)
	}
	return string(item.Data), nil
}

// Delete removes a secret from the keychain.
func (k *Keychain) Delete(key string) error {
	return k.ring.Remove(key)
}