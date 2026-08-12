// Package crypto implements AES-128-GCM encryption at rest for the 25 legacy
// encrypted fields, in the same wire format Appwrite's PHP `encrypt` filter
// produces, so ciphertext migrates without re-encryption at read time.
//
// Legacy format (app/init.php `encrypt` filter): a JSON object
//
//	{"data":"<base64 ciphertext>","method":"aes-128-gcm",
//	 "iv":"<hex nonce>","tag":"<hex gcm tag>","version":"1"}
//
// where `data` is openssl_encrypt(..., OPENSSL_RAW_DATA=0) i.e. base64 of the
// raw ciphertext (tag stored separately). See BR-MIGRAR-006.
package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
)

const (
	methodAES128GCM = "aes-128-gcm"
	formatVersion   = "1"
	keyLen          = 16 // AES-128
	nonceLen        = 12 // GCM standard IV
	tagLen          = 16 // GCM tag
)

// ErrInvalidKey is returned when the key is not exactly 16 bytes.
var ErrInvalidKey = errors.New("crypto: encryption key must be 16 bytes (AES-128)")

// Cipher encrypts and decrypts values with a fixed AES-128 key.
type Cipher struct {
	key []byte
	gcm cipher.AEAD
}

// New builds a Cipher from a 16-byte key.
func New(key []byte) (*Cipher, error) {
	if len(key) != keyLen {
		return nil, fmt.Errorf("%w (got %d bytes)", ErrInvalidKey, len(key))
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("crypto: new cipher: %w", err)
	}
	gcm, err := cipher.NewGCMWithNonceSize(block, nonceLen)
	if err != nil {
		return nil, fmt.Errorf("crypto: new gcm: %w", err)
	}
	return &Cipher{key: append([]byte(nil), key...), gcm: gcm}, nil
}

// payload is the on-disk JSON envelope (legacy-compatible).
type payload struct {
	Data    string `json:"data"`
	Method  string `json:"method"`
	IV      string `json:"iv"`
	Tag     string `json:"tag"`
	Version string `json:"version"`
}

// Encrypt returns the legacy JSON envelope string for plaintext.
func (c *Cipher) Encrypt(plaintext []byte) (string, error) {
	nonce := make([]byte, nonceLen)
	if _, err := rand.Read(nonce); err != nil {
		return "", fmt.Errorf("crypto: read nonce: %w", err)
	}
	// gcm.Seal appends the 16-byte tag to the ciphertext; split them to match
	// the legacy layout that stores ciphertext and tag separately.
	sealed := c.gcm.Seal(nil, nonce, plaintext, nil)
	ct, tag := sealed[:len(sealed)-tagLen], sealed[len(sealed)-tagLen:]

	env := payload{
		Data:    base64.StdEncoding.EncodeToString(ct),
		Method:  methodAES128GCM,
		IV:      hex.EncodeToString(nonce),
		Tag:     hex.EncodeToString(tag),
		Version: formatVersion,
	}
	out, err := json.Marshal(env)
	if err != nil {
		return "", fmt.Errorf("crypto: marshal payload: %w", err)
	}
	return string(out), nil
}

// Decrypt parses a legacy JSON envelope and returns the plaintext.
func (c *Cipher) Decrypt(encoded string) ([]byte, error) {
	var env payload
	if err := json.Unmarshal([]byte(encoded), &env); err != nil {
		return nil, fmt.Errorf("crypto: unmarshal payload: %w", err)
	}
	if env.Method != methodAES128GCM {
		return nil, fmt.Errorf("crypto: unsupported method %q", env.Method)
	}
	ct, err := base64.StdEncoding.DecodeString(env.Data)
	if err != nil {
		return nil, fmt.Errorf("crypto: decode data: %w", err)
	}
	nonce, err := hex.DecodeString(env.IV)
	if err != nil {
		return nil, fmt.Errorf("crypto: decode iv: %w", err)
	}
	tag, err := hex.DecodeString(env.Tag)
	if err != nil {
		return nil, fmt.Errorf("crypto: decode tag: %w", err)
	}
	if len(nonce) != nonceLen {
		return nil, fmt.Errorf("crypto: bad nonce length %d", len(nonce))
	}
	// Recombine ciphertext||tag for gcm.Open.
	sealed := append(append([]byte(nil), ct...), tag...)
	plaintext, err := c.gcm.Open(nil, nonce, sealed, nil)
	if err != nil {
		return nil, fmt.Errorf("crypto: authentication failed: %w", err)
	}
	return plaintext, nil
}
