// Package auth resolves a caller's Principal from an HTTP request (session
// token, JWT, or API key) and provides minimal HS256 JWT sign/verify. It
// replaces the injected $user/$project + Action scope labels of the legacy with
// explicit middleware (BR-DESCARTAR-002, BR-MIGRAR-007).
package auth

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"
)

// ErrInvalidToken is returned for a malformed, mis-signed, or expired JWT.
var ErrInvalidToken = errors.New("auth: invalid token")

// Claims is the minimal JWT payload gowrite issues.
type Claims struct {
	UserID    string   `json:"userId"`
	SessionID string   `json:"sessionId,omitempty"`
	Roles     []string `json:"roles,omitempty"`
	Expires   int64    `json:"exp"`
	IssuedAt  int64    `json:"iat"`
}

// header is a fixed HS256 JWT header.
type header struct {
	Alg string `json:"alg"`
	Typ string `json:"typ"`
}

func b64(b []byte) string { return base64.RawURLEncoding.EncodeToString(b) }

// SignJWT issues an HS256 JWT for claims valid for ttl, using secret.
// nowUnix is passed explicitly for deterministic testing.
func SignJWT(secret []byte, claims Claims, nowUnix int64, ttl time.Duration) (string, error) {
	claims.IssuedAt = nowUnix
	claims.Expires = nowUnix + int64(ttl.Seconds())

	h, err := json.Marshal(header{Alg: "HS256", Typ: "JWT"})
	if err != nil {
		return "", err
	}
	p, err := json.Marshal(claims)
	if err != nil {
		return "", err
	}
	signingInput := b64(h) + "." + b64(p)
	sig := sign(secret, signingInput)
	return signingInput + "." + b64(sig), nil
}

// VerifyJWT validates the signature and expiry and returns the claims.
func VerifyJWT(secret []byte, token string, nowUnix int64) (Claims, error) {
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return Claims{}, fmt.Errorf("%w: expected 3 parts", ErrInvalidToken)
	}
	signingInput := parts[0] + "." + parts[1]
	sig, err := base64.RawURLEncoding.DecodeString(parts[2])
	if err != nil {
		return Claims{}, fmt.Errorf("%w: bad signature encoding", ErrInvalidToken)
	}
	// Constant-time signature comparison.
	if !hmac.Equal(sig, sign(secret, signingInput)) {
		return Claims{}, fmt.Errorf("%w: signature mismatch", ErrInvalidToken)
	}
	payload, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return Claims{}, fmt.Errorf("%w: bad payload encoding", ErrInvalidToken)
	}
	var claims Claims
	if err := json.Unmarshal(payload, &claims); err != nil {
		return Claims{}, fmt.Errorf("%w: bad payload json", ErrInvalidToken)
	}
	if claims.Expires != 0 && nowUnix >= claims.Expires {
		return Claims{}, fmt.Errorf("%w: expired", ErrInvalidToken)
	}
	return claims, nil
}

func sign(secret []byte, input string) []byte {
	m := hmac.New(sha256.New, secret)
	m.Write([]byte(input))
	return m.Sum(nil)
}
