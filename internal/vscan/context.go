package vscan

import (
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
)

// Budget caps keep enhance briefs token-cheap.
const (
	MaxBriefFindings = 25
	MaxBriefPaths    = 40
	MaxBriefBytes    = 12 * 1024
)

// Pack is the programmatic context passed into enhance (and optionally shown by V).
type Pack struct {
	Root           string            `json:"root"`
	ScopeMode      string            `json:"scope_mode"` // "diff" | "full"
	ChangedPaths   []string          `json:"changed_paths"`
	LangCounts     map[string]int    `json:"lang_counts"`
	Manifests      []string          `json:"manifests"`
	Results        []Result          `json:"results"`
	Findings       []Finding         `json:"findings"`
	FindingsCapped []Finding         `json:"findings_capped"`
	BySource       map[string]int    `json:"by_source"`
	Disposition    string            `json:"disposition"`
	AdaptersRun    []string          `json:"adapters_run"`
	AdaptersSkip   map[string]string `json:"adapters_skip"`
}

// BuildPack runs builtin OWASP + requested adapters and collects scope/inventory.
func BuildPack(root string, opts Options) (*Pack, error) {
	root, _ = filepath.Abs(root)
	p := &Pack{
		Root:         root,
		LangCounts:   map[string]int{},
		BySource:     map[string]int{},
		AdaptersSkip: map[string]string{},
	}
	p.ChangedPaths, p.ScopeMode = changedPaths(root)
	p.LangCounts, p.Manifests = inventory(root)

	// Always run built-in.
	builtin, err := OWASPScanner{}.Scan(root)
	if err != nil {
		return nil, err
	}
	p.Results = append(p.Results, builtin)
	p.AdaptersRun = append(p.AdaptersRun, builtin.Source)

	want := map[string]bool{}
	for _, a := range opts.Adapters {
		want[strings.ToLower(a)] = true
	}
	if opts.AutoAdapters || os.Getenv("FRONTIER_V_AUTO") == "1" {
		for _, s := range Registry() {
			if !s.Builtin() && s.Available() {
				want[s.Name()] = true
			}
		}
	}

	for _, s := range Registry() {
		if s.Builtin() {
			continue
		}
		if !want[s.Name()] {
			continue
		}
		if !s.Available() {
			res, _ := s.Scan(root) // stub/skip reason
			p.AdaptersSkip[s.Name()] = res.SkipWhy
			p.Results = append(p.Results, res)
			continue
		}
		res, err := s.Scan(root)
		if err != nil {
			p.AdaptersSkip[s.Name()] = err.Error()
			continue
		}
		if res.Skipped {
			p.AdaptersSkip[s.Name()] = res.SkipWhy
		} else {
			p.AdaptersRun = append(p.AdaptersRun, s.Name())
		}
		p.Results = append(p.Results, res)
	}

	for _, r := range p.Results {
		for _, f := range r.Findings {
			p.Findings = append(p.Findings, f)
			p.BySource[f.Source]++
		}
	}
	p.Findings = clusterFindings(p.Findings)
	p.FindingsCapped = capFindings(p.Findings, MaxBriefFindings)
	p.Disposition = Disposition(p.Findings)
	if len(p.ChangedPaths) > MaxBriefPaths {
		p.ChangedPaths = append([]string{}, p.ChangedPaths[:MaxBriefPaths]...)
	}
	return p, nil
}

func changedPaths(root string) ([]string, string) {
	for _, base := range []string{"main", "master", "origin/main", "origin/master"} {
		cmd := exec.Command("git", "diff", "--name-only", base+"...HEAD")
		cmd.Dir = root
		out, err := cmd.Output()
		if err != nil {
			continue
		}
		lines := splitNonEmpty(string(out))
		if len(lines) == 0 {
			continue
		}
		return lines, "diff:"+base
	}
	return nil, "full"
}

func inventory(root string) (map[string]int, []string) {
	langs := map[string]int{}
	var manifests []string
	manifestNames := map[string]bool{
		"go.mod": true, "package.json": true, "requirements.txt": true,
		"pipfile": true, "cargo.toml": true, "composer.json": true,
		"gemfile": true, "pom.xml": true, "build.gradle": true,
	}
	_ = filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}
		if info.IsDir() {
			name := info.Name()
			if name == ".git" || name == "node_modules" || name == "vendor" || name == ".frontier" ||
				name == "dist" || name == "build" || name == "testdata" || name == "__pycache__" ||
				name == "venv" || name == ".venv" || name == "site-packages" {
				return filepath.SkipDir
			}
			return nil
		}
		base := strings.ToLower(filepath.Base(path))
		rel, _ := filepath.Rel(root, path)
		rel = filepath.ToSlash(rel)
		ext := strings.ToLower(filepath.Ext(base))
		switch ext {
		case ".go", ".py", ".js", ".ts", ".tsx", ".jsx", ".rb", ".php", ".java", ".cs",
			".tf", ".yml", ".yaml", ".json", ".sh", ".ps1":
			langs[ext]++
		}
		if base == "dockerfile" || strings.HasPrefix(base, "dockerfile.") {
			manifests = append(manifests, rel)
		}
		if manifestNames[base] || strings.HasSuffix(base, ".tf") || strings.HasPrefix(base, ".env") {
			manifests = append(manifests, rel)
		}
		return nil
	})
	sort.Strings(manifests)
	if len(manifests) > 30 {
		manifests = manifests[:30]
	}
	return langs, manifests
}

func clusterFindings(in []Finding) []Finding {
	type key struct{ src, rule, path string }
	seen := map[key]Finding{}
	order := []key{}
	for _, f := range in {
		k := key{f.Source, f.RuleID, f.Path}
		if _, ok := seen[k]; ok {
			continue
		}
		seen[k] = f
		order = append(order, k)
	}
	out := make([]Finding, 0, len(order))
	for _, k := range order {
		out = append(out, seen[k])
	}
	return out
}

func capFindings(in []Finding, n int) []Finding {
	if len(in) <= n {
		return in
	}
	// Prefer High/Critical first.
	sorted := append([]Finding{}, in...)
	sort.SliceStable(sorted, func(i, j int) bool {
		return sevRank(sorted[i].Severity) > sevRank(sorted[j].Severity)
	})
	return sorted[:n]
}

func sevRank(s string) int {
	switch strings.ToLower(s) {
	case "critical":
		return 4
	case "high":
		return 3
	case "medium":
		return 2
	case "low":
		return 1
	default:
		return 0
	}
}

func splitNonEmpty(s string) []string {
	var out []string
	for _, line := range strings.Split(s, "\n") {
		line = strings.TrimSpace(line)
		if line != "" {
			out = append(out, filepath.ToSlash(line))
		}
	}
	return out
}
