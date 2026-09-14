package sessionbudget

import "testing"

func TestCheck(t *testing.T) {
	if v := Check(95.0); v.OK || !v.RefuseNewRuns || !v.RequireHandoff || !v.RequireNewSession {
		t.Errorf("95%% must soft-refuse with handoff + new session: %+v", v)
	}
	if v := Check(94.99); !v.OK || v.RefuseNewRuns {
		t.Errorf("94.99%% must stay open: %+v", v)
	}
	if v := Check(99.9); !v.RequireHandoff {
		t.Error("even a full cache insists on the handoff (soft, never silent)")
	}
}
