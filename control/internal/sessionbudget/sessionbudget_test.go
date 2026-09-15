package sessionbudget

import "testing"

func TestSoftRefuse(t *testing.T) {
	const limit int64 = 1000
	if SoftRefuse(0, limit) {
		t.Error("0/1000 should allow")
	}
	if SoftRefuse(949, limit) {
		t.Error("949/1000 should allow (under 95%)")
	}
	if !SoftRefuse(950, limit) {
		t.Error("950/1000 should soft-refuse (95%)")
	}
	if !SoftRefuse(1000, limit) {
		t.Error("full budget should soft-refuse")
	}
	if !SoftRefuse(0, 0) {
		t.Error("zero limit should soft-refuse")
	}
	if !SoftRefuse(10, -1) {
		t.Error("negative limit should soft-refuse")
	}
}

func TestRemaining(t *testing.T) {
	const limit int64 = 1000
	if got := Remaining(0, limit); got != 950 {
		t.Errorf("Remaining(0,1000)=%d want 950", got)
	}
	if got := Remaining(950, limit); got != 0 {
		t.Errorf("Remaining(950,1000)=%d want 0", got)
	}
	if got := Remaining(960, limit); got != -10 {
		t.Errorf("Remaining(960,1000)=%d want -10", got)
	}
}
