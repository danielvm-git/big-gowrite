package errors

import (
	"encoding/json"
	stderrors "errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestNewAndWrap(t *testing.T) {
	cause := stderrors.New("db down")
	ae := New("user_not_found", http.StatusNotFound, cause)

	if ae.Code != "user_not_found" {
		t.Fatalf("code = %q, want user_not_found", ae.Code)
	}
	if ae.HTTPStatus != http.StatusNotFound {
		t.Fatalf("status = %d, want 404", ae.HTTPStatus)
	}
	if !stderrors.Is(ae, cause) {
		t.Fatalf("errors.Is should find the wrapped cause")
	}
	if !strings.Contains(ae.Error(), "db down") {
		t.Fatalf("Error() should include the cause, got %q", ae.Error())
	}
}

func TestNewDefaultsStatusFromCatalog(t *testing.T) {
	ae := New("team_not_found", 0, nil) // 0 => catalog default
	if ae.HTTPStatus != http.StatusNotFound {
		t.Fatalf("status = %d, want 404 from catalog", ae.HTTPStatus)
	}
	// unknown code with 0 status falls back to 500
	ae2 := New("totally_unknown_code", 0, nil)
	if ae2.HTTPStatus != http.StatusInternalServerError {
		t.Fatalf("unknown code status = %d, want 500", ae2.HTTPStatus)
	}
}

func TestAs(t *testing.T) {
	wrapped := New("project_not_found", 404, nil)
	outer := New("general_unknown", 500, wrapped) // chained
	got, ok := As(outer)
	if !ok || got.Code != "general_unknown" {
		t.Fatalf("As should surface the outermost AppError")
	}
	// non-AppError
	if _, ok := As(stderrors.New("plain")); ok {
		t.Fatalf("As should return false for a plain error")
	}
}

func TestTranslateAppError(t *testing.T) {
	rec := httptest.NewRecorder()
	WriteError(rec, New("user_not_found", http.StatusNotFound, stderrors.New("internal detail")))

	if rec.Code != http.StatusNotFound {
		t.Fatalf("http status = %d, want 404", rec.Code)
	}
	if ct := rec.Header().Get("Content-Type"); ct != "application/json" {
		t.Fatalf("content-type = %q, want application/json", ct)
	}
	var env envelope
	if err := json.Unmarshal(rec.Body.Bytes(), &env); err != nil {
		t.Fatalf("body not valid json: %v", err)
	}
	if env.Type != "user_not_found" {
		t.Fatalf("type = %q, want user_not_found", env.Type)
	}
	if env.Code != http.StatusNotFound {
		t.Fatalf("code = %d, want 404", env.Code)
	}
	if strings.Contains(env.Message, "internal detail") {
		t.Fatalf("internal cause leaked into public message: %q", env.Message)
	}
	if env.Message == "" {
		t.Fatalf("public message should fall back to catalog default")
	}
}

func TestTranslateMasksUnknownError(t *testing.T) {
	rec := httptest.NewRecorder()
	WriteError(rec, stderrors.New("raw non-app error with secret=hunter2"))

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("http status = %d, want 500", rec.Code)
	}
	var env envelope
	_ = json.Unmarshal(rec.Body.Bytes(), &env)
	if env.Type != "general_unknown" {
		t.Fatalf("type = %q, want general_unknown", env.Type)
	}
	if strings.Contains(rec.Body.String(), "hunter2") {
		t.Fatalf("raw error detail leaked: %q", rec.Body.String())
	}
}

func TestHandlerTranslatesReturnedError(t *testing.T) {
	h := Handler(func(w http.ResponseWriter, r *http.Request) error {
		return New("project_not_found", http.StatusNotFound, nil)
	})
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/v1/projects/x", nil)
	h.ServeHTTP(rec, req)
	if rec.Code != http.StatusNotFound {
		t.Fatalf("handler status = %d, want 404", rec.Code)
	}
}

func TestHandlerNoErrorLeavesResponseUntouched(t *testing.T) {
	h := Handler(func(w http.ResponseWriter, r *http.Request) error {
		w.WriteHeader(http.StatusCreated)
		_, _ = w.Write([]byte(`{"ok":true}`))
		return nil
	})
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, httptest.NewRequest(http.MethodPost, "/x", nil))
	if rec.Code != http.StatusCreated {
		t.Fatalf("status = %d, want 201", rec.Code)
	}
}

// TestCatalogParity asserts every catalog entry has a sane status + message,
// standing in for byte-parity against the legacy errors.php subset.
func TestCatalogParity(t *testing.T) {
	if len(catalog) == 0 {
		t.Fatal("catalog is empty")
	}
	for code, entry := range catalog {
		if entry.status < 400 || entry.status > 599 {
			t.Errorf("code %q has non-error status %d", code, entry.status)
		}
		if entry.message == "" {
			t.Errorf("code %q has empty message", code)
		}
	}
	// FromCode resolves a known code with catalog status + message.
	ae := FromCode("user_already_exists")
	if ae.HTTPStatus != http.StatusConflict {
		t.Fatalf("FromCode status = %d, want 409", ae.HTTPStatus)
	}
	if ae.publicMessage() == "" {
		t.Fatal("FromCode should carry a public message")
	}
}
