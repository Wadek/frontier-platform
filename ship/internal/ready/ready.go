// Package ready detects whether an app tree already follows the standard
// merge/deploy process. It does not mutate the tree. No LLM.
package ready

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// Status is pass, fail, or advise.
const (
	Pass   = "pass"
	Fail   = "fail"
	Advise = "advise"
)

// Check is one checklist item.
type Check struct {
	ID     string `json:"id"`
	Title  string `json:"title"`
	Status string `json:"status"`
	Fix    string `json:"fix,omitempty"`
}

// Report is the machine-readable first-deploy detector output.
type Report struct {
	Root    string  `json:"root"`
	Ready   bool    `json:"ready"`
	Failed  int     `json:"failed"`
	Advised int     `json:"advised"`
	Checks  []Check `json:"checks"`
	Next    string  `json:"next"`
}

// Inspect walks root read-only.
func Inspect(root string) (*Report, error) {
	root, err := filepath.Abs(root)
	if err != nil {
		return nil, err
	}
	st, err := os.Stat(root)
	if err != nil {
		return nil, err
	}
	if !st.IsDir() {
		return nil, fmt.Errorf("not a directory: %s", root)
	}

	r := &Report{Root: root, Checks: make([]Check, 0, 8)}
	r.add(checkGit(root))
	r.add(checkTests(root))
	r.add(checkCI(root))
	r.add(checkScanners(root))
	r.add(checkDeployScript(root))
	r.add(checkDeployWorkflow(root))
	r.add(checkCompose(root))
	r.add(checkGitHubProtectionHint())

	for _, c := range r.Checks {
		switch c.Status {
		case Fail:
			r.Failed++
		case Advise:
			r.Advised++
		}
	}
	r.Ready = r.Failed == 0
	if r.Ready {
		r.Next = "plan → apply → git push on feat/* ; human PR dev → main ; then deploy script + release-check"
	} else {
		r.Next = "fix each fail hint, then run frontier ready again"
	}
	return r, nil
}

func (r *Report) add(c Check) { r.Checks = append(r.Checks, c) }

func exists(root string, rel string) bool {
	_, err := os.Stat(filepath.Join(root, filepath.FromSlash(rel)))
	return err == nil
}

func checkGit(root string) Check {
	if exists(root, ".git") {
		return Check{ID: "git", Title: "Git repository", Status: Pass}
	}
	return Check{ID: "git", Title: "Git repository", Status: Fail,
		Fix: "git init and add a GitHub remote; work on feat/<slug>, never main"}
}

func checkTests(root string) Check {
	if exists(root, "pytest.ini") || exists(root, "tests") || exists(root, "go.mod") {
		return Check{ID: "tests", Title: "Automated tests", Status: Pass}
	}
	pkg := filepath.Join(root, "package.json")
	if b, err := os.ReadFile(pkg); err == nil && strings.Contains(string(b), `"test"`) {
		return Check{ID: "tests", Title: "Automated tests", Status: Pass}
	}
	return Check{ID: "tests", Title: "Automated tests", Status: Fail,
		Fix: "add the stack test command (pytest -v, npm test, or go test ./...) and a tests/ directory"}
}

func workflowBlob(root string) string {
	dir := filepath.Join(root, ".github", "workflows")
	ents, err := os.ReadDir(dir)
	if err != nil {
		return ""
	}
	var b strings.Builder
	for _, e := range ents {
		if e.IsDir() {
			continue
		}
		name := strings.ToLower(e.Name())
		if !strings.HasSuffix(name, ".yml") && !strings.HasSuffix(name, ".yaml") {
			continue
		}
		raw, err := os.ReadFile(filepath.Join(dir, e.Name()))
		if err != nil {
			continue
		}
		b.Write(raw)
		b.WriteByte('\n')
	}
	return strings.ToLower(b.String())
}

func checkCI(root string) Check {
	blob := workflowBlob(root)
	if blob == "" {
		return Check{ID: "ci_workflow", Title: "GitHub Actions verify/ci workflow", Status: Fail,
			Fix: "copy ship/templates/release/verify.yml.example to .github/workflows/verify.yml"}
	}
	return Check{ID: "ci_workflow", Title: "GitHub Actions verify/ci workflow", Status: Pass}
}

