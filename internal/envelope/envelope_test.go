package envelope_test

import (
	"testing"

	"filippo.io/age"

	"github.com/yourorg/envseal/internal/envelope"
)

func newIdentity(t *testing.T) (*age.X25519Identity, *age.X25519Recipient) {
	t.Helper()
	id, err := age.GenerateX25519Identity()
	if err != nil {
		t.Fatalf("generate identity: %v", err)
	}
	return id, id.Recipient()
}

func TestSealOpen_RoundTrip(t *testing.T) {
	id, rec := newIdentity(t)
	plaintext := []byte("DB_PASSWORD=supersecret\nAPI_KEY=abc123\n")

	e, err := envelope.Seal("production", 1, plaintext, rec)
	if err != nil {
		t.Fatalf("Seal: %v", err)
	}

	if e.Header.Environment != "production" {
		t.Errorf("expected environment 'production', got %q", e.Header.Environment)
	}
	if e.Header.Version != 1 {
		t.Errorf("expected version 1, got %d", e.Header.Version)
	}

	got, err := envelope.Open(e, id)
	if err != nil {
		t.Fatalf("Open: %v", err)
	}
	if string(got) != string(plaintext) {
		t.Errorf("plaintext mismatch: got %q, want %q", got, plaintext)
	}
}

func TestMarshalUnmarshal_RoundTrip(t *testing.T) {
	_, rec := newIdentity(t)
	plaintext := []byte("SECRET=hello\n")

	e, err := envelope.Seal("staging", 2, plaintext, rec)
	if err != nil {
		t.Fatalf("Seal: %v", err)
	}

	data, err := envelope.Marshal(e)
	if err != nil {
		t.Fatalf("Marshal: %v", err)
	}

	e2, err := envelope.Unmarshal(data)
	if err != nil {
		t.Fatalf("Unmarshal: %v", err)
	}

	if e2.Header.Environment != e.Header.Environment {
		t.Errorf("environment mismatch after marshal round-trip")
	}
	if e2.Header.Version != e.Header.Version {
		t.Errorf("version mismatch after marshal round-trip")
	}
	if string(e2.Ciphertext) != string(e.Ciphertext) {
		t.Errorf("ciphertext mismatch after marshal round-trip")
	}
}

func TestOpen_WrongIdentity(t *testing.T) {
	_, rec := newIdentity(t)
	wrongID, _ := newIdentity(t)

	e, err := envelope.Seal("dev", 1, []byte("X=1\n"), rec)
	if err != nil {
		t.Fatalf("Seal: %v", err)
	}

	_, err = envelope.Open(e, wrongID)
	if err == nil {
		t.Fatal("expected error opening with wrong identity, got nil")
	}
}
