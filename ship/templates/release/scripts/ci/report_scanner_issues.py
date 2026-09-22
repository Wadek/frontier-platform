"""Turn Stage 1 scanner JSON into deduped GitHub issues (runner-owned feedback).

Copy to the app as `scripts/ci/report_scanner_issues.py` and call from verify.yml
after scanners write gitleaks.json / trivy.json / checkov.json / semgrep.json.
"""

from __future__ import annotations

import json
import os
import subprocess
import sys
import urllib.error
import urllib.parse
import urllib.request
from pathlib import Path

ROOT = Path(os.environ.get("GITHUB_WORKSPACE") or Path.cwd())
RUN_URL = os.environ.get("GITHUB_SERVER_URL", "https://github.com").rstrip("/")
REPO = os.environ.get("GITHUB_REPOSITORY", "")
RUN_ID = os.environ.get("GITHUB_RUN_ID", "")
TOKEN = os.environ.get("GITHUB_TOKEN") or os.environ.get("GH_TOKEN") or ""


def _api(method: str, path: str, body: dict | None = None) -> dict | list | None:
    if not TOKEN or not REPO:
        print("skip_issues: missing GITHUB_TOKEN or GITHUB_REPOSITORY", file=sys.stderr)
        return None
    url = f"https://api.github.com/repos/{REPO}{path}"
    data = None if body is None else json.dumps(body).encode("utf-8")
    req = urllib.request.Request(
        url,
        data=data,
        method=method,
        headers={
            "Authorization": f"Bearer {TOKEN}",
            "Accept": "application/vnd.github+json",
            "X-GitHub-Api-Version": "2022-11-28",
            "User-Agent": "frontier-scanner-feedback",
            "Content-Type": "application/json",
        },
    )
    try:
        with urllib.request.urlopen(req, timeout=60) as resp:
            raw = resp.read().decode("utf-8")
            return json.loads(raw) if raw else None
    except urllib.error.HTTPError as exc:
        err = exc.read().decode("utf-8", errors="replace")
        print(f"api_error {method} {path}: {exc.code} {err[:400]}", file=sys.stderr)
        return None


def _ensure_labels(names: list[str]) -> None:
    for name in names:
        _api("POST", "/labels", {"name": name, "color": "5319e7"})


def _load(path: Path):
    if not path.is_file() or path.stat().st_size == 0:
        return None
    raw = path.read_bytes()
    # Semgrep on Windows may emit cp1252 (e.g. 0x97 em dash). Only try utf-16 with BOM.
    if raw.startswith(b"\xff\xfe") or raw.startswith(b"\xfe\xff"):
        text = raw.decode("utf-16")
    else:
        text = None
        for enc in ("utf-8-sig", "utf-8", "cp1252", "latin-1"):
            try:
                text = raw.decode(enc)
                break
            except UnicodeDecodeError:
                continue
        if text is None:
            text = raw.decode("utf-8", errors="replace")
    try:
        return json.loads(text)
    except json.JSONDecodeError as exc:
        print(f"bad_json {path}: {exc}", file=sys.stderr)
        return None


