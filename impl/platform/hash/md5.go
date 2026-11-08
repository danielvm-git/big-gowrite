package hash

import (
	"crypto/md5"
	"crypto/subtle"
	"encoding/hex"
)

func verifyMD5(plain, storedHash string) (bool, error) {
	sum := md5.Sum([]byte(plain))
	got := hex.EncodeToString(sum[:])
	return subtle.ConstantTimeCompare([]byte(got), []byte(storedHash)) == 1, nil
}
