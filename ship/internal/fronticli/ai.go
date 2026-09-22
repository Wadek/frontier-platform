package fronticli

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/Wadek/frontier-platform/ship/internal/frontierai"
	"github.com/Wadek/frontier-platform/ship/internal/ledger"
)

func runAI(cwd string, args []string) {
	if len(args) == 0 {
		fmt.Println(`ai commands (local-first cascade → DeepSeek):
  frontier ai route              probe ollama → habitat(mock) → deepseek
  frontier ai complete [text]    complete via cascade; record tokens
  frontier ai complete -f PATH   complete from file contents
  frontier ai tokens             print PR token-report table from .frontier/tokens.jsonl

Env: DEEPSEEK_API_KEY  DEEPSEEK_MODEL  OLLAMA_HOST  OLLAMA_MODEL
     FRONTIER_HABITAT_BASE  FRONTIER_HABITAT_MOCK=1 (default)`)
		return
	}
	switch strings.ToLower(args[0]) {
	case "route", "probe":
		runAIRoute(cwd)
	case "complete", "ask", "chat":
		runAIComplete(cwd, args[1:])
	case "tokens", "report", "ledger":
		runAITokens(cwd)
	default:
		fmt.Fprintf(os.Stderr, "unknown ai subcommand %q\n", args[0])
		os.Exit(2)
	}
}

func runAIRoute(cwd string) {
	c := frontierai.NewFromEnv()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	route, detail := c.Probe(ctx)
	fmt.Printf("route=%s\n", route)
	fmt.Printf("detail=%s\n", detail)
	led, err := frontierai.OpenLedger(cwd)
	if err == nil {
		_ = led.Append(frontierai.Entry{Route: frontierai.RouteProcess, TotalTokens: 0, Purpose: "ai.route", Model: "probe"})
	}
	axiom("F0", "ai.route", route)
}

func runAIComplete(cwd string, args []string) {
	user := ""
	if len(args) >= 2 && (args[0] == "-f" || args[0] == "--file") {
		b, err := os.ReadFile(args[1])
		if err != nil {
			fail(err)
			return
		}
		user = string(b)
	} else {
		user = strings.TrimSpace(strings.Join(args, " "))
	}
	if user == "" {
		fail(fmt.Errorf("usage: frontier ai complete [-f PATH] <prompt>"))
		return
	}
	c := frontierai.NewFromEnv()
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Minute)
	defer cancel()
	sys := "You are Frontier AI. Prefer short, actionable answers. Do not invent secrets."
	res, err := c.Complete(ctx, sys, user)
	if err != nil {
		fail(err)
		return
	}
	fmt.Printf("route=%s model=%s tokens=%d (prompt=%d completion=%d)\n\n",
		res.Route, res.Model, res.Usage.TotalTokens, res.Usage.PromptTokens, res.Usage.CompletionTokens)
	fmt.Println(res.Content)

	tled, err := frontierai.OpenLedger(cwd)
	if err != nil {
		fail(err)
		return
	}
	if err := tled.Record(res, "ai.complete"); err != nil {
		fail(err)
		return
	}
	led, err := ledger.Open(findLedger(cwd))
	if err == nil {
		_, _ = led.Append("frontier-git", "ai.completed", map[string]any{
			"route":             res.Route,
			"model":             res.Model,
			"prompt_tokens":     res.Usage.PromptTokens,
			"completion_tokens": res.Usage.CompletionTokens,
			"total_tokens":      res.Usage.TotalTokens,
			"token_ledger":      tled.Path(),
		})
	}
	axiom("F0", "ai.completed", res.Route)
}

func runAITokens(cwd string) {
	tled, err := frontierai.OpenLedger(cwd)
	if err != nil {
		fail(err)
		return
	}
	sum, err := tled.Summarize()
	if err != nil {
		fail(err)
		return
	}
	fmt.Printf("token_ledger=%s entries=%d total=%d\n", tled.Path(), sum.Entries, sum.Total)
	fmt.Print(sum.PRTable())
}
