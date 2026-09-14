# WEATHER — pricing policy for the fleet

Official source (checked 2026-09-14, re-checked at go-time): [DeepSeek Models & Pricing](https://api-docs.deepseek.com/quick_start/pricing)

## The window

- **Peak: Monday–Friday 01:00–04:00 and 06:00–10:00 UTC.**
- **Everything else is off-peak** (including all weekend). Off-peak rates are half of peak.
- `deepseek-v4-pro` API service continues past 2026-09-14 with unchanged billing.

## Prices (per 1M tokens)

| Model | Input (cache miss) | Output |
|---|---|---|
| deepseek-flash | $0.15 off-peak / $0.30 peak | $0.60 / $1.20 |
| deepseek-v4-pro | $0.66 off-peak / $1.32 peak | $1.98 / $3.96 |

Vision: **flash ✓**, v4-pro ✗ — all image work flies on flash.

## Fleet rules

1. Pro runway never at peak (`tower workflow review` → `holding`; job service schedules off-peak).
2. Daily one-shot check of the pricing page (Task Scheduler; alerts on change — DeepSeek
   reserves the right to adjust prices and says to check regularly).
3. Local first: the local GPU flies before any cloud runway.
4. All cloud calls egress-filtered first; refused trees never cross.
5. Spend ledgers per plane (tokens + USD), fed into the data index.

Config: `control/weather.yaml` (window + rates, editable without code).
