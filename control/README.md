# frontier-control (codename Milkcow)

The fleet's control tower. It coordinates; it never does the work itself.

```
English  →  english/          what we mean
Haskell  →  haskell/          what is true (proof surface)
Go       →  cmd/, internal/   what runs (authored here; builds where Go exists)
Python   →  reference/        stdlib reference of the same logic — runs on any host,
                              tests pass right now (SSC precedent)
```

- Doctrine: [english/CONTROL_TOWER.md](english/CONTROL_TOWER.md)
- Vocabulary: [english/TERMS.md](english/TERMS.md) (Red Hat standard terms adopted)
- Architecture: [english/ARCHITECTURE_ASCII.md](english/ARCHITECTURE_ASCII.md)
- Pricing: [english/WEATHER.md](english/WEATHER.md)
- Learning format: [english/LEARNING.md](english/LEARNING.md)

## Command surface

```
tower taxonomy init|diff|add|approve|retire   tower data generate --plane <p>  (teacher pass, off-peak)
tower workflow run|status|review|list         tower roster
tower ai audit|improve <plane>                tower ops health|docker|ip|decommission|jobs
tower runways list|test                       tower weather show|check
tower monitor                                 tower prove
```

## Status matrix (honesty first)

| Layer | State |
|---|---|
| English doctrine | complete (this tree) |
| Python reference (`reference/tower_ref.py`) | complete, tests pass on Python 3.14 stdlib |
| Go runtime | authored; `go test ./...` runs where Go is installed (not on this host yet) |
| Haskell witnesses | authored; builds where GHC/cabal exist |
| Runway adapters (`runways.yaml`) | config complete; waka-cli + deepseek verified against this habitat |
| Weather (`weather.yaml`) | live, matching the official DeepSeek window |
