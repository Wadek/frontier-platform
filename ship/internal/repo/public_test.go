// Package repo holds tests that assert invariants about this repository's own
// tree rather than about runtime behaviour.
//
// English is the policy; these tests are the machine-checked witness (F5). The
// policy here is the public-tree rule: this repository is published, so it must
// not leak the author's local filesystem, must not carry stale repository
// identities, and must not ship cp1252 mojibake into user-facing output.
//
// These defects are not hypothetical. All three were present in the public tree
// and are exactly what this test now prevents from returning.
package repo

import (
	"bytes"
	"fmt"
	"go/parser"
	"go/token"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
	"unicode/utf8"
)

// skip lists tracked paths exempt from a rule, with the reason. Keep this short
// and justified; an unexplained entry means the rule is wrong.
var skip = map[string]string{
	// This file necessarily contains the forbidden patterns as literals.
	"ship/internal/repo/public_test.go": "contains the patterns by construction",
}

type rule struct {
	name     string
	patterns []string
	why      string
}

var rules = []rule{
	{
		name:     "stale repository identity",
		patterns: []string{"frontier-ship"},
		why: "the repository is Wadek/frontier-platform; 'frontier-ship' is the retired " +
			"private name and produces clone URLs and release links that 404",
	},
	{
		name: "habitat-local path",
		patterns: []string{
			`D:\wakalabs`, `D:/wakalabs`, `D:\\wakalabs`,
			`C:\Users\waka`, `C:/Users/waka`, `C:\\Users\\waka`,
		},
		why: "this tree is public; absolute paths from the author's machine are meaningless " +
			"to a reader and leak local layout",
	},
	{
		name:     "cp1252 mojibake",
		patterns: []string{"\u00e2\u20ac", "\u00e2\u2020", "\u00c3\u00a2", "\u00c2\u00b7"},
		why: "UTF-8 bytes were decoded as cp1252 and re-saved; the result renders as " +
			"'â€”' in CLI help and docs",
	},
}

// textExt are the kinds of file whose contents are expected to be UTF-8 text.
var textExt = map[string]bool{
	".md": true, ".txt": true, ".go": true, ".py": true, ".ps1": true,
	".sh": true, ".yml": true, ".yaml": true, ".json": true, ".puml": true,
	".mod": true, ".example": true, ".js": true,
}

func repoRoot(t *testing.T) string {
	t.Helper()
	dir, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			t.Skip("not inside a module; skipping repository policy test")
		}
		dir = parent
	}
}

// trackedFiles returns every git-tracked path, relative to the repository root.
// Only tracked files are checked: local scratch is not part of the published tree.
func trackedFiles(t *testing.T, root string) []string {
	t.Helper()
	cmd := exec.Command("git", "-C", root, "ls-files", "-z")
	out, err := cmd.Output()
	if err != nil {
		t.Skipf("git ls-files unavailable (%v); skipping repository policy test", err)
	}
	var paths []string
	for _, p := range bytes.Split(out, []byte{0}) {
		if len(p) > 0 {
			paths = append(paths, string(p))
		}
	}
	if len(paths) == 0 {
		t.Skip("no tracked files; skipping repository policy test")
	}
	return paths
}

func isText(rel string) bool {
	if textExt[strings.ToLower(filepath.Ext(rel))] {
		return true
	}
	base := filepath.Base(rel)
	return base == ".gitignore" || base == ".dockerignore"
}

func TestPublicTreePolicy(t *testing.T) {
	root := repoRoot(t)
	files := trackedFiles(t, root)

	type violation struct {
		rule   string
		file   string
		line   int
		text   string
		reason string
	}
	var found []violation

	for _, rel := range files {
		if reason, ok := skip[rel]; ok {
			_ = reason
			continue
		}
		if !isText(rel) {
			continue
		}
		raw, err := os.ReadFile(filepath.Join(root, rel))
		if err != nil {
			t.Fatalf("read %s: %v", rel, err)
		}
		if bytes.IndexByte(raw, 0) >= 0 {
			continue // binary
		}
		body := string(raw)
		if !utf8.ValidString(body) {
			found = append(found, violation{
				rule: "not valid UTF-8", file: rel, line: 0,
				text: "", reason: "tracked text file is not valid UTF-8",
			})
			continue
		}
		for _, r := range rules {
			for _, pat := range r.patterns {
				idx := strings.Index(body, pat)
				if idx < 0 {
					continue
				}
				line := 1 + strings.Count(body[:idx], "\n")
				found = append(found, violation{
					rule: r.name, file: rel, line: line,
					text: snippet(body, idx), reason: r.why,
				})
			}
		}
	}

	if len(found) == 0 {
		return
	}
	for _, v := range found {
		t.Errorf("public-tree policy violation [%s]\n  %s:%d\n  found: %q\n  why:   %s",
			v.rule, v.file, v.line, v.text, v.reason)
	}
	t.Fatalf("%d public-tree policy violation(s) in %d tracked files", len(found), len(files))
}

