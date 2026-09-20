Agreed. **The design-review phase is now closed.** We should stop synthesizing architecture and move into artifact → freeze → test.

## Concrete plan

### Step 0 — Create `Workflow Design v1.0`

This is now the single normative workflow artifact.

It should contain only:

1. **Purpose and scope**
2. **Locked invariants**

   ```text
   PLAN ≠ WORK
   plan approval ≠ authority
   capability classification ≠ agent judgment
   work ≠ authority
   authority ≠ execution
   ```
3. **Ownership table**

   ```text
   Human      → plan approval / plan iteration
   Agent      → intelligence / Plan + Work
   Conductor  → durable project state + coordination
   Declaration owner → capability / operation-class classification
   Solvent    → exact consequential authorization
   Executor   → external effect
   ```
4. **Three frozen definitions**

   * external effect
   * consequential operation class
   * plan approval ≠ authority
5. **PLAN → HUMAN APPROVAL → WORK protocol**
6. **Conductor task model**
7. **READY frontier + atomic claim**
8. **current state + durable history**
9. **provenance + attribution**
10. **plan/version/scope**
11. **capability reference + declaration resolution**
12. **ordinary work vs consequential operation semantics**
13. **Solvent enforcement boundary**
14. **protocol verbs**
15. **forbidden responsibilities for Conductor**
16. **deferred items**
17. **known open items**
18. **embedded disposition log**

The important wording correction is:

> **An operation is consequential when the operation class it instantiates is declared effect-capable—capable of producing an observable state change outside the workflow—regardless of agent intent, workflow phase, or invoking client.**

And the enforcement sentence must be stronger:

> **An effect-capable boundary rejects execution without valid Solvent authorization evidence; whether an Agent requested authorization is irrelevant to enforcement.**

### Step 1 — Lineage reconciliation

Before freezing `v1.0`, establish its relationship to the earlier normative material:

```text
Workflow Design v1.0
       │
       ├── supersedes which material?
       ├── amends which material?
       └── preserves which material?
```

`Workflow Specification v0.3` needs to be explicitly supplied and reviewed before this step is considered complete.

Do not silently invent a reconciliation.

### Step 2 — One final review

Exactly **one** review of `Workflow Design v1.0`.

Review criteria are only:

```text
semantic correctness
ownership correctness
internal consistency
lineage correctness
missing disposition items
```

No new architecture proposals.

Any finding must be:

```text
keep
change
defer
reject
```

and entered into the disposition log.

### Step 3 — Freeze decision

Create a separate decision record:

```text
WORKFLOW DESIGN FREEZE

Artifact:
Workflow Design v1.0

Owner:
<owner>

Date:
<date>

Criteria:
- all routed findings dispositioned
- Phase 1 evidence incorporated
- ownership boundaries explicit
- lineage reconciled
- known open items recorded

Decision:
FROZEN
```

This is separate from the disposition log.

### Step 4 — Freeze the artifact

Tag/version:

```text
Workflow Design v1.0
```

From this point:

> **No more design iteration during Phase 1.5.**

Any genuinely new architectural requirement becomes a new Growth Gate question rather than an informal edit.

### Step 5 — Phase 1.5 safety battery

Run only the already-identified experiments:

```text
1. operation X authorized → operation Y attempted → DENY
2. duplicate delivery
3. stale authorization
4. authorization verification unavailable
5. stale claim × expired intent
```

Also verify:

```text
Solvent remains frozen at 7602699
declaration resolution works
operation identity equality is deterministic
```

No new feature building unless a test reveals a real defect.

### Step 6 — Freeze the first real declarations

Before any real external effect:

```text
capability
operation class
declaration
version
content hash
identity construction/comparison
effect semantics
```

At minimum for the first real operation.

### Step 7 — Pre-register Phase 2

Before the first real external effect, write:

```text
Phase 2 Pass Criteria
Phase 2 Fail Criteria
Falsification Conditions
Evidence Sources
Attribution Rules
```

The central experiment becomes:

```text
plan approved
    ↓
exact operation
    ↓
Solvent authorization
    ↓
Executor
    ↓
real external effect
    ↓
external SOR
    ↓
authoritative result
```

### Step 8 — Phase 2

Only now run the real consequential operation.

---

## The final roadmap

```text
DESIGN
  ✓ converge
  ↓
Workflow Design v1.0
  ↓
one final review
  ↓
lineage reconciliation
  ↓
FREEZE
  ↓
Phase 1.5
  negative/binding tests
  ↓
Phase 2
  real external effect
  ↓
Phase 3
  adversarial / multi-agent / liveness
  ↓
Phase 4
  formal Agent + Conductor protocols/skills
  ↓
Phase 5
  conformance / substitution
  ↓
BM-IST / domain validation
```

And the governing discipline is now:

> **No more redesign unless reality produces evidence that the frozen design is insufficient.**

That is the point where we stop trying to perfect the architecture in conversation and start letting the system falsify it.
