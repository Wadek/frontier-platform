package learning

import "testing"

func TestValidateMethod(t *testing.T) {
	if err := ValidateMethod("retrieval"); err != nil {
		t.Fatalf("retrieval rejected: %v", err)
	}
	if err := ValidateMethod("youtube"); err == nil {
		t.Fatal("youtube accepted; YouTube is a source, not a method")
	}
}

func TestSelect(t *testing.T) {
	got := Select(KindSpeech)
	if len(got) != 2 || got[0] != MethodListen || got[1] != MethodDual {
		t.Fatalf("speech select = %v", got)
	}
	if Select(Kind("nope")) != nil {
		t.Fatal("unknown kind should return nil")
	}
	if Select(KindSession)[0] != MethodDebrief {
		t.Fatal("session must debrief")
	}
}
