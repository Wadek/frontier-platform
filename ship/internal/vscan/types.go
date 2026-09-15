// Package vscan is programmatic V: built-in OWASP + optional tool adapters.
// Enhance briefs consume this pack so host models burn fewer tokens.
package vscan

import "strings"

// Finding is one programmatic hit under V (any source).
type Finding struct {
	Source   string `json:"source"`
	RuleID   string `json:"rule_id"`
	OWASP    string `json:"owasp,omitempty"`
	Severity string `json:"severity"`
	Path     string `json:"path"`
	Line     int    `json:"line,omitempty"`
	Snippet  string `json:"snippet,omitempty"`
}

// Result is one scanner's output.
type Result struct {
	Source   string         `json:"source"`
	Findings []Finding      `json:"findings"`
	Meta     map[string]any `json:"meta,omitempty"`
	Skipped  bool           `json:"skipped,omitempty"`
	SkipWhy  string         `json:"skip_why,omitempty"`
}

// Options controls programmatic scans.
type Options struct {
	// Adapters lists optional scanner names to run (e.g. "checkov"). Empty = none.
	Adapters []string
	// AutoAdapters runs every Available() non-builtin adapter.
	AutoAdapters bool
}

// Scanner is a programmatic V backend (no LLM).
type Scanner interface {
	Name() string
	Available() bool
	Builtin() bool
	Scan(root string) (Result, error)
}

// BlocksGate is true when any High/Critical finding is present.
func BlocksGate(findings []Finding) bool {
	for _, f := range findings {
		sev := strings.ToLower(f.Severity)
		if sev == "high" || sev == "critical" {
			return true
		}
	}
	return false
}

// Disposition returns block | advise | record for a finding set.
func Disposition(findings []Finding) string {
	if BlocksGate(findings) {
		return "block"
	}
	if len(findings) > 0 {
		return "advise"
	}
	return "record"
}
