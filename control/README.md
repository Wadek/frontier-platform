# frontier-control

Coordination layer for **frontier-platform**. It routes work, enforces review gates, and records outcomes. It never performs the work itself.

```
English  →  english/          policy (what we mean)
Go       →  cmd/, internal/   runtime + `_test.go` witness (what is enforced)
```

- Doctrine: [english/CONTROL.md](english/CONTROL.md)
- Vocabulary: [english/TERMS.md](english/TERMS.md) (standard taxonomy/orchestrator terms)
- Architecture: [english/ARCHITECTURE.md](english/ARCHITECTURE.md)
- Pricing: [english/PRICING.md](english/PRICING.md)
- Learning format: [english/LEARNING.md](english/LEARNING.md)

## Command surface

Implemented today:

```
control pricing show|verify                      control roster [ROOT]
control transfer validate <file|->               control check
```

Reserved — these exit non-zero rather than pretend to run:

```
control taxonomy init|diff|add|approve|retire   control data generate --project <p>
control workflow run|status|review|list         control ai audit|improve <project>
control ops health|docker|ip|decommission|jobs   control providers list|test
control monitor                                 control prove
```

## Status matrix

| Layer | State |
|---|---|
| English doctrine | complete (this tree) |
| Go runtime | authored; `go test ./...` from the repo root |
| Provider adapters (`providers.yaml`) | config complete; local + cloud providers documented |
| Pricing (`pricing.yaml`) | peak/off-peak windows aligned with the official DeepSeek schedule |
