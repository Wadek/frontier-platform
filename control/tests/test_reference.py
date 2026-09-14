"""Tests for the tower Python reference — the deterministic core of frontier-control.

Run: python -m unittest discover -s control/tests -v   (from the repo root)
"""
import json
import os
import sys
import tempfile
import unittest
from datetime import datetime, timezone
from pathlib import Path

sys.path.insert(0, str(Path(__file__).resolve().parent.parent / "reference"))

import tower_ref as t


def utc(y, mo, d, h, mi=0):
    return datetime(y, mo, d, h, mi, tzinfo=timezone.utc)


class TestWeather(unittest.TestCase):
    # 2026-09-14 is a Monday; 2026-09-12/13 are Saturday/Sunday.
    def test_peak_windows(self):
        self.assertTrue(t.is_peak(utc(2026, 9, 14, 1, 0)))    # Mon 01:00
        self.assertTrue(t.is_peak(utc(2026, 9, 14, 3, 59)))    # Mon 03:59
        self.assertFalse(t.is_peak(utc(2026, 9, 14, 4, 0)))    # Mon 04:00 gap
        self.assertTrue(t.is_peak(utc(2026, 9, 14, 6, 0)))     # Mon 06:00
        self.assertTrue(t.is_peak(utc(2026, 9, 14, 9, 59)))    # Mon 09:59
        self.assertFalse(t.is_peak(utc(2026, 9, 14, 10, 0)))   # Mon 10:00 off
        self.assertFalse(t.is_peak(utc(2026, 9, 15, 5, 0)))    # Tue 05:00 off
        self.assertFalse(t.is_peak(utc(2026, 9, 12, 8, 0)))    # Saturday all off
        self.assertFalse(t.is_peak(utc(2026, 9, 13, 23, 59)))  # Sunday all off

    def test_now_matches_official_clock(self):
        # Clock check, not a window check: the official page said the window;
        # this only proves is_peak() answers consistently with route().
        now = datetime.now(timezone.utc)
        self.assertEqual(t.is_peak(now), t.route("teacher", now)["peak"])


class TestRoute(unittest.TestCase):
    def test_local_first(self):
        for kind in ("atomic", "audit", "ops"):
            d = t.route(kind, utc(2026, 9, 14, 7, 0))  # peak hour
            self.assertEqual(d["target"], "local")
            self.assertFalse(d["egress"])

    def test_flash_anytime(self):
        d = t.route("coach", utc(2026, 9, 14, 7, 0))   # peak hour
        self.assertEqual(d["target"], "deepseek-flash")
        self.assertTrue(d["egress"])

    def test_pro_off_peak_only(self):
        peak = t.route("teacher", utc(2026, 9, 14, 7, 0))
        self.assertEqual(peak["target"], "deepseek-flash")
        self.assertIn("pro forbidden", peak["reason"])
        off = t.route("teacher", utc(2026, 9, 14, 10, 0))
        self.assertEqual(off["target"], "deepseek-v4-pro")

    def test_unknown_kind_defaults_local(self):
        self.assertEqual(t.route("nonsense", utc(2026, 9, 14, 10, 0))["target"], "local")


class TestFlightPlan(unittest.TestCase):
    def test_happy_path(self):
        path = ["queued", "review", "taxiing", "airborne", "landed", "debriefed", "filed"]
        for a, b in zip(path, path[1:]):
            self.assertTrue(t.next_ok(a, b), f"{a} -> {b}")

    def test_holding_path(self):
        self.assertTrue(t.next_ok("review", "holding"))
        self.assertTrue(t.next_ok("holding", "review"))

    def test_handoff_path(self):
        self.assertTrue(t.next_ok("airborne", "handed_off"))
        self.assertTrue(t.next_ok("handed_off", "airborne"))

    def test_rejections_and_retry(self):
        self.assertTrue(t.next_ok("review", "rejected"))
        self.assertTrue(t.next_ok("airborne", "failed"))
        self.assertTrue(t.next_ok("failed", "review"))

    def test_illegal_transitions_denied(self):
        self.assertFalse(t.next_ok("queued", "airborne"))      # no clearance skip
        self.assertFalse(t.next_ok("landed", "filed"))         # debrief mandatory (T4)
        self.assertFalse(t.next_ok("taxiing", "filed"))
        self.assertFalse(t.next_ok("filed", "review"))         # terminal
        self.assertFalse(t.next_ok("rejected", "airborne"))

    def test_terminal(self):
        self.assertIn("filed", t.TERMINAL)
        self.assertIn("rejected", t.TERMINAL)


