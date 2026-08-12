// Package errors defines gowrite's single error type, AppError, and the HTTP
// translation layer that renders it as the legacy Appwrite JSON error envelope.
//
// Paradigm note (specs/migration/target_business_rules.md BR-DESCARTAR-001):
// this replaces PHP try/catch + AppwriteException with explicit Go error values.
package errors

import (
	"encoding/json"
	stderrors "errors"
	"fmt"
	"net/http"
)

// AppError carries a public error type/code and the HTTP status the API surface
// should return, while wrapping the underlying cause for logs.
type AppError struct {
	// Code is the public, stable error type string (e.g. "user_not_found").
	Code string
	// HTTPStatus is the HTTP status code the client receives.
	HTTPStatus int
	// Message is the public, human-readable message. If empty, the catalog
	// default (or the wrapped error) is used.
	Message string
	// Err is the wrapped underlying cause (never exposed to clients).
	Err error
}

// Error implements the error interface.
func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %v", e.Code, e.Err)
	}
	if e.Message != "" {
		return fmt.Sprintf("%s: %s", e.Code, e.Message)
	}
	return e.Code
}

// Unwrap exposes the wrapped cause for errors.Is / errors.As.
func (e *AppError) Unwrap() error { return e.Err }

// New builds an AppError from a code, HTTP status, and optional wrapped cause.
// If status is 0, the catalog default for the code is used (falling back to 500).
func New(code string, status int, err error) *AppError {
	if status == 0 {
		if entry, ok := catalog[code]; ok {
			status = entry.status
		} else {
			status = http.StatusInternalServerError
		}
	}
	return &AppError{Code: code, HTTPStatus: status, Err: err}
}

// FromCode builds an AppError purely from a known catalog code, using the
// catalog's status and default message.
func FromCode(code string) *AppError {
	entry, ok := catalog[code]
	status := http.StatusInternalServerError
	msg := "server error"
	if ok {
		status = entry.status
		msg = entry.message
	}
	return &AppError{Code: code, HTTPStatus: status, Message: msg}
}

// WithMessage returns a copy of the AppError with an overridden public message.
func (e *AppError) WithMessage(msg string) *AppError {
	clone := *e
	clone.Message = msg
	return &clone
}

// publicMessage resolves the message shown to clients: explicit message, then
// catalog default, then a generic fallback.
func (e *AppError) publicMessage() string {
	if e.Message != "" {
		return e.Message
	}
	if entry, ok := catalog[e.Code]; ok && entry.message != "" {
		return entry.message
	}
	return "server error"
}

// envelope is the legacy Appwrite JSON error shape.
type envelope struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
	Type    string `json:"type"`
}

// As extracts an *AppError from err, if present.
func As(err error) (*AppError, bool) {
	var ae *AppError
	if stderrors.As(err, &ae) {
		return ae, true
	}
	return nil, false
}

// WriteError renders err as the legacy JSON error envelope. A non-AppError is
// masked as a generic 500 so internal detail never leaks to clients.
func WriteError(w http.ResponseWriter, err error) {
	ae, ok := As(err)
	if !ok {
		ae = &AppError{Code: "general_unknown", HTTPStatus: http.StatusInternalServerError, Err: err}
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(ae.HTTPStatus)
	_ = json.NewEncoder(w).Encode(envelope{
		Message: ae.publicMessage(),
		Code:    ae.HTTPStatus,
		Type:    ae.Code,
	})
}

// HandlerFunc is an HTTP handler that returns an error, letting the Handler
// middleware centralize error translation (replaces per-handler try/catch).
type HandlerFunc func(http.ResponseWriter, *http.Request) error

// Handler adapts a HandlerFunc into an http.Handler, translating any returned
// error into the legacy JSON envelope.
func Handler(fn HandlerFunc) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := fn(w, r); err != nil {
			WriteError(w, err)
		}
	})
}