def _findings() -> list[dict]:
    out: list[dict] = []
    run_link = f"{RUN_URL}/{REPO}/actions/runs/{RUN_ID}" if RUN_ID else ""

    gitleaks = _load(ROOT / "gitleaks.json")
    if isinstance(gitleaks, list):
        for row in gitleaks:
            rule = row.get("RuleID") or row.get("Description") or "leak"
            path = row.get("File") or row.get("Path") or "?"
            line = row.get("StartLine") or row.get("Line") or ""
            title = f"[gitleaks] {rule} in {path}"
            body = (
                f"**Tool:** gitleaks\n**File:** `{path}`\n**Line:** {line}\n"
                f"**Rule:** `{rule}`\n\n```\n{(row.get('Match') or row.get('Secret') or '')[:200]}\n```\n"
                f"\nWorkflow: {run_link}\n\n"
                f"## Autofix queue\n"
                f"Label `autofix:queued`. A **local-first** Frontier agent should open "
                f"`fix/vuln-<issue>` from `dev`, remove/rotate the secret, PR into `dev`, "
                f"and comment on any blocked PR that it depends on this hotfix.\n"
            )
            out.append({"tool": "gitleaks", "title": title[:240], "body": body, "autofix": True})

    trivy = _load(ROOT / "trivy.json")
    if isinstance(trivy, dict):
        for result in trivy.get("Results") or []:
            target = result.get("Target") or "?"
            for vuln in result.get("Vulnerabilities") or []:
                vid = vuln.get("VulnerabilityID") or "CVE"
                pkg = vuln.get("PkgName") or "?"
                sev = vuln.get("Severity") or "?"
                title = f"[trivy] {sev} {vid} in {pkg} ({target})"
                body = (
                    f"**Tool:** trivy\n**Severity:** {sev}\n**Package:** `{pkg}`\n"
                    f"**ID:** `{vid}`\n**Target:** `{target}`\n"
                    f"**Title:** {vuln.get('Title') or ''}\n\nWorkflow: {run_link}\n"
                )
                out.append({"tool": "trivy", "title": title[:240], "body": body})

    checkov = _load(ROOT / "checkov.json")
    failed = []
    if isinstance(checkov, dict):
        results = checkov.get("results") or {}
        failed = results.get("failed_checks") or checkov.get("failed_checks") or []
        if not failed and isinstance(checkov.get("summary"), dict):
            for value in checkov.values():
                if isinstance(value, dict) and "failed_checks" in value:
                    failed.extend(value.get("failed_checks") or [])
    elif isinstance(checkov, list):
        for block in checkov:
            if isinstance(block, dict):
                results = block.get("results") or {}
                failed.extend(results.get("failed_checks") or [])
    for row in failed:
        cid = row.get("check_id") or row.get("id") or "check"
        path = row.get("file_path") or row.get("repo_file_path") or "?"
        title = f"[checkov] {cid} in {path}"
        body = (
            f"**Tool:** checkov\n**Check:** `{cid}`\n**File:** `{path}`\n"
            f"**Name:** {row.get('check_name') or ''}\n"
            f"**Guideline:** {row.get('guideline') or ''}\n\nWorkflow: {run_link}\n"
        )
        out.append({"tool": "checkov", "title": title[:240], "body": body})

    semgrep = _load(ROOT / "semgrep.json")
    if isinstance(semgrep, dict):
        for row in semgrep.get("results") or []:
            path = row.get("path") or "?"
            check = (row.get("check_id") or "rule").split(".")[-1]
            start = ((row.get("start") or {}).get("line")) or ""
            title = f"[semgrep] {check} in {path}:{start}"
            body = (
                f"**Tool:** semgrep\n**Check:** `{row.get('check_id')}`\n"
                f"**File:** `{path}`\n**Line:** {start}\n"
                f"**Message:** {(row.get('extra') or {}).get('message') or ''}\n\n"
                f"Workflow: {run_link}\n"
            )
            out.append({"tool": "semgrep", "title": title[:240], "body": body})

    return out


def _existing_titles() -> set[str]:
    titles: set[str] = set()
    q = urllib.parse.urlencode({"state": "open", "labels": "scanner", "per_page": "100"})
    data = _api("GET", f"/issues?{q}")
    if isinstance(data, list):
        for issue in data:
            if isinstance(issue, dict) and "pull_request" not in issue:
                titles.add(issue.get("title") or "")
    return titles


def main() -> int:
    findings = _findings()
    print(f"findings={len(findings)}")
    if not findings:
        return 0
    _ensure_labels(
        ["scanner", "needs-triage", "gitleaks", "trivy", "checkov", "semgrep", "autofix:queued", "hotfix"]
    )
    existing = _existing_titles()
    created = 0
    for item in findings:
        if item["title"] in existing:
            print(f"exists: {item['title']}")
            continue
        labels = ["scanner", "needs-triage", item["tool"]]
        if item.get("autofix") or item["tool"] in {"gitleaks", "trivy"}:
            labels.append("autofix:queued")
        res = _api("POST", "/issues", {"title": item["title"], "body": item["body"], "labels": labels})
        if isinstance(res, dict) and res.get("number"):
            created += 1
            existing.add(item["title"])
            print(f"created #{res['number']}: {item['title']}")
        else:
            try:
                subprocess.run(
                    [
                        "gh",
                        "issue",
                        "create",
                        "--title",
                        item["title"],
                        "--body",
                        item["body"],
                        "--label",
                        ",".join(labels),
                    ],
                    check=False,
                )
                created += 1
            except FileNotFoundError:
                print(f"failed: {item['title']}", file=sys.stderr)
    print(f"created={created}")
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
