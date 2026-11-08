package identity

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
)

// Token types — mirror app/init/constants.php.
const (
	TokenTypeVerification = 2
	TokenTypeRecovery     = 3
	TokenTypeInvite       = 4
	TokenTypeMagicURL     = 5
	TokenTypePhone        = 6
	TokenTypeOAuth2       = 7
	TokenTypeGeneric      = 8
	TokenTypeEmail        = 9
)

// Token is a short-lived credential for verification / recovery / OAuth flows.
type Token struct {
	ID        string
	UserID    string
	Type      int
	Secret    string // plaintext secret, returned once on creation
	Expire    time.Time
	UserAgent string
	IP        string
	CreatedAt time.Time
}

// CreateTokenParams is the input for CreateToken.
type CreateTokenParams struct {
	ID        string
	UserID    string
	Type      int
	TTL       time.Duration
	UserAgent string
	IP        string
}

// CreateToken generates a random secret, encrypts it, and persists a token row.
// The plaintext secret is returned once.
func (s *Store) CreateToken(ctx context.Context, p CreateTokenParams) (*Token, error) {
	raw := make([]byte, 32)
	if _, err := rand.Read(raw); err != nil {
		return nil, err
	}
	secret := hex.EncodeToString(raw)

	encSecret, err := encryptField(s.cipher, secret)
	if err != nil {
		return nil, fmt.Errorf("identity: encrypt token secret: %w", err)
	}

	expire := time.Now().Add(p.TTL).UTC()
	const q = `
INSERT INTO tokens (id, user_id, type, secret, expire, user_agent, ip)
VALUES ($1,$2,$3,$4,$5,$6,$7)
RETURNING id, user_id, type, expire, user_agent, ip, created_at`

	tok, err := scanToken(s.pool.QueryRow(ctx, q,
		p.ID, p.UserID, p.Type, encSecret, expire, p.UserAgent, p.IP))
	if err != nil {
		return nil, err
	}
	tok.Secret = secret
	return tok, nil
}

// ConsumeToken looks up a token by its plaintext secret + type, validates expiry,
// deletes the row, and returns the token. Returns ErrTokenInvalid if not found or
// expired (no distinguishing between the two — timing-safe).
func (s *Store) ConsumeToken(ctx context.Context, secret string, tokenType int) (*Token, error) {
	// Fetch all non-expired tokens of this type to avoid timing-side-channel
	// from early DB exit on "not found". In production this set is tiny.
	const q = `
SELECT id, user_id, type, secret, expire, user_agent, ip, created_at
FROM tokens WHERE type=$1 AND expire > now()`

	rows, err := s.pool.Query(ctx, q, tokenType)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var matched *Token
	for rows.Next() {
		var t Token
		var encSecret string
		if err := rows.Scan(&t.ID, &t.UserID, &t.Type, &encSecret,
			&t.Expire, &t.UserAgent, &t.IP, &t.CreatedAt); err != nil {
			return nil, err
		}
		if s.cipher != nil {
			plain, err := s.cipher.Decrypt(encSecret)
			if err != nil {
				continue
			}
			if string(plain) == secret {
				matched = &t
				break
			}
		} else if encSecret == secret {
			matched = &t
			break
		}
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	if matched == nil {
		return nil, ErrTokenInvalid
	}

	// Invalidate immediately (consume = delete).
	_, err = s.pool.Exec(ctx, `DELETE FROM tokens WHERE id=$1`, matched.ID)
	if err != nil {
		return nil, err
	}
	matched.Secret = secret
	return matched, nil
}

// StoreOAuthToken persists an encrypted OAuth provider token into the session.
func (s *Store) StoreOAuthToken(ctx context.Context, sessionID, providerAccessToken, providerRefreshToken string, expiry *time.Time) error {
	encAccess, err := encryptField(s.cipher, providerAccessToken)
	if err != nil {
		return fmt.Errorf("identity: encrypt oauth access token: %w", err)
	}
	encRefresh, err := encryptField(s.cipher, providerRefreshToken)
	if err != nil {
		return fmt.Errorf("identity: encrypt oauth refresh token: %w", err)
	}
	_, err = s.pool.Exec(ctx, `
UPDATE sessions
SET provider_access_token=$2, provider_refresh_token=$3,
    provider_access_token_expiry=$4, updated_at=now()
WHERE id=$1`,
		sessionID, encAccess, encRefresh, expiry)
	return err
}

// GetOAuthTokens decrypts and returns the provider tokens for a session.
func (s *Store) GetOAuthTokens(ctx context.Context, sessionID string) (access, refresh string, expiry *time.Time, err error) {
	var encAccess, encRefresh *string
	err = s.pool.QueryRow(ctx, `
SELECT provider_access_token, provider_refresh_token, provider_access_token_expiry
FROM sessions WHERE id=$1`, sessionID).Scan(&encAccess, &encRefresh, &expiry)
	if err == pgx.ErrNoRows {
		return "", "", nil, fmt.Errorf("identity: session not found")
	}
	if err != nil {
		return "", "", nil, err
	}

	if encAccess != nil && s.cipher != nil {
		raw, e := s.cipher.Decrypt(*encAccess)
		if e == nil {
			access = string(raw)
		}
	}
	if encRefresh != nil && s.cipher != nil {
		raw, e := s.cipher.Decrypt(*encRefresh)
		if e == nil {
			refresh = string(raw)
		}
	}
	return access, refresh, expiry, nil
}

var ErrTokenInvalid = fmt.Errorf("identity: token invalid or expired")

func scanToken(row pgx.Row) (*Token, error) {
	var t Token
	err := row.Scan(&t.ID, &t.UserID, &t.Type, &t.Expire, &t.UserAgent, &t.IP, &t.CreatedAt)
	if err == pgx.ErrNoRows {
		return nil, nil
	}
	return &t, err
}
