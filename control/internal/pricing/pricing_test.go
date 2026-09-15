package pricing

import (
	"testing"
	"time"
)

func utc(y, mo, d, h, mi int) time.Time {
	return time.Date(y, time.Month(mo), d, h, mi, 0, 0, time.UTC)
}

// 2026-09-14 is a Monday; 2026-09-12/13 are Saturday/Sunday.
func TestIsPeakWindows(t *testing.T) {
	cases := []struct {
		in   time.Time
		want bool
	}{
		{utc(2026, 9, 14, 1, 0), true},
		{utc(2026, 9, 14, 3, 59), true},
		{utc(2026, 9, 14, 4, 0), false},
		{utc(2026, 9, 14, 6, 0), true},
		{utc(2026, 9, 14, 9, 59), true},
		{utc(2026, 9, 14, 10, 0), false},
		{utc(2026, 9, 15, 5, 0), false},
		{utc(2026, 9, 12, 8, 0), false},
		{utc(2026, 9, 13, 23, 59), false},
	}
	for _, c := range cases {
		if got := IsPeak(c.in); got != c.want {
			t.Errorf("IsPeak(%v) = %v, want %v", c.in, got, c.want)
		}
	}
}

func TestRoute(t *testing.T) {
	peak := utc(2026, 9, 14, 7, 0)
	off := utc(2026, 9, 14, 10, 0)

	if d := Route("atomic", peak); d.Target != "local" || d.Egress {
		t.Errorf("atomic at peak routed to %v (want local, no egress)", d.Target)
	}
	if d := Route("coach", peak); d.Target != "deepseek-flash" || !d.Egress {
		t.Errorf("coach at peak routed to %v (want deepseek-flash, egress)", d.Target)
	}
	if d := Route("teacher", peak); d.Target != "deepseek-flash" {
		t.Errorf("teacher at peak routed to %v (pro must be grounded)", d.Target)
	}
	if d := Route("teacher", off); d.Target != "deepseek-v4-pro" {
		t.Errorf("teacher off-peak routed to %v (want deepseek-v4-pro)", d.Target)
	}
	if d := Route("nonsense", off); d.Target != "local" {
		t.Errorf("unknown kind routed to %v (want local default)", d.Target)
	}
}
