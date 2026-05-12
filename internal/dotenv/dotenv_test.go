package dotenv_test

import (
	"strings"
	"testing"

	"github.com/yourorg/envseal/internal/dotenv"
)

const sampleEnv = `# Database config
DB_HOST=localhost
DB_PORT=5432
DB_NAME=myapp

# App settings
APP_SECRET="hello world"
DEBUG=true
`

func TestParse_Basic(t *testing.T) {
	env, err := dotenv.Parse(strings.NewReader(sampleEnv))
	if err != nil {
		t.Fatalf("Parse: %v", err)
	}
	m := env.Map()
	cases := map[string]string{
		"DB_HOST":    "localhost",
		"DB_PORT":    "5432",
		"DB_NAME":    "myapp",
		"APP_SECRET": "hello world",
		"DEBUG":      "true",
	}
	for k, want := range cases {
		if got := m[k]; got != want {
			t.Errorf("key %q: got %q, want %q", k, got, want)
		}
	}
}

func TestMarshal_RoundTrip(t *testing.T) {
	env, err := dotenv.Parse(strings.NewReader(sampleEnv))
	if err != nil {
		t.Fatalf("Parse: %v", err)
	}
	b, err := dotenv.Marshal(env)
	if err != nil {
		t.Fatalf("Marshal: %v", err)
	}
	env2, err := dotenv.Parse(strings.NewReader(string(b)))
	if err != nil {
		t.Fatalf("Parse after marshal: %v", err)
	}
	m1, m2 := env.Map(), env2.Map()
	if len(m1) != len(m2) {
		t.Fatalf("map length mismatch: %d vs %d", len(m1), len(m2))
	}
	for k, v := range m1 {
		if m2[k] != v {
			t.Errorf("key %q: got %q, want %q", k, m2[k], v)
		}
	}
}

func TestParse_InvalidLine(t *testing.T) {
	_, err := dotenv.Parse(strings.NewReader("NOEQUALSIGN\n"))
	if err == nil {
		t.Fatal("expected error for invalid line, got nil")
	}
}

func TestMap_SkipsComments(t *testing.T) {
	env, _ := dotenv.Parse(strings.NewReader("# comment\nFOO=bar\n"))
	m := env.Map()
	if _, ok := m[""]; ok {
		t.Error("map should not contain empty-key entries from comments")
	}
	if m["FOO"] != "bar" {
		t.Errorf("expected FOO=bar, got %q", m["FOO"])
	}
}
