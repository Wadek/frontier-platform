# Release templates

Generic deploy/rollback for apps that use Frontier + Compose.

1. Copy these scripts into the app `scripts/` (or call them with `-Config`).
2. Fill `env-table.example.ps1` for your envs.
3. Set `RUNTIME=docker` or `RUNTIME=podman`.
4. Prod: `frontier release-check` then this deploy script. GitHub protects `main`; this is not a push of `main`.
5. CD: GitHub Actions runner on the deploy host, or run the script by hand after merge. No desktop Task Scheduler.
6. Copy `verify.yml.example` to `.github/workflows/verify.yml` and replace the Stage 2 echo with the stack test command.
