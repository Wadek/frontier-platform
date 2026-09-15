# Hygiene service

Optional local service for the ship's **H** family (`frontier hygiene`) and monitor directive D6. Default listen address: `http://127.0.0.1:8765`. It inspects AI-provenance marks (Unicode / C2PA / metadata). It is **advisory-only** — the gate does not hard-block when the service is down.

Re-homed from a retired watermark-detection service; research ML harnesses from that project are not part of this tree.

```powershell
docker build -t frontier-hygiene service/
docker run -d -p 127.0.0.1:8765:8765 frontier-hygiene
```

Keep the port loopback-bound. Start the container when you want hygiene checks to return live results.
