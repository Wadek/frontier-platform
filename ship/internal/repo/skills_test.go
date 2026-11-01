package repo

// Skill-tree policy: the English for this rule is ship/skills/grok/README.md
// ("Every SKILL.md has the same shape"). This test is its witness. It walks the
// working tree (not only tracked files) so a new skill is checked before its
// first commit.

import (
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"testing"
)

const skillsDir = "ship/skills/grok"

// pstackRoots are the pstack skills a Frontier skill may name as its root.
// A backticked principle-* name outside this list is a typo or a root that
// does not exist.
var pstackRoots = map[string]bool{
	"architect": true, "arena": true, "automate-me": true, "benchmark-checklist": true,
	"blast-radius": true, "correct": true, "figure-it-out": true, "how": true,
	"interrogate": true, "poteto-mode": true, "recall": true, "reflect": true,
	"sharpen": true, "show-me-your-work": true, "statechart": true, "swarm": true,
	"tdd": true, "teach": true, "technical-writing": true, "unslop": true, "why": true,
	"principle-attack-the-premise":                       true,
	"principle-boundary-discipline":                      true,
	"principle-build-the-lever":                          true,
	"principle-encode-lessons-in-structure":              true,
	"principle-exhaust-the-design-space":                 true,
	"principle-experience-first":                         true,
	"principle-explain-the-number":                       true,
	"principle-fix-root-causes":                          true,
	"principle-foundational-thinking":                    true,
	"principle-guard-the-context-window":                 true,
	"principle-laziness-protocol":                        true,
	"principle-make-operations-idempotent":               true,
	"principle-migrate-callers-then-delete-legacy-apis":  true,
	"principle-minimize-reader-load":                     true,
	"principle-model-the-domain":                         true,
	"principle-never-block-on-the-human":                 true,
	"principle-outcome-oriented-execution":               true,
	"principle-prove-it-works":                           true,
	"principle-redesign-from-first-principles":           true,
	"principle-separate-before-serializing-shared-state": true,
	"principle-sequence-verifiable-units":                true,
	"principle-subtract-before-you-add":                  true,
	"principle-test-behavior-not-implementation":         true,
	"principle-type-system-discipline":                   true,
}

var allowedFrontmatter = map[string]bool{
	"name": true, "description": true, "when-to-use": true, "argument-hint": true,
	"user-invocable": true, "disable-model-invocation": true,
}

var (
	kebab       = regexp.MustCompile(`^[a-z0-9]+(-[a-z0-9]+)*$`)
	fmKey       = regexp.MustCompile(`^([a-z][a-z0-9-]*):\s*(.*)$`)
	mdLink      = regexp.MustCompile(`\]\(([^)\s]+)\)`)
	backticked  = regexp.MustCompile("`([^`\n]+)`")
	repoPathPfx = []string{"ship/", "control/", "scripts/", ".github/"}
)

// frontmatter parses the leading YAML block of a skill. It supports the subset
// the skills use (scalar keys and folded/literal block scalars); the module has
// no YAML dependency and should not gain one for this.
func frontmatter(src string) (map[string]string, string, bool) {
	if !strings.HasPrefix(src, "---\n") {
		return nil, src, false
	}
	end := strings.Index(src[4:], "\n---\n")
	if end < 0 {
		return nil, src, false
	}
	block, body := src[4:4+end], src[4+end+5:]
	out := map[string]string{}
	var key string
	for _, line := range strings.Split(block, "\n") {
		if m := fmKey.FindStringSubmatch(line); m != nil {
			key = m[1]
			v := strings.TrimSpace(m[2])
			if strings.HasPrefix(v, ">") || strings.HasPrefix(v, "|") {
				v = ""
			}
			out[key] = strings.Trim(v, `"`)
			continue
		}
		if key != "" && strings.HasPrefix(line, " ") {
			out[key] = strings.TrimSpace(out[key] + " " + strings.TrimSpace(line))
		}
	}
	return out, body, true
}

// section returns the text under a "## title" heading up to the next "## ".
func section(body, title string) (string, bool) {
	h := "\n## " + title + "\n"
	i := strings.Index("\n"+body, h)
	if i < 0 {
		return "", false
	}
	rest := ("\n" + body)[i+len(h):]
	if j := strings.Index(rest, "\n## "); j >= 0 {
		rest = rest[:j]
	}
	return rest, true
}

// checkRefs reports markdown links and backticked repo paths that do not exist.
func checkRefs(t *testing.T, root, file, body string) {
	t.Helper()
	dir := filepath.Dir(file)
	for _, m := range mdLink.FindAllStringSubmatch(body, -1) {
		target := m[1]
		if strings.Contains(target, "://") || strings.HasPrefix(target, "#") || strings.HasPrefix(target, "mailto:") {
			continue
		}
		target = strings.SplitN(target, "#", 2)[0]
		if _, err := os.Stat(filepath.Join(dir, filepath.FromSlash(target))); err != nil {
			t.Errorf("%s: dead link %q", rel(root, file), m[1])
		}
	}
	for _, m := range backticked.FindAllStringSubmatch(body, -1) {
		p := m[1]
		if strings.ContainsAny(p, " <>*…$") {
			continue
		}
		for _, pfx := range repoPathPfx {
			if strings.HasPrefix(p, pfx) {
				if _, err := os.Stat(filepath.Join(root, filepath.FromSlash(strings.TrimSuffix(p, "/")))); err != nil {
					t.Errorf("%s: repo path `%s` does not exist", rel(root, file), p)
				}
				break
			}
		}
	}
}

