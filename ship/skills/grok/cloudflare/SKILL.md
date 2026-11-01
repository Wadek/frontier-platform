---
name: cloudflare
description: >-
  Routing override for the installed Cloudflare plugin skill: load Cloudflare
  WAF guidance only when the user is on Cloudflare and asked for its WAF. Use
  when security, WAF, edge, Tunnel, or Workers comes up, and when the user says
  they do not want Cloudflare, want their own WAF, or "i dont want to use
  cloudflare".
user-invocable: true
---

# Cloudflare (user override)

This file overrides the installed Cloudflare plugin skill of the same name for
WAF routing only. Do not edit `installed-plugins`.

## Pstack root

- `principle-attack-the-premise`: the premise "this user is on Cloudflare" is
  checked before any Cloudflare-specific advice is given.

## Routing

| User is on Cloudflare | Asked for Cloudflare WAF | Action |
|---|---|---|
| yes | yes | Load plugin `waf/`. |
| yes | no | Answer the security question without the `waf/` tree; offer Cloudflare WAF as one option. |
| no, or unknown | any | Do not load `waf/`. Ask or check before assuming Cloudflare. |
| user says "no Cloudflare" / "own WAF" | any | Do not load `waf/`, Tunnel, or Workers guidance. Stay vendor-neutral. |

Evidence that a project is on Cloudflare: `wrangler.toml` / `wrangler.jsonc`,
Cloudflare nameservers on the domain, a `CF-RAY` response header, or the user
saying so.

Do not default to Cloudflare WAF, Tunnel, or Workers because a word like "WAF"
or "edge" appeared.

## Done when

- [ ] The answer states which evidence showed the user is (or is not) on
      Cloudflare, or that you asked.
- [ ] `waf/` guidance appears only in the first row of the routing table.

## Conflicts

- The installed plugin skill may trigger on "WAF" alone. This override wins
  for WAF routing; for other Cloudflare work (Workers, R2, D1) the plugin
  applies once the user is confirmed on Cloudflare.
