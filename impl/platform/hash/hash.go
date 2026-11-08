package hash

import (
	"encoding/json"
	"fmt"
)

// Algorithm names match the legacy PHP utopia/auth getName() strings.
type Algorithm = string

const (
	AlgArgon2    Algorithm = "argon2"
	AlgBcrypt    Algorithm = "bcrypt"
	AlgMD5       Algorithm = "md5"
	AlgSHA       Algorithm = "sha"
	AlgPHPass    Algorithm = "phpass"
	AlgScrypt    Algorithm = "scrypt"
	AlgScryptMod Algorithm = "scryptMod"
	AlgPlaintext Algorithm = "plaintext"
)

// Options is the JSON blob stored in users.hash_options.
type Options map[string]any

// ParseOptions unmarshals a raw JSON options string into Options.
func ParseOptions(raw string) (Options, error) {
	if raw == "" || raw == "{}" {
		return Options{}, nil
	}
	var o Options
	if err := json.Unmarshal([]byte(raw), &o); err != nil {
		return nil, err
	}
	return o, nil
}

// Verify checks whether plain matches the stored hash string using the given
// algorithm and options. Returns ErrUnknownAlgorithm for unrecognised names.
func Verify(plain, storedHash string, algo Algorithm, opts Options) (bool, error) {
	switch algo {
	case AlgArgon2:
		return verifyArgon2(plain, storedHash)
	case AlgBcrypt:
		return verifyBcrypt(plain, storedHash)
	case AlgMD5:
		return verifyMD5(plain, storedHash)
	case AlgSHA:
		return verifySHA(plain, storedHash, opts)
	case AlgPHPass:
		return verifyPHPass(plain, storedHash)
	case AlgScrypt:
		return verifyScrypt(plain, storedHash, opts)
	case AlgScryptMod:
		return verifyScryptMod(plain, storedHash, opts)
	case AlgPlaintext:
		return plain == storedHash, nil
	default:
		return false, fmt.Errorf("hash: unknown algorithm %q", algo)
	}
}
