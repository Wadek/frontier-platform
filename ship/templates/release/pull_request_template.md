## Summary

<!-- What changed and why. Link scanner issues with Closes #n when this PR fixes them. -->

Closes #

## Token report

Inference cascade is always local-first. Record actual consumption (0 is fine when process-only).

| Route | Tokens |
|-------|--------|
| Local tools / process | 0 |
| Laptop Ollama | 0 (fallthrough) |
| Habitat Qwen (mock until live) | 0 (mocked fallthrough) |
| Frontier AI (DeepSeek) | 0 |

## Checklist

- [ ] Feature branch (not `dev` / `main`)
- [ ] `frontier plan` + `frontier apply` before push
- [ ] Stage 1 scanners addressed or issues linked
- [ ] After merge to `dev`: close/delete the issue and delete this branch
