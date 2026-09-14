// frontier — standalone CLI (same engine as `git frontier …`).
//
// Examples:
//
//	frontier guard
//	frontier plan
//	frontier apply
//	frontier S
//
// Not `go frontier` — that is not how the Go toolchain works.
// For development you may run: go run ./cmd/frontier guard
package main

import (
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
	fronticli.Run(os.Args[1:])
}
