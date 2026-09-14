"""tower_ref.py — Python-stdlib reference of the tower's deterministic core.

Same logic as internal/weather, internal/flightplan, internal/learning (Go).
Runs anywhere Python 3.8+ is installed; tests in ../tests prove it right now.
SSC precedent: stdlib reference implementation first, runtime ports later.
"""
from __future__ import annotations

import json
from datetime import datetime, timezone
from pathlib import Path

# ---------------------------------------------------------------- weather
# Official DeepSeek window (api-docs.deepseek.com/quick_start/pricing):
# peak Mon-Fri 01:00-04:00 and 06:00-10:00 UTC; everything else off-peak.
PEAK_WINDOWS = ((1, 4), (6, 10))

def is_peak(now: datetime | None = None) -> bool:
    now = now or datetime.now(timezone.utc)
    if now.weekday() >= 5:  # Saturday/Sunday: all off-peak
        return False
    h = now.hour + now.minute / 60 + now.second / 3600
    return any(lo <= h < hi for lo, hi in PEAK_WINDOWS)

LOCAL_KINDS = {"atomic", "audit", "ops", "ingest", "train"}
FLASH_KINDS = {"plan", "label", "census-draft", "tips", "vision", "coach"}
PRO_KINDS = {"census-deep", "verify", "teacher", "plan-pro", "label-pro"}

def route(kind: str, now: datetime | None = None) -> dict:
    """Deterministic runway sizing. A model never decides which model to use."""
    peak = is_peak(now)
    if kind in LOCAL_KINDS:
        return {"target": "local", "model": "qwen2.5-coder:64k", "egress": False,
                "peak": peak, "reason": "local-first"}
    if kind in FLASH_KINDS:
        return {"target": "deepseek-flash", "model": "deepseek-flash", "egress": True,
                "peak": peak, "reason": f"{kind} after egress; flash"}
    if kind in PRO_KINDS:
        if peak:
            return {"target": "deepseek-flash", "model": "deepseek-flash", "egress": True,
                    "peak": True, "reason": "pro forbidden at peak"}
        return {"target": "deepseek-v4-pro", "model": "deepseek-v4-pro", "egress": True,
                "peak": False, "reason": "pro off-peak"}
    return {"target": "local", "model": "qwen2.5-coder:64k", "egress": False,
            "peak": peak, "reason": "default local"}

# ---------------------------------------------------------------- flight plan FSM
FSM = {
    "queued": {"review"},
    "review": {"taxiing", "holding", "rejected"},
    "holding": {"review"},
    "taxiing": {"airborne", "rejected"},
    "airborne": {"landed", "handed_off", "failed"},
    "handed_off": {"airborne", "failed"},
    "failed": {"review"},
    "landed": {"debriefed", "failed"},
    "debriefed": {"filed", "failed"},
    "filed": set(),
    "rejected": set(),
}
TERMINAL = {"filed", "rejected"}

def next_ok(state: str, nxt: str) -> bool:
    return nxt in FSM.get(state, set())

# ---------------------------------------------------------------- learning ledger
REQUIRED_KEYS = {"session", "plane", "pilot", "runway", "ts", "tldr"}
RUNWAYS = {"local", "flash", "pro", "dsh"}
LIST_KEYS = ("learned", "skills_proposed", "mappings_updated", "next")

def validate_record(rec: dict) -> tuple[bool, str]:
    if not isinstance(rec, dict):
        return False, "not an object"
    missing = REQUIRED_KEYS - set(rec)
    if missing:
        return False, "missing %s" % sorted(missing)
    if rec.get("runway") not in RUNWAYS:
        return False, "bad runway"
    try:
        datetime.fromisoformat(str(rec["ts"]).replace("Z", "+00:00"))
    except ValueError:
        return False, "bad ts"
    for key in LIST_KEYS:
        if key in rec and not isinstance(rec[key], list):
            return False, f"{key} not a list"
    if "attribution" in rec and not isinstance(rec["attribution"], dict):
        return False, "attribution not an object"
    return True, "ok"

def append_record(path: str | Path, rec: dict) -> str:
    ok, why = validate_record(rec)
    if not ok:
        raise ValueError(why)
    line = json.dumps(rec, sort_keys=True, ensure_ascii=False)
    with open(path, "a", encoding="utf-8") as f:
        f.write(line + "\n")
    return line

def read_records(path: str | Path) -> list[dict]:
    out = []
    with open(path, encoding="utf-8") as f:
        for line in f:
            line = line.strip()
            if line:
                out.append(json.loads(line))
    return out

# ---------------------------------------------------------------- registry
DEFAULT_EXCLUDE = {"frontier-fleet", "deepseek_harness", "ollama-data", "grok skills"}

def scan_planes(root: str | Path, exclude: set | None = None) -> list[dict]:
    exclude = exclude or DEFAULT_EXCLUDE
    root = Path(root)
    out = []
    for p in sorted(root.iterdir()):
        if not p.is_dir() or p.name.startswith(".") or p.name in exclude:
            continue
        pilot = root / (".agent_" + p.name)
        out.append({"plane": p.name, "pilot_dir": str(pilot), "has_pilot": pilot.is_dir()})
    return out

# ---------------------------------------------------------------- handoff cards
CARD_KEYS = {"specversion", "type", "source", "id", "data"}

def valid_card(env: dict) -> bool:
    return isinstance(env, dict) and CARD_KEYS <= set(env) and env.get("specversion") == "1.0"
