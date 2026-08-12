package crypto

import (
	"encoding/hex"
	"errors"
	"fmt"
	"os"
)

// EnvKey is the environment variable holding the hex-encoded 16-byte AES key.
const EnvKey = "GOWRITE_ENCRYPTION_KEY"

// ErrMissingKey is returned when the key env var is unset.
var ErrMissingKey = errors.New("crypto: " + EnvKey + " is not set")

// FromEnv builds a Cipher from a hex-encoded key in GOWRITE_ENCRYPTION_KEY.
func FromEnv() (*Cipher, error) {
	raw := os.Getenv(EnvKey)
	if raw == "" {
		return nil, ErrMissingKey
	}
	return FromHexKey(raw)
}

// FromHexKey builds a Cipher from a hex-encoded 16-byte key.
func FromHexKey(hexKey string) (*Cipher, error) {
	key, err := hex.DecodeString(hexKey)
	if err != nil {
		return nil, fmt.Errorf("crypto: key is not valid hex: %w", err)
	}
	return New(key)
}
