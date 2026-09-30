# RFC 001 -- WakaBot authority, learn namespaces, and reserved Google AI provider

Status: draft (docs + stubs only)
Applies to: frontier-platform doctrine and provider config
Non-goals for this change: runtime Google/DeepSeek test calls, Herdr install, secrets

## 1. Authority (public face vs operator)

Outsiders and sibling agents never address Wade/user by name. The control face is **WakaBot**.

```text
admin `waka`
    └─ WakaBot          master on D: (wakalabs habitat); public / control face
         └─ slave bots / agents
              ├─ Grok Bot (tech-lead)     slave -- schematics / triage / approvals only
              ├─ Frontier / waka-cli      slave -- plan -> apply -> push execution
              └─ project agents           slave -- scoped `.agent_<project>/` planes
```

| Role | Who | What outsiders see |
|------|-----|--------------------|
| Admin | `waka` (operator account / admin handle) | Never exposed as a personal identity |
| Master | **WakaBot** | Public and control face for the habitat |
| Slaves | Grok Bot tech-lead, Frontier runners, project agents | Named agents under WakaBot -- not the human |

Rules:

- Grok Bot tech-lead is a **slave**: capture intent, triage, short Q&A, approvals, schematics. It does not own the ship path.
- Product code still ships only through Frontier (`hygiene` -> `plan` -> `apply` -> `git push` -> PR into `dev`).
- No slave speaks as Wade/user outwardly; WakaBot is the face.

## 2. Herdr (optional ops multiplexer -- out of ship path)

**Herdr** is an optional operations multiplexer (wait / attention semantics for multi-agent ops).

| In scope later | Not in this RFC / not in Frontier ship path |
|----------------|-----------------------------------------------|
| Borrow wait/attention semantics into control or habitat watchers when useful | Installing Herdr |
| Document non-goals so agents do not invent a mega-bot | Making Herdr a gate, hook, or `plan`/`apply` dependency |

Non-goals:

- No mega-bot that collapses WakaBot + Grok + Frontier into one process.
- No Herdr install required to use frontier-platform.
- Herdr is **not** on the Frontier ship path (`PIPELINE.md` / `plan` -> `apply` -> push).

## 3. Learn namespaces: Grok `/learn` vs Frontier `L`

These are different words that must not be conflated.

| Name | Where | Job |
|------|-------|-----|
| Grok `/learn` | Grok harness / TUI skill | Change harness skills from **traces** (skill distillation from session traces) |
| `waka-cli learn` / local `/learn` | waka-cli (local twin) | Habitat awareness: ingest -> train -> apply briefings (file/dir/root). Same *shape* as Grok `/learn`, local, no cloud API by default |
| Frontier `L` (`frontier learn`) | frontier-platform ship | **Landscape classify** before change: kind/compose/langs/topology; seal `learn.classified`. Read-only. Spec: `ship/english/L_LEARN.md` |
| Control learning ledger | frontier-control | Session **debrief** methods + `.agent_learning.jsonl`. Spec: `control/english/LEARNING.md` |

Optional ingest under WakaBot: habitat/session traces may be ingested into WakaBot-owned awareness paths for later use. That ingest is **not** Frontier `L`, and it is **not** a substitute for `frontier learn classify` before a ship change.

## 4. Providers (DeepSeek live path; Google reserved stub)

Live cloud path today remains DeepSeek (`api.deepseek.com`) -- existing alternate / test path for Frontier AI cascade and control providers.

Google AI Pro is introduced only as a **reserved fail-closed stub** in `control/providers.yaml`:

- `kind`: `google-ai` (id e.g. `google-ai-pro`)
- No secrets in tree or config
- No live calls; adapters must refuse / exit non-zero until explicitly enabled by a later doctrine + wiring change
- Must not be smuggled in by renaming banned vendors or piggybacking C7 bans

See C7 / C8 in `control/english/CONTROL.md`.

## 5. Acceptance for this PR

- [x] Authority hierarchy documented
- [x] Herdr non-goals and out-of-ship-path stated
- [x] Learn namespaces distinguished
- [x] DeepSeek entries kept; Google AI reserved stub added with no secrets
- [x] C7/C8 doctrine notes aligned (Google as explicit new kind)
- [ ] Runtime Google enablement (future PR; not this one)
