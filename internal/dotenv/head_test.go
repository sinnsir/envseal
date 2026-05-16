package dotenv

import (
	"strings"
	"testing"
)

func TestHead_BasicN(t *testing.T) {
	src := map[string]string{
		"ALPHA": "1",
		"BETA":  "2",
		"GAMMA": "3",
		"DELTA": "4",
	}
	res, err := Head(src, 2)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Shown != 2 {
		t.Errorf("expected Shown=2, got %d", res.Shown)
	}
	if res.Total != 4 {
		t.Errorf("expected Total=4, got %d", res.Total)
	}
	// Keys should be alphabetically first two: ALPHA, BETA
	if res.Keys[0] != "ALPHA" || res.Keys[1] != "BETA" {
		t.Errorf("unexpected keys: %v", res.Keys)
	}
}

func TestHead_NGreaterThanLen(t *testing.T) {
	src := map[string]string{"A": "1", "B": "2"}
	res, err := Head(src, 100)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Shown != 2 {
		t.Errorf("expected Shown=2, got %d", res.Shown)
	}
	if res.Total != 2 {
		t.Errorf("expected Total=2, got %d", res.Total)
	}
}

func TestHead_NZeroReturnsAll(t *testing.T) {
	src := map[string]string{"X": "10", "Y": "20", "Z": "30"}
	res, err := Head(src, 0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Shown != 3 {
		t.Errorf("expected Shown=3, got %d", res.Shown)
	}
}

func TestHead_NilSource(t *testing.T) {
	_, err := Head(nil, 3)
	if err == nil {
		t.Fatal("expected error for nil source")
	}
}

func TestHead_EmptyMap(t *testing.T) {
	res, err := Head(map[string]string{}, 5)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Shown != 0 || res.Total != 0 {
		t.Errorf("expected empty result, got shown=%d total=%d", res.Shown, res.Total)
	}
}

func TestFormatHead_ShowsOmittedCount(t *testing.T) {
	src := map[string]string{
		"ALPHA": "1",
		"BETA":  "2",
		"GAMMA": "3",
	}
	res, _ := Head(src, 2)
	out := FormatHead(res)
	if !strings.Contains(out, "1 more key") {
		t.Errorf("expected omitted count in output, got:\n%s", out)
	}
}

func TestFormatHead_NoOmittedWhenAll(t *testing.T) {
	src := map[string]string{"A": "1", "B": "2"}
	res, _ := Head(src, 10)
	out := FormatHead(res)
	if strings.Contains(out, "more key") {
		t.Errorf("did not expect omitted count when all keys shown, got:\n%s", out)
	}
}

func TestFormatHead_NilResult(t *testing.T) {
	out := FormatHead(nil)
	if out != "" {
		t.Errorf("expected empty string for nil result, got %q", out)
	}
}
