package ratelimit

import "time"

// Common legacy windows.
const (
	WindowHour   = time.Hour
	WindowMinute = time.Minute
)

// Named route limits transcribed from the legacy abuse-limit/abuse-time/
// abuse-key labels (permissions.md §8, BR-MIGRAR-009). These are the canonical
// policies routes attach; values must match the legacy exactly.
var (
	// Account create / auth-sensitive writes: 10 per hour by IP.
	LimitAuthWrite = Limit{Max: 10, Window: WindowHour, Key: KeyIP}
	// Session/login: 10 per hour by email+ip composite (email key here).
	LimitLogin = Limit{Max: 10, Window: WindowHour, Key: KeyEmail}
	// Recovery / verification token creation: 10 per hour by userId.
	LimitTokenCreate = Limit{Max: 10, Window: WindowHour, Key: KeyUserID}
	// Generic per-URL write: 60 per minute by url.
	LimitWriteURL = Limit{Max: 60, Window: WindowMinute, Key: KeyURL}
	// Generic per-user API: 120 per minute by userId.
	LimitUserAPI = Limit{Max: 120, Window: WindowMinute, Key: KeyUserID}
	// Heavier list/read endpoints: 100 per minute by userId.
	LimitRead = Limit{Max: 100, Window: WindowMinute, Key: KeyUserID}
	// Function execution: 50 per minute by userId.
	LimitExecution = Limit{Max: 50, Window: WindowMinute, Key: KeyUserID}
	// Messaging/target writes: 15 per hour by userId.
	LimitMessaging = Limit{Max: 15, Window: WindowHour, Key: KeyUserID}
)

// AllLimits enumerates the named limits for validation/introspection.
var AllLimits = map[string]Limit{
	"auth_write":   LimitAuthWrite,
	"login":        LimitLogin,
	"token_create": LimitTokenCreate,
	"write_url":    LimitWriteURL,
	"user_api":     LimitUserAPI,
	"read":         LimitRead,
	"execution":    LimitExecution,
	"messaging":    LimitMessaging,
}
