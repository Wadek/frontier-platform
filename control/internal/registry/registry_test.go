package registry

import (
	"os"
	"path/filepath"
	"testing"
)

func mkdirs(t *testing.T, root string, names ...string) {
	t.Helper()
	for _, n := range names {
		if err := os.MkdirAll(filepath.Join(root, n), 0o755); err != nil {
			t.Fatal(err)
		}
	}
}

func TestScan(t *testing.T) {
	root := t.TempDir()
	mkdirs(t, root, "alpha", "beta", ".agent_alpha", ".hidden", "runtime")

	projects, err := Scan(root, nil)
	if err != nil {
		t.Fatal(err)
	}
	got := map[string]Project{}
	for _, p := range projects {
		got[p.Name] = p
	}
	if len(got) != 2 {
		t.Fatalf("scanned %d projects, want 2 (alpha,beta): %v", len(got), got)
	}
	if _, ok := got["alpha"]; !ok {
		t.Error("alpha missing")
	}
	if _, ok := got["beta"]; !ok {
		t.Error("beta missing")
	}
	if _, ok := got["runtime"]; ok {
		t.Error("runtime should be excluded by DefaultExclude")
	}
	if _, ok := got[".hidden"]; ok {
		t.Error("hidden dir should be skipped")
	}
	if !got["alpha"].HasAgent {
		t.Error("alpha should report HasAgent")
	}
	if got["beta"].HasAgent {
		t.Error("beta should not report HasAgent")
	}
}

func TestScanSorted(t *testing.T) {
	root := t.TempDir()
	mkdirs(t, root, "zeta", "alpha", "mid")
	projects, err := Scan(root, nil)
	if err != nil {
		t.Fatal(err)
	}
	want := []string{"alpha", "mid", "zeta"}
	for i, p := range projects {
		if p.Name != want[i] {
			t.Fatalf("position %d = %q, want %q", i, p.Name, want[i])
		}
	}
}

func TestScanMissingRoot(t *testing.T) {
	if _, err := Scan(filepath.Join(t.TempDir(), "nope"), nil); err == nil {
		t.Error("expected error for missing root")
	}
}

func TestScanCustomExclude(t *testing.T) {
	root := t.TempDir()
	mkdirs(t, root, "alpha", "beta")
	projects, err := Scan(root, map[string]bool{"beta": true})
	if err != nil {
		t.Fatal(err)
	}
	if len(projects) != 1 || projects[0].Name != "alpha" {
		t.Fatalf("custom exclude ignored: %+v", projects)
	}
}
