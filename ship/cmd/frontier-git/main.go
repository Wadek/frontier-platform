// frontier-git — drop-in `git` shim. High-blast verbs are gated;
// `git frontier …` delegates to the same CLI as the standalone `frontier` binary.
package main

import (
	"fmt"
	"os"

	"github.com/Wadek/frontier-ship/internal/fronticli"
)

// Set by release ldflags.
var (
	version = "dev"
	commit  = "none"
)

func main() {
	fronticli.Version = version
	fronticli.Commit = commit

	args := os.Args[1:]
	realGit := fronticli.FindRealGit()

	if len(args) == 0 {
		os.Exit(fronticli.RunGitPassthrough(realGit, args))
	}

	verb := args[0]
	soft := os.Getenv("FRONTIER_SOFT") == "1"
	strict := os.Getenv("FRONTIER_STRICT") == "1"

	switch verb {
	case "push":
		if err := fronticli.GuardPush(soft); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	case "commit":
		if err := fronticli.GuardCommit(args, soft, strict); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	case "frontier":
		fronticli.Run(args[1:])
		return
	}

	os.Exit(fronticli.RunGitPassthrough(realGit, args))
}
