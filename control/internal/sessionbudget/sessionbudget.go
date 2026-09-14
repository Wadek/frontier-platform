// Package sessionbudget implements the 95% rule (T9):
// at 95% session context/cache, soft-refuse new runs — update the handoff
// and open a new session. Soft: the refusal is an insistence on artifacts,
// never a silent stop.
package sessionbudget

// SoftLimit is the fleet-wide context/cache threshold.
const SoftLimit = 95.0

// Verdict is the deterministic budget decision.
type Verdict struct {
	OK               bool
	RefuseNewRuns    bool
	RequireHandoff   bool
	RequireNewSession bool
	Reason           string
}

// Check decides whether a session may start new work at this usage percentage.
func Check(usagePct float64) Verdict {
	if usagePct >= SoftLimit {
		return Verdict{
			RefuseNewRuns: true, RequireHandoff: true, RequireNewSession: true,
			Reason: "session context/cache >= 95% - update the handoff and open a new session",
		}
	}
	return Verdict{OK: true, Reason: "budget ok"}
}
