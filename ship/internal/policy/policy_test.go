package policy

import "testing"

func TestEvaluatePushGate_BlocksMain(t *testing.T) {
	g := EvaluatePushGate("main", "abc", "", false)
	if g.OK {
		t.Fatal("expected main push blocked")
	}
	if len(g.Codes) != 1 || g.Codes[0] != CodeRefusePushMain {
		t.Fatalf("codes=%v want [%s]", g.Codes, CodeRefusePushMain)
	}
}

func TestEvaluatePushGate_BlocksDev(t *testing.T) {
	g := EvaluatePushGate("dev", "abc", "", false)
	if g.OK {
		t.Fatal("expected dev push blocked")
	}
	if len(g.Codes) != 1 || g.Codes[0] != CodeRefusePushDev {
		t.Fatalf("codes=%v want [%s]", g.Codes, CodeRefusePushDev)
	}
}

func TestReleaseAuthorized_DevRefusalFails(t *testing.T) {
	g := EvaluatePushGate("dev", "abc", "", false)
	if ReleaseAuthorized(g) {
		t.Fatal("dev is integration, not production; release-check must fail")
	}
}

func TestReleaseAuthorized_MainRefusalOnly(t *testing.T) {
	g := EvaluatePushGate("main", "abc", "", false)
	if !ReleaseAuthorized(g) {
		t.Fatal("main refusal alone should authorize release")
	}
	g.AddCode(CodeOWASPBlock, MsgOWASPBlock)
	if ReleaseAuthorized(g) {
		t.Fatal("OWASP block must fail release-check")
	}
}

func TestReleaseAuthorized_DirtyFails(t *testing.T) {
	g := EvaluatePushGate("main", "abc", " M x.go", false)
	if ReleaseAuthorized(g) {
		t.Fatal("dirty tree must fail release-check")
	}
}

func TestReleaseAuthorized_CleanFeatureOK(t *testing.T) {
	g := EvaluatePushGate("feat/x", "abc", "", false)
	if !ReleaseAuthorized(g) {
		t.Fatal("no blocking codes should authorize release")
	}
}

func TestDeriveCodes_LegacyMainString(t *testing.T) {
	codes := DeriveCodes([]string{MsgRefusePushMain})
	if len(codes) != 1 || codes[0] != CodeRefusePushMain {
		t.Fatalf("got %v", codes)
	}
	g := GateResult{Reasons: []string{MsgRefusePushMain}}
	if !ReleaseAuthorized(g) {
		t.Fatal("legacy reason string should still authorize release")
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
