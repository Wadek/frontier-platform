// Package optimize implements Optimize (O): behavior-preserving speed reports.
package optimize

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

const (
	MaxFindings   = 25
	MaxBriefBytes = 24 * 1024
	MaxFuncLines  = 80 // hotspot threshold for programmatic candidates
)

// Finding is one Opt-* review unit for a developer / PR.
type Finding struct {
	ID                 string `json:"id"` // Opt-001
	Title              string `json:"title"`
	Path               string `json:"path"`
	Function           string `json:"function,omitempty"`
	StartLine          int    `json:"start_line,omitempty"`
	EndLine            int    `json:"end_line,omitempty"`
	IntendedBehavior   string `json:"intended_behavior"`
	WhyWasteful        string `json:"why_wasteful"`
	SuggestedChange    string `json:"suggested_change"`
	SnippetBefore      string `json:"snippet_before,omitempty"`
	SnippetAfter       string `json:"snippet_after,omitempty"`
	RiskEquivalence    string `json:"risk_equivalence"`
	Disposition        string `json:"disposition"` // advise default
	Source             string `json:"source"`      // hotspot|enhance|manual
}

// Report is one Optimize pass over a project root.
type Report struct {
	Root      string    `json:"root"`
	Name      string    `json:"name"`
	Stamp     string    `json:"stamp"`
	Findings  []Finding `json:"findings"`
	Notes     []string  `json:"notes,omitempty"`
	PRProcess string    `json:"pr_process"`
}

// Artifacts written under .frontier/optimize/.
type Artifacts struct {
	Dir      string
	JSON     string
	Markdown string
	Latest   string
	Stamp    string
}

// BuildReport scans root for programmatic Optimize candidates and builds a report.
func BuildReport(root string) (*Report, error) {
	root, err := filepath.Abs(root)
	if err != nil {
		return nil, err
	}
	info, err := os.Stat(root)
	if err != nil {
		return nil, err
	}
	if !info.IsDir() {
		return nil, fmt.Errorf("not a directory: %s", root)
	}

	r := &Report{
		Root:  root,
		Name:  filepath.Base(root),
		Stamp: time.Now().UTC().Format("20060102T150405Z"),
		PRProcess: "One Opt-ID per branch (frontier/opt-<id>-…). PR body = this report section. " +
			"Behavior must stay equivalent. Advise-only until promoted.",
		Notes: []string{
			"Programmatic pass: large-function hotspots only (heuristic).",
			"Use frontier enhance optimize for CS-depth residual fills.",
			"Do not change intended behavior.",
		},
	}

	hots, err := scanLargeFunctions(root)
	if err != nil {
		return nil, err
	}
	for i, h := range hots {
		if i >= MaxFindings {
			r.Notes = append(r.Notes, fmt.Sprintf("truncated hotspots at %d", MaxFindings))
			break
		}
		id := fmt.Sprintf("Opt-%03d", i+1)
		r.Findings = append(r.Findings, Finding{
			ID:        id,
			Title:     fmt.Sprintf("Large function hotspot: %s", h.FuncName),
			Path:      h.Path,
			Function:  h.FuncName,
			StartLine: h.Start,
			EndLine:   h.End,
			IntendedBehavior: "Preserve current observable results and API contract for callers of " +
				h.FuncName + ".",
			WhyWasteful: fmt.Sprintf(
				"Function spans ~%d lines (threshold %d). Large routines often hide repeated work, "+
					"poor data-structure fit, or mixed hot/cold paths (CS: locality, complexity, SRP).",
				h.End-h.Start+1, MaxFuncLines),
			SuggestedChange: "Read for algorithmic waste; extract cold paths; replace nested scans with "+
				"appropriate structures only when equality/ordering semantics match. Prefer the simplest "+
				"equivalent transform.",
			SnippetBefore:   h.Snippet,
			SnippetAfter:    "(propose in PR after human/enhance review — no auto-rewrite in Phase B)",
			RiskEquivalence: "Any extract/refactor must keep inputs→outputs identical; add/adjust tests covering this function.",
			Disposition:     "advise",
			Source:          "hotspot",
		})
	}
	if len(r.Findings) == 0 {
		r.Notes = append(r.Notes, "No large-function hotspots under threshold — enhance optimize may still add Opt-* items.")
	}
	return r, nil
}

// WriteArtifacts writes O-<stamp>.{md,json} and LATEST under .frontier/optimize/.
func WriteArtifacts(r *Report) (*Artifacts, error) {
	dir := filepath.Join(r.Root, ".frontier", "optimize")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return nil, err
	}
	base := "O-" + r.Stamp
	mdPath := filepath.Join(dir, base+".md")
	jsPath := filepath.Join(dir, base+".json")
	latest := filepath.Join(dir, "LATEST")

	raw, err := json.MarshalIndent(r, "", "  ")
	if err != nil {
		return nil, err
	}
	if err := os.WriteFile(jsPath, raw, 0o644); err != nil {
		return nil, err
	}
	md := RenderMarkdown(r)
	if len(md) > MaxBriefBytes {
		md = md[:MaxBriefBytes] + "\n\n...[truncated for token budget]\n"
	}
	if err := os.WriteFile(mdPath, []byte(md), 0o644); err != nil {
		return nil, err
	}
	_ = os.WriteFile(latest, []byte(base+"\n"), 0o644)
	return &Artifacts{Dir: dir, JSON: jsPath, Markdown: mdPath, Latest: latest, Stamp: r.Stamp}, nil
}

