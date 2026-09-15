// control — frontier-platform control binary.
package main

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/Wadek/frontier-platform/control/internal/pricing"
	"github.com/Wadek/frontier-platform/control/internal/workflow"
)

func usage() {
	fmt.Print(`control — frontier-platform coordinator

  control taxonomy init|diff|add|approve|retire
  control data generate --project <p>      (teacher pass, off-peak)
  control workflow run|status|review|list
  control roster [ROOT]                    (FRONTIER_WORKSPACE or arg; projects + agents)
  control ai audit|improve <project>
  control ops health|docker|ip|decommission|jobs
  control providers list|test
  control pricing show                     (window + current status)
  control monitor
  control prove
`)
}

func main() {
	if len(os.Args) < 2 {
		usage()
		os.Exit(2)
	}
	switch os.Args[1] {
	case "pricing":
		now := time.Now().UTC()
		d := pricing.Route("teacher", now)
		fmt.Printf("utc=%s peak=%v off-peak provider=%s (%s)\n",
			now.Format("2006-01-02T15:04Z"), pricing.IsPeak(now), d.Target, d.Reason)
	case "roster":
		root := os.Getenv("FRONTIER_WORKSPACE")
		if len(os.Args) >= 3 && os.Args[2] != "" {
			root = os.Args[2]
		}
		if root == "" {
			fmt.Fprintln(os.Stderr, "control roster: set FRONTIER_WORKSPACE or pass ROOT as first arg")
			os.Exit(1)
		}
		entries, err := os.ReadDir(root)
		if err != nil {
			fmt.Fprintln(os.Stderr, "control roster:", err)
			os.Exit(1)
		}
		for _, e := range entries {
			if !e.IsDir() || e.Name()[0] == '.' {
				continue
			}
			agent := filepath.Join(root, ".agent_"+e.Name())
			if _, err := os.Stat(agent); err == nil {
				fmt.Printf("%-24s agent: .agent_%s\n", e.Name(), e.Name())
			} else {
				fmt.Printf("%-24s (no agent yet)\n", e.Name())
			}
		}
	case "check":
		// cheap smoke: prove the FSM holds and the pro provider is gated at peak
		if workflow.NextOk("completed", "closed") {
			fmt.Println("FAIL: debrief gate missing")
			os.Exit(1)
		}
		peak := time.Date(2026, 9, 14, 7, 0, 0, 0, time.UTC) // Mon peak
		if pricing.Route("teacher", peak).Target != "deepseek-flash" {
			fmt.Println("FAIL: pro routed at peak")
			os.Exit(1)
		}
		fmt.Println("control check: FSM + pricing gates hold")
	default:
		fmt.Printf("control: %s not implemented yet\n", os.Args[1])
		usage()
		os.Exit(1)
	}
}
