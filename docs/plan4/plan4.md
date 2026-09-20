The latest two reviews sharpen the design considerably. After consolidating them with my prior position and applying the same **less-is-more** constraint, I would make **one major correction to our current lock** and preserve several small rules that are actually semantic necessities.

The result is still a very small architecture.

## 1. The two-mode model survives

This remains the core:

```text
PLAN
  ↓
HUMAN APPROVAL
  ↓
WORK
```

Exactly two Agent modes. No planning engine in Conductor.

The human controls plan iteration. Conductor records the resulting plan/decision; it does not orchestrate the iterations.

And:

```text
Human plan approval ≠ Solvent authority
```

The reviews agree this is a major improvement. 

---

# 2. Correct the consequentiality definition

This is the most important correction.

Our previous wording:

> "anything whose successful completion intentionally mutates authoritative state outside the workflow"

should **not** be locked.

The word **"intentionally"** gives the Agent a classification escape hatch. "Authoritative state" is also too narrow and ambiguous. The latest review correctly identifies both problems. 

Use:

> **An operation is consequential when its declared capability is capable of producing an external effect—an observable state change outside the workflow—regardless of agent intent, workflow phase, or invoking client.**

Therefore:

```text
declaration → classification
Solvent      → authorization
Executor     → effect
```

Not:

```text
Agent → decides whether consequential
Conductor → decides whether consequential
```

This is a critical semantic boundary.

---

# 3. Do not add a task-level ordinary/consequential fork

Both reviews converge here.

The task is just work.

```text
Task
  ↓
capability_ref
```

The capability/declaration determines what the operation can do.

And the **operation** is classified at the boundary.

So there is no:

```text
kind = ordinary
kind = consequential
```

and no:

```text
READY_FOR_AUTHORIZATION
```

The task lifecycle remains identical.

The consequential branch occurs only when the Agent actually reaches an effect-capable operation.

That preserves the smallest Conductor.

---

# 4. Human approval still needs one small amount of structure

This is where the latest review adds something genuinely useful.

We previously said:

> human approves the plan, but does not enumerate every exact future operation.

Keep that.

However, a plan must still have a **machine-readable scope at the operation-class level**.

Not exact operations.

For example:

```text
PLAN v3

Allowed operation class:
    DeployWorkflow
    target class:
        repository = pithomlabs/reference-loop-test
```

Then later:

```text
operation_id =
    deploy:pithomlabs/reference-loop-test:ref-loop.yml:main:run-42
```

is formed at runtime and checked by Solvent.

So we get three grains:

```text
Human:
    approves operation class / scope

Agent:
    constructs exact operation

Solvent:
    authorizes exact operation
```

This is the right compromise.

The latest review correctly identifies that `plan_id + operation_id` alone is insufficient to determine whether an operation was actually within the approved plan. 

---

# 5. But do not build a plan-enforcement engine

The above does **not** mean Conductor should enforce plan scope.

Conductor records:

```text
plan_id
plan_version
scope
approval
```

The evidence later lets us compare:

```text
approved class
        vs
actual operation class
```

But automatic prevention is not required yet.

This keeps:

```text
Conductor = memory / coordination
Solvent = exact authority
```

rather than:

```text
Conductor = second authority system
```

The review itself recognizes that plan/operation prevention should remain a future Growth Gate question. 

---

# 6. Four small rules that must be restored

The latest reviews identify four things that were silently dropped. I agree they are not "features"; they are **semantic fences**.

### A. Capability resolution

`capability_ref` must resolve to an actual declaration.

A task cannot become actionable if the referenced capability has no valid declaration.

Keep this minimal:

```text
capability_ref
    ↓
declaration
```

No capability registry inside Conductor is required.

The owner of the declaration serves it; Conductor merely records the reference.

The review correctly identifies this as load-bearing. 

### B. Acceptance ownership

Conductor must not decide whether submitted work is correct.

It records:

```text
ACCEPT
REJECT
```

from the authority that actually makes the decision.

That authority may be:

```text
test system
human
domain verifier
```

`verification_ref` identifies that authority.

This prevents "DONE because the Agent said so." 

### C. History fence

Conductor's history contains:

```text
coordination facts
attribution
decisions
references
```

It must **not duplicate authoritative Solvent or Executor facts as if they were its own authority**.

Most importantly:

```text
decision ≠ authorization
```

This is a small but important invariant. 

### D. Attribution

Tasks and decisions need minimal provenance:

```text
created_by
created_at
```

The task text is data, not trusted instructions.

