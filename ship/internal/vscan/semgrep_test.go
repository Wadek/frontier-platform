package vscan

import "testing"

func TestParseSemgrepJSON(t *testing.T) {
	raw := []byte(`{
	  "results": [
	    {
	      "check_id": "rules.python.security.injection",
	      "path": "app.py",
	      "start": { "line": 42 },
	      "extra": { "message": "Potential SQL injection", "severity": "ERROR" }
	    },
	    {
	      "check_id": "rules.python.style.unused-import",
	      "path": "app.py",
	      "start": { "line": 1 },
	      "extra": { "message": "Unused import os", "severity": "WARNING" }
	    }
	  ]
	}`)
	fs, meta, err := ParseSemgrepJSON(raw)
	if err != nil {
		t.Fatal(err)
	}
	if meta["count"] != 2 || len(fs) != 2 {
		t.Fatalf("want 2 got meta=%v n=%d", meta, len(fs))
	}
	if fs[0].Severity != "High" || fs[0].Line != 42 || fs[0].RuleID != "rules.python.security.injection" {
		t.Fatalf("%+v", fs[0])
	}
	if fs[1].Severity != "Medium" {
		t.Fatalf("%+v", fs[1])
	}
}

func TestParseSemgrepJSONEmpty(t *testing.T) {
	fs, meta, err := ParseSemgrepJSON([]byte(`{"results":[]}`))
	if err != nil || len(fs) != 0 || meta["count"] != 0 {
		t.Fatalf("%v %v %v", fs, meta, err)
	}
}
