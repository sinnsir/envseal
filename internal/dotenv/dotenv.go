// Package dotenv provides parsing and serialization of .env files.
package dotenv

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"strings"
)

// Env represents a parsed .env file as an ordered list of entries.
type Env struct {
	Entries []Entry
}

// Entry is a single line in a .env file.
type Entry struct {
	Key     string
	Value   string
	Comment string // non-empty if the line is a comment or blank
}

// Parse reads a .env file from r and returns an Env.
func Parse(r io.Reader) (*Env, error) {
	env := &Env{}
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := scanner.Text()
		trimmed := strings.TrimSpace(line)
		if trimmed == "" || strings.HasPrefix(trimmed, "#") {
			env.Entries = append(env.Entries, Entry{Comment: line})
			continue
		}
		idx := strings.IndexByte(trimmed, '=')
		if idx < 0 {
			return nil, fmt.Errorf("dotenv: invalid line %q", line)
		}
		key := strings.TrimSpace(trimmed[:idx])
		val := strings.TrimSpace(trimmed[idx+1:])
		val = unquote(val)
		env.Entries = append(env.Entries, Entry{Key: key, Value: val})
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return env, nil
}

// Marshal serializes an Env back to .env format.
func Marshal(env *Env) ([]byte, error) {
	var buf bytes.Buffer
	for _, e := range env.Entries {
		if e.Comment != "" || (e.Key == "" && e.Value == "") {
			fmt.Fprintln(&buf, e.Comment)
			continue
		}
		val := e.Value
		if strings.ContainsAny(val, " \t\n#") {
			val = fmt.Sprintf("%q", val)
		}
		fmt.Fprintf(&buf, "%s=%s\n", e.Key, val)
	}
	return buf.Bytes(), nil
}

// Map returns the env entries as a key→value map (comments excluded).
func (e *Env) Map() map[string]string {
	m := make(map[string]string, len(e.Entries))
	for _, entry := range e.Entries {
		if entry.Key != "" {
			m[entry.Key] = entry.Value
		}
	}
	return m
}

func unquote(s string) string {
	if len(s) >= 2 {
		if (s[0] == '"' && s[len(s)-1] == '"') ||
			(s[0] == '\'' && s[len(s)-1] == '\'') {
			return s[1 : len(s)-1]
		}
	}
	return s
}
