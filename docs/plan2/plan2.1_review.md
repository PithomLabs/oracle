Approved. This is now a sufficiently constrained execution order for Agent B.

The key improvement is that the plan makes the **forensic matrix a hard prerequisite** rather than allowing Agent B to immediately start modifying Agent A's code. That preserves the empirical-pivot discipline while avoiding a second rewrite.

I would hand it to Agent B exactly as written, with one tiny operational clarification added:

> **When the forensic matrix identifies an item as COMPLETE, do not modify it unless a failing test or direct evidence shows it is incorrect. When an item is PARTIAL, fix only the missing portion. When an item is BLOCKED, stop at that boundary and report the evidence.**

That gives Agent B an explicit preservation rule:

```text
COMPLETE  → preserve
PARTIAL   → minimal correction
INCORRECT → minimal correction
MISSING   → implement
BLOCKED   → stop/report
```

The resulting execution sequence is clean:

```text
Agent A work
    ↓
Forensic gap/status matrix
    ↓
Contract closure
    ↓
Strict SQLite compatibility
    ↓
Correct MCP/REST interfaces
    ↓
Happy-path proof
    ↓
Focused tests
    ↓
STOP + report
```

The critical scope boundaries are also now explicit: SQLite compatibility is permitted, but **Solvent authority semantics are frozen**, and broad architectural surgery is a hard stop. The GitHub deploy remains the first consequential operation, while BM-IST remains explicitly deferred.

This is ready for Agent B execution.

