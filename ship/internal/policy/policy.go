package policy

import (
	"fmt"
	"strings"
	"time"

	"github.com/Wadek/frontier-platform/ship/internal/ledger"
	"github.com/Wadek/frontier-platform/ship/internal/role"
)

// Deny > Ask > Allow evaluated in Require.
func Require(have role.Role, min role.Role, tool string) error {
	if !have.Can(min) {
		return fmt.Errorf("deny: tool %s needs role >= %s (have %s)", tool, min, have)
	}
	return nil
}

// Stable reason codes for plan --json and release-check (do not reword casually).
const (
	CodeDetachedHEAD   = "detached_head"
	CodeMissingHEAD    = "missing_head"
	CodeDirtyTree      = "dirty_tree"
	CodeRefusePushMain = "refuse_push_main"
	CodeRefusePushDev  = "refuse_push_dev"
	CodeOWASPBlock     = "owasp_high_or_critical"
	CodeHygieneBlock   = "hygiene_block"
)

// Legacy human reason strings (still emitted for display / older classifiers).
const (
	MsgRefusePushMain = "refusing direct push to main/master (use a feature branch)"
	MsgRefusePushDev  = "refusing direct push to dev (open a PR from feat/* into dev)"
	MsgDirtyTree      = "working tree dirty; commit or clean before push"
	MsgDetachedHEAD   = "detached HEAD or empty branch"
	MsgMissingHEAD    = "missing HEAD"
	MsgOWASPBlock     = "OWASP V: untriaged High/Critical finding(s)"
	MsgHygieneBlock   = "Hygiene H: untriaged provenance marks (FRONTIER_HYGIENE_BLOCK=1)"
)

// GateResult is a sealed pre-push policy check.
type GateResult struct {
	OK        bool           `json:"ok"`
	Codes     []string       `json:"codes"`
	Reasons   []string       `json:"reasons"`
	Branch    string         `json:"branch"`
	Head      string         `json:"head"`
	Dirty     bool           `json:"dirty"`
	SealHash  string         `json:"seal,omitempty"`
	ExpiresAt string         `json:"expires_at,omitempty"`
	Extras    map[string]any `json:"extras,omitempty"`
}

// GateTTL is how long a plan.passed / gate.passed seal stays fresh for push.
// Exported so the monitor can re-verify freshness from ledger timestamps.
const GateTTL = 15 * time.Minute

// DirtyPorcelain reports whether porcelain output has real changes,
// ignoring Frontier's own ledger/metadata under .frontier/.
func DirtyPorcelain(porcelain string) bool {
	for _, line := range strings.Split(porcelain, "\n") {
		line = strings.TrimRight(line, "\r")
		if strings.TrimSpace(line) == "" {
			continue
		}
		// porcelain: XY<space>path (path starts at index 3 when standard)
		path := line
		if len(line) >= 3 {
			path = strings.TrimSpace(line[3:])
		}
		path = strings.TrimPrefix(path, "\"")
		if strings.HasPrefix(path, ".frontier/") || path == ".frontier" {
			continue
		}
		return true
	}
	return false
}

// EvaluatePushGate is local-first and cheap: no model calls.
func EvaluatePushGate(branch, head, porcelain string, allowDirty bool) GateResult {
	var codes, reasons []string
	dirty := DirtyPorcelain(porcelain)
	if branch == "" {
		codes = append(codes, CodeDetachedHEAD)
		reasons = append(reasons, MsgDetachedHEAD)
	}
	if head == "" {
		codes = append(codes, CodeMissingHEAD)
		reasons = append(reasons, MsgMissingHEAD)
	}
	if dirty && !allowDirty {
		codes = append(codes, CodeDirtyTree)
		reasons = append(reasons, MsgDirtyTree)
	}
	if IsProduction(branch) {
		codes = append(codes, CodeRefusePushMain)
		reasons = append(reasons, MsgRefusePushMain)
	} else if IsIntegration(branch) {
		codes = append(codes, CodeRefusePushDev)
		reasons = append(reasons, MsgRefusePushDev)
	}
	ok := len(codes) == 0
	return GateResult{
		OK:      ok,
		Codes:   codes,
		Reasons: reasons,
		Branch:  branch,
		Head:    head,
		Dirty:   dirty,
	}
}

// AddCode appends a stable code and matching human reason if not already present.
func (g *GateResult) AddCode(code, msg string) {
	for _, c := range g.Codes {
		if c == code {
			return
		}
	}
	g.Codes = append(g.Codes, code)
	g.Reasons = append(g.Reasons, msg)
	g.OK = false
}

