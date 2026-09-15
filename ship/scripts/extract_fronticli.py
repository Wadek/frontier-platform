"""One-shot helper historically used to split fronticli out of cmd/frontier-git.

Paths are relative to the ship/ tree. Prefer editing internal/fronticli directly.
"""
from pathlib import Path

root = Path(__file__).resolve().parents[1]
src = root / "cmd" / "frontier-git" / "main.go"
dst = root / "internal" / "fronticli" / "cli.go"
print(f"src={src}")
print(f"dst={dst}")
print("This script is a stub; fronticli already lives at internal/fronticli.")
