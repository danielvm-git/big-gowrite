package hash

import (
	"crypto/subtle"
	"encoding/base64"
	"fmt"
	"strconv"
	"strings"

	"golang.org/x/crypto/argon2"
)

// verifyArgon2 verifies a plain password against a PHC-format argon2id hash
// produced by PHP's password_hash($value, PASSWORD_ARGON2ID, ...).
// PHC format: $argon2id$v=19$m=65536,t=4,p=3$<b64salt>$<b64hash>
func verifyArgon2(plain, phcHash string) (bool, error) {
	parts := strings.Split(phcHash, "$")
	// parts[0]="" parts[1]="argon2id" parts[2]="v=19" parts[3]="m=...,t=...,p=..."
	// parts[4]=b64salt parts[5]=b64hash
	if len(parts) != 6 {
		return false, fmt.Errorf("hash: invalid argon2 PHC format")
	}
	if parts[1] != "argon2id" {
		return false, fmt.Errorf("hash: unsupported argon2 variant %q", parts[1])
	}

	params := make(map[string]uint32)
	for _, kv := range strings.Split(parts[3], ",") {
		pair := strings.SplitN(kv, "=", 2)
		if len(pair) != 2 {
			continue
		}
		v, err := strconv.ParseUint(pair[1], 10, 32)
		if err != nil {
			return false, fmt.Errorf("hash: bad argon2 param %q: %w", kv, err)
		}
		params[pair[0]] = uint32(v)
	}

	salt, err := base64.RawStdEncoding.DecodeString(parts[4])
	if err != nil {
		return false, fmt.Errorf("hash: argon2 salt decode: %w", err)
	}
	wantHash, err := base64.RawStdEncoding.DecodeString(parts[5])
	if err != nil {
		return false, fmt.Errorf("hash: argon2 hash decode: %w", err)
	}

	m := params["m"]
	t := params["t"]
	p := uint8(params["p"])
	keyLen := uint32(len(wantHash))

	got := argon2.IDKey([]byte(plain), salt, t, m, p, keyLen)
	return subtle.ConstantTimeCompare(got, wantHash) == 1, nil
}
