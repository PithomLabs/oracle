I reviewed both the amended prompt and the generated implementation plan. The **pivot direction is correct, but I would not approve this plan for implementation yet**. There are several concrete issues where the plan contradicts either the amended specification or the architectural invariants.

### The biggest issue: the plan claims Phase 0 is complete, but the evidence is too weak

The plan says the majority of complexity is deployment and concludes GO, and says the Solvent kernel is directly reusable and all Conductor semantics are load-bearing.  

That can be valid **only if those findings came from actual repository inspection**. The amended prompt explicitly makes reconnaissance the foundation of the decision and says code is the source of truth. 

The document should therefore include concrete code-level evidence for the key findings, not just conclusions.

### Major technical problem: Solvent `internal/view` cannot simply be imported

The plan says:

> “Solvent internal (`internal/view`, `internal/belief`, etc.) — KEEP — import as module”

and later proposes re-exporting `internal/view` from ARGUS.  

That is a Go module boundary problem. A package under another module's `internal/` tree is not importable from ARGUS.

So the plan needs one of these decisions:

```text
A. Solvent exposes the required projections as public packages
B. Copy/extract the required projection code into ARGUS
C. Define a public Solvent kernel/view interface and keep implementation private
```

This is not a cosmetic correction; it affects the central “import Solvent as a module” strategy.

### Major semantic problem: packet submission must not retire debt

The proposed packet flow contains:

```text
AddEvidence
RetireDebt
CreateTask
```

inside `app.SubmitPacket()`. 

That conflicts with the authority model we've been protecting:

```text
Agent → propose work/evidence
Human → discharge debt / promote / retract / authorize
```

The amended prompt explicitly says human consequential operations include discharge and promotion, while agent operations are only `get_context` and `submit_packet`. 

So **remove `RetireDebt` from the agent submission path**. Evidence may qualify a debt for later human discharge, but submission itself must not discharge it.

### Another important issue: “single transaction” is currently only asserted

The plan says:

> “single DB transaction where possible”

while calling separate methods such as `epistemic.EnterBelief`, `epistemic.AddEvidence`, `work.CreateTask`, etc. 

If those methods each open their own transaction, this is not atomic.

The implementation plan needs a definite transaction strategy:

```text
Application transaction
        ↓
single *sql.Tx
   ┌────┴────┐
epistemic   work
   └────┬────┘
      commit
```

Either subsystem methods accept a transaction abstraction, or the application layer owns the transaction and repositories operate against it.

This is one of the genuine benefits of the pivot; we should actually capture it rather than merely remove HTTP.

### The boundary test proposal needs refinement

The plan specifies both `depguard` and a “compile-time negative test.” 

`depguard` is straightforward.

A literal compile-time negative test is awkward because the test intentionally contains code that should not compile. The implementation plan should instead define a mechanical dependency-graph check, for example:

```text
depguard
+
go/packages / go list import-graph assertion
```

That gives CI an executable structural check without maintaining intentionally uncompilable source.

The architectural intent is right; the proposed mechanism needs to be implementable.

### The Conductor conclusion is internally inconsistent

The plan says:

> “All Conductor semantics are load-bearing and small. Nothing is discarded.”

But its own inventory says dependency propagation and blocked/ready behavior are simplified/deferred. 

The correct statement is:

> **All semantics required by the POC demonstration are preserved; nonessential Conductor features are simplified or deferred.**

That is much more consistent with the pivot's “delete-first” objective.

### Migration strategy needs one more decision

The plan currently proposes importing Solvent as a module while also incorporating its schema and migrations into the ARGUS application. 

The missing question is **migration ownership**.

Who owns the canonical Solvent migrations after the pivot?

I would strongly prefer:

```text
Solvent repository
    ↓
canonical epistemic schema

ARGUS
    ↓
references/version-controls the required schema
```

rather than silently duplicating ten SQL migrations into ARGUS. Otherwise we risk creating two authorities over the same schema.

This should be resolved before implementation.

### `task dev` is too destructive

The proposed Taskfile does:

```text
task dev
 → cockroach demo
 → argus reset
 → argus serve
```

which resets the database every time. 

That is not the developer experience we originally wanted.

I would retain:

```text
task dev     # non-destructive
task fresh   # destructive
task test
task verify
```

and make `argus reset` explicitly a destructive development command.

### The plan's complexity ledger needs factual cleanup

The proposed “current” module count says three modules: Oracle, Trust UI, reference-loop. 

But the actual current environment also involves Solvent and Conductor as separate Go modules/repositories. The metric needs a precise definition such as:

```text
ARGUS-owned Go modules
External Go modules/repositories
```

Otherwise the before/after measurement is misleading.

---

## What I would change before implementation

I would send the coding agent a **plan-revision prompt**, not an implementation prompt.

The revision should require exactly these corrections:

```text
1. Provide repository evidence for every Phase 0 GO finding.

2. Resolve the impossible cross-module import of Solvent internal packages.

3. Remove RetireDebt from agent packet submission.
   Only human decision paths may discharge debt.

4. Define a real single-transaction strategy across epistemic + work writes.

5. Replace the literal compile-time-negative-import concept with a practical
   mechanically executable dependency-graph check, while retaining depguard.

6. Correct the Conductor finding:
   preserve all POC-required semantics, simplify/defer nonessential behavior.

7. Explicitly define schema/migration ownership between Solvent and ARGUS.

8. Make task dev non-destructive; move reset to task fresh / explicit reset.

9. Correct the complexity ledger's module/process counts and definitions.

10. Re-check the 17-step demo against the authority invariant after these fixes.
```

The rest of the plan is directionally strong. In particular, the overall structure—four phases, explicit GO/NO-GO gate, single process, one DB, domain-pack portability, and delete-first migration—is aligned with the amended prompt. 

The **two things I would absolutely block implementation on** are the Solvent `internal/` import issue and the accidental agent-driven `RetireDebt`. Those are concrete architectural defects, not stylistic concerns.
