package identity

import (
	"context"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha1"
	"encoding/base32"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
)

// AuthenticatorType identifies the kind of MFA factor.
const AuthenticatorTypeTOTP = "totp"

// Authenticator holds an enrolled MFA factor.
type Authenticator struct {
	ID        string
	UserID    string
	Type      string
	Verified  bool
	CreatedAt time.Time
	UpdatedAt time.Time
}

// totpData is the JSON persisted in authenticators.data (encrypted).
type totpData struct {
	Secret string `json:"secret"` // base32-encoded TOTP secret
}

// EnrollTOTP creates an unverified TOTP authenticator for a user.
// Returns the base32-encoded secret the client must display as a QR code.
// Returns ErrAlreadyEnrolled if a verified TOTP factor already exists.
func (s *Store) EnrollTOTP(ctx context.Context, id, userID string) (secret string, err error) {
	// Block re-enroll if a verified TOTP already exists (BR-MIGRAR-001).
	var exists bool
	err = s.pool.QueryRow(ctx,
		`SELECT EXISTS(SELECT 1 FROM authenticators WHERE user_id=$1 AND type=$2 AND verified=true)`,
		userID, AuthenticatorTypeTOTP,
	).Scan(&exists)
	if err != nil {
		return "", err
	}
	if exists {
		return "", ErrAlreadyEnrolled
	}

	rawSecret := make([]byte, 20)
	if _, err := rand.Read(rawSecret); err != nil {
		return "", err
	}
	secret = base32.StdEncoding.WithPadding(base32.NoPadding).EncodeToString(rawSecret)

	data, err := json.Marshal(totpData{Secret: secret})
	if err != nil {
		return "", err
	}
	encData, err := encryptField(s.cipher, string(data))
	if err != nil {
		return "", fmt.Errorf("identity: encrypt totp data: %w", err)
	}

	_, err = s.pool.Exec(ctx,
		`INSERT INTO authenticators (id, user_id, type, verified, data)
         VALUES ($1,$2,$3,false,$4)
         ON CONFLICT (user_id, type) DO UPDATE SET data=$4, verified=false, updated_at=now()`,
		id, userID, AuthenticatorTypeTOTP, encData,
	)
	return secret, err
}

// ConfirmTOTP verifies an OTP code against the stored secret and marks the
// authenticator as verified. Returns ErrInvalidOTP if the code is wrong.
func (s *Store) ConfirmTOTP(ctx context.Context, userID, code string) error {
	var authID string
	var encData string
	err := s.pool.QueryRow(ctx,
		`SELECT id, data FROM authenticators WHERE user_id=$1 AND type=$2 AND verified=false`,
		userID, AuthenticatorTypeTOTP,
	).Scan(&authID, &encData)
	if err == pgx.ErrNoRows {
		return fmt.Errorf("identity: no pending TOTP enrollment")
	}
	if err != nil {
		return err
	}

	rawData, err := s.cipher.Decrypt(encData)
	if err != nil {
		return fmt.Errorf("identity: decrypt totp data: %w", err)
	}
	var td totpData
	if err := json.Unmarshal(rawData, &td); err != nil {
		return err
	}

	if !verifyTOTP(td.Secret, code, time.Now()) {
		return ErrInvalidOTP
	}

	_, err = s.pool.Exec(ctx,
		`UPDATE authenticators SET verified=true, updated_at=now() WHERE id=$1`, authID)
	return err
}

// ListAuthenticators returns all authenticators for a user.
func (s *Store) ListAuthenticators(ctx context.Context, userID string) ([]*Authenticator, error) {
	rows, err := s.pool.Query(ctx,
		`SELECT id, user_id, type, verified, created_at, updated_at
         FROM authenticators WHERE user_id=$1 ORDER BY created_at`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []*Authenticator
	for rows.Next() {
		var a Authenticator
		if err := rows.Scan(&a.ID, &a.UserID, &a.Type, &a.Verified, &a.CreatedAt, &a.UpdatedAt); err != nil {
			return nil, err
		}
		out = append(out, &a)
	}
	return out, rows.Err()
}

// RecomputeFactors returns the deduplicated factor list for a session based on
// its authenticators and provider. Currently: "email" if provider == "email",
// "totp" if a verified TOTP exists.
func (s *Store) RecomputeFactors(ctx context.Context, userID, provider string) ([]string, error) {
	rows, err := s.pool.Query(ctx,
		`SELECT type FROM authenticators WHERE user_id=$1 AND verified=true`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	seen := map[string]struct{}{}
	if provider != "" {
		seen[provider] = struct{}{}
	}
	for rows.Next() {
		var typ string
		if err := rows.Scan(&typ); err != nil {
			return nil, err
		}
		seen[typ] = struct{}{}
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	out := make([]string, 0, len(seen))
	for k := range seen {
		out = append(out, k)
	}
	return out, nil
}

// verifyTOTP checks a 6-digit OTP code against the base32-encoded secret using
// RFC 6238 (TOTP). Accepts ±1 time step (30 s) to allow clock skew.
func verifyTOTP(secret, code string, now time.Time) bool {
	key, err := base32.StdEncoding.WithPadding(base32.NoPadding).DecodeString(secret)
	if err != nil {
		return false
	}
	t := now.Unix() / 30
	for delta := int64(-1); delta <= 1; delta++ {
		if totpCode(key, t+delta) == code {
			return true
		}
	}
	return false
}

func totpCode(key []byte, counter int64) string {
	msg := make([]byte, 8)
	binary.BigEndian.PutUint64(msg, uint64(counter))
	mac := hmac.New(sha1.New, key)
	mac.Write(msg)
	h := mac.Sum(nil)
	offset := h[len(h)-1] & 0x0F
	code := (binary.BigEndian.Uint32(h[offset:offset+4]) & 0x7FFFFFFF) % 1_000_000
	return fmt.Sprintf("%06d", code)
}

// Sentinel errors.
var (
	ErrAlreadyEnrolled = fmt.Errorf("identity: totp already enrolled and verified")
	ErrInvalidOTP      = fmt.Errorf("identity: invalid OTP code")
)
