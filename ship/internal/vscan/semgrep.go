package vscan

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"
)

// SemgrepScanner runs Semgrep when installed (SAST — programmatic, no tokens).
type SemgrepScanner struct{}

func (SemgrepScanner) Name() string  { return "semgrep" }
func (SemgrepScanner) Builtin() bool { return false }
func (SemgrepScanner) Available() bool {
	_, err := exec.LookPath("semgrep")
	return err == nil
}

func (s SemgrepScanner) Scan(root string) (Result, error) {
	if !s.Available() {
		return Result{Source: "semgrep", Skipped: true, SkipWhy: "semgrep not on PATH"}, nil
	}
	// Match the Stage 1 CI contract: p/ci ruleset, JSON to stdout.
	cmd := exec.Command("semgrep", "scan", "--config", "p/ci", "--json", "--quiet")
	cmd.Dir = root
	out, err := cmd.Output()
	if len(out) == 0 && err != nil {
		return Result{Source: "semgrep", Skipped: true, SkipWhy: fmt.Sprintf("semgrep run failed: %v", err)}, nil
	}
	findings, meta, parseErr := ParseSemgrepJSON(out)
	if parseErr != nil {
		return Result{Source: "semgrep", Skipped: true, SkipWhy: parseErr.Error()}, nil
	}
	findings = DropSkippedPaths(findings)
	if meta == nil {
		meta = map[string]any{}
	}
	meta["count"] = len(findings)
	return Result{Source: "semgrep", Findings: findings, Meta: meta}, nil
}

type semgrepReport struct {
	Results []semgrepResult `json:"results"`
}

type semgrepResult struct {
	CheckID string `json:"check_id"`
	Path    string `json:"path"`
	Start   struct {
		Line int `json:"line"`
	} `json:"start"`
	Extra struct {
		Message  string `json:"message"`
		Severity string `json:"severity"`
	} `json:"extra"`
}

// ParseSemgrepJSON maps `semgrep scan --json` to Findings.
func ParseSemgrepJSON(raw []byte) ([]Finding, map[string]any, error) {
	s := strings.TrimSpace(string(raw))
	if s == "" || s == "null" {
		return nil, map[string]any{"count": 0}, nil
	}
	var report semgrepReport
	if err := json.Unmarshal([]byte(s), &report); err != nil {
		return nil, nil, fmt.Errorf("semgrep json: %w", err)
	}
	out := make([]Finding, 0, len(report.Results))
	for _, r := range report.Results {
		out = append(out, Finding{
			Source:   "semgrep",
			RuleID:   nz(r.CheckID, "semgrep.unknown"),
			Severity: mapSemgrepSeverity(r.Extra.Severity),
			Path:     r.Path,
			Line:     r.Start.Line,
			Snippet:  trimSnippet(r.Extra.Message),
		})
	}
	return out, map[string]any{"count": len(out)}, nil
}

func mapSemgrepSeverity(s string) string {
	switch strings.ToUpper(strings.TrimSpace(s)) {
	case "ERROR", "HIGH", "CRITICAL":
		return "High"
	case "WARNING", "MEDIUM":
		return "Medium"
	case "INFO", "LOW":
		return "Low"
	default:
		return "Medium"
	}
}
