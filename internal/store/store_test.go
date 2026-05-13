package store_test

import (
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/yourorg/envseal/internal/store"
)

func newStore(t *testing.T) *store.Store {
	t.Helper()
	dir := t.TempDir()
	return store.New(filepath.Join(dir, ".envseal"))
}

func TestWrite_Read_RoundTrip(t *testing.T) {
	s := newStore(t)
	data := []byte("sealed-payload-bytes")

	if err := s.Write("production", data); err != nil {
		t.Fatalf("Write: %v", err)
	}

	got, err := s.Read("production")
	if err != nil {
		t.Fatalf("Read: %v", err)
	}
	if string(got) != string(data) {
		t.Errorf("got %q, want %q", got, data)
	}
}

func TestRead_NotFound(t *testing.T) {
	s := newStore(t)
	_, err := s.Read("staging")
	if !errors.Is(err, store.ErrNotFound) {
		t.Errorf("expected ErrNotFound, got %v", err)
	}
}

func TestList(t *testing.T) {
	s := newStore(t)

	for _, env := range []string{"dev", "staging", "production"} {
		if err := s.Write(env, []byte(env)); err != nil {
			t.Fatalf("Write(%s): %v", env, err)
		}
	}

	envs, err := s.List()
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if len(envs) != 3 {
		t.Errorf("expected 3 environments, got %d: %v", len(envs), envs)
	}
}

func TestList_EmptyDir(t *testing.T) {
	s := newStore(t)
	envs, err := s.List()
	if err != nil {
		t.Fatalf("List on missing dir: %v", err)
	}
	if len(envs) != 0 {
		t.Errorf("expected empty list, got %v", envs)
	}
}

func TestRemove(t *testing.T) {
	s := newStore(t)
	if err := s.Write("dev", []byte("data")); err != nil {
		t.Fatal(err)
	}
	if err := s.Remove("dev"); err != nil {
		t.Fatalf("Remove: %v", err)
	}
	if _, err := os.Stat(s.Path("dev")); !errors.Is(err, os.ErrNotExist) {
		t.Error("file should have been removed")
	}
}

func TestRemove_NotFound(t *testing.T) {
	s := newStore(t)
	err := s.Remove("ghost")
	if !errors.Is(err, store.ErrNotFound) {
		t.Errorf("expected ErrNotFound, got %v", err)
	}
}
