# Customer — Satokori (farm network grocer)

**Local:** `D:\wakalabs\satokori`  
**Remote:** `github.com/Wadek/Farm` (rebranded in-tree as Satokori)  
**Public:** `https://satokori.wakalabs.net`  
**Loopback:** `http://127.0.0.1:8791`  
**Date:** 2026-08-21  
**Playbooks:** [stacks/python-web-static.md](stacks/python-web-static.md), [stacks/habitat-edge.md](stacks/habitat-edge.md)

## Story

Five prior attempts (kariniemi-farms, perinnepelto, mycelium-network, Farm, Offerit) lived on C: and Firebase/PocketBase. Habitat rule: **one tree on D:**, Docker on `waka-net`, chalkboard not a shop.

Value is the **farmer network**. Organizer field-visits onboard a farm in one form. Customers browse who has what / when / where. Cash at the gate decrements stock; no euro sales row.

## What we ran

```text
cd D:\wakalabs\satokori
frontier learn classify
frontier guard
git checkout -b frontier/satokori-launch
frontier plan
frontier apply
git push -u origin HEAD
docker compose up -d --build
```

## Axioms observed

| When | What |
|------|------|
| Learn | `app_compose` (compose + `data/` + Python) |
| 1033 | Tunnel down because Docker Desktop was off — not an app bug |
| Guard | Do not bake `SECRET_KEY` into the image |
| Edge | Public hostname, no GitHub oauth2-proxy |
| Push | Fail-closed plan → apply → push on a `frontier/…` branch, not `main` |

## Demo logins (local seed only)

| Role | Email | Password |
|------|-------|----------|
| Organizer | `wade@kariniemi.farm` | `farmgate` |
| Farmer | `maija@naapuri.fi` | `farmgate` |
| Customer | `anna@hyvinkaa.fi` | `market` |
