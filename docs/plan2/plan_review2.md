Approved. The five decisions and the added evidence-collector constraint are consistent with the implementation direction we locked.

The reference-loop environment is now:

```text
Agent
 ├─ MCP ──→ Conductor ──→ SQLite
 │
 └─ REST ─→ Solvent ────→ SQLite
                  │
                  ↓
             Executor
          RecordingFunc / GitHub
                  │
                  ↓
             GitHub Actions
                  │
                  ↓
        GitHub = authoritative SOR
```

For the coding agent, the immediate scope should now be **only**:

```text
1. Create/prepare pithomlabs/reference-loop-test
2. Create the isolated SQLite test environment using modernc.org/sqlite
3. Implement the happy_path reference loop
4. Capture participant-owned evidence after execution
5. Demonstrate the complete trace
6. Stop and report
```

Do **not** begin the adversarial scenarios yet. Do **not** add the Conformance Test Matrix yet. Do **not** modify Conductor, Solvent, or Executor source unless reconnaissance proves a genuine implementation defect requiring a minimal fix.

One final requirement I would preserve in the agent prompt:

> The first successful run must work with the real component interfaces and produce evidence that can distinguish an actual GitHub workflow execution from a simulated RecordingFunc execution.

That keeps the first milestone firmly in the empirical-proof phase rather than turning it into another mock architecture exercise.
