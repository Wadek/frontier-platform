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
	write(t, root, ".github/workflows/verify.yml", "name: verify\n# gitleaks trivy semgrep checkov\n")
	write(t, root, ".github/workflows/deploy.yml", "name: deploy\n")
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
