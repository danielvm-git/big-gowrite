package hash

import (
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"crypto/subtle"
	"encoding/hex"
	"fmt"
	"hash"
)

// verifySHA verifies plain against a hex-encoded SHA digest.
// opts["version"] selects the algorithm (default: "sha256").
func verifySHA(plain, storedHash string, opts Options) (bool, error) {
	version, _ := opts["version"].(string)
	if version == "" {
		version = "sha256"
	}

	var h hash.Hash
	switch version {
	case "sha1":
		h = sha1.New()
	case "sha224":
		h = sha256.New224()
	case "sha256":
		h = sha256.New()
	case "sha384":
		h = sha512.New384()
	case "sha512":
		h = sha512.New()
	default:
		return false, fmt.Errorf("hash: unsupported sha version %q", version)
	}
	h.Write([]byte(plain))
	got := hex.EncodeToString(h.Sum(nil))
	return subtle.ConstantTimeCompare([]byte(got), []byte(storedHash)) == 1, nil
}
