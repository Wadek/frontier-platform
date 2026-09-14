# Stack: Go HTTP API (stub)

**Status:** stub — fill after the first Go Opt cycle on a real customer/dogfood app.

## Detect

- `go.mod`, `main.go` / `cmd/`, `net/http` or chi/echo/gin  

## Equivalence (Layer A)

```text
go test ./... -v
```

Add httptest golden tests for public routes.

## Browser (Layer B)

- If no UI: skip with note in PR  
- If serves templates/SPA: Playwright against `go run` or compose URL  

## GitHub Actions

- `go test ./... -v` on PR  

## Pitfalls learned

*(none yet — append after first Opt)*
