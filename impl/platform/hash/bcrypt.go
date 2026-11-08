package hash

import "golang.org/x/crypto/bcrypt"

// verifyBcrypt verifies plain against a standard bcrypt hash produced by PHP's
// password_hash($value, PASSWORD_BCRYPT, ...).
func verifyBcrypt(plain, storedHash string) (bool, error) {
	err := bcrypt.CompareHashAndPassword([]byte(storedHash), []byte(plain))
	if err == bcrypt.ErrMismatchedHashAndPassword {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil
}
