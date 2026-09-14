"""Daily pricing-window check (Task Scheduler: WakaWeatherCheck, daily 11:05 local).

One-shot: fetches the official DeepSeek pricing page, verifies the peak-window line
still matches the fleet constant, appends one line to the log, exits. Never polls.

Log: D:\\wakalabs\\frontier-fleet\\runtime\\weather-check.log
"""
import datetime as _dt
import re
import sys
import urllib.request

URL = "https://api-docs.deepseek.com/quick_start/pricing"
EXPECTED = "01:00 - 04:00 and 06:00 - 10:00 UTC"
LOG = r"D:\wakalabs\frontier-fleet\runtime\weather-check.log"


def is_peak(now=None):
    now = now or _dt.datetime.now(_dt.timezone.utc)
    if now.weekday() >= 5:
        return False
    h = now.hour + now.minute / 60
    return (1 <= h < 4) or (6 <= h < 10)


def main():
    ok = True
    found = ""
    try:
        req = urllib.request.Request(URL, headers={"User-Agent": "frontier-fleet/1.0"})
        with urllib.request.urlopen(req, timeout=30) as r:
            text = r.read().decode("utf-8", errors="replace")
        m = re.search(r"\d{2}:00\s*-\s*\d{2}:00\s*and\s*\d{2}:00\s*-\s*\d{2}:00\s*UTC", text)
        found = m.group(0) if m else ""
        if not found or found != EXPECTED:
            ok = False
    except Exception as exc:  # noqa: BLE001 — log and exit non-zero
        ok = False
        found = "fetch error: %s" % exc

    now = _dt.datetime.now(_dt.timezone.utc)
    line = "%s | window=%s | expected_ok=%s | peak_now=%s\n" % (
        now.isoformat(), found or "missing", ok, is_peak(now))
    with open(LOG, "a", encoding="utf-8") as f:
        f.write(line)
    print(line.strip())
    sys.exit(0 if ok else 1)


if __name__ == "__main__":
    main()
