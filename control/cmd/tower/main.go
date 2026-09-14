// tower — the frontier-fleet control tower (frontier-control · Milkcow).
package main

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/Wadek/frontier-fleet/control/internal/flightplan"
	"github.com/Wadek/frontier-fleet/control/internal/weather"
)

func usage() {
	fmt.Print(`tower — the fleet control tower (frontier-control · Milkcow)

  tower taxonomy init|diff|add|approve|retire
  tower data generate --plane <p>        (teacher pass, off-peak)
  tower workflow run|status|review|list
  tower roster                           (scan D:\wakalabs for planes + pilots)
  tower ai audit|improve <plane>
  tower ops health|docker|ip|decommission|jobs
  tower runways list|test
  tower weather show                     (window + current status)
  tower monitor
  tower prove
`)
}

func main() {
	if len(os.Args) < 2 {
		usage()
		os.Exit(2)
	}
	switch os.Args[1] {
	case "weather":
		now := time.Now().UTC()
		d := weather.Route("teacher", now)
		fmt.Printf("utc=%s peak=%v off-peak runway=%s (%s)\n",
			now.Format("2006-01-02T15:04Z"), weather.IsPeak(now), d.Target, d.Reason)
	case "roster":
		root := `D:\wakalabs`
		entries, err := os.ReadDir(root)
		if err != nil {
			fmt.Fprintln(os.Stderr, "tower roster:", err)
			os.Exit(1)
		}
		exclude := map[string]bool{"frontier-fleet": true, "deepseek_harness": true,
			"ollama-data": true, "grok skills": true}
		for _, e := range entries {
			if !e.IsDir() || exclude[e.Name()] || e.Name()[0] == '.' {
				continue
			}
			pilot := filepath.Join(root, ".agent_"+e.Name())
			if _, err := os.Stat(pilot); err == nil {
				fmt.Printf("%-24s pilot: .agent_%s\n", e.Name(), e.Name())
			} else {
				fmt.Printf("%-24s (no pilot yet)\n", e.Name())
			}
		}
	case "check":
		// cheap smoke: prove the FSM holds and the pro runway is gated at peak
		if !flightplan.NextOk("landed", "filed") {
			fmt.Println("FAIL: debrief gate missing")
			os.Exit(1)
		}
		peak := time.Date(2026, 9, 14, 7, 0, 0, 0, time.UTC) // Mon peak
		if weather.Route("teacher", peak).Target != "deepseek-flash" {
			fmt.Println("FAIL: pro flew at peak")
			os.Exit(1)
		}
		fmt.Println("tower check: FSM + weather gates hold")
	default:
		fmt.Printf("tower: %s not implemented yet (see english/CONTROL_TOWER.md)\n", os.Args[1])
		usage()
		os.Exit(1)
	}
}
