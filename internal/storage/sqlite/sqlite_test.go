package sqlite

import (
	"errors"
	"testing"
	"url-shortener/internal/storage"
)

func TestSaveURLReturnsAliasConflict(t *testing.T) {
	db, err := New(":memory:")
	if err != nil {
		t.Fatalf("create storage: %v", err)
	}
	t.Cleanup(func() {
		_ = db.db.Close()
	})

	if err := db.SaveURL("https://example.com/first", "alias"); err != nil {
		t.Fatalf("save first URL: %v", err)
	}

	err = db.SaveURL("https://example.com/second", "alias")
	if !errors.Is(err, storage.ErrAliasExists) {
		t.Fatalf("expected ErrAliasExists, got %v", err)
	}
}
