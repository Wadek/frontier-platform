// Package learning implements the .agent_learning.jsonl ledger.
// Append-only: records are validated, then appended in time order. Never rewrite (F0).
// Spec: english/LEARNING.md
package learning

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"
)

// Record is one session debrief. Field names are the platform standard.
type Record struct {
	Session         string            `json:"session"`
	Project         string            `json:"project"`
	Agent           string            `json:"agent"`
	Provider        string            `json:"provider"`
	TS              string            `json:"ts"`
	TLDR            string            `json:"tldr"`
	Learned         []string          `json:"learned,omitempty"`
	Methods         []string          `json:"methods,omitempty"`
	SkillsProposed  []string          `json:"skills_proposed,omitempty"`
	MappingsUpdated []string          `json:"mappings_updated,omitempty"`
	Attribution     map[string]string `json:"attribution,omitempty"`
	Next            []string          `json:"next,omitempty"`
}

var providers = map[string]bool{"local": true, "flash": true, "pro": true, "dsh": true}

// Validate checks a record against the schema in english/LEARNING.md.
func Validate(r Record) error {
	if r.Session == "" || r.Project == "" || r.Agent == "" || r.Provider == "" || r.TS == "" || r.TLDR == "" {
		return fmt.Errorf("missing required field (session/project/agent/provider/ts/tldr)")
	}
	if !providers[r.Provider] {
		return fmt.Errorf("bad provider %q (local|flash|pro|dsh)", r.Provider)
	}
	if _, err := time.Parse(time.RFC3339, r.TS); err != nil {
		return fmt.Errorf("bad ts %q: %w", r.TS, err)
	}
	for _, m := range r.Methods {
		if err := ValidateMethod(m); err != nil {
			return err
		}
	}
	return nil
}

// Append validates and appends one record as a single JSON line.
func Append(path string, r Record) error {
	if err := Validate(r); err != nil {
		return err
	}
	line, err := json.Marshal(r)
	if err != nil {
		return err
	}
	f, err := os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o644)
	if err != nil {
		return err
	}
	defer f.Close()
	if _, err := f.Write(append(line, '\n')); err != nil {
		return err
	}
	return nil
}

// Read parses every line of a learning ledger.
func Read(path string) ([]Record, error) {
	raw, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}
	var out []Record
	for _, ln := range strings.Split(string(raw), "\n") {
		ln = strings.TrimSpace(ln)
		if ln == "" {
			continue
		}
		var r Record
		if err := json.Unmarshal([]byte(ln), &r); err != nil {
			return nil, fmt.Errorf("bad record: %w", err)
		}
		out = append(out, r)
	}
	return out, nil
}
