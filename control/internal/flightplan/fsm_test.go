package flightplan

import "testing"

func TestHappyPath(t *testing.T) {
	path := []string{"queued", "review", "taxiing", "airborne", "landed", "debriefed", "filed"}
	for i := 0; i < len(path)-1; i++ {
		if !NextOk(path[i], path[i+1]) {
			t.Fatalf("legal transition denied: %s -> %s", path[i], path[i+1])
		}
	}
}

func TestGates(t *testing.T) {
	illegal := [][2]string{
		{"queued", "airborne"}, // no clearance skip
		{"landed", "filed"},    // debrief mandatory (T4)
		{"taxiing", "filed"},
		{"filed", "review"},    // terminal
		{"rejected", "airborne"},
	}
	for _, p := range illegal {
		if NextOk(p[0], p[1]) {
			t.Errorf("illegal transition allowed: %s -> %s", p[0], p[1])
		}
	}
}

func TestBranches(t *testing.T) {
	for _, p := range [][2]string{
		{"review", "holding"}, {"holding", "review"},
		{"airborne", "handed_off"}, {"handed_off", "airborne"},
		{"airborne", "failed"}, {"failed", "review"},
		{"review", "rejected"},
	} {
		if !NextOk(p[0], p[1]) {
			t.Errorf("legal branch denied: %s -> %s", p[0], p[1])
		}
	}
}
