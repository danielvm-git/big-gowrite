package identity

import (
	"context"
	"encoding/base32"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/danielvm-git/gowrite/platform/crypto"
	"github.com/danielvm-git/gowrite/platform/hash"
)

const testKey = "000102030405060708090a0b0c0d0e0f"

func testStore(t *testing.T) *Store {
	t.Helper()
	dsn := os.Getenv("GOWRITE_TEST_DATABASE_URL")
	if dsn == "" {
		t.Skip("GOWRITE_TEST_DATABASE_URL not set; skipping Postgres integration test")
	}
	pool, err := pgxpool.New(context.Background(), dsn)
	if err != nil {
		t.Fatalf("connect: %v", err)
	}
	t.Cleanup(pool.Close)

	cipher, err := crypto.FromHexKey(testKey)
	if err != nil {
		t.Fatalf("cipher: %v", err)
	}

	for _, tbl := range []string{"tokens", "challenges", "authenticators", "sessions", "users"} {
		if _, err := pool.Exec(context.Background(),
			fmt.Sprintf("DELETE FROM %s", tbl)); err != nil {
			t.Fatalf("clean %s: %v", tbl, err)
		}
	}
	return NewStore(pool, cipher)
}

// e02s01 — Users & accounts

func TestUserCRUD(t *testing.T) {
	s := testStore(t)
	ctx := context.Background()

	u, err := s.CreateUser(ctx, CreateUserParams{
		ID:       "user-1",
		Name:     "Alice",
		Email:    ptr("alice@example.com"),
		PlainPwd: "s3cr3t!",
	})
	if err != nil || u == nil {
		t.Fatalf("CreateUser: %v", err)
	}
	if u.ID != "user-1" || u.Name != "Alice" {
		t.Fatalf("unexpected user: %+v", u)
	}

	got, err := s.GetUser(ctx, "user-1")
	if err != nil || got == nil {
		t.Fatalf("GetUser: %v", err)
	}
	if got.Email == nil || *got.Email != "alice@example.com" {
		t.Fatalf("email mismatch: %v", got.Email)
	}

	if err := s.UpdateUserName(ctx, "user-1", "Alice Smith"); err != nil {
		t.Fatalf("UpdateUserName: %v", err)
	}
	got2, _ := s.GetUser(ctx, "user-1")
	if got2.Name != "Alice Smith" {
		t.Fatalf("name not updated: %s", got2.Name)
	}

	if err := s.DeleteUser(ctx, "user-1"); err != nil {
		t.Fatalf("DeleteUser: %v", err)
	}
	gone, _ := s.GetUser(ctx, "user-1")
	if gone != nil {
		t.Fatal("user should be deleted")
	}
}

func TestHashParity(t *testing.T) {
	s := testStore(t)
	ctx := context.Background()

	_, err := s.CreateUser(ctx, CreateUserParams{
		ID:       "user-hash",
		Name:     "Bob",
		Email:    ptr("bob@example.com"),
		PlainPwd: "hunter2",
	})
	if err != nil {
		t.Fatalf("CreateUser: %v", err)
	}

	ok, err := s.VerifyPassword(ctx, "user-hash", "hunter2")
	if err != nil || !ok {
		t.Fatalf("correct password should verify: ok=%v err=%v", ok, err)
	}
	ok2, _ := s.VerifyPassword(ctx, "user-hash", "wrong")
	if ok2 {
		t.Fatal("wrong password should not verify")
	}
}

func TestUniqueness(t *testing.T) {
	s := testStore(t)
	ctx := context.Background()

	_, err := s.CreateUser(ctx, CreateUserParams{ID: "u1", Name: "Carol", Email: ptr("carol@example.com")})
	if err != nil {
		t.Fatalf("first create: %v", err)
	}
	_, err2 := s.CreateUser(ctx, CreateUserParams{ID: "u2", Name: "Carol2", Email: ptr("carol@example.com")})
	if err2 == nil {
		t.Fatal("duplicate email should be rejected")
	}
}

