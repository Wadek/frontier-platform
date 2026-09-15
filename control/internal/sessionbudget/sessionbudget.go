// Package sessionbudget implements session budget policy C9 (session budget rule).
// Soft-refuse further work when usage reaches SoftRefuseRatio of the session limit.
package sessionbudget

// SoftRefuseRatio is the fraction of the session budget at which new work is soft-refused.
const SoftRefuseRatio = 0.95

// SoftRefuse reports whether used tokens have reached the soft-refuse threshold of limit.
// A non-positive limit is treated as exhausted (refuse).
func SoftRefuse(used, limit int64) bool {
	if limit <= 0 {
		return true
	}
	return float64(used) >= SoftRefuseRatio*float64(limit)
}

// Remaining returns tokens left before the soft-refuse threshold.
// Negative means already at or past the soft-refuse line.
func Remaining(used, limit int64) int64 {
	if limit <= 0 {
		return 0
	}
	thresh := int64(SoftRefuseRatio * float64(limit))
	return thresh - used
}
