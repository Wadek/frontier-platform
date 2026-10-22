# CPU offload — deterministic work leaves the model

Companion to [`TOKEN_THRIFT.md`](TOKEN_THRIFT.md) (prose thrift) and
[`SECURITY_POLICY_LLM.md`](SECURITY_POLICY_LLM.md) LLM06 / LLM10.

**Doctrine:** deterministic agent tasks run as plain CPU code. A model is
invoked only when real reasoning is required.

This is both a **token-economy** rule and an **LLM06 containment** measure:
fewer tools in the model's hands, fewer iterations, less coordination
surface for covert channels ([`COVERT_CHANNELS.md`](COVERT_CHANNELS.md)).

---

## Split

| Kind of step | Runs as | Model call? |
|--------------|---------|-------------|
| Fetch, poll, list, copy | CPU job | no |
| Format, assemble, render | CPU job | no |
| Schema validate, gate, seal | CPU job / Frontier binary | no |
| Summarize, plan, decide under ambiguity | Model | **yes** (counted) |

If a step has a deterministic answer given its inputs, it is a CPU job.
If a step needs judgment the inputs do not determine, it is a model call.
Narrating a deterministic step inside a model turn is a thrift failure
([`TOKEN_THRIFT.md`](TOKEN_THRIFT.md) rule 1) and an agency failure (LLM06).

---

## Rules

1. **Default off the model.** New agent pipelines start with every step as
   CPU; promote a step to a model call only with a one-line reason in the
   plan.
2. **Count model invocations.** Every model call seals a ledger row. CPU
   jobs do not pretend to be model calls and model stubs are labelled.
3. **Least tool surface on the model path.** The model reaches only the
   tools the reasoning step needs. Fetch / poll / format tools stay on
   the CPU side.
4. **Same security floor as any other service.** CPU jobs that talk to
   each other use mTLS (or equivalent), dedicated ports, non-root,
   loopback- or NetworkPolicy-isolated listeners, and no public publish.
5. **Logs are platform-owned.** CPU jobs do not select individual log
   entries to drop; retention is declarative. See D8 /
   [`COVERT_CHANNELS.md`](COVERT_CHANNELS.md).
6. **Budgets still bind.** Offloading does not exempt the remaining model
   calls from LLM10 quotas.

---

## Why this is security, not only cost

| Without offload | With offload |
|-----------------|--------------|
| Model holds fetch / schedule / delete tools "just in case" | Those tools never enter the model context |
| Every morning brief is 6–12 model turns of polling and formatting | One model turn for the summarize step; the rest is CPU |
| Repeated tool use + shared artefacts = coordination surface (Rippin et al.) | Smaller surface, fewer shared focal points |
| Unbounded iteration while "assembling" | Assembly has no model loop to unbounded-consume |

Savings figures are **estimates** until an operator compares against their
own usage dashboard. Frontier does not claim visibility into external
token meters. The architectural split is what this policy requires; the
meter reading is the operator's.

---

## Relation to other English

```
  TOKEN_THRIFT.md     prose: stop narrating, act
  CPU_OFFLOAD.md      architecture: deterministic → CPU, reason → model
  SECURITY_POLICY_LLM.md   LLM06 least agency, LLM10 budgets
  COVERT_CHANNELS.md  smaller tool/coordination surface as defense
  MONITOR.md          D8 + sealed model-invocation rows
```

## See also

- [`TOKEN_THRIFT.md`](TOKEN_THRIFT.md)
- [`SECURITY_POLICY_LLM.md`](SECURITY_POLICY_LLM.md)
- [`COVERT_CHANNELS.md`](COVERT_CHANNELS.md)
- [`MONITOR.md`](MONITOR.md)
- [`R_RUNTIME.md`](R_RUNTIME.md)