// e02s02 — Sessions & MFA/TOTP

func TestSession(t *testing.T) {
	s := testStore(t)
	ctx := context.Background()

	_, _ = s.CreateUser(ctx, CreateUserParams{ID: "sess-user", Name: "Dan"})

	sess, err := s.CreateSession(ctx, CreateSessionParams{
		ID:          "sess-1",
		UserID:      "sess-user",
		Provider:    "email",
		Expire:      time.Now().Add(time.Hour),
		UserAgent:   "Go-test/1.0",
		IP:          "127.0.0.1",
		CountryCode: "US",
		Factors:     []string{"email"},
	})
	if err != nil || sess == nil {
		t.Fatalf("CreateSession: %v", err)
	}
	if sess.Secret == "" {
		t.Fatal("secret should be returned on creation")
	}

	list, err := s.ListSessions(ctx, "sess-user")
	if err != nil || len(list) != 1 {
		t.Fatalf("ListSessions: got %d, err=%v", len(list), err)
	}

	if err := s.DeleteSession(ctx, "sess-1"); err != nil {
		t.Fatalf("DeleteSession: %v", err)
	}
	list2, _ := s.ListSessions(ctx, "sess-user")
	if len(list2) != 0 {
		t.Fatal("session should be deleted")
	}
}

func TestTOTP(t *testing.T) {
	s := testStore(t)
	ctx := context.Background()

	_, _ = s.CreateUser(ctx, CreateUserParams{ID: "totp-user", Name: "Eve"})

	secret, err := s.EnrollTOTP(ctx, "auth-1", "totp-user")
	if err != nil || secret == "" {
		t.Fatalf("EnrollTOTP: %v", err)
	}

	key, decErr := base32.StdEncoding.WithPadding(base32.NoPadding).DecodeString(secret)
	if decErr != nil {
		t.Fatalf("decode totp secret: %v", decErr)
	}
	code := totpCode(key, time.Now().Unix()/30)
	if err := s.ConfirmTOTP(ctx, "totp-user", code); err != nil {
		t.Fatalf("ConfirmTOTP: %v", err)
	}

	// Re-enroll a verified factor must return ErrAlreadyEnrolled (BR-MIGRAR-001).
	_, err = s.EnrollTOTP(ctx, "auth-2", "totp-user")
	if err != ErrAlreadyEnrolled {
		t.Fatalf("expected ErrAlreadyEnrolled, got: %v", err)
	}
}

func TestFactors(t *testing.T) {
	s := testStore(t)
	ctx := context.Background()

	_, _ = s.CreateUser(ctx, CreateUserParams{ID: "factor-user", Name: "Frank"})

	secret, _ := s.EnrollTOTP(ctx, "auth-f", "factor-user")
	key, _ := base32.StdEncoding.WithPadding(base32.NoPadding).DecodeString(secret)
	code := totpCode(key, time.Now().Unix()/30)
	_ = s.ConfirmTOTP(ctx, "factor-user", code)

	sess, err := s.CreateSession(ctx, CreateSessionParams{
		ID:       "sess-f",
		UserID:   "factor-user",
		Provider: "email",
		Expire:   time.Now().Add(time.Hour),
	})
	if err != nil || sess == nil {
		t.Fatalf("CreateSession: %v", err)
	}

	factors, err := s.RecomputeFactors(ctx, "factor-user", sess.Provider)
	if err != nil {
		t.Fatalf("RecomputeFactors: %v", err)
	}

	hasEmail, hasTOTP := false, false
	for _, f := range factors {
		switch f {
		case "email":
			hasEmail = true
		case "totp":
			hasTOTP = true
		}
	}
	if !hasEmail || !hasTOTP {
		t.Fatalf("expected email+totp factors, got: %v", factors)
	}

	if err := s.UpdateSessionFactors(ctx, sess.ID, factors); err != nil {
		t.Fatalf("UpdateSessionFactors: %v", err)
	}
}

// e02s03 — Tokens

