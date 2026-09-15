// Package workflow implements the control package's workflow state machine.
// Every transition is checked against this map; every transition seals a journal row (F0).
// No debrief, no "closed" (C4).
package workflow

var fsm = map[string]map[string]bool{
	"queued":      {"review": true},
	"review":      {"preparing": true, "deferred": true, "rejected": true},
	"deferred":    {"review": true},
	"preparing":   {"running": true, "rejected": true},
	"running":     {"completed": true, "transferred": true, "failed": true},
	"transferred": {"running": true, "failed": true},
	"failed":      {"review": true},
	"completed":   {"debriefed": true, "failed": true},
	"debriefed":   {"closed": true, "failed": true},
	"closed":      {},
	"rejected":    {},
}

// Terminal states end a workflow instance.
var Terminal = map[string]bool{"closed": true, "rejected": true}

// NextOk reports whether the transition is legal.
func NextOk(state, next string) bool {
	return fsm[state][next]
}

// States returns the ordered lifecycle for humans to read.
func States() []string {
	return []string{"queued", "review", "deferred", "preparing", "running",
		"transferred", "failed", "completed", "debriefed", "closed", "rejected"}
}
