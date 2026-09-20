## prompt

just to be absolutely clear: lock in this principle, at all times our Conductor + Solvent workflow must be tool-agnostic and domain-agnostic. Even Conductor and Solvent must remain tool-agnostic and domain-agnostic as well. EBP and Pizza Bot are just implementation details and replaceable across any domain and tools, got it?



Got it. **Locked in as a permanent architectural principle.**

### Core invariant

```text
Conductor + Solvent workflow
        ↓
TOOL-AGNOSTIC
DOMAIN-AGNOSTIC
```

This applies to the components individually as well:

```text
Conductor = tool-agnostic + domain-agnostic
Solvent   = tool-agnostic + domain-agnostic
```

### Replaceable implementation layer

EBP and Pizza Bot are **not architectural dependencies**. They are replaceable implementations around the invariant core.

```text
                 DOMAIN / WORKLOAD
                       │
             ┌─────────┴─────────┐
             │ replaceable layer │
             │                   │
             │ EBP / Pizza Bot   │
             │ or alternatives   │
             └─────────┬─────────┘
                       │
                  Conductor
                       │
                    Solvent
                       │
                    Executor
                       │
                 External SOR
```

Therefore:

* **EBP is one possible epistemic/research implementation**, not part of Conductor.
* **Pizza Bot is one possible Agent/runtime implementation**, not part of the workflow architecture.
* Another domain must be able to replace EBP entirely.
* Another Agent framework must be able to replace Pizza Bot entirely.
* Another toolchain, executor, SOR, research methodology, or domain must be able to plug in without changing Conductor or Solvent's conceptual model.
* **No BM-IST concepts, EBP concepts, Pizza Bot concepts, or specific tool semantics may leak into Conductor or Solvent.**

### Hard design test

For every proposed change, ask:

> **Would this still make architectural sense if EBP, Pizza Bot, BM-IST, GitHub, and the current toolchain were all replaced tomorrow?**

If **no**, it does not belong in Conductor or Solvent.

And the stronger version:

> **Can this requirement be satisfied outside Conductor/Solvent?**

If **yes**, keep it outside the kernel.

This reinforces the existing separation:

```text
AGENCY       → replaceable
METHODOLOGY  → replaceable
TOOLS        → replaceable
DOMAIN       → replaceable
COORDINATION → Conductor
AUTHORITY    → Solvent
```

**Conductor and Solvent are the durable, domain/tool-independent core. Everything around them is replaceable.**
