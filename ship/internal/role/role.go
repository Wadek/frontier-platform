package role

import "fmt"

// Role is a frontier ladder rung. Lower = less blast radius.
type Role int

const (
	Observer Role = iota
	Analyst
	Operator
	Executor
)

func Parse(s string) (Role, error) {
	switch s {
	case "observer", "Observer", "OBSERVER":
		return Observer, nil
	case "analyst", "Analyst", "ANALYST":
		return Analyst, nil
	case "operator", "Operator", "OPERATOR":
		return Operator, nil
	case "executor", "Executor", "EXECUTOR":
		return Executor, nil
	default:
		return Observer, fmt.Errorf("unknown role %q (use observer|analyst|operator|executor)", s)
	}
}

func (r Role) String() string {
	switch r {
	case Observer:
		return "observer"
	case Analyst:
		return "analyst"
	case Operator:
		return "operator"
	case Executor:
		return "executor"
	default:
		return "observer"
	}
}

func (r Role) Can(min Role) bool { return r >= min }

// Elevate steps up exactly one rung.
func (r Role) Elevate() (Role, error) {
	if r >= Executor {
		return r, fmt.Errorf("already at executor")
	}
	return r + 1, nil
}
