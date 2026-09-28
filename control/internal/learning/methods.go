package learning

import "fmt"

// Method is a platform learning-method ID. Catalog: english/LEARNING.md.
type Method string

const (
	MethodDebrief       Method = "debrief"
	MethodGenerate      Method = "generate"
	MethodRetrieval     Method = "retrieval"
	MethodSpaced        Method = "spaced"
	MethodInterleave    Method = "interleave"
	MethodElaborate     Method = "elaborate"
	MethodFadedExample  Method = "faded_example"
	MethodDeliberate    Method = "deliberate"
	MethodTeachback     Method = "teachback"
	MethodDual          Method = "dual"
	MethodListen        Method = "listen"
	MethodSeed          Method = "seed"
)

var methods = map[Method]bool{
	MethodDebrief:      true,
	MethodGenerate:     true,
	MethodRetrieval:    true,
	MethodSpaced:       true,
	MethodInterleave:   true,
	MethodElaborate:    true,
	MethodFadedExample: true,
	MethodDeliberate:   true,
	MethodTeachback:    true,
	MethodDual:         true,
	MethodListen:       true,
	MethodSeed:         true,
}

// Kind is the work-shape used to pick methods. Domain skills map tasks onto a kind.
type Kind string

const (
	KindSession    Kind = "session"
	KindFact       Kind = "fact"
	KindProcedure  Kind = "procedure"
	KindSimilar    Kind = "similar"
	KindWhy        Kind = "why"
	KindSpeech     Kind = "speech"
	KindTeacher    Kind = "teacher"
)

// ValidateMethod rejects unknown IDs.
func ValidateMethod(id string) error {
	if !methods[Method(id)] {
		return fmt.Errorf("unknown learning method %q", id)
	}
	return nil
}

// Select returns the default method stack for a kind. Always includes debrief
// except when the kind is already the teacher pass (seed is the product; debrief
// still required at session end via KindSession).
func Select(kind Kind) []Method {
	switch kind {
	case KindSession:
		return []Method{MethodDebrief}
	case KindFact:
		return []Method{MethodGenerate, MethodRetrieval, MethodSpaced}
	case KindProcedure:
		return []Method{MethodFadedExample, MethodDeliberate}
	case KindSimilar:
		return []Method{MethodInterleave}
	case KindWhy:
		return []Method{MethodElaborate, MethodTeachback}
	case KindSpeech:
		return []Method{MethodListen, MethodDual}
	case KindTeacher:
		return []Method{MethodSeed}
	default:
		return nil
	}
}
