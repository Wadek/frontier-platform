# Stack: Habitat edge (waka-net / cloudflared)

**Learned while launching Pigeon** (`pigeon.wakalabs.net`) and auditing the wakalabs habitat.

## Detect

- `cloudflared/config.yml` with `ingress:` hostnames  
- External Docker network `waka-net`  
- Optional `scripts/manage-tunnel-routes.ps1`  

## Standard publish flow (new public app)

```text
1. App joins waka-net; prefer NO host ports (like tasks/books)
2. App has its own auth (or goes through gateway oauth2-proxy)
3. Edit cloudflared/config.yml — add hostname + service
4. Publish DNS:
   cloudflared tunnel route dns --overwrite-dns <tunnel-id> <hostname>
   (or fixed manage-tunnel-routes.ps1 once CLI args are correct)
5. docker compose restart cloudflared-waka
6. Probe: https://<hostname>/api/health (or /)
7. Frontier: learn | guard on the app repo before inviting users
```

## Exposure audit (before first public hostname)

Read-only checklist (see also Pigeon `AUDIT-2026-08-19.md`):

| Check | Why |
|-------|-----|
| `docker ps` Ports column | Spot `0.0.0.0` binds (LAN; dangerous if WAN-forwarded) |
| cloudflared ingress list | What is already on the public Internet |
| Direct vs gateway | `home`/`vault`/`immich` style direct routes need strong app auth |
| Router port-forwards | Human confirms — agents cannot see ISP NAT |
| Cookie `Secure` | `Secure=true` breaks `http://127.0.0.1`; use env `COOKIE_SECURE` |

## Pitfalls learned

| Pitfall | Fix |
|---------|-----|
| `express-session` `secure: NODE_ENV===production` | Blocks sessions on localhost HTTP after register — loft never loads |
| manage-tunnel-routes wrong CLI order | Use `cloudflared tunnel route dns --overwrite-dns <id> <hostname>` |
| Config path still `D:\waka-net\...` in old scripts | Use `D:\wakalabs\waka-net\...` after habitat move |
| Committing sqlite under `data/` | gitignore runtime data |
| Local folder still named `frontier-push-mcp` | Rename to `frontier-ship` to match GitHub/module |

## Pigeon reference

- Repo: `github.com/Wadek/pigeon`  
- Internal: `pigeon-web:80`, `pigeon-api:8787` on `waka-net`  
- Public: `https://pigeon.wakalabs.net` → `http://pigeon-web:80`  
- Loopback tryout: `127.0.0.1:8790` only  

## Guard / Optimize

- Run `frontier learn` + `frontier guard` on the app before enabling the tunnel.  
- Add stack playbook link in `.frontier/optimize/stack.md` when Optimize starts (likely `node-express` + this edge playbook).

## Satokori (2026-08-21)

Second public habitat customer: FastAPI grocer at `D:\wakalabs\satokori`, hostname `satokori.wakalabs.net`.

**Error 1033** means the Cloudflare connector is gone, not that nginx is wrong. Typical cause on this PC: Docker Desktop stopped → `docker-desktop` WSL distro Stopped → `cloudflared-waka` not running. Fix: start Docker Desktop, wait until `docker info` works, confirm `cloudflared-waka` logs `Registered tunnel connection`. Local loopback still works while the tunnel is down.

| Pitfall | Fix |
|---------|-----|
| Cloudflare **Error 1033** on `*.wakalabs.net` | Docker Desktop must stay running. `com.docker.service` is Manual; closing Desktop kills the tunnel. |
| `python:3.12-alpine` + `cryptography`/`bcrypt` | Use `python:3.12-slim`. Alpine musl wheels stall the build. |
| `entrypoint.sh` CRLF from Windows | `sed -i 's/\\r$//' /app/entrypoint.sh` in the Dockerfile. |
| Public farm app behind GitHub oauth2-proxy | Farmers on phones cannot log in. Public hostname, **app-owned auth**, like `info`/`gaze`. |
| urllib/curl **403** on the public host | Cloudflare bot fight. Send a browser User-Agent; browsers are fine. |
| Project on `C:\Users\…\src` | Habitat rule: apps live under `D:\wakalabs\<app>`. C is not for projects. |
| Five parallel farm repos (kariniemi-farms, perinnepelto, mycelium-network, Farm, Offerit) | One tree. Learn the old ones; ship in the habitat copy. |
| Hardcoded `SECRET_KEY` in the image | Compose/env only. Never bake demo secrets into `Dockerfile`. |
| Binding `0.0.0.0:8791` | Loopback only: `127.0.0.1:8791` (Pigeon uses `8790`). |
| Nginx/ingress edited, still 1033 or old vhost | Restart **both** `gateway` and `cloudflared-waka`. DNS: `manage-tunnel-routes.ps1`. |

### Satokori publish recipe

```text
1. D:\wakalabs\satokori  (not C:)
2. docker-compose joins waka-net; ports: "127.0.0.1:8791:8080"
3. gateway nginx: public server_name satokori.wakalabs.net → satokori:8080 (no oauth2)
4. cloudflared ingress: hostname satokori.wakalabs.net → http://gateway:80
5. manage-tunnel-routes.ps1 ; docker restart cloudflared-waka gateway
6. Probe: http://127.0.0.1:8791/health  then  https://satokori.wakalabs.net/health
7. frontier learn | guard | plan | apply | push before treating it as shipped
```

Product lesson (not a till): cash at the farm gate decrements stock only. No euro row, no buyer. Basket is localStorage. Do not build a second books or a “hide from Vero” mode.

Playbooks: this file **and** [python-web-static.md](python-web-static.md).
