I reviewed the Phase 2 plan. It is **close, but not READY as written**. Two corrections are necessary before implementation; neither requires architectural expansion.

### 1. Fix the `run_id` ambiguity

The plan uses `run_id` in two different senses:

```text
operation identity:
deploy:repo:workflow:ref:run_id
```

while GitHub also has its own numeric workflow-run `id`. The plan even uses both concepts in the SOR section.  

Make them explicitly distinct:

```text
execution_run_id
    = immutable Reference Loop correlation/run identity

github_run_id
    = numeric GitHub Actions run ID assigned by GitHub
```

Then:

```text
operation_id =
deploy:repo:workflow:ref:execution_run_id
```

and:

```text
execution_run_id
    ↔ github_run_id
```

is evidence correlation, not operation identity.

This is important because GitHub does not know the operation identity when the dispatch is issued.

---

### 2. Do not claim the Executor independently verifies authorization

The plan says:

> `Executor independently verifies via kernel.Authorize` 

But the actual architecture described later is:

```text
POST AuthorizeAction
      ↓
POST ExecuteAction
      ↓
PrepareForAction
      ↓
kernel.Authorize
      ↓
ClaimIntent
      ↓
github executor
```

So the **fresh authorization enforcement occurs in Solvent's execution path**, immediately before the Executor produces the effect. The Executor adapter itself is not independently becoming an authority verifier.

Rewrite PASS-5 as:

> **The effect-capable execution boundary performs fresh Solvent authorization verification before invoking the external-effect Executor.**

And preserve:

```text
Executor ≠ authority
```

This is more faithful to the frozen Solvent architecture.

---

### 3. Declaration resolution is still the only real readiness gap

The plan now proposes:

```text
declaration_github_deploy_v1.json
        ↓
Agent reads it
        ↓
Executor uses it
```



That is sufficient for **the first experiment's declaration artifact**, but it does not fully satisfy Workflow Design v1.0's declaration-resolution contract unless the system actually verifies:

```text
owner
version
content_hash
effective_reference
retrievable_content
```

before the task becomes actionable.

Since the plan already declares this a single-operation declaration rather than a general registry, keep it minimal. Add a tiny resolver/integrity check rather than building a capability registry.

That means:

```text
capability_ref
    ↓
load declaration
    ↓
verify version + hash + effective reference
    ↓
READY
```

not merely:

```text
file exists
    ↓
READY
```

---

### 4. UNKNOWN is correctly non-blocking for this experiment

I agree with the plan's decision here.

The plan explicitly records that frozen Solvent currently collapses UNKNOWN/UNAVAILABLE into fail-closed `Allowed=false`, while the workflow semantics distinguish UNKNOWN from DENIED. 

For a **positive-path Phase 2 experiment**, that does not block execution because the scenario requires:

```text
AUTHORIZED
```

rather than UNKNOWN.

Keep the limitation explicit. Do not modify Solvent.

---

## Final verdict

After these two corrections:

```text
Phase 2 preparation = READY
```

The experiment itself is very well scoped:

```text
Human approves Plan v1
        ↓
Agent WORK
        ↓
Conductor
        ↓
exact operation
        ↓
Solvent fresh authorization
        ↓
execution boundary
        ↓
GitHub Actions
        ↓
GitHub SOR
```

The external effect is also properly separated from invocation: GitHub is the SOR, and success requires `completed + success`; invocation alone is not effect confirmation. 

**Do those two edits, then stop planning and implement Phase 2.**
