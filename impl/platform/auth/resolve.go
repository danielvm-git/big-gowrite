package auth

import (
	"errors"
	"net/http"
	"strings"

	"github.com/danielvm-git/gowrite/platform/permissions"
)

// Credential kinds carried by a request.
type credKind int

const (
	credNone credKind = iota
	credSession
	credJWT
	credAPIKey
)

// ErrNoCredential means the request carried no recognizable credential; the
// caller is treated as a guest by ResolveOrGuest.
var ErrNoCredential = errors.New("auth: no credential present")

// Resolver looks up the backing store for a credential. It is implemented by
// the identity context; auth stays storage-agnostic (no DB import here).
type Resolver interface {
	// SessionPrincipal resolves an opaque session secret to a Principal.
	SessionPrincipal(secret string) (permissions.Principal, error)
	// KeyPrincipal resolves an API key secret to a Principal.
	KeyPrincipal(secret string) (permissions.Principal, error)
	// JWTSecret returns the signing secret for JWT verification.
	JWTSecret() []byte
}

// extractCredential pulls the first recognized credential from the request.
// Precedence matches the legacy: API key header, then JWT, then session.
func extractCredential(r *http.Request) (credKind, string) {
	if k := r.Header.Get("X-Appwrite-Key"); k != "" {
		return credAPIKey, k
	}
	if j := r.Header.Get("X-Appwrite-JWT"); j != "" {
		return credJWT, j
	}
	if a := r.Header.Get("Authorization"); a != "" {
		if v, ok := strings.CutPrefix(a, "Bearer "); ok {
			return credJWT, v
		}
	}
	// Session cookie (console uses a cookie session, DEV-005).
	if c, err := r.Cookie("a_session"); err == nil && c.Value != "" {
		return credSession, c.Value
	}
	if s := r.Header.Get("X-Appwrite-Session"); s != "" {
		return credSession, s
	}
	return credNone, ""
}

// Resolve resolves the request's Principal using the resolver, or returns
// ErrNoCredential if none is present.
func Resolve(r *http.Request, res Resolver, nowUnix int64) (permissions.Principal, error) {
	kind, secret := extractCredential(r)
	switch kind {
	case credAPIKey:
		return res.KeyPrincipal(secret)
	case credJWT:
		claims, err := VerifyJWT(res.JWTSecret(), secret, nowUnix)
		if err != nil {
			return permissions.Principal{}, err
		}
		aclRoles := append([]string{"user:" + claims.UserID, "users"}, claims.Roles...)
		return permissions.NewPrincipal(claims.UserID, []string{permissions.RoleUsers}, aclRoles), nil
	case credSession:
		return res.SessionPrincipal(secret)
	default:
		return permissions.Principal{}, ErrNoCredential
	}
}

// ResolveOrGuest is like Resolve but returns a guest Principal when no
// credential is present, instead of an error.
func ResolveOrGuest(r *http.Request, res Resolver, nowUnix int64) permissions.Principal {
	p, err := Resolve(r, res, nowUnix)
	if err != nil {
		return permissions.NewPrincipal("", []string{permissions.RoleGuests}, []string{"guests"})
	}
	return p
}
