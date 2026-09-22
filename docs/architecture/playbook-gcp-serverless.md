# Goal: Extract Satokori Learnings into Frontier-Platform

The goal is to distill the technical hurdles, architectural pivots, and operational gotchas encountered during the Satokori GCP migration into reusable patterns, documentation, and scaffolding for `frontier-platform`.

## User Review Required

> [!IMPORTANT]
> Before executing, I need you to clarify **where** in the `frontier-platform` repository these learnings should be applied. 
> 
> * Should I update documentation (e.g., `docs/playbooks/` or `architecture.md`)?
> * Should I update Python backend scaffolding/templates (e.g., standardizing `config.py` and `Dockerfile` templates)?
> * Should I update infrastructure-as-code modules (e.g., Terraform for Cloud Run)?

## Proposed Changes (Extracted Learnings)

The following core learnings have been extracted from the past 48 hours and will be integrated into the platform:

### 1. Serverless Container Deployments (Cloud Run)
#### [NEW] Standardized Secret Manager Mounting Pattern
*   **The Trap:** Mounting a GCP Secret Manager volume directly to the working directory (`/app/.env`) completely shadows the underlying Linux directory, silently wiping out the application codebase on boot.
*   **The Standard:** All `frontier-platform` Docker templates must mount secrets to an isolated `/secrets/` volume, and securely symlink them into the app path (`RUN ln -s /secrets/.env /app/.env`).

#### [MODIFY] Pydantic Configuration Standards
*   **The Trap:** Cloud Run injects system environment variables (e.g., `PORT`, `K_SERVICE`) at runtime. If a FastAPI app's `SettingsConfigDict` uses `extra="forbid"`, the app will crash instantly on boot.
*   **The Standard:** Update `config.py` scaffolding across `frontier-platform` to strictly enforce `extra="ignore"` for cloud environments.

#### [MODIFY] Docker Entrypoint Hardening
*   **The Trap:** Legacy bootstrap logic (e.g., `if [ ! -f /data/db ]; then seed()`) in `entrypoint.sh` causes fatal crashes in immutable serverless environments. Encoding mismatches (Windows CRLF in shell scripts) also trigger opaque `no such file or directory` errors.
*   **The Standard:** Strip shell entrypoints for production. Default `frontier-platform` Dockerfiles must use direct execution arrays: `CMD ["uvicorn", "main:app", "--host", "0.0.0.0", "--port", "8080"]`.

### 2. Background Processing
#### [MODIFY] Deprecation of FastAPI BackgroundTasks for Cloud
*   **The Trap:** FastAPI's native `BackgroundTasks` execute in the web server thread. On GCP Cloud Run, CPU allocation drops to near-zero the millisecond the HTTP response is returned, freezing or killing the background task.
*   **The Standard:** For serverless deployments, `frontier-platform` must default to the **Webhook / Task Queue pattern** (e.g., Google Cloud Tasks). Background jobs must be decoupled into separate HTTP endpoints (`POST /worker/execute`) triggered by the queue.

### 3. Domain Routing & Networking
#### [NEW] Firebase Hosting as Default Ingress
*   **The Trap:** Native Cloud Run Custom Domain Mappings are deprecated or entirely unsupported in newer/specialized GCP regions (e.g., `europe-north1`).
*   **The Standard:** Update networking playbooks to use **Firebase Hosting** as the default reverse proxy for Cloud Run. It is region-agnostic, free, automatically handles edge SSL provisioning, and easily rewrites traffic to the container via `firebase.json`.

### 4. AI & SDK Optimization
#### [MODIFY] The 3-Tier AI Cascade & Lightweight Clients
*   **The Pattern:** For resilient AI features, implement a strict fallback cascade: Fast/Cheap LLM (DeepSeek) -> Reliable LLM (Gemini Flash) -> Deterministic fallback (hardcoded text).
*   **The Standard:** Avoid heavy vendor SDKs (like `openai` or `google-genai`) in edge services where possible to reduce Docker image bloat and cold start times. Standardize on raw `urllib` or `httpx` for lightweight inference calls.

### 5. Frontend & UX Patterns
#### [NEW] The `dry_run` Pre-flight Pattern
*   **The Pattern:** When building batch APIs that trigger irreversible external actions (like sending SMS via 46elks or Twilio), standardizing a `?dry_run=true` parameter.
*   **The Standard:** `frontier-platform` API templates should include a dry-run middleware/pattern that returns the exact payload/action preview to the frontend for a final user confirmation modal before execution.

## Verification Plan
1. Switch to the `frontier-platform` workspace.
2. Search for existing Dockerfiles, FastAPI configs, and architecture documentation.
3. Apply the listed learnings to the scaffolding.
4. Run `frontier pre-commit` or relevant checks to ensure the platform scaffolding remains valid.
