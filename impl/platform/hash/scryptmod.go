package hash

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/subtle"
	"encoding/base64"
	"fmt"

	xscrypt "golang.org/x/crypto/scrypt"
)

// verifyScryptMod verifies plain against a Firebase-style scryptMod hash.
// opts must contain: "salt" (b64), "saltSeparator" (b64), "signerKey" (b64).
// storedHash is base64-encoded.
func verifyScryptMod(plain, storedHash string, opts Options) (bool, error) {
	salt, err := b64Decode(opts, "salt")
	if err != nil {
		return false, fmt.Errorf("hash: scryptMod salt: %w", err)
	}
	saltSep, err := b64Decode(opts, "saltSeparator")
	if err != nil {
		return false, fmt.Errorf("hash: scryptMod saltSeparator: %w", err)
	}
	signerKey, err := b64Decode(opts, "signerKey")
	if err != nil {
		return false, fmt.Errorf("hash: scryptMod signerKey: %w", err)
	}

	derived, err := xscrypt.Key([]byte(plain), append(salt, saltSep...), 16384, 8, 1, 64)
	if err != nil {
		return false, fmt.Errorf("hash: scryptMod scrypt: %w", err)
	}

	key := derived[:32]
	iv := make([]byte, 16) // all zeros

	block, err := aes.NewCipher(key)
	if err != nil {
		return false, fmt.Errorf("hash: scryptMod aes: %w", err)
	}
	got := make([]byte, len(signerKey))
	cipher.NewCTR(block, iv).XORKeyStream(got, signerKey)

	want, err := base64.StdEncoding.DecodeString(storedHash)
	if err != nil {
		return false, fmt.Errorf("hash: scryptMod hash decode: %w", err)
	}
	return subtle.ConstantTimeCompare(got, want) == 1, nil
}

func b64Decode(opts Options, key string) ([]byte, error) {
	s, _ := opts[key].(string)
	if s == "" {
		return nil, fmt.Errorf("missing %q", key)
	}
	return base64.StdEncoding.DecodeString(s)
}
