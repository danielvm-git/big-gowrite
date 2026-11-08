package hash

import (
	"crypto/md5"
	"strings"
)

// itoa64 is the PHPass alphabet used for encode64 and the salt alphabet.
const phpassItoa64 = "./0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"

// verifyPHPass verifies plain against a PHPass hash (WordPress/Drupal style).
// Supports both portable $P$ hashes and bcrypt $2a$ hashes (which bcrypt handles).
func verifyPHPass(plain, storedHash string) (bool, error) {
	// $2a$ hashes are bcrypt; delegate.
	if strings.HasPrefix(storedHash, "$2") {
		return verifyBcrypt(plain, storedHash)
	}
	got := phpassCryptPrivate(plain, storedHash)
	if got[0] == '*' {
		return false, nil
	}
	return got == storedHash, nil
}

func phpassCryptPrivate(password, setting string) string {
	output := "*0"
	if strings.HasPrefix(setting, output) {
		output = "*1"
	}

	id := setting[:3]
	if id != "$P$" && id != "$H$" {
		return output
	}

	countLog2 := strings.IndexByte(phpassItoa64, setting[3])
	if countLog2 < 7 || countLog2 > 30 {
		return output
	}

	count := 1 << countLog2
	salt := setting[4:12]
	if len(salt) != 8 {
		return output
	}

	hashBytes := md5.Sum([]byte(salt + password))
	h := hashBytes[:]
	for count--; count > 0; count-- {
		next := md5.Sum(append(h, []byte(password)...))
		h = next[:]
	}

	out := setting[:12] + phpassEncode64(h, 16)
	return out
}

func phpassEncode64(input []byte, count int) string {
	var output []byte
	i := 0
	for i < count {
		value := int(input[i])
		i++
		output = append(output, phpassItoa64[value&0x3F])
		if i < count {
			value |= int(input[i]) << 8
		}
		output = append(output, phpassItoa64[(value>>6)&0x3F])
		if i >= count {
			break
		}
		i++
		if i < count {
			value |= int(input[i]) << 16
		}
		output = append(output, phpassItoa64[(value>>12)&0x3F])
		if i >= count {
			break
		}
		i++
		output = append(output, phpassItoa64[(value>>18)&0x3F])
	}
	return string(output)
}