func rel(root, p string) string {
	r, err := filepath.Rel(root, p)
	if err != nil {
		return p
	}
	return filepath.ToSlash(r)
}

func requireSections(t *testing.T, root, file, body string, names ...string) {
	t.Helper()
	for _, s := range names {
		text, ok := section(body, s)
		if !ok {
			t.Errorf("%s: missing section \"## %s\"", rel(root, file), s)
			continue
		}
		switch s {
		case "Done when":
			if !strings.Contains(text, "- [ ]") {
				t.Errorf("%s: \"## Done when\" has no checklist item (\"- [ ]\")", rel(root, file))
			}
		case "Pstack root":
			known := 0
			for _, m := range backticked.FindAllStringSubmatch(text, -1) {
				n := m[1]
				if pstackRoots[n] {
					known++
				} else if strings.HasPrefix(n, "principle-") {
					t.Errorf("%s: unknown pstack root `%s`", rel(root, file), n)
				}
			}
			if known == 0 {
				t.Errorf("%s: \"## Pstack root\" names no known pstack skill", rel(root, file))
			}
		}
	}
}

func TestSkillsTree(t *testing.T) {
	root := repoRoot(t)
	base := filepath.Join(root, filepath.FromSlash(skillsDir))
	entries, err := os.ReadDir(base)
	if err != nil {
		t.Skipf("%s not present: %v", skillsDir, err)
	}
	readme, err := os.ReadFile(filepath.Join(base, "README.md"))
	if err != nil {
		t.Fatalf("%s/README.md (skill index) is required: %v", skillsDir, err)
	}
	checkRefs(t, root, filepath.Join(base, "README.md"), string(readme))

	var skills []string
	for _, e := range entries {
		if !e.IsDir() || e.Name() == "rules" {
			continue
		}
		name := e.Name()
		skills = append(skills, name)
		file := filepath.Join(base, name, "SKILL.md")
		raw, err := os.ReadFile(file)
		if err != nil {
			t.Errorf("%s/%s: missing SKILL.md", skillsDir, name)
			continue
		}
		src := strings.ReplaceAll(string(raw), "\r\n", "\n")
		fm, body, ok := frontmatter(src)
		if !ok {
			t.Errorf("%s: missing or unterminated frontmatter", rel(root, file))
			continue
		}
		for k := range fm {
			if !allowedFrontmatter[k] {
				t.Errorf("%s: unknown frontmatter key %q", rel(root, file), k)
			}
		}
		if !kebab.MatchString(name) {
			t.Errorf("%s: folder name is not kebab-case", rel(root, file))
		}
		if fm["name"] != name {
			t.Errorf("%s: frontmatter name %q != folder %q", rel(root, file), fm["name"], name)
		}
		d := fm["description"]
		switch {
		case d == "":
			t.Errorf("%s: empty description", rel(root, file))
		case len(d) > 1024:
			t.Errorf("%s: description is %d chars; limit 1024", rel(root, file), len(d))
		case !strings.Contains(strings.ToLower(d), "use when"):
			t.Errorf("%s: description must say when to use the skill (\"Use when ...\")", rel(root, file))
		}
		if !strings.Contains(body, "\n# ") && !strings.HasPrefix(strings.TrimLeft(body, "\n"), "# ") {
			t.Errorf("%s: missing H1 title", rel(root, file))
		}
		requireSections(t, root, file, body, "Pstack root", "Done when", "Conflicts")
		checkRefs(t, root, file, body)
		if !strings.Contains(string(readme), "("+name+"/SKILL.md)") {
			t.Errorf("%s/README.md: index does not link %s/SKILL.md", skillsDir, name)
		}

		// references/*.md are part of the skill: their links must resolve too.
		refs, _ := filepath.Glob(filepath.Join(base, name, "references", "*.md"))
		for _, r := range refs {
			b, err := os.ReadFile(r)
			if err != nil {
				t.Error(err)
				continue
			}
			checkRefs(t, root, r, string(b))
		}
	}
	if len(skills) == 0 {
		t.Fatalf("%s: no skills found", skillsDir)
	}

	rules, _ := filepath.Glob(filepath.Join(base, "rules", "*.md"))
	sort.Strings(rules)
	for _, r := range rules {
		b, err := os.ReadFile(r)
		if err != nil {
			t.Error(err)
			continue
		}
		body := strings.ReplaceAll(string(b), "\r\n", "\n")
		requireSections(t, root, r, body, "Pstack root", "Done when")
		checkRefs(t, root, r, body)
		if !strings.Contains(string(readme), "(rules/"+filepath.Base(r)+")") {
			t.Errorf("%s/README.md: index does not link rules/%s", skillsDir, filepath.Base(r))
		}
	}
}
