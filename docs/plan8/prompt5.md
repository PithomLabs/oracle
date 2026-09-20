I would **not switch to build mode yet**. The checklist is close, but there are four concrete issues to fix first.

### 1. Phase 0 is still implementing code

The biggest issue.

The approved gate says Phase 0 is reconnaissance **before code changes**, but the checklist currently has:

```text
0.1 import Solvent + write go.mod
0.2 add GetSnapshot / ExplainSnapshot
0.2 create migration package
0.3 define Pack + implement BM-IST
```

Those are implementation changes, not reconnaissance.

Make Phase 0:

```text
0.1 Verify Solvent reusability
    inspect/import feasibility, existing constraints, dependency graph

0.2 Verify Solvent projection + migration requirements
    confirm exactly two bounded Solvent changes are sufficient

0.3 Verify Domain Pack contract
    inspect current BM-IST pack and confirm contract coverage

GO / REVISE / NO-GO
```

Then put the actual Solvent/Pack modifications into Phase 1.

### 2. Explicitly verify Solvent uniqueness constraints before `ON CONFLICT`

Step 1.5 assumes:

```text
(scenario_id, claim_hash)
(scenario_id, content_sha256)
(scenario_id, from_id, to_id, kind)
```

already have matching unique constraints.

That must be a **Phase 0 checkpoint**. If any are absent, stop and reopen the Solvent change budget rather than silently modifying Solvent schema.

### 3. Fix the acceptance-count mismatch

The checklist says:

```text
"all 25 acceptance criteria pass"
```

but the frozen plan has **21 acceptance criteria**. That should be corrected before execution.

### 4. `argus reset` must be explicitly destructive

Step 1.1 currently says:

> `reset` drops + recreates all tables (idempotent)

Use:

```text
reset = destructive, repeatable development reset
```

`task dev` remains non-destructive; `task fresh` invokes the destructive reset.

### One additional cleanup

The checklist compresses the approved **4-phase plan into only Phase 0 + Phase 1** by putting end-to-end integration, portability, ledger verification, and cleanup into Step 1.8.

That is acceptable as an **execution checklist**, but label it clearly as:

```text
Phase 0/1 implementation checklist
```

rather than implying the architecture plan itself has only two phases.

Once those corrections are made, **switch to build mode**. The checklist itself is otherwise well-shaped: eight coherent implementation steps, explicit verification checkpoints, refusal tests early, and entity-level idempotency.
