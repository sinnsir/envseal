package dotenv

import (
	"os"
	"path/filepath"
	"testing"
)

func TestReadFile_Basic(t *testing.T) {
	tmp := t.TempDir()
	path := filepath.Join(tmp, ".env")
	_ = os.WriteFile(path, []byte("FOO=bar\nBAZ=qux\n"), 0600)

	kvs, err := ReadFile(path)
	if err != nil {
		t.Fatalf("ReadFile: %v", err)
	}
	if kvs["FOO"] != "bar" {
		t.Errorf("expected FOO=bar, got %q", kvs["FOO"])
	}
	if kvs["BAZ"] != "qux" {
		t.Errorf("expected BAZ=qux, got %q", kvs["BAZ"])
	}
}

func TestReadFile_NotFound(t *testing.T) {
	_, err := ReadFile("/nonexistent/.env")
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}

func TestWriteFile_RoundTrip(t *testing.T) {
	tmp := t.TempDir()
	path := filepath.Join(tmp, ".env")

	input := map[string]string{"KEY": "value", "OTHER": "123"}
	if err := WriteFile(path, input); err != nil {
		t.Fatalf("WriteFile: %v", err)
	}

	kvs, err := ReadFile(path)
	if err != nil {
		t.Fatalf("ReadFile: %v", err)
	}
	if kvs["KEY"] != "value" {
		t.Errorf("expected KEY=value, got %q", kvs["KEY"])
	}
	if kvs["OTHER"] != "123" {
		t.Errorf("expected OTHER=123, got %q", kvs["OTHER"])
	}
}

func TestMerge_OverrideTakesPrecedence(t *testing.T) {
	base := map[string]string{"A": "1", "B": "2"}
	override := map[string]string{"B": "99", "C": "3"}

	out := Merge(base, override)
	if out["A"] != "1" {
		t.Errorf("expected A=1, got %q", out["A"])
	}
	if out["B"] != "99" {
		t.Errorf("expected B=99 (override), got %q", out["B"])
	}
	if out["C"] != "3" {
		t.Errorf("expected C=3, got %q", out["C"])
	}
}

func TestMerge_DoesNotMutateInputs(t *testing.T) {
	base := map[string]string{"A": "1"}
	override := map[string]string{"A": "2"}
	Merge(base, override)
	if base["A"] != "1" {
		t.Error("base map was mutated")
	}
}

func TestFilterKeys(t *testing.T) {
	kvs := map[string]string{"FOO": "1", "BAR": "2", "BAZ": "3"}
	out := FilterKeys(kvs, []string{"FOO", "BAZ"})
	if len(out) != 2 {
		t.Errorf("expected 2 keys, got %d", len(out))
	}
	if _, ok := out["BAR"]; ok {
		t.Error("BAR should have been filtered out")
	}
}
