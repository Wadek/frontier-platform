# frontier-platform

A local, fail-closed pre-push gate. AI-authored code does not reach a remote
without a machine-checked exam and a human merge path.

Nothing is pushed until `plan` and `apply` both exit 0. `plan` refuses `main`,
`master`, and this repository's default branch, refuses a dirty tree, and refuses
a changeset with High/Critical OWASP findings. `apply` seals evidence into a
hash-chained ledger outside the work tree. The `pre-push` hook re-runs both, so
`--no-verify` is the only way around it and that is a deliberate act.

## Layout

| Path | Binary | Role |
|---|---|---|
| `ship/` | `frontier` | The gate: `plan` / `apply` / `release-check`, Guard (OWASP v0), hygiene client, ledger, monitor, optimize |
| `ship/cmd/frontier-git/` | `frontier-git` | The `git frontier …` interface |
| `control/` | `control` | Coordinator: provider pricing, workspace roster, handoff cards |
| `ship/english/` | — | The policy, in English. This is the source of truth for intent |
| `.github/workflows/` | — | CI (build, vet, test, cross-build, public-tree policy) |

## One language

The platform is **Go**. There is one module (`go.mod` at the root) and one
toolchain.

English states the policy; **Go tests are the machine-checked witness** that the
implementation actually obeys it. That pairing is the whole point: a rule that
exists only in prose drifts, and a rule implemented several times drifts faster.
This repository previously carried the same invariants as Go *and* Haskell *and*
Python. That was quadratic maintenance for no additional proof, and the Haskell
layer could not be built on the machine it was written on. It is gone.

Go was chosen over TypeScript because the gate is a **git hook**: it runs on
every commit, on a contributor's machine, and must be fast with no runtime to
install. `frontier` is a single static binary.

## Build

Requires Go 1.22 or newer.

```sh
go build ./...
go test ./...
go vet ./...
```

Install the binaries where the hooks expect them (`$FRONTIER_RUNTIME/bin`,
defaulting to `./runtime/bin`):

```sh
go build -o runtime/bin/frontier        ./ship/cmd/frontier
go build -o runtime/bin/frontier-git    ./ship/cmd/frontier-git
go build -o runtime/bin/control         ./control/cmd/control
```

On Windows add the `.exe` suffix. `runtime/` is gitignored.

## Use

```sh
git checkout -b frontier/topic
# edit
git add -A && git commit -m "msg"

frontier plan            # must exit 0   (or: frontier plan --json)
frontier apply           # must exit 0
git push -u origin HEAD
gh pr create --fill
```

Production deploys are authorized with `frontier release-check`, then the app's
own deploy script. Process and branch model: `ship/english/RELEASE_PROCESS.md`.
The permanent pipeline: `PIPELINE.md`.

## Status, stated plainly

**Implemented and used:** the gate (`plan`, `apply`, `release-check`), Guard
OWASP v0, the ledger, the hooks, `learn`, `guard`, `optimize`, `monitor`,
`runtime`, `skills`/`agents`.

**Reserved but not implemented — these exit non-zero rather than pretend to
run:** most of `control`'s surface (`taxonomy`, `data`, `workflow`, `ai`, `ops`,
`providers`, `prove`), and `frontier slim`.

`control` currently implements: `pricing show|verify`, `roster`,
`transfer validate`, `check`.

## Deliberately not in this tree

| Thing | Where it lives | Why |
|---|---|---|
| Haskell proof layer | removed | Unbuildable on the authoring host, never executed in CI, and duplicated rules Go already enforced with tests |
| `watermarks-remover` service and its agent skill | upstream [`guillaumemeyer/watermarks-remover`](https://github.com/guillaumemeyer/watermarks-remover) | It was vendored here byte-for-byte. Hygiene (H) is advisory-only, and frontier needs nothing from it beyond the Go client in `ship/internal/hygiene/` |

## License

MIT — see [LICENSE](LICENSE).
