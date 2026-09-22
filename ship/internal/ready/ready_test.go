package ready

import (
	"os"
	"path/filepath"
	"testing"
)

func write(t *testing.T, root, rel, body string) {
	t.Helper()
	p := filepath.Join(root, filepath.FromSlash(rel))
	if err := os.MkdirAll(filepath.Dir(p), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(p, []byte(body), 0o644); err != nil {
		t.Fatal(err)
	}
}

func TestMissingAppFails(t *testing.T) {
	root := t.TempDir()
	_ = os.Mkdir(filepath.Join(root, ".git"), 0o755)
	write(t, root, "docker-compose.yml", "services:\n  web:\n    image: nginx\n")
	rep, err := Inspect(root)
	if err != nil {
		t.Fatal(err)
	}
	if rep.Ready {
		t.Fatal("want not ready")
	}
	got := map[string]string{}
	for _, c := range rep.Checks {
		got[c.ID] = c.Status
	}
	if got["tests"] != Fail || got["ci_workflow"] != Fail || got["deploy_script"] != Fail {
		t.Fatalf("%+v", got)
	}
	if got["git"] != Pass {
		t.Fatalf("git: %s", got["git"])
	}
	if got["github_protection"] != Advise {
		t.Fatalf("protection should be advise")
	}
}

func TestReadyAppPassesFails(t *testing.T) {
	root := t.TempDir()
	_ = os.Mkdir(filepath.Join(root, ".git"), 0o755)
	write(t, root, "pytest.ini", "[pytest]\n")
	write(t, root, "docker-compose.yml", "services:\n  web:\n    image: nginx\n")
	write(t, root, "scripts/deploy.ps1", "# deploy\n")
	write(t, root, "scripts/ci/report_scanner_issues.py", "# reporter\n")
	write(t, root, "ops/github-runner/Dockerfile", "FROM scratch\n")
	write(t, root, ".github/pull_request_template.md", "## Token report\n| Frontier AI (DeepSeek) | 0 |\n")
	write(t, root, ".github/workflows/verify.yml", "name: verify\n# gitleaks trivy semgrep checkov\n# Stage 1 gate\nscanner-status\nuses: actions/checkout@11d5960a326750d5838078e36cf38b85af677262 # v4\n")
	write(t, root, ".github/workflows/deploy.yml", "name: deploy\n")
	write(t, root, ".github/workflows/cleanup-merged.yml", "name: cleanup-merged\n")
	rep, err := Inspect(root)
	if err != nil {
		t.Fatal(err)
	}
	if !rep.Ready {
		t.Fatalf("want ready, failed=%d %+v", rep.Failed, rep.Checks)
	}
	got := map[string]string{}
	for _, c := range rep.Checks {
		got[c.ID] = c.Status
	}
	if got["tests"] != Pass || got["ci_workflow"] != Pass || got["stage1_scanners"] != Pass {
		t.Fatalf("%+v", got)
	}
	if got["deploy_script"] != Pass || got["deploy_workflow"] != Pass {
		t.Fatalf("deploy %+v", got)
	}
	if got["scanner_reporter"] != Pass || got["stage1_gate"] != Pass || got["actions_pinned"] != Pass {
		t.Fatalf("new checks %+v", got)
	}
	if got["pr_token_report"] != Pass || got["runner_docker"] != Pass || got["cleanup_merged"] != Pass {
		t.Fatalf("template checks %+v", got)
	}
}

func TestPinnedActionsAdviseOnMutableTags(t *testing.T) {
	root := t.TempDir()
	_ = os.Mkdir(filepath.Join(root, ".git"), 0o755)
	write(t, root, ".github/workflows/verify.yml", "name: verify\njobs:\n  x:\n    steps:\n      - uses: actions/checkout@v4\n")
	c := checkPinnedActions(root)
	if c.Status != Advise {
		t.Fatalf("%+v", c)
	}
}

func TestPackageJSONTestsCount(t *testing.T) {
	root := t.TempDir()
	_ = os.Mkdir(filepath.Join(root, ".git"), 0o755)
	write(t, root, "package.json", `{"scripts":{"test":"node --test"}}`)
	rep, err := Inspect(root)
	if err != nil {
		t.Fatal(err)
	}
	for _, c := range rep.Checks {
		if c.ID == "tests" && c.Status != Pass {
			t.Fatalf("package.json test script should pass tests check: %+v", c)
		}
	}
}
