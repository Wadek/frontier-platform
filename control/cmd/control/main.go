// control — frontier-platform coordinator.
//
// The coordinator routes work, enforces review gates, and records outcomes. It
// never performs the work itself. Only the verbs listed by usage() are
// implemented; everything else exits non-zero rather than pretending to run.
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/Wadek/frontier-platform/control/internal/pricing"
	"github.com/Wadek/frontier-platform/control/internal/registry"
	"github.com/Wadek/frontier-platform/control/internal/transfer"
	"github.com/Wadek/frontier-platform/control/internal/workflow"
)

func usage() {
	fmt.Print(`control — frontier-platform coordinator

  control pricing show|verify              (window, routing, official drift check)
  control roster [ROOT]                    (FRONTIER_WORKSPACE or arg; projects + agents)
  control transfer validate <file|->       (handoff card envelope)
  control check                            (smoke: FSM + pricing gates hold)

Not implemented yet (reserved): taxonomy, data, workflow, ai, ops, providers,
monitor, prove. These exit 1 rather than run.
`)
}

func main() {
	if len(os.Args) < 2 {
		usage()
		os.Exit(2)
	}
	switch os.Args[1] {
	case "pricing":
		sub := ""
		if len(os.Args) >= 3 {
			sub = os.Args[2]
		}
		switch sub {
		case "", "show":
			showPricing()
		case "verify":
			verifyPricing()
		default:
			fmt.Fprintf(os.Stderr, "control pricing: unknown %q (show|verify)\n", sub)
			os.Exit(2)
		}
	case "roster":
		roster()
	case "transfer":
		if len(os.Args) < 4 || os.Args[2] != "validate" {
			fmt.Fprintln(os.Stderr, "usage: control transfer validate <file|->")
			os.Exit(2)
		}
		validateTransfer(os.Args[3])
	case "check":
		check()
	default:
		fmt.Printf("control: %s not implemented yet\n", os.Args[1])
		usage()
		os.Exit(1)
	}
}

func showPricing() {
	now := time.Now().UTC()
	d := pricing.Route("teacher", now)
	fmt.Printf("utc=%s peak=%v provider=%s (%s)\n",
		now.Format("2006-01-02T15:04Z"), pricing.IsPeak(now), d.Target, d.Reason)
	fmt.Printf("peak window (official): %s\n", pricing.ExpectedWindow)
}

func verifyPricing() {
	now := time.Now().UTC()
	window, err := pricing.Check(context.Background(), nil, pricing.OfficialURL)
	line := pricing.LogLine(now, window, err)
	if logErr := pricing.AppendLog(pricing.LogPath(), line); logErr != nil {
		fmt.Fprintln(os.Stderr, "control pricing verify: log:", logErr)
	}
	fmt.Println(line)
	if err != nil {
		fmt.Fprintln(os.Stderr, "control pricing verify:", err)
		os.Exit(1)
	}
}

func roster() {
	root := os.Getenv("FRONTIER_WORKSPACE")
	if len(os.Args) >= 3 && os.Args[2] != "" {
		root = os.Args[2]
	}
	if root == "" {
		fmt.Fprintln(os.Stderr, "control roster: set FRONTIER_WORKSPACE or pass ROOT as first arg")
		os.Exit(1)
	}
	projects, err := registry.Scan(root, nil)
	if err != nil {
		fmt.Fprintln(os.Stderr, "control roster:", err)
		os.Exit(1)
	}
	for _, p := range projects {
		if p.HasAgent {
			fmt.Printf("%-24s agent: %s%s\n", p.Name, registry.AgentPrefix, p.Name)
		} else {
			fmt.Printf("%-24s (no agent yet)\n", p.Name)
		}
	}
}

func validateTransfer(path string) {
	raw, err := readInput(path)
	if err != nil {
		fmt.Fprintln(os.Stderr, "control transfer validate:", err)
		os.Exit(1)
	}
	card, err := transfer.Parse(raw)
	if err != nil {
		fmt.Fprintf(os.Stderr, "control transfer validate: invalid card: %v\n", err)
		os.Exit(1)
	}
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	_ = enc.Encode(card)
}

func readInput(path string) ([]byte, error) {
	if path == "-" {
		return io.ReadAll(os.Stdin)
	}
	return os.ReadFile(path)
}

func check() {
	// Cheap smoke: prove the debrief gate and the peak pricing gate both hold.
	if workflow.NextOk("completed", "closed") {
		fmt.Println("FAIL: debrief gate missing")
		os.Exit(1)
	}
	peak := time.Date(2026, 9, 14, 7, 0, 0, 0, time.UTC) // Monday peak
	if pricing.Route("teacher", peak).Target != "deepseek-flash" {
		fmt.Println("FAIL: pro routed at peak")
		os.Exit(1)
	}
	if transfer.Validate(transfer.Card{}) == nil {
		fmt.Println("FAIL: empty transfer card accepted")
		os.Exit(1)
	}
	fmt.Println("control check: FSM + pricing + transfer gates hold")
}