That matters once Agent X, Y, Z can join the same project. 

---

# 7. What I reject from the adversarial review

The first review proposes three "minimal primitives":

```text
claim TTLs
bounded plan scopes
proof tokens
```

I would **not add all three**.

### Bounded plan scope

**Yes**, but only as the small operation-class scope described above—not as a new plan engine.

### Claim TTL

**Not yet.**

Stale claims are a real liveness issue, but we've not observed the failure in the reference environment yet.

Keep:

```text
claim + release
```

first.

When real multi-agent runs demonstrate abandoned claims, add the smallest recovery mechanism then.

### Proof tokens

**No.**

This risks turning Conductor into a verification/attestation system.

Use:

```text
verification_ref
+
external acceptance decision
```

and record the result.

If that proves insufficient in practice, then add something stronger.

That is exactly the subtraction discipline we're trying to preserve.

The critical review's concerns about liveness and verification are valid; its proposed machinery is not yet earned. 

---

# 8. External-effect enforcement remains outside Conductor

This concern **is** real:

> What stops Work Mode from calling a raw external tool?

The answer cannot be "Conductor."

It must be:

```text
effect-capable tool
        ↓
enforcement boundary
        ↓
Solvent
        ↓
Executor
```

For a deployment using shell/python/general-purpose tools, the environment must ensure external-effect paths are mediated appropriately.

This is a **deployment/conformance requirement**, not a new Conductor subsystem.

The review's point here is correct. 

---

# 9. Tactical vs structural plan changes

We should retain the concept, but make it explicitly **advisory**.

During Work Mode:

```text
tactical adjustment
    → Agent may continue

structural plan change
    → Agent returns to Plan Mode
    → human approves revision
```

But Conductor does not attempt to algorithmically classify the distinction.

The important safety fact is:

```text
external effect
    → always reaches Solvent
```

So a bad Agent classification cannot bypass the authority boundary.

That is the right containment strategy.

---

# 10. The resulting minimal data model

I would now lock this:

```text
PROJECT
├── current state
├── durable history
├── plans
│   ├── plan_id
│   ├── version
│   ├── scope (operation classes)
│   └── approval decision
│
└── TASK
    ├── task_id
    ├── objective
    ├── capability_ref
    ├── inputs
    ├── outputs
    ├── verification_ref
    ├── dependencies
    ├── claim
    ├── lifecycle
    ├── created_by
    └── created_at
```

Derived:

```text
READY
```

Not stored.

And:

```text
ordinary/consequential
```

Not stored.

```text
authorization/execution/effect
```

Not stored as Conductor lifecycle states.

---

# 11. The resulting workflow

```text
                  HUMAN
                    │
                    ▼
              ┌───────────┐
              │   AGENT   │
              │ PLAN MODE │
              └─────┬─────┘
                    │
              plan iteration
              stays here
                    │
             HUMAN APPROVES
                    │
                    ▼
              ┌───────────┐
              │   AGENT   │
              │ WORK MODE │
              └─────┬─────┘
                    │
             work protocol
                    │
                    ▼
              ┌───────────┐
              │ CONDUCTOR │
              │           │
              │ state     │
              │ history   │
              │ tasks     │
              │ deps      │
              │ claims    │
              │ READY     │
              └───────────┘

                    │
              external effect
                    │
                    ▼
                 SOLVENT
                    │
                authorize
                    │
                    ▼
                 EXECUTOR
                    │
                    ▼
               EXTERNAL SOR
```

And the key rule is:

> **Conductor never routes on consequence. The boundary does.**

---

# 12. Final "less is more" verdict

After this round, I would lock these and only these:

### Keep

```text
PLAN / WORK
Human plan approval
Plan version + operation-class scope
Conductor durable project state
Conductor durable history
Task graph
Dependencies
READY frontier
Atomic claim
Provenance
Attribution
Capability reference
Verification reference
Solvent exact-operation checkpoint
Executor effect
```

### Explicitly defer

```text
claim TTL / heartbeat
proof-token machinery
Conductor-side policy enforcement
cryptographic Conductor↔Solvent binding
automatic plan-scope enforcement
workflow templates
scheduler
knowledge graph
DAG engine
```

The most important final correction is this:

> **The human approves the scope of intended work; the Agent determines the exact implementation; the declaration determines whether an operation is consequential; Solvent determines whether the exact consequential operation is authorized; Executor produces the effect; Conductor remembers and coordinates the whole project.**

That is the smallest coherent design I see after this review cycle.
