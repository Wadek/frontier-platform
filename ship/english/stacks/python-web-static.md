# Stack: Python + static web UI

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

- Env for app dir, data file, and port so tests use temp data + real UI files
- Avoid hard-coded `/app` and `/data` only

## Browser (Layer B)

```text
playwright install chromium
pytest tests/test_browser_smoke.py -v --tb=short
```

**Smoke flows**

- Load UI → add → search → toggle → edit → delete

**Env**

- Default: ephemeral server from `conftest.py` (no token)
- Optional gateway base URL via env only, never commit secrets

## GitHub Actions

- Workflow: `.github/workflows/verify.yml`
- Job runs `pytest -v` after `playwright install --with-deps chromium`
- Must appear as a check on every Opt PR

## Pitfalls learned (do not repeat)

| Pitfall | Fix |
|---------|-----|
| `pytest -q` hides which tests passed | Default `addopts = -v --tb=short` in `pytest.ini`; paste verbose log in PR/chat |
| Fixture named `base_url` clashes with `pytest-base-url` | Use `app_base` (or disable plugin) |
| Python hotspot scanner blamed whole file on one `def` | End function at dedent (see frontier `internal/optimize/hotspot.go`) |
| Live habitat used `ui/` while GitHub used `app/` | Align names early; Learn/Optimize paths drift otherwise |
| Agent said “5 passed” without showing output | O_VERIFY **Visibility** — always print named results |
| `datetime.utcnow()` deprecation noise | Prefer `datetime.now(datetime.UTC)` on next touch |

## Optimize notes for this stack

| Pattern | Opt angle |
|---------|-----------|
| Linear scan by id in list | Dict/index |
| Re-read static files every request | Startup byte cache |
| Full list refetch after every mutation | Patch local state from response |
| Pretty-print JSON every save | Compact dump if schema allows (advise) |

## FastAPI + compose variant

### Detect extra

- `main.py` + `app/routes/` + `static/*.html` (not only `server.py` in `app/`)
- `uvicorn main:app`
- Learn still `app_compose` when `docker-compose.yml` + `data/` exist

### Equivalence extras

```text
pytest tests/ -v --tb=short
```

Must stay equal: `/health` plus the app’s documented domain invariants (write those into this playbook copy in the customer repo, not here).

### Docker

- Image: `python:3.12-slim` (not Alpine)
- SQLite (or similar) on a bind volume; URL must match the mount
- Seed only if the db file is missing
- Strip CRLF on `entrypoint.sh`

### Pitfalls

| Pitfall | Fix |
|---------|-----|
| Walking `venv/` during `frontier learn` | Learn skipDir must include `venv` and `.venv` |
| Demo passwords in the HTML header | Fine for local tryout; never the production `SECRET_KEY` |
| Baking secrets into the image | Compose/env only |
| GCP Cloud Run shadowing app code on boot | Never mount Secret Manager volumes to `/app/` directly. Mount to `/secrets/.env` and symlink via `RUN ln -s /secrets/.env /app/.env` |
| GCP Cloud Run crashing on system env vars (e.g. `PORT`) | Pydantic `SettingsConfigDict` must explicitly set `extra="ignore"` for cloud environments |
| FastAPI `BackgroundTasks` freezing on Cloud Run | Serverless drops CPU allocation to zero after HTTP response. Use Google Cloud Tasks Webhooks for background jobs instead |
| Cloud Run domain mappings failing in new regions | Do not rely on native domain mappings. Standardize on Firebase Hosting `firebase.json` reverse proxies for universal custom domain support |
| Brittle `.sh` entrypoints breaking deployment | Deprecate `entrypoint.sh` for production containers. Embed execution directly as `CMD ["uvicorn", "main:app", "--host", "0.0.0.0", "--port", "8080"]` |

## Reference layout

```text
app/           # server.py, index.html, main.js
data/          # persisted JSON (not always in git)
tests/         # Layer A + B
.github/workflows/verify.yml
requirements-dev.txt
pytest.ini
```
