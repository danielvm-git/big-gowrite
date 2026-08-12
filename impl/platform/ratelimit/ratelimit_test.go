package ratelimit

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestAllowUnderLimit(t *testing.T) {
	l := New()
	lim := Limit{Max: 5, Window: time.Minute, Key: KeyIP}
	for i := 0; i < 5; i++ {
		ok, _ := l.Allow("ip:1.2.3.4", lim)
		if !ok {
			t.Fatalf("request %d should be allowed", i+1)
		}
	}
}

func TestAllowOverLimit(t *testing.T) {
	l := New()
	lim := Limit{Max: 3, Window: time.Hour, Key: KeyIP}
	for i := 0; i < 3; i++ {
		if ok, _ := l.Allow("k", lim); !ok {
			t.Fatalf("request %d should be allowed", i+1)
		}
	}
	ok, retry := l.Allow("k", lim)
	if ok {
		t.Fatal("4th request should be denied")
	}
	if retry <= 0 {
		t.Fatal("denied request should report a positive Retry-After")
	}
}

func TestSeparateKeysAreIndependent(t *testing.T) {
	l := New()
	lim := Limit{Max: 1, Window: time.Hour, Key: KeyIP}
	if ok, _ := l.Allow("a", lim); !ok {
		t.Fatal("first key should pass")
	}
	if ok, _ := l.Allow("b", lim); !ok {
		t.Fatal("different key should have its own bucket")
	}
	if ok, _ := l.Allow("a", lim); ok {
		t.Fatal("first key should now be exhausted")
	}
}

func TestRefillOverTime(t *testing.T) {
	now := time.Unix(0, 0)
	l := newWithClock(func() time.Time { return now })
	lim := Limit{Max: 2, Window: 2 * time.Second, Key: KeyIP} // 1 token/sec
	l.Allow("k", lim)
	l.Allow("k", lim)
	if ok, _ := l.Allow("k", lim); ok {
		t.Fatal("should be exhausted")
	}
	now = now.Add(1100 * time.Millisecond) // ~1.1 tokens refilled
	if ok, _ := l.Allow("k", lim); !ok {
		t.Fatal("should allow after refill")
	}
}

func TestUnlimitedWhenNoLabel(t *testing.T) {
	l := New()
	for i := 0; i < 1000; i++ {
		if ok, _ := l.Allow("k", Limit{Max: 0}); !ok {
			t.Fatal("Max<=0 means unlimited")
		}
	}
}

// TestLabels asserts the transcribed legacy limits use the documented values
// and windows (permissions.md §8).
func TestLabels(t *testing.T) {
	allowedMax := map[int]bool{10: true, 15: true, 50: true, 60: true, 100: true, 120: true}
	allowedWin := map[time.Duration]bool{WindowHour: true, WindowMinute: true}
	if len(AllLimits) == 0 {
		t.Fatal("no limits transcribed")
	}
	for name, lim := range AllLimits {
		if !allowedMax[lim.Max] {
			t.Errorf("%s: Max %d not in legacy set {10,15,50,60,100,120}", name, lim.Max)
		}
		if !allowedWin[lim.Window] {
			t.Errorf("%s: Window %v not in legacy set {3600s,60s}", name, lim.Window)
		}
		if lim.Key == "" {
			t.Errorf("%s: missing abuse key", name)
		}
	}
}

func TestMiddleware429WithRetryAfter(t *testing.T) {
	l := New()
	lim := Limit{Max: 1, Window: time.Hour, Key: KeyIP}
	ok := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(http.StatusOK) })
	h := l.Middleware(lim, DefaultKeyFunc, ok)

	// first request passes
	rec1 := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/v1/account", nil)
	req.RemoteAddr = "9.9.9.9:1234"
	h.ServeHTTP(rec1, req)
	if rec1.Code != http.StatusOK {
		t.Fatalf("first request = %d, want 200", rec1.Code)
	}

	// second is throttled
	rec2 := httptest.NewRecorder()
	h.ServeHTTP(rec2, req)
	if rec2.Code != http.StatusTooManyRequests {
		t.Fatalf("second request = %d, want 429", rec2.Code)
	}
	if rec2.Header().Get("Retry-After") == "" {
		t.Fatal("429 must include Retry-After")
	}
	if got := rec2.Body.String(); got == "" ||
		!contains(got, "general_rate_limit_exceeded") {
		t.Fatalf("429 body should carry the legacy error type, got %q", got)
	}
}

func contains(s, sub string) bool {
	return len(s) >= len(sub) && (func() bool {
		for i := 0; i+len(sub) <= len(s); i++ {
			if s[i:i+len(sub)] == sub {
				return true
			}
		}
		return false
	})()
}
