# Learn git *through* Frontier (userland first)

You do **not** need to rebuild the Linux kernel to get value.  
First we wrap **git** (the tool everyone already pretends to know). Kernels / “linux frontier” come later.

## What you have locally

| Path | What |
|------|------|
| `D:\frontier\src\git` | Upstream **git/git** source (study copy, shallow clone) |
| `frontier-git` | Our shim: real git + F0–F4 guards on push / main commits |
| System git | `C:\Program Files\Git\cmd\git.exe` (engine under the hood) |

Reading `builtin/push.c` in the source tree teaches how push works.  
Using `frontier-git` teaches **when you are allowed to push**.

## Install — the interface **is** `git`

Frontier installs a shim named `git.exe` ahead of system git on your PATH.  
You never need a special command name. Same muscle memory as everyone else.

```powershell
$env:Path = "C:\Users\waka\sdk\go\bin;C:\Program Files\Git\cmd;" + $env:Path
cd C:\Users\waka\src\frontier-ship
go build -o D:\frontier\bin\frontier-git.exe ./cmd/frontier-git
Copy-Item -Force D:\frontier\bin\frontier-git.exe D:\frontier\bin\git.exe

# real engine
[Environment]::SetEnvironmentVariable("FRONTIER_GIT_BIN", "C:\Program Files\Git\cmd\git.exe", "User")
# put D:\frontier\bin first on User PATH (one-time)
# then open a NEW terminal

$env:FRONTIER_SOFT = "1"   # learn mode; remove when ready
git frontier explain
git --version
Get-Command git   # should show D:\frontier\bin\git.exe
```

## Git basics as a Frontier ladder

The interface is ordinary git. Frontier adds seals on the dangerous verbs.

```
  git status              ≈  observer
  git diff                ≈  analyst
  git commit              ≈  operator   (blocked on main/master)
  git push                ≈  executor   (needs: git frontier gate)
  git frontier status|gate|ledger|explain
```

### Drill 1 — observe

```powershell
cd <your-repo>
git status
git frontier status
```

### Drill 2 — branch (never commit on main)

```powershell
git checkout -b frontier/learn-1
# edit a file
git add -A
git commit -m "frontier: learning commit"
```

### Drill 3 — gate then push

```powershell
git frontier gate          # seals gate.passed or fails
$env:FRONTIER_SOFT = "0"   # real deny mode
git push -u origin HEAD
```

If gate fails, read the reasons (dirty tree, on main, no seal). Fix; re-gate; push.

### Drill 4 — peek at upstream source (optional)

```powershell
cd D:\frontier\src\git
notepad builtin\push.c
# or
findstr /n "push" Documentation\git-push.txt
```

You are learning **two** gits: the C program, and the Frontier *policy* around the same `git` command.

## Roadmap (so kernels don’t scare you)

```
  NOW     frontier-git wrapper + study clone of git/git
  NEXT    logic/ Datalog core + Python/Rust witnesses (F5)
  LATER   WSL build of your own git binary from source
  MUCH    "linux frontier" — kernel/policy deep dive only after
          userland habits (status→commit→push+gate) are boring
```

**Simple is better:** master the shim before compiling kernels.
