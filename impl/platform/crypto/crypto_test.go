package crypto

import (
	"bytes"
	"os"
	"strings"
	"testing"
)

var testKey = mustHex("000102030405060708090a0b0c0d0e0f")

func mustHex(s string) []byte {
	c, err := FromHexKey(s)
	if err != nil {
		panic(err)
	}
	return c.key
}

func newCipher(t *testing.T) *Cipher {
	t.Helper()
	c, err := New(testKey)
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	return c
}

func TestRoundTrip(t *testing.T) {
	c := newCipher(t)
	for _, plain := range []string{"", "hunter2", "a longer secret value with spaces", "unicode: café ☕"} {
		enc, err := c.Encrypt([]byte(plain))
		if err != nil {
			t.Fatalf("Encrypt(%q): %v", plain, err)
		}
		got, err := c.Decrypt(enc)
		if err != nil {
			t.Fatalf("Decrypt(%q): %v", plain, err)
		}
		if string(got) != plain {
			t.Fatalf("round-trip = %q, want %q", got, plain)
		}
	}
}

func TestCiphertextIsNotPlaintext(t *testing.T) {
	c := newCipher(t)
	secret := "super-secret-token"
	enc, err := c.Encrypt([]byte(secret))
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(enc, secret) {
		t.Fatalf("plaintext leaked into ciphertext envelope: %s", enc)
	}
	// Envelope must carry iv + tag (AEAD).
	for _, field := range []string{`"iv"`, `"tag"`, `"aes-128-gcm"`} {
		if !strings.Contains(enc, field) {
			t.Fatalf("envelope missing %s: %s", field, enc)
		}
	}
}

func TestNonceIsRandomPerEncryption(t *testing.T) {
	c := newCipher(t)
	a, _ := c.Encrypt([]byte("same"))
	b, _ := c.Encrypt([]byte("same"))
	if a == b {
		t.Fatal("two encryptions of the same plaintext produced identical output (nonce reuse)")
	}
}

func TestTamperDetection(t *testing.T) {
	c := newCipher(t)
	enc, _ := c.Encrypt([]byte("integrity"))
	// Flip a base64 char in the data field.
	tampered := strings.Replace(enc, `"data":"`, `"data":"X`, 1)
	if _, err := c.Decrypt(tampered); err == nil {
		t.Fatal("Decrypt should fail on tampered ciphertext")
	}
}

func TestKeyValidation(t *testing.T) {
	if _, err := New([]byte("short")); err == nil {
		t.Fatal("New should reject a non-16-byte key")
	}
	if _, err := FromHexKey("zzzz"); err == nil {
		t.Fatal("FromHexKey should reject invalid hex")
	}
}

func TestFromEnv(t *testing.T) {
	t.Setenv(EnvKey, "000102030405060708090a0b0c0d0e0f")
	if _, err := FromEnv(); err != nil {
		t.Fatalf("FromEnv with valid key: %v", err)
	}
	if err := os.Unsetenv(EnvKey); err != nil {
		t.Fatalf("unset %s: %v", EnvKey, err)
	}
	if _, err := FromEnv(); err == nil {
		t.Fatal("FromEnv should error when key is unset")
	}
}

// TestLegacyParity decrypts a ciphertext produced INDEPENDENTLY by PHP's
// openssl_encrypt (testdata/legacy_fixture.json) with the same key, proving the
// Go decoder is byte-compatible with the legacy `encrypt` filter wire format.
func TestLegacyParity(t *testing.T) {
	raw, err := os.ReadFile("testdata/legacy_fixture.json")
	if err != nil {
		t.Skipf("no legacy fixture present: %v", err)
	}
	c, err := FromHexKey("000102030405060708090a0b0c0d0e0f")
	if err != nil {
		t.Fatal(err)
	}
	got, err := c.Decrypt(strings.TrimSpace(string(raw)))
	if err != nil {
		t.Fatalf("Decrypt legacy fixture: %v", err)
	}
	const want = "oauth-access-token-42" // from testdata/legacy_fixture.meta
	if !bytes.Equal(got, []byte(want)) {
		t.Fatalf("legacy decrypt = %q, want %q", got, want)
	}
}

func TestFieldRegistry(t *testing.T) {
	if FieldCount() != 25 {
		t.Fatalf("registry has %d fields, want 25 (BR-MIGRAR-006)", FieldCount())
	}
	// no duplicates
	seen := map[EncryptedField]bool{}
	for _, f := range EncryptedFields {
		if seen[f] {
			t.Fatalf("duplicate field in registry: %+v", f)
		}
		seen[f] = true
		if f.Collection == "" || f.Attribute == "" {
			t.Fatalf("registry field has empty part: %+v", f)
		}
	}
}

func TestReEncryptRotatesKey(t *testing.T) {
	src, _ := New(testKey)
	dstKey := mustHex("0f0e0d0c0b0a09080706050403020100")
	dst, _ := New(dstKey)

	enc, _ := src.Encrypt([]byte("rotate-me"))
	rotated, err := ReEncrypt(src, dst, enc)
	if err != nil {
		t.Fatalf("ReEncrypt: %v", err)
	}
	// Old key can no longer read it; new key can.
	if _, err := src.Decrypt(rotated); err == nil {
		t.Fatal("source key should not decrypt after rotation")
	}
	got, err := dst.Decrypt(rotated)
	if err != nil || string(got) != "rotate-me" {
		t.Fatalf("dst.Decrypt = %q, %v; want rotate-me", got, err)
	}
}
