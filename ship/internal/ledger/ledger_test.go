package ledger

import (
	"path/filepath"
	"testing"
)

func TestAppendAndTail(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "ledger.jsonl")
	l, err := Open(path)
	if err != nil {
		t.Fatal(err)
	}
	e1, err := l.Append("observer", "observe", map[string]any{"ok": true})
	if err != nil {
		t.Fatal(err)
	}
	if e1.PrevHash != "genesis" {
		t.Fatalf("prev=%s", e1.PrevHash)
	}
	e2, err := l.Append("analyst", "analyze", map[string]any{"n": 1})
	if err != nil {
		t.Fatal(err)
	}
	if e2.PrevHash != e1.EntryHash {
		t.Fatal("chain broken")
	}
	tail, err := l.Tail(10)
	if err != nil || len(tail) != 2 {
		t.Fatalf("tail=%v err=%v", tail, err)
	}
	last, err := l.LastAction("observe")
	if err != nil || last == nil || last.Action != "observe" {
		t.Fatalf("last=%v err=%v", last, err)
	}
}
