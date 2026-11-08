package identity

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"

	"github.com/danielvm-git/gowrite/platform/crypto"
)

// Session represents an authenticated session.
type Session struct {
	ID                        string
	UserID                    string
	Provider                  string
	ProviderUID               string
	ProviderAccessTokenExpiry *time.Time
	Secret                    string // plaintext — caller must hash before storing elsewhere
	UserAgent                 string
	IP                        string
	CountryCode               string
	Factors                   []string
	Expire                    time.Time
	MFAUpdatedAt              *time.Time
	CreatedAt                 time.Time
	UpdatedAt                 time.Time
}

// CreateSessionParams is the input for CreateSession.
type CreateSessionParams struct {
	ID                        string
	UserID                    string
	Provider                  string
	ProviderUID               string
	ProviderAccessToken       *string
	ProviderAccessTokenExpiry *time.Time
	ProviderRefreshToken      *string
	UserAgent                 string
	IP                        string
	CountryCode               string
	Factors                   []string
	Expire                    time.Time
}

// sessionStore is embedded in Store — no separate type needed.

// CreateSession generates a random secret token, encrypts it, and persists a
// new session row. The plaintext secret is returned once for the caller to
// hand to the client.
func (s *Store) CreateSession(ctx context.Context, p CreateSessionParams) (*Session, error) {
	raw := make([]byte, 32)
	if _, err := rand.Read(raw); err != nil {
		return nil, err
	}
	secret := hex.EncodeToString(raw)

	encSecret, err := encryptField(s.cipher, secret)
	if err != nil {
		return nil, fmt.Errorf("identity: encrypt session secret: %w", err)
	}
	encAccess, err := encryptFieldPtr(s.cipher, p.ProviderAccessToken)
	if err != nil {
		return nil, fmt.Errorf("identity: encrypt providerAccessToken: %w", err)
	}
	encRefresh, err := encryptFieldPtr(s.cipher, p.ProviderRefreshToken)
	if err != nil {
		return nil, fmt.Errorf("identity: encrypt providerRefreshToken: %w", err)
	}

	if p.Factors == nil {
		p.Factors = []string{}
	}

	const q = `
INSERT INTO sessions
    (id, user_id, provider, provider_uid,
     provider_access_token, provider_access_token_expiry,
     provider_refresh_token,
     secret, user_agent, ip, country_code, factors, expire)
VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13)
RETURNING id, user_id, provider, provider_uid,
          provider_access_token_expiry,
          user_agent, ip, country_code, factors, expire, mfa_updated_at,
          created_at, updated_at`

	row := s.pool.QueryRow(ctx, q,
		p.ID, p.UserID, p.Provider, p.ProviderUID,
		encAccess, p.ProviderAccessTokenExpiry, encRefresh,
		encSecret, p.UserAgent, p.IP, p.CountryCode, p.Factors, p.Expire,
	)

	sess, err := scanSession(row)
	if err != nil {
		return nil, err
	}
	sess.Secret = secret // plaintext for the caller
	return sess, nil
}

// DeleteSession removes a session by ID.
func (s *Store) DeleteSession(ctx context.Context, id string) error {
	_, err := s.pool.Exec(ctx, `DELETE FROM sessions WHERE id=$1`, id)
	return err
}

// DeleteUserSessions removes all sessions for a user.
func (s *Store) DeleteUserSessions(ctx context.Context, userID string) error {
	_, err := s.pool.Exec(ctx, `DELETE FROM sessions WHERE user_id=$1`, userID)
	return err
}

// ListSessions returns all non-expired sessions for a user, most recent first.
func (s *Store) ListSessions(ctx context.Context, userID string) ([]*Session, error) {
	const q = `
SELECT id, user_id, provider, provider_uid,
       provider_access_token_expiry,
       user_agent, ip, country_code, factors, expire, mfa_updated_at,
       created_at, updated_at
FROM sessions WHERE user_id=$1 AND expire > now()
ORDER BY created_at DESC`

	rows, err := s.pool.Query(ctx, q, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []*Session
	for rows.Next() {
		sess, err := scanSession(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, sess)
	}
	return out, rows.Err()
}

// UpdateSessionFactors rewrites the factors array for a session (factor recompute).
func (s *Store) UpdateSessionFactors(ctx context.Context, sessionID string, factors []string) error {
	if factors == nil {
		factors = []string{}
	}
	_, err := s.pool.Exec(ctx,
		`UPDATE sessions SET factors=$2, mfa_updated_at=now(), updated_at=now() WHERE id=$1`,
		sessionID, factors)
	return err
}

func scanSession(row pgx.Row) (*Session, error) {
	var sess Session
	err := row.Scan(
		&sess.ID, &sess.UserID, &sess.Provider, &sess.ProviderUID,
		&sess.ProviderAccessTokenExpiry,
		&sess.UserAgent, &sess.IP, &sess.CountryCode,
		&sess.Factors, &sess.Expire, &sess.MFAUpdatedAt,
		&sess.CreatedAt, &sess.UpdatedAt,
	)
	if err == pgx.ErrNoRows {
		return nil, nil
	}
	return &sess, err
}

func encryptField(c *crypto.Cipher, v string) (string, error) {
	if c == nil {
		return v, nil
	}
	return c.Encrypt([]byte(v))
}

func encryptFieldPtr(c *crypto.Cipher, v *string) (*string, error) {
	if v == nil {
		return nil, nil
	}
	enc, err := encryptField(c, *v)
	if err != nil {
		return nil, err
	}
	return &enc, nil
}
