package identity

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/argon2"

	"github.com/danielvm-git/gowrite/platform/crypto"
	"github.com/danielvm-git/gowrite/platform/hash"
)

// User represents a gowrite user.
type User struct {
	ID                string
	Name              string
	Email             *string
	Phone             *string
	Status            bool
	Labels            []string
	Hash              string
	HashOptions       hash.Options
	PasswordUpdate    *time.Time
	Prefs             map[string]any
	Registration      time.Time
	EmailVerification bool
	PhoneVerification bool
	Reset             bool
	MFA               bool
	AccessedAt        *time.Time
	CreatedAt         time.Time
	UpdatedAt         time.Time
}

// CreateUserParams is the input for CreateUser.
type CreateUserParams struct {
	ID          string
	Name        string
	Email       *string
	Phone       *string
	PlainPwd    string // empty for OAuth users
	Hash        string // empty defaults to argon2
	HashOptions hash.Options
}

// Store is a thin data-access layer for identity tables.
type Store struct {
	pool   *pgxpool.Pool
	cipher *crypto.Cipher // may be nil in unit tests that skip encrypt
}

// NewStore returns a Store.
func NewStore(pool *pgxpool.Pool, c *crypto.Cipher) *Store {
	return &Store{pool: pool, cipher: c}
}

// CreateUser inserts a new user. PlainPwd is encrypted and stored as an
// AES-128-GCM envelope; Hash/HashOptions track the algorithm used.
func (s *Store) CreateUser(ctx context.Context, p CreateUserParams) (*User, error) {
	if p.Hash == "" {
		p.Hash = hash.AlgArgon2
	}

	var encPwd *string
	if p.PlainPwd != "" && s.cipher != nil {
		phc, err := argon2PHC(p.PlainPwd)
		if err != nil {
			return nil, fmt.Errorf("identity: hash password: %w", err)
		}
		enc, err := s.cipher.Encrypt([]byte(phc))
		if err != nil {
			return nil, fmt.Errorf("identity: encrypt password: %w", err)
		}
		encPwd = &enc
	}

	optsJSON, err := json.Marshal(p.HashOptions)
	if err != nil {
		return nil, err
	}

	now := time.Now().UTC()
	const q = `
INSERT INTO users
    (id, name, email, phone, password, hash, hash_options, registration, created_at, updated_at)
VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$8,$8)
RETURNING id, name, email, phone, status, labels, hash, hash_options,
          password_update, prefs, registration,
          email_verification, phone_verification, reset, mfa,
          accessed_at, created_at, updated_at`

	row := s.pool.QueryRow(ctx, q,
		p.ID, p.Name, p.Email, p.Phone, encPwd,
		p.Hash, string(optsJSON), now)
	return scanUser(row)
}

// GetUser returns a user by ID.
func (s *Store) GetUser(ctx context.Context, id string) (*User, error) {
	const q = `
SELECT id, name, email, phone, status, labels, hash, hash_options,
       password_update, prefs, registration,
       email_verification, phone_verification, reset, mfa,
       accessed_at, created_at, updated_at
FROM users WHERE id = $1`
	return scanUser(s.pool.QueryRow(ctx, q, id))
}

// UpdateUserName changes the user's display name.
func (s *Store) UpdateUserName(ctx context.Context, id, name string) error {
	_, err := s.pool.Exec(ctx,
		`UPDATE users SET name=$2, updated_at=now() WHERE id=$1`, id, name)
	return err
}

// DeleteUser removes a user and cascades to all child rows.
func (s *Store) DeleteUser(ctx context.Context, id string) error {
	_, err := s.pool.Exec(ctx, `DELETE FROM users WHERE id=$1`, id)
	return err
}

// VerifyPassword checks plain against the stored encrypted password.
func (s *Store) VerifyPassword(ctx context.Context, userID, plain string) (bool, error) {
	var encPwd *string
	var alg string
	var optsRaw []byte

	err := s.pool.QueryRow(ctx,
		`SELECT password, hash, hash_options FROM users WHERE id=$1`, userID,
	).Scan(&encPwd, &alg, &optsRaw)
	if err == pgx.ErrNoRows {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	if encPwd == nil || s.cipher == nil {
		return false, nil
	}

	storedPlain, err := s.cipher.Decrypt(*encPwd)
	if err != nil {
		return false, fmt.Errorf("identity: decrypt password: %w", err)
	}

	var opts hash.Options
	_ = json.Unmarshal(optsRaw, &opts)

	return hash.Verify(plain, string(storedPlain), alg, opts)
}

// argon2PHC hashes plain with argon2id and returns a PHC-format string.
func argon2PHC(plain string) (string, error) {
	salt := make([]byte, 16)
	if _, err := rand.Read(salt); err != nil {
		return "", err
	}
	h := argon2.IDKey([]byte(plain), salt, 4, 65536, 3, 32)
	return fmt.Sprintf("$argon2id$v=19$m=65536,t=4,p=3$%s$%s",
		base64.RawStdEncoding.EncodeToString(salt),
		base64.RawStdEncoding.EncodeToString(h)), nil
}

func scanUser(row pgx.Row) (*User, error) {
	var u User
	var optsRaw []byte
	var prefsRaw []byte

	err := row.Scan(
		&u.ID, &u.Name, &u.Email, &u.Phone,
		&u.Status, &u.Labels,
		&u.Hash, &optsRaw,
		&u.PasswordUpdate, &prefsRaw, &u.Registration,
		&u.EmailVerification, &u.PhoneVerification, &u.Reset, &u.MFA,
		&u.AccessedAt, &u.CreatedAt, &u.UpdatedAt,
	)
	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	_ = json.Unmarshal(optsRaw, &u.HashOptions)
	_ = json.Unmarshal(prefsRaw, &u.Prefs)
	return &u, nil
}