// DeriveCodes maps legacy reason strings to stable codes (for older ledgers/tests).
func DeriveCodes(reasons []string) []string {
	var out []string
	for _, r := range reasons {
		switch {
		case r == MsgRefusePushMain || strings.Contains(r, "refusing direct push to main"):
			out = append(out, CodeRefusePushMain)
		case r == MsgRefusePushDev || strings.Contains(r, "refusing direct push to dev"):
			out = append(out, CodeRefusePushDev)
		case r == MsgDirtyTree || strings.Contains(r, "working tree dirty"):
			out = append(out, CodeDirtyTree)
		case r == MsgDetachedHEAD || strings.Contains(r, "detached HEAD"):
			out = append(out, CodeDetachedHEAD)
		case r == MsgMissingHEAD || strings.Contains(r, "missing HEAD"):
			out = append(out, CodeMissingHEAD)
		case r == MsgOWASPBlock || strings.Contains(r, "OWASP V:"):
			out = append(out, CodeOWASPBlock)
		case r == MsgHygieneBlock || strings.Contains(r, "Hygiene H:"):
			out = append(out, CodeHygieneBlock)
		default:
			out = append(out, "unknown:"+r)
		}
	}
	return out
}

// IsProduction is main/master: production deploy is allowed, direct push is not.
func IsProduction(branch string) bool {
	b := strings.ToLower(strings.TrimSpace(branch))
	return b == "main" || b == "master"
}

// IsIntegration is the integration branch. Code arrives only by PR from feat/*.
func IsIntegration(branch string) bool {
	return strings.EqualFold(strings.TrimSpace(branch), "dev")
}

// IsProtected must only move via pull request (feat -> dev -> main).
func IsProtected(branch string) bool {
	return IsProduction(branch) || IsIntegration(branch)
}

// ReleaseAuthorized reports whether a GateResult may authorize a production deploy.
// Push to main remains refused; that single refusal is tolerated for release mode
// (deploy runs from a checkout of main). A refusal to push to dev is not a
// deploy authorization: integration is not production.
// Any other code (dirty, OWASP, hygiene, unknown) fails closed.
func ReleaseAuthorized(g GateResult) bool {
	codes := g.Codes
	if len(codes) == 0 && len(g.Reasons) > 0 {
		codes = DeriveCodes(g.Reasons)
	}
	for _, c := range codes {
		if c == CodeRefusePushMain {
			continue
		}
		return false
	}
	return true
}

// SealGate writes gate.passed or gate.failed to the ledger.
func SealGate(l *ledger.Ledger, actor string, g GateResult) (*GateResult, error) {
	action := "gate.failed"
	if g.OK {
		action = "gate.passed"
		g.ExpiresAt = time.Now().UTC().Add(GateTTL).Format(time.RFC3339)
	}
	payload := map[string]any{
		"ok":         g.OK,
		"codes":      g.Codes,
		"reasons":    g.Reasons,
		"branch":     g.Branch,
		"head":       g.Head,
		"dirty":      g.Dirty,
		"expires_at": g.ExpiresAt,
	}
	e, err := l.Append(actor, action, payload)
	if err != nil {
		return nil, err
	}
	g.SealHash = e.EntryHash
	return &g, nil
}

// FreshGateOK reports whether a recent gate.passed matches current head/branch.
func FreshGateOK(l *ledger.Ledger, branch, head string) (bool, string) {
	return freshAction(l, "gate.passed", branch, head, "run: git frontier apply  (or gate)")
}

// FreshPlanOK reports whether a recent plan.passed matches current head/branch.
func FreshPlanOK(l *ledger.Ledger, branch, head string) (bool, string) {
	return freshAction(l, "plan.passed", branch, head, "run: git frontier plan")
}

func freshAction(l *ledger.Ledger, action, branch, head, hint string) (bool, string) {
	e, err := l.LastAction(action)
	if err != nil || e == nil {
		return false, "no " + action + " in ledger; " + hint
	}
	exp, _ := e.Payload["expires_at"].(string)
	if exp != "" {
		t, err := time.Parse(time.RFC3339, exp)
		if err == nil && time.Now().UTC().After(t) {
			return false, action + " expired; " + hint
		}
	}
	b, _ := e.Payload["branch"].(string)
	h, _ := e.Payload["head"].(string)
	if b != branch || h != head {
		return false, action + " was for different branch/HEAD; " + hint
	}
	return true, e.EntryHash
}

// SealPlan writes plan.passed or plan.failed (Terraform-like preview state).
func SealPlan(l *ledger.Ledger, actor string, g GateResult) (*GateResult, error) {
	action := "plan.failed"
	if g.OK {
		action = "plan.passed"
		g.ExpiresAt = time.Now().UTC().Add(GateTTL).Format(time.RFC3339)
	}
	payload := map[string]any{
		"ok":         g.OK,
		"codes":      g.Codes,
		"reasons":    g.Reasons,
		"branch":     g.Branch,
		"head":       g.Head,
		"dirty":      g.Dirty,
		"expires_at": g.ExpiresAt,
		"V":          "OWASP-Top10-2021-v0",
		"S":          "not_enforced_yet",
	}
	e, err := l.Append(actor, action, payload)
	if err != nil {
		return nil, err
	}
	g.SealHash = e.EntryHash
	return &g, nil
}
