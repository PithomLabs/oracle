Agreed. This review identifies the **last substantive semantic defect** in the workflow design, and I agree with the proposed repair.

The critical correction is:

> **An operation is consequential when the operation class it instantiates is declared effect-capable—capable of producing an observable state change outside the workflow—regardless of agent intent, workflow phase, or invoking client.**

That is the correct grain. It resolves the universal-capability problem with `shell`, `python`, etc. without making every invocation consequential and without allowing those tools to become an escape hatch.

So the final semantics become:

```text
Capability
    ↓
Operation class
    ↓
operation instantiated with exact inputs/run_id
    ↓
if class is effect-capable
    → Solvent checkpoint
```

not:

```text
Capability itself = consequential
```

and not:

```text
Agent decides = consequential
```

### What I would lock now

The consolidated design is:

```text
PLAN MODE
    Agent reasons / decomposes / iterates
        ↓
    Human approves plan + operation-class scope

WORK MODE
    Agent works
        ↓
    Conductor maintains durable project state/history
        ↓
    READY frontier + atomic claims
        ↓
    Agents X/Y/Z can participate

Consequential operation
    ↓
operation class
    ↓
declaration
    ↓
Solvent exact authorization
    ↓
Executor
    ↓
external effect
```

And the ownership boundaries remain:

```text
Human      = plan approval
Agent      = intelligence
Conductor  = project state / coordination
Declaration= operation classification
Solvent    = exact consequential authorization
Executor   = external effect
```

### Two things still need to be carried into implementation

**Declaration pinning.** `capability_ref` must resolve to a specific declaration **version + content hash**, with a named resolution check before work can become actionable. The earlier simplification dropped that security property; restore it. This is especially important because the declaration now determines whether an operation reaches Solvent. The review explicitly identifies this as an integrity requirement.

**Operation-class vocabulary.** The plan's approved classes and declaration-defined classes must use the same identifiers. Otherwise plan scope comparison becomes ad hoc string matching.

### One important procedural change

The **disposition log** should now become mandatory.

Not another subsystem—just a tiny review artifact:

```text
finding | disposition | reason | source
```

Every consolidation round records what was kept, rejected, or deferred.

That prevents exactly what has happened repeatedly: a useful finding disappears in the next synthesis and has to be rediscovered.

### What remains deliberately deferred

I agree with the review that we should **not** add:

```text
claim TTL / heartbeat
proof-token subsystem
automatic plan-scope enforcement
Conductor policy routing
cryptographic Conductor↔Solvent ledger
workflow engine
scheduler
```

Those are future experiments, not current architecture.

The one operational test that should be added to Phase 1.5 is:

```text
stale claim
+
expired Solvent intent
→
defined rejection behavior
```

That is a semantic interaction worth proving, without prematurely adding lease machinery.

### Final verdict

**No remaining workflow-design objection after these edits.**

At this point I would stop redesigning and move to implementation/document lock:

1. Fix consequentiality to **operation-class grain**.
2. Restore **declaration version + content-hash pinning** and the resolution checkpoint.
3. Add the **disposition log**.
4. Correct the diagram so Conductor is not visually placed in the consequence path.
5. Add the stale-claim/expired-intent scenario to Phase 1.5.
6. Then **freeze the workflow design and test it**.

The review's own conclusion is right: after C1, the remaining risk is concentrated in declaration integrity, operation-identity semantics, and the real enforcement boundary in Agent runtimes—not in the Conductor architecture itself.
