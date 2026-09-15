package policy

import "testing"

func TestEvaluatePushGate_BlocksMain(t *testing.T) {
	g := EvaluatePushGate("main", "abc", "", false)
	if g.OK {
		t.Fatal("expected main push blocked")
	}
}

func TestEvaluatePushGate_Dirty(t *testing.T) {
	g := EvaluatePushGate("feat/x", "abc", " M file.go", false)
	if g.OK {
		t.Fatal("expected dirty blocked")
	}
	g2 := EvaluatePushGate("feat/x", "abc", " M file.go", true)
	if !g2.OK {
		t.Fatal("expected allow_dirty to pass")
	}
}

func TestEvaluatePushGate_OK(t *testing.T) {
	g := EvaluatePushGate("feat/x", "abc123", "", false)
	if !g.OK {
		t.Fatalf("expected ok, got %v", g.Reasons)
	}
}

func TestDirtyPorcelain_IgnoresFrontier(t *testing.T) {
	if DirtyPorcelain("?? .frontier/ledger.jsonl\n") {
		t.Fatal(".frontier should be ignored")
	}
	if !DirtyPorcelain(" M README.md\n?? .frontier/ledger.jsonl\n") {
		t.Fatal("real dirty file should count")
	}
}
