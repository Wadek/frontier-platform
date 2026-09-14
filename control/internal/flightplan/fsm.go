// Package flightplan implements the tower's workflow state machine.
// Every transition is checked against this map; every transition seals a journal row (F0).
// No debrief, no "filed" (T4).
package flightplan

var fsm = map[string]map[string]bool{
	"queued":     {"review": true},
	"review":     {"taxiing": true, "holding": true, "rejected": true},
	"holding":    {"review": true},
	"taxiing":    {"airborne": true, "rejected": true},
	"airborne":   {"landed": true, "handed_off": true, "failed": true},
	"handed_off": {"airborne": true, "failed": true},
	"failed":     {"review": true},
	"landed":     {"debriefed": true, "failed": true},
	"debriefed":  {"filed": true, "failed": true},
	"filed":      {},
	"rejected":   {},
}

// Terminal states end a workflow instance.
var Terminal = map[string]bool{"filed": true, "rejected": true}

// NextOk reports whether the transition is legal.
func NextOk(state, next string) bool {
	return fsm[state][next]
}

// States returns the ordered lifecycle for humans to read.
func States() []string {
	return []string{"queued", "review", "holding", "taxiing", "airborne",
		"handed_off", "failed", "landed", "debriefed", "filed", "rejected"}
}
