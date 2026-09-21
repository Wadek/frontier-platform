package vscan

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestParseCheckovJSON(t *testing.T) {
	raw := []byte(`{
  "summary": {"passed": 2, "failed": 1},
  "results": {
    "failed_checks": [{
      "check_id": "CKV_AWS_20",
      "check_name": "S3 Bucket has an ACL defined",
      "file_path": "/Dockerfile",
      "file_line_range": [1, 3],
      "severity": "HIGH"
    }]
  }
}`)
	fs, meta, err := ParseCheckovJSON(raw)
	if err != nil {
		t.Fatal(err)
	}
	if len(fs) != 1 {
		t.Fatalf("want 1 finding, got %d", len(fs))
	}
	if fs[0].RuleID != "CKV_AWS_20" || fs[0].Severity != "High" || fs[0].Source != "checkov" {
		t.Fatalf("unexpected finding: %+v", fs[0])
	}
	if meta["failed"] != 1 {
		t.Fatalf("meta failed=%v", meta["failed"])
	}
}

func TestClusterAndCap(t *testing.T) {
	in := []Finding{
		{Source: "owasp-v0", RuleID: "a", Path: "x.go", Severity: "Low"},
		{Source: "owasp-v0", RuleID: "a", Path: "x.go", Severity: "Low"},
		{Source: "owasp-v0", RuleID: "b", Path: "y.go", Severity: "Critical"},
		{Source: "checkov", RuleID: "c", Path: "z.tf", Severity: "Medium"},
	}
	c := clusterFindings(in)
	if len(c) != 3 {
		t.Fatalf("cluster want 3 got %d", len(c))
	}
	cap := capFindings(c, 2)
	if len(cap) != 2 {
		t.Fatalf("cap want 2 got %d", len(cap))
	}
	if cap[0].Severity != "Critical" {
		t.Fatalf("expected Critical first, got %s", cap[0].Severity)
	}
}

func TestWriteEnhanceBriefBudget(t *testing.T) {
	dir := t.TempDir()
	// minimal fake tree
	_ = os.WriteFile(filepath.Join(dir, "main.go"), []byte("package main\n"), 0o644)
	p := &Pack{
		Root:         dir,
		ScopeMode:    "full",
		Disposition:  "record",
		AdaptersRun:  []string{"owasp-v0"},
		BySource:     map[string]int{},
		AdaptersSkip: map[string]string{"checkov": "not on PATH"},
		LangCounts:   map[string]int{".go": 1},
	}
	art, err := WriteEnhanceBrief(dir, p)
	if err != nil {
		t.Fatal(err)
	}
	b, err := os.ReadFile(art.Markdown)
	if err != nil {
		t.Fatal(err)
	}
	if len(b) == 0 || len(b) > MaxBriefBytes {
		t.Fatalf("brief size %d", len(b))
	}
	if !strings.Contains(string(b), "Do **not** re-litigate") {
		t.Fatal("missing thrift mission text")
	}
}
