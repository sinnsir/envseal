package dotenv

import (
	"testing"
)

func TestPatch_Set(t *testing.T) {
	src := map[string]string{"FOO": "bar", "BAZ": "qux"}
	ops := []PatchOp{
		{Key: "FOO", Value: "newbar", Strategy: PatchSet},
		{Key: "EXTRA", Value: "added", Strategy: PatchSet},
	}
	out, err := Patch(src, ops)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["FOO"] != "newbar" {
		t.Errorf("expected FOO=newbar, got %q", out["FOO"])
	}
	if out["EXTRA"] != "added" {
		t.Errorf("expected EXTRA=added, got %q", out["EXTRA"])
	}
	if out["BAZ"] != "qux" {
		t.Errorf("expected BAZ=qux unchanged, got %q", out["BAZ"])
	}
}

func TestPatch_Delete(t *testing.T) {
	src := map[string]string{"FOO": "bar", "BAZ": "qux"}
	ops := []PatchOp{
		{Key: "FOO", Strategy: PatchDelete},
	}
	out, err := Patch(src, ops)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := out["FOO"]; ok {
		t.Error("expected FOO to be deleted")
	}
	if out["BAZ"] != "qux" {
		t.Errorf("expected BAZ=qux unchanged, got %q", out["BAZ"])
	}
}

func TestPatch_DoesNotMutateSrc(t *testing.T) {
	src := map[string]string{"FOO": "bar"}
	_, err := Patch(src, []PatchOp{{Key: "FOO", Value: "changed", Strategy: PatchSet}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if src["FOO"] != "bar" {
		t.Error("Patch mutated the source map")
	}
}

func TestPatch_InvalidKey(t *testing.T) {
	src := map[string]string{}
	ops := []PatchOp{{Key: "1INVALID", Value: "v", Strategy: PatchSet}}
	_, err := Patch(src, ops)
	if err == nil {
		t.Error("expected error for invalid key, got nil")
	}
}

func TestParsePatchOps_Basic(t *testing.T) {
	raw := "FOO=hello\nBAR=world\n"
	ops, err := ParsePatchOps(raw)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(ops) != 2 {
		t.Fatalf("expected 2 ops, got %d", len(ops))
	}
	for _, op := range ops {
		if op.Strategy != PatchSet {
			t.Errorf("expected PatchSet strategy, got %v", op.Strategy)
		}
	}
}

func TestParsePatchOps_InvalidLine(t *testing.T) {
	_, err := ParsePatchOps("!!!invalid")
	if err == nil {
		t.Error("expected error for invalid dotenv input")
	}
}
