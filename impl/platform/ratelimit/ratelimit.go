// Package ratelimit implements a per-instance token-bucket limiter that
// preserves the legacy per-route abuse limits and keys (BR-MIGRAR-009) without
// Redis. A single-VPS target makes per-instance state acceptable (red-team).
package ratelimit

import (
	"sync"
	"time"
)

// Limit describes a route's abuse policy: Max requests per Window, keyed by Key.
// Mirrors the legacy abuse-limit / abuse-time / abuse-key route labels.
type Limit struct {
	Max    int           // e.g. 10, 15, 50, 60, 100, 120
	Window time.Duration // e.g. 3600s or 60s
	Key    KeyKind       // which request attribute forms the bucket key
}

// KeyKind is the abuse-key dimension.
type KeyKind string

const (
	KeyIP     KeyKind = "ip"
	KeyUserID KeyKind = "userId"
	KeyURL    KeyKind = "url"
	KeyEmail  KeyKind = "email"
)

// bucket is a single token bucket.
type bucket struct {
	tokens   float64
	lastFill time.Time
}

// Limiter is a concurrency-safe token-bucket limiter. Time is injectable via
// `now` for deterministic tests.
type Limiter struct {
	mu      sync.Mutex
	buckets map[string]*bucket
	now     func() time.Time
}

// New builds a Limiter using the real clock.
func New() *Limiter {
	return &Limiter{buckets: map[string]*bucket{}, now: time.Now}
}

// newWithClock builds a Limiter with an injectable clock (tests).
func newWithClock(now func() time.Time) *Limiter {
	return &Limiter{buckets: map[string]*bucket{}, now: now}
}

// Allow reports whether a request under identity (the resolved abuse key) is
// permitted for the given limit. When denied, retryAfter is the wait until one
// token refills. A Max <= 0 limit means "unlimited" (no abuse label).
func (l *Limiter) Allow(identity string, lim Limit) (allowed bool, retryAfter time.Duration) {
	if lim.Max <= 0 || lim.Window <= 0 {
		return true, 0
	}
	l.mu.Lock()
	defer l.mu.Unlock()

	now := l.now()
	refillPerSec := float64(lim.Max) / lim.Window.Seconds()

	b, ok := l.buckets[identity]
	if !ok {
		b = &bucket{tokens: float64(lim.Max), lastFill: now}
		l.buckets[identity] = b
	} else {
		elapsed := now.Sub(b.lastFill).Seconds()
		if elapsed > 0 {
			b.tokens += elapsed * refillPerSec
			if b.tokens > float64(lim.Max) {
				b.tokens = float64(lim.Max)
			}
			b.lastFill = now
		}
	}

	if b.tokens >= 1 {
		b.tokens--
		return true, 0
	}
	// Tokens needed to reach 1, converted to wait time.
	deficit := 1 - b.tokens
	wait := time.Duration(deficit/refillPerSec*float64(time.Second)) * 1
	return false, wait
}

// Reset clears all buckets (test/maintenance helper).
func (l *Limiter) Reset() {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.buckets = map[string]*bucket{}
}
