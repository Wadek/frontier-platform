package owasp

import (
	"path/filepath"
	"runtime"
	"testing"
)

func testdataRoot(t *testing.T) string {
	t.Helper()
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("caller")
	}
	// internal/owasp → repo root/testdata/owasp
	return filepath.Clean(filepath.Join(filepath.Dir(file), "..", "..", "testdata", "owasp"))
}

func TestScorecard_PositivesBlockOrMatch(t *testing.T) {
	root := testdataRoot(t)
	cases := []string{
		"a01_pos", "a02_pos", "a03_pos", "a04_pos", "a05_pos",
		"a06_pos", "a07_pos", "a08_pos", "a09_pos", "a10_pos",
	}
	var tp, fn int
	for _, c := range cases {
		fs, err := ScanTree(filepath.Join(root, c))
		if err != nil {
			t.Fatal(err)
		}
		if len(fs) == 0 {
			fn++
			t.Errorf("FN: %s expected findings", c)
			continue
		}
		tp++
	}
	t.Logf("positive recall cases: TP=%d FN=%d (want FN=0)", tp, fn)
}

func TestScorecard_NegativesClean(t *testing.T) {
	root := testdataRoot(t)
	cases := []string{
		"a01_neg", "a02_neg", "a03_neg", "a04_neg", "a05_neg",
		"a06_neg", "a07_neg", "a08_neg", "a09_neg", "a10_neg",
	}
	var tn, fp int
	for _, c := range cases {
		fs, err := ScanTree(filepath.Join(root, c))
		if err != nil {
			t.Fatal(err)
		}
		if BlocksGate(fs) || len(fs) > 0 {
			// a04/a05/a06/a09/a10 may match medium — treat any finding as FP for scorecard v0
			fp++
			t.Errorf("FP: %s unexpected %#v", c, fs)
			continue
		}
		tn++
	}
	t.Logf("negative precision cases: TN=%d FP=%d (want FP=0)", tn, fp)
}
