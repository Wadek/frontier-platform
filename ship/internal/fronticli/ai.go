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
  frontier ai route                 probe ollama → habitat(mock) → deepseek
  frontier ai complete [text]       complete via cascade; record tokens
  frontier ai complete -f PATH      complete from file contents
  frontier ai tokens                print PR token-report table from .frontier/tokens.jsonl
  frontier ai secrets status        where the DeepSeek key would come from (no value printed)
  frontier ai secrets set deepseek  store key in local keystore (paste on stdin; do not commit)
  frontier ai secrets delete deepseek

Key resolution (same on laptop or Actions runner):
  1) process env DEEPSEEK_API_KEY  ← GitHub Actions secrets land here when a workflow sets them
  2) local keystore                ← Windows DPAPI file / Unix ~/.config/frontier/secrets

DeepSeek does not provide a key — you create one at platform.deepseek.com and store it once.
GitHub:  gh secret set DEEPSEEK_API_KEY
Local:   frontier ai secrets set deepseek

Optional env: DEEPSEEK_MODEL  OLLAMA_HOST  OLLAMA_MODEL  FRONTIER_HABITAT_BASE  FRONTIER_HABITAT_MOCK`)
		return
	}
	switch strings.ToLower(args[0]) {
	case "route", "probe":
		runAIRoute(cwd)
	case "complete", "ask", "chat":
		runAIComplete(cwd, args[1:])
	case "tokens", "report", "ledger":
		runAITokens(cwd)
	case "secrets", "secret", "keys", "key":
		runAISecrets(args[1:])
	default:
		fmt.Fprintf(os.Stderr, "unknown ai subcommand %q\n", args[0])
		os.Exit(2)
	}
}

func runAISecrets(args []string) {
	if len(args) == 0 {
		fmt.Println(`secrets:
  frontier ai secrets status
  frontier ai secrets set deepseek     # read key from stdin (one line)
  frontier ai secrets delete deepseek`)
		return
	}
	switch strings.ToLower(args[0]) {
	case "status", "st", "where":
		st := frontierai.StatusDeepSeek()
		fmt.Printf("resolved_from=%s\n", st.ResolvedFrom)
		fmt.Printf("env_set=%v\n", st.EnvSet)
		fmt.Printf("keystore_set=%v\n", st.KeystoreSet)
		fmt.Printf("keystore_path=%s\n", st.KeystorePath)
		fmt.Printf("in_github_actions=%v\n", st.InActions)
		fmt.Printf("hint=%s\n", st.Hint)
	case "set", "store", "put":
		name := DeepSeekName(args[1:])
		fmt.Fprintln(os.Stderr, "Paste DeepSeek API key, then Enter (input is not echoed to the ledger):")
		key, err := readSecretLine()
		if err != nil {
			fail(err)
			return
		}
		if err := frontierai.SetDeepSeekKey(key); err != nil {
			fail(err)
			return
		}
		// Prefer keystore over a stale User env for this process.
		_ = os.Unsetenv(frontierai.DeepSeekSecretName)
		st := frontierai.StatusDeepSeek()
		fmt.Printf("stored name=%s source=keystore path=%s\n", name, st.KeystorePath)
		fmt.Println("Next: frontier ai secrets status   then   frontier ai route")
		fmt.Println("For GitHub Actions runners: gh secret set DEEPSEEK_API_KEY")
	case "delete", "rm", "remove":
		_ = DeepSeekName(args[1:])
		if err := frontierai.DeleteDeepSeekKey(); err != nil {
			fail(err)
			return
		}
		fmt.Println("deleted local keystore entry for DEEPSEEK_API_KEY")
	default:
		fmt.Fprintf(os.Stderr, "unknown secrets subcommand %q\n", args[0])
		os.Exit(2)
	}
}

func DeepSeekName(args []string) string {
	if len(args) == 0 {
		return frontierai.DeepSeekSecretName
	}
	switch strings.ToLower(args[0]) {
	case "deepseek", "ds", frontierai.DeepSeekSecretName, "api", "key":
		return frontierai.DeepSeekSecretName
	default:
		return frontierai.DeepSeekSecretName
	}
}

func readSecretLine() (string, error) {
	// Read from stdin so `echo key | frontier ai secrets set` and interactive paste both work.
	// Do not print the value.
	buf := make([]byte, 0, 128)
	tmp := make([]byte, 1)
	for {
		n, err := os.Stdin.Read(tmp)
		if n > 0 {
			if tmp[0] == '\n' || tmp[0] == '\r' {
				break
			}
			buf = append(buf, tmp[0])
		}
		if err != nil {
			if len(buf) > 0 {
				break
			}
			return "", err
		}
	}
	key := strings.TrimSpace(string(buf))
	if key == "" {
		return "", fmt.Errorf("empty key — paste the DeepSeek API key from https://platform.deepseek.com")
	}
	return key, nil
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