func TestToken(t *testing.T) {
	s := testStore(t)
	ctx := context.Background()

	_, _ = s.CreateUser(ctx, CreateUserParams{ID: "tok-user", Name: "Grace"})

	tok, err := s.CreateToken(ctx, CreateTokenParams{
		ID:     "tok-1",
		UserID: "tok-user",
		Type:   TokenTypeVerification,
		TTL:    time.Hour,
	})
	if err != nil || tok == nil {
		t.Fatalf("CreateToken: %v", err)
	}
	if tok.Secret == "" {
		t.Fatal("secret should be returned on creation")
	}

	consumed, err := s.ConsumeToken(ctx, tok.Secret, TokenTypeVerification)
	if err != nil || consumed == nil {
		t.Fatalf("ConsumeToken (valid): %v", err)
	}
	if consumed.UserID != "tok-user" {
		t.Fatalf("wrong user: %s", consumed.UserID)
	}

	_, err2 := s.ConsumeToken(ctx, tok.Secret, TokenTypeVerification)
	if err2 != ErrTokenInvalid {
		t.Fatalf("expected ErrTokenInvalid on second consume, got: %v", err2)
	}
}

func TestExpiredToken(t *testing.T) {
	s := testStore(t)
	ctx := context.Background()

	_, _ = s.CreateUser(ctx, CreateUserParams{ID: "exp-user", Name: "Hank"})

	secret := "deadbeef"
	encSecret, err := s.cipher.Encrypt([]byte(secret))
	if err != nil {
		t.Fatalf("encrypt: %v", err)
	}
	_, err = s.pool.Exec(ctx, `
INSERT INTO tokens (id, user_id, type, secret, expire)
VALUES ('tok-exp','exp-user',$1,$2, now()-interval '1 second')`,
		TokenTypeRecovery, encSecret)
	if err != nil {
		t.Fatalf("insert expired token: %v", err)
	}

	_, err = s.ConsumeToken(ctx, secret, TokenTypeRecovery)
	if err != ErrTokenInvalid {
		t.Fatalf("expected ErrTokenInvalid for expired token, got: %v", err)
	}
}

func TestOAuthTokenCrypto(t *testing.T) {
	s := testStore(t)
	ctx := context.Background()

	_, _ = s.CreateUser(ctx, CreateUserParams{ID: "oauth-user", Name: "Iris"})
	expiry := time.Now().Add(time.Hour)
	sess, _ := s.CreateSession(ctx, CreateSessionParams{
		ID:     "oauth-sess",
		UserID: "oauth-user",
		Expire: time.Now().Add(24 * time.Hour),
	})

	if err := s.StoreOAuthToken(ctx, sess.ID, "access-token-abc", "refresh-token-xyz", &expiry); err != nil {
		t.Fatalf("StoreOAuthToken: %v", err)
	}

	access, refresh, gotExpiry, err := s.GetOAuthTokens(ctx, sess.ID)
	if err != nil {
		t.Fatalf("GetOAuthTokens: %v", err)
	}
	if access != "access-token-abc" {
		t.Fatalf("access token mismatch: %q", access)
	}
	if refresh != "refresh-token-xyz" {
		t.Fatalf("refresh token mismatch: %q", refresh)
	}
	if gotExpiry == nil {
		t.Fatal("expiry should not be nil")
	}

	// Raw value must be encrypted, not plaintext.
	var raw *string
	_ = s.pool.QueryRow(ctx,
		`SELECT provider_access_token FROM sessions WHERE id=$1`, sess.ID).Scan(&raw)
	if raw == nil || *raw == "access-token-abc" {
		t.Fatal("OAuth token should be stored encrypted, not as plaintext")
	}
}

// TestVerifyHashDispatch is a unit-level smoke test for the hash package —
// runs without a database.
func TestVerifyHashDispatch(t *testing.T) {
	ok, err := hash.Verify("test", "test", hash.AlgPlaintext, nil)
	if err != nil || !ok {
		t.Fatalf("hash dispatch: ok=%v err=%v", ok, err)
	}
}

func ptr(s string) *string { return &s }
