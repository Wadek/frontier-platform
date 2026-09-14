# frontier-fleet — HANDOFF (session 2026-09-14, ~95% context)

Read this first in the next session. The 95% rule (T9) applies: this session soft-refuses
further runs and hands off here.

## State as of handoff

- **Repo**: `Wadek/frontier-fleet` (private). `main` = merged PR #1 (monorepo, tower scaffold).
  Local `main` tracks `origin/main` (needed for the lefthook pre-push gate).
- **Open fleet PRs** (gate-sealed, awaiting Wade's merge):
  - `#2 frontier/finalize` — hygiene service re-home + daily weather-check script.
  - `#3 frontier/session-budget` — the 95% rule (T9): Python/Go/Haskell triple + doctrine.
- **Improvement PRs** (single-vendor migration, all gate-sealed):
  - `Wadek/wakagym#14` → dev · `Wadek/Farm#31` → main · `Wadek/polly-engine#2` → main ·
    `Wadek/food#2` → main.
- **Decommissioned** (containers + gateway routes removed, gateway reloaded, nginx -t OK):
  caravans, human-atlas, pigeon, phylax, watermarks-remover. Data preserved on disk;
  salvage briefs in each `.agent_<plane>\brief.md`. waka-net nginx change committed locally
  on `feat/decommission-routes` (no remote); its frontier seal was blocked by pre-existing
  key files in that tree (`phylax\rsa_key.pem`, `pihole\tls.pem`) — by design, not a bug.
- **Daily job**: `WakaWeatherCheck` (11:05 local) runs
  `frontier-fleet\control\scripts\weather_check.py` → logs `runtime\weather-check.log`.

## Wade's action list

1. Merge fleet PRs #2 and #3; then the four improvement PRs (wakagym → dev first).
2. opengym `.env`: add `DEEPSEEK_COACH_BUDGET_USD`, delete `XAI_*`; satokori: delete
   `ANTHROPIC_API_KEY`.
3. Cloudflare dashboard: remove `pigeon.wakalabs.net` + `vault.wakalabs.net` tunnel entries.
4. Back up the Vaultwarden data (`D:\wakalabs\waka-net\phylax`) before any future deletion.
5. Optionally start the re-homed hygiene service (`frontier-fleet\hygiene\README.md`) —
   advisory-only while down.

## Next work (new session)

1. **DSH 95% hook** — the harness has pressure-aware compaction (`packages/compaction/
   compaction-basic`); wire `session_budget(95%)` as a soft-refusal surface there (the
   fleet-side policy is done; this makes it automatic for this GUI).
2. `ship-scan-gitignore-aware` improvement (the OWASP walker scans ignored trees — the
   waka-net key-file block is the concrete case).
3. Optional: opengym `coach-verify.mjs` scheduled via the job service; food catalog fallback
   polish; retire-tree cleanup when ready.
