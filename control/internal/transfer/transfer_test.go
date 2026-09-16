package transfer

import (
	"encoding/json"
	"testing"
)

const validCard = `{"specversion":"1.0","type":"handoff.request",` +
	`"source":".agent_beta","id":"h1","data":{"task":"x"}}`

func TestParseValid(t *testing.T) {
	c, err := Parse([]byte(validCard))
	if err != nil {
		t.Fatalf("valid card rejected: %v", err)
	}
	if !c.IsHandoff() {
		t.Error("handoff.request not recognised as a handoff")
	}
	if c.Source != ".agent_beta" {
		t.Errorf("source = %q", c.Source)
	}
}

func TestParseEmpty(t *testing.T) {
	if _, err := Parse([]byte(`{}`)); err == nil {
		t.Error("empty object accepted")
	}
}

func TestParseWrongSpecVersion(t *testing.T) {
	bad := `{"specversion":"0.9","type":"x","source":"s","id":"i","data":{}}`
	if _, err := Parse([]byte(bad)); err == nil {
		t.Error("specversion 0.9 accepted")
	}
}

func TestParseMissingKey(t *testing.T) {
	for _, missing := range requiredKeys {
		card := map[string]any{
			"specversion": "1.0", "type": "handoff.request",
			"source": "s", "id": "i", "data": map[string]any{},
		}
		delete(card, missing)
		raw := mustJSON(t, card)
		if _, err := Parse(raw); err == nil {
			t.Errorf("card missing %q accepted", missing)
		}
	}
}

func TestParseNotAnObject(t *testing.T) {
	if _, err := Parse([]byte(`["nope"]`)); err == nil {
		t.Error("array accepted as a card")
	}
}

func TestParseEmptySource(t *testing.T) {
	bad := `{"specversion":"1.0","type":"handoff.request","source":"","id":"i","data":{}}`
	if _, err := Parse([]byte(bad)); err == nil {
		t.Error("empty source accepted")
	}
}

func TestParseNullData(t *testing.T) {
	bad := `{"specversion":"1.0","type":"handoff.request","source":"s","id":"i","data":null}`
	if _, err := Parse([]byte(bad)); err == nil {
		t.Error("null data accepted")
	}
}

func mustJSON(t *testing.T, v any) []byte {
	t.Helper()
	raw, err := json.Marshal(v)
	if err != nil {
		t.Fatal(err)
	}
	return raw
}
