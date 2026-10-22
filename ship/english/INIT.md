# How to start Frontier

Read this in order. Simple English only here.

## What you are installing

1. **Laws** (English) — what is allowed
2. **`git` + `frontier`** (Go) — how you talk to the machine every day
3. **Go tests** — the witness that the laws and the code agree (F5)

You do **not** install a new operating system.

---

## Step 0 — Tools

- Git (real engine on your PATH)
- Go 1.22+ (to build the runtime and run the tests)

Linux, macOS, and Windows are all fine. The examples below show both shells.

---

## Step 1 — Get the code

```sh
git clone https://github.com/Wadek/frontier-platform.git
cd frontier-platform
```

If a Frontier `git` shim is already ahead of the real engine on your PATH and
blocks the clone, call the real `git` by absolute path once (Windows example:
`"C:\Program Files\Git\cmd\git.exe" clone …`).

---

## Step 2 — Build the Go runtime

From the repository root:

```sh
mkdir -p runtime/bin
go build -o runtime/bin/frontier     ./ship/cmd/frontier
go build -o runtime/bin/frontier-git ./ship/cmd/frontier-git
go build -o runtime/bin/control      ./control/cmd/control
```

On Windows, add the `.exe` suffix to each `-o` path, or run the helper:

```powershell
powershell -File .\ship\scripts\install-git-interface.ps1
```

Put `runtime/bin` on your `PATH`, or set `FRONTIER_RUNTIME` to this `runtime`
directory. Open a **new** terminal, then check:

```sh
frontier version
git frontier explain   # if you installed the git shim
go test ./...
```

---

## Step 3 — Learn mode (first week)

Warn instead of blocking while you learn the flow:

```sh
export FRONTIER_SOFT=1          # bash / zsh
# $env:FRONTIER_SOFT = "1"      # PowerShell
```

Turn it off when ready (`FRONTIER_SOFT=0` / unset). Soft mode is for learning
only — never for a real ship.

---

## Step 4 — First safe drill (Terraform-like)

```sh
cd <some-repo>
git checkout -b feat/first
# edit a file
git add -A
git commit -m "frontier: first sealed change"

frontier learn         # classify project (L)
frontier guard         # security exam (G)
frontier plan          # must succeed — fail closed
frontier apply         # only works after plan.passed
git push -u origin HEAD
```

(`git frontier …` is the same engine as bare `frontier …`.)

If `plan` fails, fix the reasons; do not expect apply/push to work.
The ledger is state (like Terraform state).

From an **application** tree the first time you intend to deploy:

```sh
frontier ready
```

That command only reports gaps (tests, CI, deploy script, …). It does not
edit the app. See [SETUP.md](../../SETUP.md) and [RELEASE_PROCESS.md](RELEASE_PROCESS.md).

---

## Step 5 — Agents use the same CLI

Coding agents run `frontier` and `git` like a human. There is no separate host protocol.

Corporate laptop with no install: paste `teach/CORPORATE_AI_TRIAL_PROMPT.md`.

---

## Where to read next

| File | Why |
|------|-----|
| `english/LANGUAGE.md` | Why English + Go |
| `english/AXIOMS.md` | The laws in English |
| `teach/LEARN_GIT.md` | Practice with `git` |
| `internal/` | The laws as Go code + `_test.go` witness |
| `english/WHY_NOT_A_CUSTOM_OS.md` | Why we do not write an OS |
