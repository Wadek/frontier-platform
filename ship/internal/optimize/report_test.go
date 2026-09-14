package optimize

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestBuildAndWriteReport(t *testing.T) {
	dir := t.TempDir()
	// one large python function
	var b strings.Builder
	b.WriteString("def big():\n")
	for i := 0; i < 100; i++ {
		b.WriteString("    x = 1\n")
	}
	_ = os.WriteFile(filepath.Join(dir, "app.py"), []byte(b.String()), 0o644)
	_ = os.WriteFile(filepath.Join(dir, "small.py"), []byte("def tiny():\n    return 1\n"), 0o644)

	r, err := BuildReport(dir)
	if err != nil {
		t.Fatal(err)
	}
	if len(r.Findings) < 1 {
		t.Fatal("expected at least one hotspot finding")
	}
	if r.Findings[0].ID != "Opt-001" {
		t.Fatalf("id=%s", r.Findings[0].ID)
	}
	art, err := WriteArtifacts(r)
	if err != nil {
		t.Fatal(err)
	}
	md, err := os.ReadFile(art.Markdown)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(md), "### Opt-001") {
		t.Fatal("markdown missing Opt-001 section")
	}

	loaded, err := LoadLatest(dir)
	if err != nil {
		t.Fatal(err)
	}
	f, err := FindingByID(loaded, "Opt-1")
	if err != nil {
		t.Fatal(err)
	}
	body := FormatFinding(*f)
	if !strings.Contains(body, "Intended behavior") {
		t.Fatal("pr-body fragment incomplete")
	}
}
