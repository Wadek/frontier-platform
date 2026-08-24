# Stack: Python + static web UI

**Learned from:** `github.com/Wadek/tasks` (`D:\wakalabs\tasks`) — Opt-001…003 cycle.

## Detect

- `docker-compose.yml` + Python server (`server.py` / Flask / FastAPI) serving HTML/JS  
- Layout often `app/` (or `ui/`) + `data/`  
- Learn kind: often `app_compose`

## Equivalence (Layer A)

```text
pip install -r requirements-dev.txt
pytest tests/test_api_equivalence.py -v --tb=short
```

**Must stay equal**

- API status codes and JSON fields for create/toggle/update/delete/list  
- Static `/` and `/main.js` (or equivalent) **byte-identical** to files on disk when claiming cache opts  

**Testability knobs (bake into the app)**

- `TASKS_APP_DIR`, `TASKS_DATA_FILE`, `TASKS_PORT` (or stack-specific env) so tests use temp data + real UI files  
- Avoid hard-coded `/app` and `/data` only  

## Browser (Layer B)

```text
playwright install chromium
pytest tests/test_browser_smoke.py -v --tb=short
```

**Smoke flows (tasks)**

- Load UI → add → search → toggle → edit → delete  

**Env**

- Default: ephemeral server from `conftest.py` (no token)  
- Optional gateway: `TASKS_E2E_BASE=https://…?token=…` (env only, never commit)

## GitHub Actions

- Workflow: `.github/workflows/verify.yml`  
- Job runs `pytest -v` after `playwright install --with-deps chromium`  
- Must appear as a check on every Opt PR  

## Pitfalls learned (do not repeat)

| Pitfall | Fix |
|---------|-----|
| `pytest -q` hides which tests passed | Default `addopts = -v --tb=short` in `pytest.ini`; paste verbose log in PR/chat |
| Fixture named `base_url` clashes with `pytest-base-url` | Use `tasks_base` (or disable plugin) |
| Python hotspot scanner blamed whole file on one `def` | End function at dedent (see frontier `internal/optimize/hotspot.go`) |
| Live habitat used `ui/` while GitHub used `app/` | Align names early; Learn/Optimize paths drift otherwise |
| Agent said “5 passed” without showing output | O_VERIFY **Visibility** — always print named results |
| `datetime.utcnow()` deprecation noise | Prefer `datetime.now(datetime.UTC)` on next touch |

## Optimize notes for this stack

| Pattern | Opt angle |
|---------|-----------|
| Linear scan by id in list | Dict/index (Opt-001) |
| Re-read static files every request | Startup byte cache (Opt-002) |
| Full list refetch after every mutation | Patch local state from response (Opt-003) |
| Pretty-print JSON every save | Compact dump if schema allows (advise) |

## FastAPI habitat variant (Satokori)

**Learned from:** `D:\wakalabs\satokori` (cloned from `Wadek/Farm`, relaunched as Satokori).

### Detect extra

- `main.py` + `app/routes/` + `static/*.html` (not only `server.py` in `app/`)
- `uvicorn main:app`
- Learn still `app_compose` when `docker-compose.yml` + `data/` exist

### Equivalence extras

```text
pytest tests/ -v --tb=short
```

Must stay equal: `/health`, `/catalog` lot count, `/onboard` (organizer-only), `/listings/{id}/gate` drops kg **without** a `/ledger` row.

### Docker

- Image: `python:3.12-slim` (not Alpine)
- `DATABASE_URL=sqlite:////data/farm.db` on a bind volume
- Seed only if the db file is missing
- Strip CRLF on `entrypoint.sh`

### Pitfalls

| Pitfall | Fix |
|---------|-----|
| Walking `venv/` during `frontier learn` | Learn skipDir must include `venv` and `.venv` |
| Demo passwords in the HTML header | Fine for local tryout; never the production `SECRET_KEY` |
| Treating `/transactions/complete` as the default checkout | Default path is cash-at-gate; MYC ledger is opt-in ecology, not a kassa |

## Reference layout

```text
app/           # server.py, index.html, main.js
data/          # persisted JSON (not always in git)
tests/         # Layer A + B
.github/workflows/verify.yml
requirements-dev.txt
pytest.ini
```
