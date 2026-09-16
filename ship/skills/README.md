# Skills

Bundled skills for agents working through frontier-platform.

| Tree | Purpose |
|------|---------|
| `grok/` | General workflow skills (dogfood, local-first, test-as-you-go, token-thrift, …) |

AI-provenance (Hygiene / H) tooling and its agent skill now live in the separate upstream
`watermarks-remover` project rather than in this tree; frontier keeps only the Go client
`ship/internal/hygiene/`.

List and show via the ship CLI:

```powershell
frontier skills
frontier skills show dogfood
```

Host-specific skill packs do not belong in this public tree.
