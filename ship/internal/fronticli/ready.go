package fronticli

import (
	"fmt"
	"os"

	"github.com/Wadek/frontier-platform/ship/internal/ready"
)

func runReady(cwd string, args []string) {
	root := cwd
	asJSON := false
	for _, a := range args {
		if a == "--json" {
			asJSON = true
			continue
		}
		if a == "-h" || a == "--help" || a == "help" {
			fmt.Print(`frontier ready — first-deploy detector (does not mutate the tree)

  frontier ready           text report; exit 1 if any check fails
  frontier ready --json    JSON report; same exit code

Copy templates from ship/templates/release/ when a check says so.
GitHub branch protection is configured in the GitHub UI, not by this command.
`)
			return
		}
		if a != "" && a[0] != '-' {
			root = a
		}
	}
	rep, err := ready.Inspect(root)
	if err != nil {
		fail(err)
	}
	if asJSON {
		b, err := ready.Marshal(rep)
		if err != nil {
			fail(err)
		}
		fmt.Println(string(b))
	} else {
		fmt.Print(ready.RenderText(rep))
	}
	if !rep.Ready {
		os.Exit(1)
	}
}
