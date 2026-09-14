# Hygiene service (re-homed from watermarks-remover, 2026-09-14)

The ship's **H** family (`frontier hygiene`) and monitor directive D6 call this service at
`http://127.0.0.1:8765`. It strips AI-provenance marks (Unicode / C2PA / metadata). It is
**advisory-only** — the gate never hard-blocks on it being down.

- Provenance: `D:\wakalabs\watermarks-remover\service\` (project retired 2026-09-14;
  the ML research harnesses stayed behind and are retired with it).
- Run: `docker build -t fleet-hygiene service/ && docker run -d -p 127.0.0.1:8765:8765 fleet-hygiene`
  (port must stay loopback-bound — it runs on the host, not behind the WAF).
- Status today: service is **down** (advisory). Start it from here whenever hygiene
  checks should go green again.
