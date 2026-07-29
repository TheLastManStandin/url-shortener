package save

import (
	"bytes"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"
	"url-shortener/internal/storage"
)

type urlSaverFunc func(urlToSave, alias string) error

func (f urlSaverFunc) SaveURL(urlToSave, alias string) error {
	return f(urlToSave, alias)
}

func TestUserAliasConflict(t *testing.T) {
	var calls int
	handler := New(testLogger(), urlSaverFunc(func(_, alias string) error {
		calls++
		if alias != "taken" {
			t.Fatalf("unexpected alias: %q", alias)
		}
		return storage.ErrAliasExists
	}))

	recorder := performRequest(handler, `{"url":"https://example.com","alias":"taken"}`)

	if recorder.Code != http.StatusConflict {
		t.Fatalf("expected status %d, got %d", http.StatusConflict, recorder.Code)
	}
	if calls != 1 {
		t.Fatalf("expected one save attempt, got %d", calls)
	}
}

func TestGeneratedAliasRetriesAfterConflict(t *testing.T) {
	var calls int
	handler := New(testLogger(), urlSaverFunc(func(_, alias string) error {
		calls++
		if alias == "" {
			t.Fatal("generated alias is empty")
		}
		if calls == 1 {
			return storage.ErrAliasExists
		}
		return nil
	}))

	recorder := performRequest(handler, `{"url":"https://example.com"}`)

	if recorder.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, recorder.Code)
	}
	if calls != 2 {
		t.Fatalf("expected two save attempts, got %d", calls)
	}
}

func TestGeneratedAliasAttemptsAreLimited(t *testing.T) {
	var calls int
	handler := New(testLogger(), urlSaverFunc(func(_, _ string) error {
		calls++
		return storage.ErrAliasExists
	}))

	recorder := performRequest(handler, `{"url":"https://example.com"}`)

	if recorder.Code != http.StatusInternalServerError {
		t.Fatalf("expected status %d, got %d", http.StatusInternalServerError, recorder.Code)
	}
	if calls != maxAliasSaveAttempts {
		t.Fatalf("expected %d save attempts, got %d", maxAliasSaveAttempts, calls)
	}
}

func TestSaveErrorIsNotRetried(t *testing.T) {
	var calls int
	handler := New(testLogger(), urlSaverFunc(func(_, _ string) error {
		calls++
		return errors.New("database unavailable")
	}))

	recorder := performRequest(handler, `{"url":"https://example.com"}`)

	if recorder.Code != http.StatusInternalServerError {
		t.Fatalf("expected status %d, got %d", http.StatusInternalServerError, recorder.Code)
	}
	if calls != 1 {
		t.Fatalf("expected one save attempt, got %d", calls)
	}
}

func performRequest(handler http.HandlerFunc, body string) *httptest.ResponseRecorder {
	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodPost, "/url", bytes.NewBufferString(body))
	handler.ServeHTTP(recorder, request)
	return recorder
}

func testLogger() *slog.Logger {
	return slog.New(slog.NewTextHandler(io.Discard, nil))
}
