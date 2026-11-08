package hash

import (
	"crypto/subtle"
	"encoding/hex"
	"fmt"

	"golang.org/x/crypto/scrypt"
)

// verifyScrypt verifies plain against a hex-encoded scrypt hash.
// opts must contain: "salt" (hex string), "costCpu" (N), "costMemory" (r),
// "costParallel" (p), "length" (keyLen). Defaults match PHP Scrypt defaults.
func verifyScrypt(plain, storedHash string, opts Options) (bool, error) {
	saltHex, _ := opts["salt"].(string)
	if saltHex == "" {
		return false, fmt.Errorf("hash: scrypt opts missing salt")
	}
	salt, err := hex.DecodeString(saltHex)
	if err != nil {
		return false, fmt.Errorf("hash: scrypt salt decode: %w", err)
	}

	n := optInt(opts, "costCpu", 8)
	r := optInt(opts, "costMemory", 14)
	p := optInt(opts, "costParallel", 1)
	keyLen := optInt(opts, "length", 64)

	got, err := scrypt.Key([]byte(plain), salt, 1<<n, r, p, keyLen)
	if err != nil {
		return false, fmt.Errorf("hash: scrypt: %w", err)
	}

	wantHex, err := hex.DecodeString(storedHash)
	if err != nil {
		// stored as raw hex string — compare as hex directly
		gotHex := hex.EncodeToString(got)
		return subtle.ConstantTimeCompare([]byte(gotHex), []byte(storedHash)) == 1, nil
	}
	return subtle.ConstantTimeCompare(got, wantHex) == 1, nil
}

func optInt(opts Options, key string, def int) int {
	switch v := opts[key].(type) {
	case int:
		return v
	case float64:
		return int(v)
	case int64:
		return int(v)
	}
	return def
}
