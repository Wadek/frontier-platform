package role

import "testing"

func TestElevateLadder(t *testing.T) {
	r := Observer
	for _, want := range []Role{Analyst, Operator, Executor} {
		next, err := r.Elevate()
		if err != nil {
			t.Fatal(err)
		}
		if next != want {
			t.Fatalf("got %s want %s", next, want)
		}
		r = next
	}
	if _, err := r.Elevate(); err == nil {
		t.Fatal("expected error at top")
	}
}

func TestCan(t *testing.T) {
	if !Operator.Can(Analyst) {
		t.Fatal("operator should include analyst")
	}
	if Executor.Can(Observer) != true {
		t.Fatal("executor includes observer")
	}
	if Observer.Can(Operator) {
		t.Fatal("observer cannot operator")
	}
}
