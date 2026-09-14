package vscan

import "github.com/Wadek/frontier-ship/internal/owasp"

// OWASPScanner wraps the built-in OWASP Top10 v0 regex ScanTree.
type OWASPScanner struct{}

func (OWASPScanner) Name() string     { return "owasp-v0" }
func (OWASPScanner) Available() bool  { return true }
func (OWASPScanner) Builtin() bool    { return true }

func (OWASPScanner) Scan(root string) (Result, error) {
	fs, err := owasp.ScanTree(root)
	if err != nil {
		return Result{Source: "owasp-v0"}, err
	}
	out := make([]Finding, 0, len(fs))
	for _, f := range fs {
		out = append(out, Finding{
			Source:   "owasp-v0",
			RuleID:   f.RuleID,
			OWASP:    f.OWASP,
			Severity: f.Severity,
			Path:     f.Path,
			Line:     f.Line,
			Snippet:  f.Snippet,
		})
	}
	return Result{
		Source:   "owasp-v0",
		Findings: out,
		Meta:     map[string]any{"policy": "OWASP-Top10-2021-v0", "count": len(out)},
	}, nil
}
