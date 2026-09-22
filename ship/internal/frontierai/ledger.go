package frontierai

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"
)

// Entry is one append-only token spend row (local file under .frontier/).
type Entry struct {
	TS               string `json:"ts"`
	Route            string `json:"route"`
	Model            string `json:"model"`
	PromptTokens     int    `json:"prompt_tokens"`
	CompletionTokens int    `json:"completion_tokens"`
	TotalTokens      int    `json:"total_tokens"`
	Purpose          string `json:"purpose,omitempty"`
}

// Ledger is a JSONL token spend log for PR reports.
type Ledger struct {
	path string
	mu   sync.Mutex
}

// OpenLedger opens (or creates) root/.frontier/tokens.jsonl.
func OpenLedger(root string) (*Ledger, error) {
	dir := filepath.Join(root, ".frontier")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return nil, err
	}
	path := filepath.Join(dir, "tokens.jsonl")
	f, err := os.OpenFile(path, os.O_CREATE, 0o644)
	if err != nil {
		return nil, err
	}
	_ = f.Close()
	return &Ledger{path: path}, nil
}

func (l *Ledger) Path() string {
	if l == nil {
		return ""
	}
	return l.path
}

// Record appends one spend row from a Result.
func (l *Ledger) Record(res *Result, purpose string) error {
	if l == nil || res == nil {
		return fmt.Errorf("frontierai: nil ledger or result")
	}
	e := Entry{
		TS:               time.Now().UTC().Format(time.RFC3339),
		Route:            res.Route,
		Model:            res.Model,
		PromptTokens:     res.Usage.PromptTokens,
		CompletionTokens: res.Usage.CompletionTokens,
		TotalTokens:      res.Usage.TotalTokens,
		Purpose:          purpose,
	}
	return l.Append(e)
}

// Append writes one JSON line.
func (l *Ledger) Append(e Entry) error {
	l.mu.Lock()
	defer l.mu.Unlock()
	if e.TS == "" {
		e.TS = time.Now().UTC().Format(time.RFC3339)
	}
	b, err := json.Marshal(e)
	if err != nil {
		return err
	}
	f, err := os.OpenFile(l.path, os.O_APPEND|os.O_WRONLY, 0o644)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = f.Write(append(b, '\n'))
	return err
}

// Summary aggregates totals by route for PR token reports.
type Summary struct {
	ByRoute map[string]int `json:"by_route"`
	Total   int            `json:"total"`
	Entries int            `json:"entries"`
}

// Summarize reads the ledger and sums total_tokens by route.
func (l *Ledger) Summarize() (*Summary, error) {
	l.mu.Lock()
	defer l.mu.Unlock()
	raw, err := os.ReadFile(l.path)
	if err != nil {
		return nil, err
	}
	sum := &Summary{ByRoute: map[string]int{}}
	for _, line := range splitLines(string(raw)) {
		if line == "" {
			continue
		}
		var e Entry
		if err := json.Unmarshal([]byte(line), &e); err != nil {
			continue
		}
		sum.ByRoute[e.Route] += e.TotalTokens
		sum.Total += e.TotalTokens
		sum.Entries++
	}
	return sum, nil
}

// PRTable renders a markdown token report table (local-first cascade rows).
func (s *Summary) PRTable() string {
	if s == nil {
		s = &Summary{ByRoute: map[string]int{}}
	}
	get := func(k string) int {
		if s.ByRoute == nil {
			return 0
		}
		return s.ByRoute[k]
	}
	return fmt.Sprintf(`## Token report

| Route | Tokens |
|-------|--------|
| Local tools / process | %d |
| Laptop Ollama | %d |
| Habitat Qwen (mock until live) | %d |
| Frontier AI (DeepSeek) | %d |
| **Total recorded** | **%d** |
`,
		get(RouteProcess),
		get(RouteOllama),
		get(RouteHabitat),
		get(RouteDeepSeek),
		s.Total,
	)
}

func splitLines(s string) []string {
	var out []string
	start := 0
	for i := 0; i < len(s); i++ {
		if s[i] == '\n' {
			out = append(out, s[start:i])
			start = i + 1
		}
	}
	if start < len(s) {
		out = append(out, s[start:])
	}
	return out
}
