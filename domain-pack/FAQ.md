## prompt

just to clarify, answer brief but concise, what exactly is a domain pack and how does it clarify our goal toward domain-agnostic direction of our workflow, the goal is to maximize Solvent database schema and work around it for context graph, belief system for our physics verifier, the goal is to be ruthless in our design for the user experience where the optimal design is reached once you no longer have features to subtract

## answer

A **domain pack** is the thin layer that tells the generic ARGUS/Solvent machinery **how a particular domain talks about knowledge and evidence**.

For BM-IST, it defines things like:

```text
claim types
debt vocabulary
evidence classes
retirement rules
falsifiers
authorized verifiers
```

It should **not** redefine the core database model.

That is exactly how we get toward domain-agnostic workflow:

```text
Solvent core
────────────────────────
belief
evidence
debt
edges
provenance
lifecycle
authority
context graph
        ↑
   domain pack
────────────────────────
BM-IST semantics
```

So the goal is **not** to make the whole system abstract for abstraction's sake. The goal is to maximize the reuse of Solvent's existing schema and semantics, then put the smallest possible domain-specific layer around it.

For the physics verifier:

```text
Solvent = generic epistemic substrate
BM-IST pack = physics vocabulary/rules
Verifier = physics computation
ARGUS workflow = generic orchestration
```

That means a future domain should ideally require:

```text
new domain pack
+ domain verifier
```

rather than:

```text
new database model
+ new workflow engine
+ new authority system
+ new context model
```

And your UX principle should govern the whole design:

> **The optimal design is reached when there is nothing left to subtract.**

So we should ruthlessly ask of every feature: **does this belong in Solvent's reusable substrate, in the domain pack, or does it not need to exist at all?**
