package learning

import (
	"path/filepath"
	"testing"
)

func sample() Record {
	return Record{
		Session: "2026-09-14-x", Plane: "food", Pilot: "p1",
		Runway: "local", TS: "2026-09-14T08:30:00Z", TLDR: "t",
		Learned: []string{"a"}, Next: []string{"b"},
	}
}

func TestValidate(t *testing.T) {
	if err := Validate(sample()); err != nil {
		t.Fatalf("valid record rejected: %v", err)
	}
	bad := sample()
	bad.Runway = "claude"
	if err := Validate(bad); err == nil {
		t.Error("claude runway accepted (T7: no non-DeepSeek cloud)")
	}
	bad = sample()
	bad.TS = "not-a-time"
	if err := Validate(bad); err == nil {
		t.Error("bad ts accepted")
	}
	bad = sample()
	bad.Session = ""
	if err := Validate(bad); err == nil {
		t.Error("missing session accepted")
	}
}

func TestAppendRead(t *testing.T) {
	dir := t.TempDir()
	p := filepath.Join(dir, "learn.jsonl")
	if err := Append(p, sample()); err != nil {
		t.Fatal(err)
	}
	r2 := sample()
	r2.Session = "s2"
	if err := Append(p, r2); err != nil {
		t.Fatal(err)
	}
	if err := Append(p, Record{Session: "x"}); err == nil {
		t.Error("invalid record appended")
	}
	rows, err := Read(p)
	if err != nil {
		t.Fatal(err)
	}
	if len(rows) != 2 {
		t.Fatalf("read %d rows, want 2", len(rows))
	}
}
