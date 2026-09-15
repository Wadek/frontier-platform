# frontier-control

Coordination layer for **frontier-platform**. It routes work, enforces review gates, and records outcomes. It never performs the work itself.

```
English  →  english/          doctrine and vocabulary
Haskell  →  haskell/          proof surface (witnesses)
Go       →  cmd/, internal/   runtime (builds where Go is available)
Python   →  reference/        stdlib reference of the same logic — runs on any host
```

- Doctrine: [english/CONTROL.md](english/CONTROL.md)
- Vocabulary: [english/TERMS.md](english/TERMS.md) (standard taxonomy/orchestrator terms)
- Architecture: [english/ARCHITECTURE.md](english/ARCHITECTURE.md)
- Pricing: [english/PRICING.md](english/PRICING.md)
- Learning format: [english/LEARNING.md](english/LEARNING.md)

## Command surface

```
control taxonomy init|diff|add|approve|retire   control data generate --project <p>
control workflow run|status|review|list         control roster
control ai audit|improve <project>              control ops health|docker|ip|decommission|jobs
control providers list|test                     control pricing show|check
control monitor                                 control prove
```

## Status matrix

| Layer | State |
|---|---|
| English doctrine | complete (this tree) |
| Python reference (`reference/control_ref.py`) | complete; tests pass on Python 3 stdlib |
| Go runtime | authored; `go test ./...` where Go is installed |
| Haskell witnesses | authored; builds where GHC/cabal exist |
| Provider adapters (`providers.yaml`) | config complete; local + cloud providers documented |
| Pricing (`pricing.yaml`) | peak/off-peak windows aligned with the official DeepSeek schedule |
