// Package pricing implements the platform pricing policy.
// Official source: https://api-docs.deepseek.com/quick_start/pricing
// Peak Mon-Fri 01:00-04:00 and 06:00-10:00 UTC; everything else off-peak (50% price).
package pricing

import "time"

var peakWindows = [][2]float64{{1, 4}, {6, 10}}

// IsPeak reports whether the instant is inside a weekday peak window (UTC).
func IsPeak(now time.Time) bool {
	if now.Weekday() == time.Saturday || now.Weekday() == time.Sunday {
		return false
	}
	h := float64(now.Hour()) + float64(now.Minute())/60 + float64(now.Second())/3600
	for _, w := range peakWindows {
		if h >= w[0] && h < w[1] {
			return true
		}
	}
	return false
}

// Decision is the deterministic provider verdict. A model never decides which model to use.
type Decision struct {
	Target string
	Model  string
	Egress bool
	Peak   bool
	Reason string
}

var localKinds = map[string]bool{"atomic": true, "audit": true, "ops": true, "ingest": true, "train": true}
var flashKinds = map[string]bool{"plan": true, "label": true, "census-draft": true, "tips": true, "vision": true, "coach": true}
var proKinds = map[string]bool{"census-deep": true, "verify": true, "teacher": true, "plan-pro": true, "label-pro": true}

// Route sizes a kind of work to a provider.
func Route(kind string, now time.Time) Decision {
	peak := IsPeak(now)
	switch {
	case localKinds[kind]:
		return Decision{Target: "local", Model: "qwen2.5-coder:64k", Peak: peak, Reason: "local-first"}
	case flashKinds[kind]:
		return Decision{Target: "deepseek-flash", Model: "deepseek-flash", Egress: true, Peak: peak, Reason: kind + " after egress; flash"}
	case proKinds[kind]:
		if peak {
			return Decision{Target: "deepseek-flash", Model: "deepseek-flash", Egress: true, Peak: true, Reason: "pro forbidden at peak"}
		}
		return Decision{Target: "deepseek-v4-pro", Model: "deepseek-v4-pro", Egress: true, Reason: "pro off-peak"}
	default:
		return Decision{Target: "local", Model: "qwen2.5-coder:64k", Peak: peak, Reason: "default local"}
	}
}
