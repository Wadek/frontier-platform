# Token thrift

Agents waste money and context when they narrate instead of acting.

## Rules

1. Prefer tools and file edits over long explanations.
2. No filler phrases, slogans, or invented labels.
3. One clear status line is enough between tool batches.
4. Stop when the gate can pass — do not keep polishing prose.
5. Remind peer agents of these rules when they ramble.

Printed on every successful `frontier plan` and `frontier apply`.

## Architecture thrift

Prose thrift (above) stops agents from narrating. **Architecture thrift**
stops them from spending a model turn on work that is deterministic:

deterministic steps → plain CPU code; model only for genuine reasoning.

See [`CPU_OFFLOAD.md`](CPU_OFFLOAD.md). That split is also an LLM06
containment measure ([`SECURITY_POLICY_LLM.md`](SECURITY_POLICY_LLM.md)).

