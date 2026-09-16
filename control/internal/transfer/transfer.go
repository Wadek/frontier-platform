// Package transfer validates the handoff card envelope used to move work between
// agents. Cards are CloudEvents-shaped and land in the target's inbox/.
// Spec: english/TERMS.md — {specversion, type, source, id, data}
package transfer

import (
	"encoding/json"
	"fmt"
	"strings"
)

// SpecVersion is the only accepted CloudEvents version.
const SpecVersion = "1.0"

// TypeHandoff is the card type used to request a handoff.
const TypeHandoff = "handoff.request"

// Card is the transfer envelope.
type Card struct {
	SpecVersion string          `json:"specversion"`
	Type        string          `json:"type"`
	Source      string          `json:"source"`
	ID          string          `json:"id"`
	Data        json.RawMessage `json:"data"`
}

var requiredKeys = []string{"specversion", "type", "source", "id", "data"}

// Validate checks a decoded card. All envelope fields must be non-empty and the
// spec version must be exactly 1.0.
func Validate(c Card) error {
	if c.SpecVersion != SpecVersion {
		return fmt.Errorf("specversion %q, want %q", c.SpecVersion, SpecVersion)
	}
	for _, kv := range []struct{ name, val string }{
		{"type", c.Type}, {"source", c.Source}, {"id", c.ID},
	} {
		if kv.val == "" {
			return fmt.Errorf("empty %s", kv.name)
		}
	}
	if len(c.Data) == 0 || string(c.Data) == "null" {
		return fmt.Errorf("missing data")
	}
	return nil
}

// Parse enforces key presence, then validates. Unlike Validate it can tell a
// missing key from an empty one.
func Parse(raw []byte) (Card, error) {
	var probe map[string]json.RawMessage
	if err := json.Unmarshal(raw, &probe); err != nil {
		return Card{}, fmt.Errorf("not an object: %w", err)
	}
	for _, k := range requiredKeys {
		if _, ok := probe[k]; !ok {
			return Card{}, fmt.Errorf("missing %s", k)
		}
	}
	var c Card
	if err := json.Unmarshal(raw, &c); err != nil {
		return Card{}, err
	}
	if err := Validate(c); err != nil {
		return Card{}, err
	}
	return c, nil
}

// IsHandoff reports whether the card requests a handoff.
func (c Card) IsHandoff() bool {
	return c.Type == TypeHandoff || strings.HasPrefix(c.Type, "handoff.")
}