func checkScanners(root string) Check {
	blob := workflowBlob(root)
	if blob == "" {
		return Check{ID: "stage1_scanners", Title: "Stage 1 OSS scanners in CI", Status: Fail,
			Fix: "add Gitleaks, Semgrep, Trivy fs, and Checkov to verify.yml (see verify.yml.example)"}
	}
	need := []string{"gitleaks", "trivy", "semgrep", "checkov"}
	var miss []string
	for _, n := range need {
		if !strings.Contains(blob, n) {
			miss = append(miss, n)
		}
	}
	if len(miss) == 0 {
		return Check{ID: "stage1_scanners", Title: "Stage 1 OSS scanners in CI", Status: Pass}
	}
	return Check{ID: "stage1_scanners", Title: "Stage 1 OSS scanners in CI", Status: Advise,
		Fix: "CI exists but missing: " + strings.Join(miss, ", ") + " — copy names from ship/templates/release/verify.yml.example"}
}

func checkDeployScript(root string) Check {
	for _, p := range []string{"scripts/deploy.ps1", "scripts/deploy.sh", "scripts/deploy.py"} {
		if exists(root, p) {
			return Check{ID: "deploy_script", Title: "Deploy script", Status: Pass}
		}
	}
	return Check{ID: "deploy_script", Title: "Deploy script", Status: Fail,
		Fix: "copy ship/templates/release/ into scripts/; prod runs frontier release-check then compose up"}
}

func checkDeployWorkflow(root string) Check {
	blob := workflowBlob(root)
	if strings.Contains(blob, "deploy") && (exists(root, ".github/workflows/deploy.yml") || exists(root, ".github/workflows/deploy.yaml")) {
		return Check{ID: "deploy_workflow", Title: "Deploy workflow", Status: Pass}
	}
	if exists(root, "scripts/deploy.ps1") || exists(root, "scripts/deploy.sh") {
		return Check{ID: "deploy_workflow", Title: "Deploy workflow", Status: Advise,
			Fix: "optional: copy ship/templates/release/deploy.yml.example to .github/workflows/deploy.yml"}
	}
	return Check{ID: "deploy_workflow", Title: "Deploy workflow", Status: Advise,
		Fix: "add deploy.yml after the deploy script exists (template: ship/templates/release/deploy.yml.example)"}
}

func checkCompose(root string) Check {
	for _, n := range []string{"docker-compose.yml", "docker-compose.yaml", "compose.yml", "compose.yaml"} {
		if exists(root, n) {
			return Check{ID: "compose", Title: "Compose file", Status: Pass}
		}
	}
	if exists(root, "Dockerfile") {
		return Check{ID: "compose", Title: "Compose file", Status: Advise,
			Fix: "Dockerfile present; add compose.yml if this app is deployed with Compose"}
	}
	return Check{ID: "compose", Title: "Compose file", Status: Advise,
		Fix: "if this is a container app, add compose.yml joining your private network; bind loopback only"}
}

func checkGitHubProtectionHint() Check {
	return Check{ID: "github_protection", Title: "GitHub branch protection on main", Status: Advise,
		Fix: "in the GitHub UI: require PR, require CI, deny force-push; Frontier does not replace this"}
}

// RenderText is the short English brief (token-cheap).
func RenderText(r *Report) string {
	var b strings.Builder
	fmt.Fprintf(&b, "frontier ready: root=%s ready=%v failed=%d advised=%d\n", r.Root, r.Ready, r.Failed, r.Advised)
	for _, c := range r.Checks {
		fmt.Fprintf(&b, "  [%s] %s  %s\n", c.Status, c.ID, c.Title)
		if c.Fix != "" && c.Status != Pass {
			fmt.Fprintf(&b, "      fix: %s\n", c.Fix)
		}
	}
	fmt.Fprintf(&b, "next: %s\n", r.Next)
	return b.String()
}

// MarshalJSON report helper.
func Marshal(r *Report) ([]byte, error) {
	return json.MarshalIndent(r, "", "  ")
}
