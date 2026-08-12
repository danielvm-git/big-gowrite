package auth

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/danielvm-git/gowrite/platform/permissions"
)

var secret = []byte("test-signing-secret")

func TestJWTSignVerifyRoundTrip(t *testing.T) {
	now := int64(1_000_000)
	tok, err := SignJWT(secret, Claims{UserID: "u1", SessionID: "s1", Roles: []string{"team:t1/owner"}}, now, time.Hour)
	if err != nil {
		t.Fatal(err)
	}
	claims, err := VerifyJWT(secret, tok, now+10)
	if err != nil {
		t.Fatalf("verify: %v", err)
	}
	if claims.UserID != "u1" || claims.SessionID != "s1" {
		t.Fatalf("claims = %+v", claims)
	}
}

func TestJWTExpired(t *testing.T) {
	now := int64(1_000_000)
	tok, _ := SignJWT(secret, Claims{UserID: "u1"}, now, time.Minute)
	if _, err := VerifyJWT(secret, tok, now+61); !errors.Is(err, ErrInvalidToken) {
		t.Fatalf("expected ErrInvalidToken for expired token, got %v", err)
	}
}

func TestJWTTampered(t *testing.T) {
	now := int64(1_000_000)
	tok, _ := SignJWT(secret, Claims{UserID: "u1"}, now, time.Hour)
	// wrong secret
	if _, err := VerifyJWT([]byte("wrong"), tok, now); !errors.Is(err, ErrInvalidToken) {
		t.Fatalf("expected ErrInvalidToken for wrong secret, got %v", err)
	}
	// malformed
	if _, err := VerifyJWT(secret, "a.b", now); !errors.Is(err, ErrInvalidToken) {
		t.Fatal("expected error for malformed token")
	}
}

// fakeResolver implements Resolver for tests.
type fakeResolver struct{}

func (fakeResolver) JWTSecret() []byte { return secret }
func (fakeResolver) SessionPrincipal(s string) (permissions.Principal, error) {
	if s == "good-session" {
		return permissions.NewPrincipal("sess-user", []string{permissions.RoleUsers}, []string{"user:sess-user", "users"}), nil
	}
	return permissions.Principal{}, errors.New("bad session")
}
func (fakeResolver) KeyPrincipal(s string) (permissions.Principal, error) {
	if s == "good-key" {
		return permissions.NewPrincipal("", []string{permissions.RoleKeys}, []string{"keys"}), nil
	}
	return permissions.Principal{}, errors.New("bad key")
}

func TestResolveAPIKey(t *testing.T) {
	r := httptest.NewRequest(http.MethodGet, "/v1/health", nil)
	r.Header.Set("X-Appwrite-Key", "good-key")
	p, err := Resolve(r, fakeResolver{}, 0)
	if err != nil {
		t.Fatal(err)
	}
	if !p.Can("health.read") {
		t.Fatal("key principal should have health.read")
	}
}

func TestResolveJWT(t *testing.T) {
	now := int64(2_000_000)
	tok, _ := SignJWT(secret, Claims{UserID: "jwtuser", Roles: []string{"team:t9/owner"}}, now, time.Hour)
	r := httptest.NewRequest(http.MethodGet, "/v1/account", nil)
	r.Header.Set("Authorization", "Bearer "+tok)
	p, err := Resolve(r, fakeResolver{}, now)
	if err != nil {
		t.Fatal(err)
	}
	if p.UserID != "jwtuser" {
		t.Fatalf("userID = %q", p.UserID)
	}
	// ACL role from JWT claim should be present
	if !p.HasRole("team:t9/owner") {
		t.Fatal("jwt acl role not propagated")
	}
}

func TestResolveSessionCookie(t *testing.T) {
	r := httptest.NewRequest(http.MethodGet, "/v1/account", nil)
	r.AddCookie(&http.Cookie{Name: "a_session", Value: "good-session"})
	p, err := Resolve(r, fakeResolver{}, 0)
	if err != nil || p.UserID != "sess-user" {
		t.Fatalf("session resolve = %+v, err=%v", p, err)
	}
}

func TestResolveOrGuest(t *testing.T) {
	r := httptest.NewRequest(http.MethodGet, "/v1/account", nil)
	p := ResolveOrGuest(r, fakeResolver{}, 0)
	if p.UserID != "" {
		t.Fatal("expected empty guest user")
	}
	if p.Can("users.write") {
		t.Fatal("guest must not have users.write")
	}
	if _, err := Resolve(r, fakeResolver{}, 0); !errors.Is(err, ErrNoCredential) {
		t.Fatalf("expected ErrNoCredential, got %v", err)
	}
}