class TestLearning(unittest.TestCase):
    def make(self, **kw):
        rec = {"session": "2026-09-14-x", "plane": "food", "pilot": "p1",
               "runway": "local", "ts": "2026-09-14T08:30:00Z", "tldr": "t"}
        rec.update(kw)
        return rec

    def test_valid(self):
        ok, why = t.validate_record(self.make(learned=["a"]))
        self.assertTrue(ok, why)

    def test_missing_keys(self):
        ok, why = t.validate_record({"session": "x"})
        self.assertFalse(ok)

    def test_bad_runway(self):
        self.assertFalse(t.validate_record(self.make(runway="grok"))[0])

    def test_bad_ts(self):
        self.assertFalse(t.validate_record(self.make(ts="not-a-time"))[0])

    def test_append_and_read(self):
        with tempfile.TemporaryDirectory() as d:
            p = Path(d) / "learn.jsonl"
            t.append_record(p, self.make())
            t.append_record(p, self.make(session="s2"))
            rows = t.read_records(p)
            self.assertEqual(len(rows), 2)
            self.assertEqual(rows[0]["session"], "2026-09-14-x")

    def test_append_rejects_bad(self):
        with tempfile.TemporaryDirectory() as d:
            p = Path(d) / "learn.jsonl"
            with self.assertRaises(ValueError):
                t.append_record(p, self.make(runway="claude"))


class TestRegistry(unittest.TestCase):
    def test_scan(self):
        with tempfile.TemporaryDirectory() as d:
            root = Path(d)
            (root / "food").mkdir()
            (root / "opengym").mkdir()
            (root / ".agent_food").mkdir()
            (root / ".hidden").mkdir()
            (root / "ollama-data").mkdir()
            planes = t.scan_planes(root)
            names = {p["plane"] for p in planes}
            self.assertEqual(names, {"food", "opengym"})
            by_name = {p["plane"]: p for p in planes}
            self.assertTrue(by_name["food"]["has_pilot"])
            self.assertFalse(by_name["opengym"]["has_pilot"])


class TestHandoffCard(unittest.TestCase):
    def test_valid(self):
        card = {"specversion": "1.0", "type": "handoff.request",
                "source": ".agent_opengym", "id": "h1", "data": {"task": "x"}}
        self.assertTrue(t.valid_card(card))

    def test_invalid(self):
        self.assertFalse(t.valid_card({}))
        self.assertFalse(t.valid_card({"specversion": "0.9", "type": "x",
                                       "source": "s", "id": "i", "data": {}}))


class TestSessionBudget(unittest.TestCase):
    def test_soft_refusal_at_95(self):
        v = t.session_budget(95.0)
        self.assertFalse(v["ok"])
        self.assertTrue(v["refuse_new_runs"])
        self.assertTrue(v["require_handoff"])
        self.assertTrue(v["require_new_session"])

    def test_budget_ok_below_95(self):
        v = t.session_budget(94.99)
        self.assertTrue(v["ok"])
        self.assertFalse(v["refuse_new_runs"])

    def test_far_over_95_still_soft(self):
        # Even a full cache is a SOFT refusal: handoff + new session, never hard stop
        # without the handoff artifacts.
        v = t.session_budget(99.9)
        self.assertTrue(v["require_handoff"])


if __name__ == "__main__":
    unittest.main(verbosity=2)
