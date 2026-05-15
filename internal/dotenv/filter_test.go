package dotenv

import (
	"testing"
)

var filterSrc = map[string]string{
	"DB_HOST":     "localhost",
	"DB_PORT":     "5432",
	"APP_NAME":    "envseal",
	"APP_VERSION": "1.0.0",
	"SECRET_KEY":  "s3cr3t",
}

func TestFilter_Prefix(t *testing.T) {
	out, err := Filter(filterSrc, FilterOptions{Prefix: "DB_"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out) != 2 {
		t.Fatalf("expected 2 keys, got %d", len(out))
	}
	if _, ok := out["DB_HOST"]; !ok {
		t.Error("expected DB_HOST in result")
	}
	if _, ok := out["DB_PORT"]; !ok {
		t.Error("expected DB_PORT in result")
	}
}

func TestFilter_Suffix(t *testing.T) {
	out, err := Filter(filterSrc, FilterOptions{Suffix: "_KEY"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out) != 1 {
		t.Fatalf("expected 1 key, got %d", len(out))
	}
	if _, ok := out["SECRET_KEY"]; !ok {
		t.Error("expected SECRET_KEY in result")
	}
}

func TestFilter_Pattern(t *testing.T) {
	out, err := Filter(filterSrc, FilterOptions{Pattern: `^APP_`})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out) != 2 {
		t.Fatalf("expected 2 keys, got %d", len(out))
	}
}

func TestFilter_InvalidPattern(t *testing.T) {
	_, err := Filter(filterSrc, FilterOptions{Pattern: `[invalid`})
	if err == nil {
		t.Fatal("expected error for invalid pattern")
	}
}

func TestFilter_Invert(t *testing.T) {
	out, err := Filter(filterSrc, FilterOptions{Prefix: "DB_", Invert: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out) != 3 {
		t.Fatalf("expected 3 keys, got %d", len(out))
	}
	if _, ok := out["DB_HOST"]; ok {
		t.Error("DB_HOST should have been excluded")
	}
}

func TestFilter_NoOptions(t *testing.T) {
	out, err := Filter(filterSrc, FilterOptions{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out) != len(filterSrc) {
		t.Fatalf("expected all %d keys, got %d", len(filterSrc), len(out))
	}
}

func TestFilter_DoesNotMutateSrc(t *testing.T) {
	src := map[string]string{"A": "1", "B": "2"}
	_, err := Filter(src, FilterOptions{Prefix: "A"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(src) != 2 {
		t.Error("Filter must not mutate the source map")
	}
}
