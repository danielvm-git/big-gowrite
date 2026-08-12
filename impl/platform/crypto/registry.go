package crypto

// EncryptedField identifies a legacy encrypted-at-rest field by collection and
// attribute. There are 25 such fields (BR-MIGRAR-006,
// specs/database/business-rules.md).
type EncryptedField struct {
	Collection string
	Attribute  string
}

// EncryptedFields is the registry of the 25 fields the `encrypt` filter covered
// in the legacy schema. Write sites for these fields MUST call Cipher.Encrypt.
var EncryptedFields = []EncryptedField{
	{"users", "password"},
	{"users", "hash"},
	{"users", "hashOptions"},
	{"users", "mfaRecoveryCodes"},
	{"sessions", "secret"},
	{"sessions", "providerAccessToken"},
	{"sessions", "providerRefreshToken"},
	{"tokens", "secret"},
	{"authenticators", "data"},
	{"keys", "secret"},
	{"projects", "keys"},
	{"projects", "services"},
	{"projects", "smtpPassword"},
	{"projects", "smtpUsername"},
	{"platforms", "key"},
	{"webhooks", "signatureKey"},
	{"webhooks", "httpPass"},
	{"installations", "personalAccessToken"},
	{"installations", "providerAccessToken"},
	{"vcsCommentLocks", "secret"},
	{"providers", "credentials"},
	{"providers", "options"},
	{"migrations", "credentials"},
	{"certificates", "privateKey"},
	{"functions", "vars"},
}

// FieldCount returns the number of registered encrypted fields.
func FieldCount() int { return len(EncryptedFields) }

// ReEncrypt decrypts a value produced by src and re-encrypts it under dst.
// Used by the data-migration re-crypto step (R-D3) when rotating keys.
func ReEncrypt(src, dst *Cipher, encoded string) (string, error) {
	plaintext, err := src.Decrypt(encoded)
	if err != nil {
		return "", err
	}
	return dst.Encrypt(plaintext)
}
