package owasp

import (
	"os"
	"path/filepath"
	"testing"
)

func TestScanDetectsInjection(t *testing.T) {
	dir := t.TempDir()
	p := filepath.Join(dir, "bad.py")
	if err := os.WriteFile(p, []byte("x = eval(user_input)\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	fs, err := ScanTree(dir)
	if err != nil {
		t.Fatal(err)
	}
	if !BlocksGate(fs) {
		t.Fatalf("expected block, got %#v", fs)
	}
}

func TestCleanTree(t *testing.T) {
	dir := t.TempDir()
	p := filepath.Join(dir, "ok.py")
	if err := os.WriteFile(p, []byte("print('hello')\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	fs, err := ScanTree(dir)
	if err != nil {
		t.Fatal(err)
	}
	if BlocksGate(fs) {
		t.Fatalf("unexpected %#v", fs)
	}
}
