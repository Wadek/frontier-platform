package workflow

import "testing"

func TestHappyPath(t *testing.T) {
	path := []string{"queued", "review", "preparing", "running", "completed", "debriefed", "closed"}
	for i := 0; i < len(path)-1; i++ {
		if !NextOk(path[i], path[i+1]) {
			t.Fatalf("legal transition denied: %s -> %s", path[i], path[i+1])
		}
	}
}

func TestGates(t *testing.T) {
	illegal := [][2]string{
		{"queued", "running"},    // no review skip
		{"completed", "closed"},  // debrief mandatory (C4)
		{"preparing", "closed"},
		{"closed", "review"},     // terminal
		{"rejected", "running"},
	}
	for _, p := range illegal {
		if NextOk(p[0], p[1]) {
			t.Errorf("illegal transition allowed: %s -> %s", p[0], p[1])
		}
	}
}

func TestBranches(t *testing.T) {
	for _, p := range [][2]string{
		{"review", "deferred"}, {"deferred", "review"},
		{"running", "transferred"}, {"transferred", "running"},
		{"running", "failed"}, {"failed", "review"},
		{"review", "rejected"},
	} {
		if !NextOk(p[0], p[1]) {
			t.Errorf("legal branch denied: %s -> %s", p[0], p[1])
		}
	}
}