// RenderMarkdown is the developer-facing Optimize report.
func RenderMarkdown(r *Report) string {
	var b strings.Builder
	fmt.Fprintf(&b, "# Frontier Optimize (O) — %s\n\n", r.Name)
	fmt.Fprintf(&b, "- root: `%s`\n", r.Root)
	fmt.Fprintf(&b, "- stamp: `%s`\n", r.Stamp)
	fmt.Fprintf(&b, "- findings: **%d** (advise-only)\n", len(r.Findings))
	fmt.Fprintf(&b, "- process: %s\n\n", r.PRProcess)
	if len(r.Notes) > 0 {
		b.WriteString("## Notes\n\n")
		for _, n := range r.Notes {
			fmt.Fprintf(&b, "- %s\n", n)
		}
		b.WriteString("\n")
	}
	if len(r.Findings) == 0 {
		b.WriteString("_No programmatic Optimize candidates. Run `frontier enhance optimize` for residual CS review, or raise MaxFuncLines sensitivity later._\n")
		return b.String()
	}
	b.WriteString("## Findings\n\n")
	b.WriteString("Open **one PR per Opt-ID** (or a tight cluster). Paste that section into the PR body.\n\n")
	for _, f := range r.Findings {
		b.WriteString(FormatFinding(f))
		b.WriteString("\n")
	}
	b.WriteString("## Commands\n\n")
	b.WriteString("```text\nfrontier optimize pr-body Opt-001\n# branch: frontier/opt-001-<short>\n# PR body <- printed markdown\n```\n")
	return b.String()
}

// FormatFinding renders one Opt-* section (PR body fragment).
func FormatFinding(f Finding) string {
	var b strings.Builder
	fmt.Fprintf(&b, "### %s — %s\n", f.ID, f.Title)
	loc := f.Path
	if f.Function != "" {
		loc += "  `" + f.Function + "`"
	}
	if f.StartLine > 0 {
		fmt.Fprintf(&b, "- Location: `%s`  L%d–L%d\n", loc, f.StartLine, f.EndLine)
	} else {
		fmt.Fprintf(&b, "- Location: `%s`\n", loc)
	}
	fmt.Fprintf(&b, "- Intended behavior: %s\n", f.IntendedBehavior)
	fmt.Fprintf(&b, "- Why wasteful: %s\n", f.WhyWasteful)
	fmt.Fprintf(&b, "- Suggested change: %s\n", f.SuggestedChange)
	if f.SnippetBefore != "" {
		b.WriteString("- Snippet (before):\n\n```\n")
		b.WriteString(trimSnippet(f.SnippetBefore, 40))
		b.WriteString("\n```\n")
	}
	if f.SnippetAfter != "" {
		fmt.Fprintf(&b, "- Snippet (after): %s\n", f.SnippetAfter)
	}
	fmt.Fprintf(&b, "- Risk / equivalence: %s\n", f.RiskEquivalence)
	fmt.Fprintf(&b, "- Disposition: **%s** (source: %s)\n", f.Disposition, f.Source)
	return b.String()
}

// LoadLatest reads the latest Optimize report JSON from root.
func LoadLatest(root string) (*Report, error) {
	root, _ = filepath.Abs(root)
	latest := filepath.Join(root, ".frontier", "optimize", "LATEST")
	b, err := os.ReadFile(latest)
	if err != nil {
		return nil, fmt.Errorf("no optimize report yet — run: frontier optimize report (%w)", err)
	}
	base := strings.TrimSpace(string(b))
	jsPath := filepath.Join(root, ".frontier", "optimize", base+".json")
	raw, err := os.ReadFile(jsPath)
	if err != nil {
		return nil, err
	}
	var r Report
	if err := json.Unmarshal(raw, &r); err != nil {
		return nil, err
	}
	return &r, nil
}

// FindingByID returns one finding from a report.
func FindingByID(r *Report, id string) (*Finding, error) {
	id = strings.TrimSpace(id)
	for i := range r.Findings {
		if strings.EqualFold(r.Findings[i].ID, id) {
			f := r.Findings[i]
			return &f, nil
		}
	}
	// allow Opt-1 vs Opt-001
	want := normalizeOptID(id)
	for i := range r.Findings {
		if normalizeOptID(r.Findings[i].ID) == want {
			f := r.Findings[i]
			return &f, nil
		}
	}
	return nil, fmt.Errorf("finding %q not in latest report", id)
}

func normalizeOptID(id string) string {
	id = strings.ToUpper(strings.TrimSpace(id))
	id = strings.TrimPrefix(id, "OPT-")
	id = strings.TrimPrefix(id, "OPT")
	id = strings.Trim(id, "- ")
	n := 0
	fmt.Sscanf(id, "%d", &n)
	if n > 0 {
		return fmt.Sprintf("OPT-%03d", n)
	}
	return strings.ToUpper(id)
}

func trimSnippet(s string, maxLines int) string {
	lines := strings.Split(s, "\n")
	if len(lines) <= maxLines {
		return strings.TrimRight(s, "\n")
	}
	return strings.Join(lines[:maxLines], "\n") + "\n…"
}

// ListFindingIDs sorted.
func ListFindingIDs(r *Report) []string {
	out := make([]string, 0, len(r.Findings))
	for _, f := range r.Findings {
		out = append(out, f.ID)
	}
	sort.Strings(out)
	return out
}
