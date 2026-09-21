# Stack: Node Express (or Fastify) + SPA (stub)

**Status:** stub — fill after the first Node Opt cycle.

## Detect

- `package.json`, `express` / `fastify`, `src/` or `server.js`, optional Vite/React UI  

## Equivalence (Layer A)

```text
npm test -- --verbose
# or: vitest run / jest --verbose
```

supertest (or fetch) against app without listen where possible.

## Browser (Layer B)

```text
npx playwright test
```

## GitHub Actions

- `npm ci && npm test && npx playwright install --with-deps && npx playwright test`  

## Pitfalls learned

*(none yet — append after first Opt)*
