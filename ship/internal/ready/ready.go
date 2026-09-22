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
	r.add(checkScannerReporter(root))
	r.add(checkStage1Gate(root))
	r.add(checkPinnedActions(root))
	r.add(checkPRTokenReport(root))
	r.add(checkRunnerDocker(root))
	r.add(checkCleanupWorkflow(root))
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
		Fix: "git init and add a GitHub remote; work on feat/<slug>, never dev or main"}
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
	return Check{ID: "github_protection", Title: "GitHub branch protection on dev and main", Status: Advise,
		Fix: "in the GitHub UI: require PR into both dev and main, require CI, deny force-push; Frontier does not replace this"}
}

func checkScannerReporter(root string) Check {
	if exists(root, "scripts/ci/report_scanner_issues.py") {
		return Check{ID: "scanner_reporter", Title: "Scanner → GitHub issues reporter", Status: Pass}
	}
	return Check{ID: "scanner_reporter", Title: "Scanner → GitHub issues reporter", Status: Advise,
		Fix: "copy ship/templates/release/scripts/ci/report_scanner_issues.py into scripts/ci/"}
}

func checkStage1Gate(root string) Check {
	blob := workflowBlob(root)
	if blob == "" {
		return Check{ID: "stage1_gate", Title: "Stage 1 hard-fail gate", Status: Fail,
			Fix: "copy ship/templates/release/verify.yml.example (includes Stage 1 gate after scanners)"}
	}
	if strings.Contains(blob, "stage 1 gate") || strings.Contains(blob, "scanner-status") {
		return Check{ID: "stage1_gate", Title: "Stage 1 hard-fail gate", Status: Pass}
	}
	if strings.Contains(blob, "soft-fail") || strings.Contains(blob, "--exit-code 0") {
		return Check{ID: "stage1_gate", Title: "Stage 1 hard-fail gate", Status: Advise,
			Fix: "replace soft-fail / --exit-code 0 scanners with the hard-fail verify.yml.example contract"}
	}
	return Check{ID: "stage1_gate", Title: "Stage 1 hard-fail gate", Status: Advise,
		Fix: "add a final Stage 1 gate that fails when gitleaks/trivy/checkov/semgrep exit non-zero"}
}

func checkPinnedActions(root string) Check {
	blob := workflowBlob(root)
	if blob == "" {
		return Check{ID: "actions_pinned", Title: "Actions uses: pinned to commit SHAs", Status: Fail,
			Fix: "copy verify.yml.example (SHA-pinned uses:)"}
	}
	// Heuristic: mutable tags look like @v1 / @v4 / @main without a 40-char hex nearby on the same line.
	mutable := 0
	for _, line := range strings.Split(blob, "\n") {
		line = strings.TrimSpace(line)
		if !strings.Contains(line, "uses:") {
			continue
		}
		if strings.Contains(line, "@v") || strings.HasSuffix(line, "@main") || strings.Contains(line, "@master") {
			// SHA-pinned lines still often end with " # v4" — require a 40-char hex somewhere on the line.
			if !hasActionSHA(line) {
				mutable++
			}
		}
	}
	if mutable == 0 {
		return Check{ID: "actions_pinned", Title: "Actions uses: pinned to commit SHAs", Status: Pass}
	}
	return Check{ID: "actions_pinned", Title: "Actions uses: pinned to commit SHAs", Status: Advise,
		Fix: fmt.Sprintf("%d mutable uses: tag(s) — pin to full commit SHAs (see verify.yml.example)", mutable)}
}

func hasActionSHA(line string) bool {
	hex := 0
	for _, r := range line {
		if (r >= '0' && r <= '9') || (r >= 'a' && r <= 'f') {
			hex++
			if hex >= 40 {
				return true
			}
		} else {
			hex = 0
		}
	}
	return false
}

func checkPRTokenReport(root string) Check {
	for _, p := range []string{".github/pull_request_template.md", "pull_request_template.md", ".github/PULL_REQUEST_TEMPLATE.md"} {
		b, err := os.ReadFile(filepath.Join(root, filepath.FromSlash(p)))
		if err != nil {
			continue
		}
		low := strings.ToLower(string(b))
		if strings.Contains(low, "token report") || strings.Contains(low, "frontier ai") {
			return Check{ID: "pr_token_report", Title: "PR template includes token report", Status: Pass}
		}
		return Check{ID: "pr_token_report", Title: "PR template includes token report", Status: Advise,
			Fix: "add the Token report table from ship/templates/release/pull_request_template.md"}
	}
	return Check{ID: "pr_token_report", Title: "PR template includes token report", Status: Advise,
		Fix: "copy ship/templates/release/pull_request_template.md to .github/pull_request_template.md"}
}

func checkRunnerDocker(root string) Check {
	if exists(root, "ops/github-runner/Dockerfile") || exists(root, "ops/github-runner/compose.yml") {
		return Check{ID: "runner_docker", Title: "Docker Compose self-hosted runner template", Status: Pass}
	}
	blob := workflowBlob(root)
	if strings.Contains(blob, "self-hosted") {
		return Check{ID: "runner_docker", Title: "Docker Compose self-hosted runner template", Status: Advise,
			Fix: "workflows use self-hosted; copy ship/templates/runner-docker/ to ops/github-runner/"}
	}
	return Check{ID: "runner_docker", Title: "Docker Compose self-hosted runner template", Status: Advise,
		Fix: "optional: copy ship/templates/runner-docker/ when using Linux self-hosted runners"}
}

func checkCleanupWorkflow(root string) Check {
	blob := workflowBlob(root)
	if strings.Contains(blob, "cleanup-merged") || (strings.Contains(blob, "delete-branch") && strings.Contains(blob, "closes #")) {
		return Check{ID: "cleanup_merged", Title: "Post-merge issue/branch cleanup workflow", Status: Pass}
	}
	if exists(root, ".github/workflows/cleanup-merged.yml") || exists(root, ".github/workflows/cleanup-merged.yaml") {
		return Check{ID: "cleanup_merged", Title: "Post-merge issue/branch cleanup workflow", Status: Pass}
	}
	return Check{ID: "cleanup_merged", Title: "Post-merge issue/branch cleanup workflow", Status: Advise,
		Fix: "copy ship/templates/release/cleanup-merged.yml.example to .github/workflows/cleanup-merged.yml"}
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
