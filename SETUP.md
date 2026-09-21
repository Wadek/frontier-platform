# Setup (forks and first use)

This repository is a **generic** tool. It does not know your apps, hosts, or vendors.

## 1. Build

Needs Go 1.22+.

```sh
git clone https://github.com/Wadek/frontier-platform.git
cd frontier-platform
go test ./...
go build -o runtime/bin/frontier ./ship/cmd/frontier
go build -o runtime/bin/frontier-git ./ship/cmd/frontier-git
```

On Windows, add `.exe`. Put `runtime/bin` on `PATH`, or set `FRONTIER_RUNTIME` to this `runtime` directory.

Optional: set Git `core.hooksPath` to `$FRONTIER_RUNTIME/hooks` so `git push` re-runs the gate.

## 2. Protect `dev` and `main` on GitHub (you, in the GitHub UI)

Frontier does not replace GitHub branch protection. On **this** repo and on each app repo:

- Require a pull request before merging to `dev` and to `main`
- Require the CI check (`ci` here, `verify` in apps) to pass
- Deny force-push and branch deletion
- Do not allow direct pushes to `dev` or `main`

The laptop hook will also refuse to **commit** on `dev` or `main`. That is a local safety catch, not the source of truth.

## 3. First time you ship **an application**

From the **app** checkout (not this repo):

```sh
frontier learn classify
frontier ready
```

`ready` does not change files. It prints what the standard process still needs (tests, `verify.yml`, deploy script, …) and how to copy templates. Fix the gaps, run `ready` again, then `plan` → `apply` → `git push` on a `feat/*` branch.

Templates: `ship/templates/release/`.

## 4. Secrets

API keys (if you use a model provider at all) belong in **GitHub Actions secrets** or the app’s own environment files that are gitignored. This tool has no vendor plugin and nothing to “turn on” in a coding-agent product.

## 5. Branches

```text
feat/<slug>  --PR-->  dev  --PR-->  main
```

Human merges `dev` → `main`. Then deploy from a checkout of `main` (`frontier release-check`, then the app’s deploy script).