// snippet returns a short single-line window around idx, for readable failures.
func snippet(body string, idx int) string {
	start := idx - 30
	if start < 0 {
		start = 0
	}
	end := idx + 50
	if end > len(body) {
		end = len(body)
	}
	s := body[start:end]
	if nl := strings.IndexByte(s, '\n'); nl >= 0 {
		s = s[:nl]
	}
	return strings.TrimSpace(s)
}

func TestGitignoreAnchorsHabitatRuntime(t *testing.T) {
	root := repoRoot(t)
	raw, err := os.ReadFile(filepath.Join(root, ".gitignore"))
	if err != nil {
		t.Fatal(err)
	}
	var anchored bool
	for i, line := range strings.Split(string(raw), "\n") {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" || strings.HasPrefix(trimmed, "#") {
			continue
		}
		switch trimmed {
		case "runtime/", "runtime", "**/runtime/", "**/runtime":
			t.Errorf(".gitignore:%d %q ignores ship/internal/runtime; use /runtime/", i+1, trimmed)
		case "/runtime/", "/runtime":
			anchored = true
		}
	}
	if !anchored {
		t.Fatal(".gitignore must ignore /runtime/ (habitat) without hiding ship/internal/runtime")
	}
}

func TestRuntimePackageIsTracked(t *testing.T) {
	root := repoRoot(t)
	need := "ship/internal/runtime/runtime.go"
	for _, f := range trackedFiles(t, root) {
		if filepath.ToSlash(f) == need {
			return
		}
	}
	t.Fatalf("%s is not git-tracked; CI cannot import github.com/Wadek/frontier-platform/ship/internal/runtime", need)
}

func TestTrackedGoImportsResolveInTree(t *testing.T) {
	root := repoRoot(t)
	files := trackedFiles(t, root)
	tracked := make(map[string]bool, len(files))
	for _, f := range files {
		tracked[filepath.ToSlash(f)] = true
	}

	const mod = "github.com/Wadek/frontier-platform"
	var missing []string
	for _, rel := range files {
		if !strings.HasSuffix(rel, ".go") {
			continue
		}
		slash := filepath.ToSlash(rel)
		if strings.Contains(slash, "/testdata/") {
			continue
		}
		fset := token.NewFileSet()
		f, err := parser.ParseFile(fset, filepath.Join(root, rel), nil, parser.ImportsOnly)
		if err != nil {
			t.Fatalf("parse %s: %v", rel, err)
		}
		for _, spec := range f.Imports {
			path, err := strconv.Unquote(spec.Path.Value)
			if err != nil {
				continue
			}
			if path != mod && !strings.HasPrefix(path, mod+"/") {
				continue
			}
			dir := "."
			if path != mod {
				dir = strings.TrimPrefix(path, mod+"/")
			}
			found := false
			for tpath := range tracked {
				if !strings.HasSuffix(tpath, ".go") {
					continue
				}
				pkgDir := filepath.ToSlash(filepath.Dir(tpath))
				if pkgDir == dir {
					found = true
					break
				}
			}
			if !found {
				missing = append(missing, fmt.Sprintf("%s imports %s (no tracked .go under %s)", rel, path, dir))
			}
		}
	}
	if len(missing) > 0 {
		t.Fatalf("imported packages missing from the published tree:\n  %s", strings.Join(missing, "\n  "))
	}
}

func TestCIWorkflowCoversDevAndMain(t *testing.T) {
	root := repoRoot(t)
	raw, err := os.ReadFile(filepath.Join(root, ".github/workflows/ci.yml"))
	if err != nil {
		t.Fatal(err)
	}
	body := string(raw)
	if !strings.Contains(body, "branches: [main, dev]") && !strings.Contains(body, "branches: [dev, main]") {
		t.Fatal("CI must run on push to both main and dev")
	}
	if !strings.Contains(body, "name: ci") {
		t.Fatal("CI must expose a concluding job named ci for required status checks")
	}
}
