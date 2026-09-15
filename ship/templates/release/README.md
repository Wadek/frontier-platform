# Release templates

Generic deploy/rollback for apps that use Frontier + Compose.

1. Copy these scripts into the app `scripts/` (or call them with `-Config`).
2. Fill `env-table.example.ps1` for your envs.
3. Set `RUNTIME=docker` or `RUNTIME=podman`.
4. Prod gate: `frontier release-check` (not `apply`).
5. CD: GitHub Actions self-hosted runner or manual deploy — no desktop Task Scheduler.
