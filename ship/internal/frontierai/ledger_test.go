package frontierai

import (
	"path/filepath"
	"strings"
	"testing"
)

func TestLedgerRecordAndPRTable(t *testing.T) {
	dir := t.TempDir()
	led, err := OpenLedger(dir)
	if err != nil {
		t.Fatal(err)
	}
	if err := led.Record(&Result{
		Route: RouteDeepSeek,
		Model: "deepseek-chat",
		Usage: Usage{PromptTokens: 100, CompletionTokens: 20, TotalTokens: 120},
	}, "enhance.guard"); err != nil {
		t.Fatal(err)
	}
	_ = led.Append(Entry{Route: RouteOllama, TotalTokens: 5, Model: "qwen"})
	sum, err := led.Summarize()
	if err != nil {
		t.Fatal(err)
	}
	if sum.Total != 125 || sum.ByRoute[RouteDeepSeek] != 120 || sum.ByRoute[RouteOllama] != 5 {
		t.Fatalf("%+v", sum)
	}
	table := sum.PRTable()
	if !strings.Contains(table, "Frontier AI (DeepSeek)") || !strings.Contains(table, "120") {
		t.Fatalf("table=%s", table)
	}
	if filepath.Base(led.Path()) != "tokens.jsonl" {
		t.Fatalf("path=%s", led.Path())
	}
}
