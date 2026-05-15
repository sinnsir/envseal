package dotenv

import (
	"os"
	"sort"
	"testing"
)

func TestFromEnv_AllVars(t *testing.T) {
	t.Setenv("ENVSEAL_TEST_A", "alpha")
	t.Setenv("ENVSEAL_TEST_B", "beta")

	got := FromEnv(nil)
	if got["ENVSEAL_TEST_A"] != "alpha" {
		t.Errorf("expected ENVSEAL_TEST_A=alpha, got %q", got["ENVSEAL_TEST_A"])
	}
	if got["ENVSEAL_TEST_B"] != "beta" {
		t.Errorf("expected ENVSEAL_TEST_B=beta, got %q", got["ENVSEAL_TEST_B"])
	}
}

func TestFromEnv_FilteredKeys(t *testing.T) {
	t.Setenv("ENVSEAL_TEST_X", "x-val")
	t.Setenv("ENVSEAL_TEST_Y", "y-val")

	got := FromEnv([]string{"ENVSEAL_TEST_X", "ENVSEAL_TEST_MISSING"})
	if got["ENVSEAL_TEST_X"] != "x-val" {
		t.Errorf("expected ENVSEAL_TEST_X=x-val, got %q", got["ENVSEAL_TEST_X"])
	}
	if _, ok := got["ENVSEAL_TEST_Y"]; ok {
		t.Error("expected ENVSEAL_TEST_Y to be absent")
	}
	if _, ok := got["ENVSEAL_TEST_MISSING"]; ok {
		t.Error("expected ENVSEAL_TEST_MISSING to be absent")
	}
}

func TestToEnv_Sorted(t *testing.T) {
	m := map[string]string{
		"Z": "last",
		"A": "first",
		"M": "middle",
	}
	got := ToEnv(m)
	if !sort.StringsAreSorted(got) {
		t.Errorf("expected sorted output, got %v", got)
	}
	if got[0] != "A=first" {
		t.Errorf("expected A=first, got %q", got[0])
	}
	if got[2] != "Z=last" {
		t.Errorf("expected Z=last, got %q", got[2])
	}
}

func TestToEnv_Empty(t *testing.T) {
	got := ToEnv(map[string]string{})
	if len(got) != 0 {
		t.Errorf("expected empty slice, got %v", got)
	}
}

func TestApplyToProcess(t *testing.T) {
	m := map[string]string{
		"ENVSEAL_APPLY_TEST": "applied",
	}
	if err := ApplyToProcess(m); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got := os.Getenv("ENVSEAL_APPLY_TEST"); got != "applied" {
		t.Errorf("expected applied, got %q", got)
	}
	os.Unsetenv("ENVSEAL_APPLY_TEST")
}
