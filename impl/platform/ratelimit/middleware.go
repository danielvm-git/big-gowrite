package ratelimit

import (
	"math"
	"net"
	"net/http"
	"strconv"
)

// KeyFunc derives the abuse key value from a request for a given KeyKind.
// The identity context supplies userId/email; IP and URL are request-derived.
type KeyFunc func(r *http.Request, kind KeyKind) string

// DefaultKeyFunc keys by remote IP and request URL; userId/email default to IP
// when no identity extractor is wired.
func DefaultKeyFunc(r *http.Request, kind KeyKind) string {
	switch kind {
	case KeyURL:
		return r.URL.Path
	case KeyIP:
		return clientIP(r)
	default:
		return clientIP(r)
	}
}

func clientIP(r *http.Request) string {
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}
	return host
}

// Middleware wraps next with the given limit. On denial it writes 429 with a
// Retry-After header (whole seconds, rounded up) and does not call next.
func (l *Limiter) Middleware(lim Limit, key KeyFunc, next http.Handler) http.Handler {
	if key == nil {
		key = DefaultKeyFunc
	}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		identity := string(lim.Key) + ":" + key(r, lim.Key)
		allowed, retryAfter := l.Allow(identity, lim)
		if !allowed {
			secs := int(math.Ceil(retryAfter.Seconds()))
			if secs < 1 {
				secs = 1
			}
			w.Header().Set("Retry-After", strconv.Itoa(secs))
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusTooManyRequests)
			_, _ = w.Write([]byte(`{"message":"Rate limit exceeded. Try again later.","code":429,"type":"general_rate_limit_exceeded"}`))
			return
		}
		next.ServeHTTP(w, r)
	})
}
