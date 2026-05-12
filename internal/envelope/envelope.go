package envelope

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"time"

	"filippo.io/age"
)

// Header holds metadata stored alongside encrypted content.
type Header struct {
	Environment string    `json:"environment"`
	Version     int       `json:"version"`
	CreatedAt   time.Time `json:"created_at"`
}

// Envelope wraps an encrypted .env payload with its header.
type Envelope struct {
	Header    Header `json:"header"`
	Ciphertext []byte `json:"ciphertext"`
}

// Seal encrypts plaintext for the given recipient and returns an Envelope.
func Seal(env string, version int, plaintext []byte, recipient age.Recipient) (*Envelope, error) {
	var buf bytes.Buffer
	w, err := age.Encrypt(&buf, recipient)
	if err != nil {
		return nil, fmt.Errorf("envelope: encrypt init: %w", err)
	}
	if _, err := w.Write(plaintext); err != nil {
		return nil, fmt.Errorf("envelope: write plaintext: %w", err)
	}
	if err := w.Close(); err != nil {
		return nil, fmt.Errorf("envelope: finalise encryption: %w", err)
	}
	return &Envelope{
		Header: Header{
			Environment: env,
			Version:     version,
			CreatedAt:   time.Now().UTC(),
		},
		Ciphertext: buf.Bytes(),
	}, nil
}

// Open decrypts the envelope using the provided identity and returns plaintext.
func Open(e *Envelope, identity age.Identity) ([]byte, error) {
	r, err := age.Decrypt(bytes.NewReader(e.Ciphertext), identity)
	if err != nil {
		return nil, fmt.Errorf("envelope: decrypt: %w", err)
	}
	plaintext, err := io.ReadAll(r)
	if err != nil {
		return nil, fmt.Errorf("envelope: read plaintext: %w", err)
	}
	return plaintext, nil
}

// Marshal serialises the envelope to JSON.
func Marshal(e *Envelope) ([]byte, error) {
	data, err := json.MarshalIndent(e, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("envelope: marshal: %w", err)
	}
	return data, nil
}

// Unmarshal deserialises an envelope from JSON.
func Unmarshal(data []byte) (*Envelope, error) {
	var e Envelope
	if err := json.Unmarshal(data, &e); err != nil {
		return nil, fmt.Errorf("envelope: unmarshal: %w", err)
	}
	return &e, nil
}
