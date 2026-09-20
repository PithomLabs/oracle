## Verdict: **REVISE once more. Do not implement yet.**

Plan 10r5 incorporates the requested authority, lifecycle, MCP, migration, availability, and authentication corrections. The remaining problems are now concentrated in a few places, but several are still **P0/P1** because they undermine the exact guarantees the plan claims to establish.

### 1. **P0 — PackRef is still not actually resolved from authoritative scenario metadata**

This is the biggest remaining semantic defect.

The plan says the scenario→pack relationship is authoritative, but the implementation is:

```go
func scenarioPackMapping(scenarioID string) (string, string) {
    return "bmist", "1.0.0"
}
```

That is not scenario metadata. It is a hardcoded global mapping. The text also explicitly says the current `belief` schema has no `pack_ref`, and defers runtime resolution to a future migration. 

So the claimed invariant:

```text
scenario_id → authoritative pack
```

does not yet exist.

**Required correction:** persist the scenario→pack binding when the first packet for a scenario is accepted, then resolve every human decision from that persisted binding. Enforce that subsequent packets for the same scenario cannot silently switch packs.

A minimal POC implementation is an ARGUS-owned `scenario_pack` table:

```text
scenario_id PRIMARY KEY
pack_id
pack_version
```

Set it during packet ingestion, reject conflicting bindings, and make `SubmitDecision` read only from that table.

Do not defer this to a follow-up migration.

---

### 2. **P0 — ARGUS EADDRINUSE retry does not actually re-resolve the port**

The plan claims:

```text
ARGUS startup retries on EADDRINUSE → re-resolves ports
```

but `resolveArgusPort()` merely reads `ARGUS_PORT` back out of `.argus-runtime.env`. Therefore:

```text
8080 occupied
→ retry
→ read ARGUS_PORT=8080
→ retry 8080
→ repeat
```

It is not a re-resolution. The relevant implementation is explicit about this behavior. 

**Required correction:** on EADDRINUSE, call the actual allocator, choose a new port, update `.argus-runtime.env`, and update the effective `listen` value.

---

### 3. **P1 — `cmdServe` still has a port TOCTOU race**

The code does:

```go
ln, err := net.Listen(...)
if err == nil {
    ln.Close()
    break
}
```

and then presumably starts the real HTTP server afterward. Another process can acquire the port between `Close()` and the actual server bind. So the preflight does not guarantee the subsequent server can bind. 

**Required correction:** bind the listener once and pass that listener to the HTTP server. Do not probe-then-close.

Conceptually:

```text
net.Listen()
    ↓
listener retained
    ↓
http.Server.Serve(listener)
```

Then EADDRINUSE handling happens at the actual bind boundary.

---

### 4. **P1 — the startup lock is not atomic, and `.task/` is never created**

The plan uses:

```sh
if [ -f "$LOCKFILE" ]; then ...
echo $$ > "$LOCKFILE"
```

Two simultaneous `task dev` processes can both observe no lock and both create one. That does not actually guarantee exclusive startup. 

There is also a simpler fresh-clone defect: the plan creates:

```sh
.argus-pids
.cockroach-data
```

but the lock is:

```text
.task/dev.lock
```

and `.task` is never created.

**Required correction:** use an atomic mechanism. The simplest portable choice is:

```sh
mkdir .task/dev.lock
```

because `mkdir` is atomic. Put the owner PID inside it afterward. Or use `flock` where guaranteed available.

Also explicitly:

```sh
mkdir -p .task
```

before lock acquisition.

---

### 5. **P1 — stale-lock recovery can still start a second CRDB against the same store**

The stated lock semantics only verify the PID of `task dev`. If that process died while its CRDB child survived, the lock becomes stale, gets deleted, and a new CRDB can be launched against:

```text
.cockroach-data
```

while the old CRDB still owns the store. The plan claims concurrent startup is prevented, but the recovery logic does not establish that the database process/store is actually free. 

**Required correction:** stale-lock recovery must inspect the recorded child PIDs before starting a new stack.

At minimum:

```text
stale dev.lock
    ↓
inspect .argus-pids/*
    ↓
live expected process?
    ├─ yes → refuse/adopt explicitly
    └─ no  → clean stale PID files and continue
```

Do not infer store ownership solely from the task PID.

---

### 6. **P1 — operator_asserted evidence has no actual human-attestation implementation**

The plan correctly blocks agents from creating `operator_asserted` evidence. That part is good. But the claimed human path is:

```text
Human → Solvent wizard → evidence
```

while the implementation plan contains no wizard endpoint/UI implementation. Its test is explicitly:

```text
TestOperatorAssertedAcceptedFromWizard
    (via direct SQL insert)
```

So there is currently no real human attestation path; direct SQL is a test fixture, not human adjudication. 

This matters because BM-IST explicitly has retirement rules using `operator_asserted`, including `needNullModel` and `needFaithfulnessReview`.

**Required correction:** either implement a minimal embedded Trust UI attestation action, or explicitly remove `operator_asserted` discharge from the POC acceptance surface. Do not claim the human path exists when it is only a SQL test.

Given the frozen architecture, the minimal correct implementation is an embedded Trust UI endpoint/page that:

```text
authenticated human
→ enters attestation text
→ server computes hash
→ server creates operator_asserted evidence
→ UI can discharge against that evidence
```

No separate Solvent service/UI.

---

### 7. **P1 — migration splitter is described as “reused” but is actually duplicated**

The migration package says it reuses `splitStatements` from `internal/testdb/testdb.go`, but the plan then defines another implementation inside `migrations.go`. 

That leaves two implementations of the same SQL parsing behavior.

**Required correction:** make the migration package's splitter the single implementation and have test helpers call it. This still fits the one bounded Solvent migration-packaging change.

---

### 8. **P1 — `task fresh` is asserted but not actually closed in the plan**

The plan correctly changes `cmdReset` into an explicit destructive operation and `cmdMigrate` into the non-destructive path. 

But the lifecycle section fully specifies `task dev` and never gives the corresponding revised `task fresh` implementation. Therefore the following invariant is asserted rather than demonstrated:

```text
task dev   → migrate
task fresh → reset
```

The fresh path also needs to consume the resolved runtime DB URL rather than reviving a hardcoded port.

**Required correction:** explicitly specify `task fresh` in the final plan, including:

```text
resolve-ports
CRDB startup
wait/readiness
argus reset
argus serve
cleanup
```

using the same runtime ownership model as `task dev`.

---

### 9. **P2 — semver parser is not actually strict semver**

The implementation uses:

```go
fmt.Sscanf(s, "%d", &parts[i])
```

which can accept numeric prefixes rather than validating the entire component. It also does not reject things like extra dot components robustly.

The intended behavior is fine, but the implementation should use strict integer parsing for exactly three numeric components.

This is not architecture-blocking.

---

### 10. **P2 — canonical content hashing has a task tie-order ambiguity**

Tasks are sorted only by:

```text
Title
```

If two tasks have the same title but different descriptions/governance references, their relative order depends on input order. The resulting "canonical" hash is therefore not guaranteed canonical.

Use:

```text
Title
Description
GovernanceRef
```

as the complete sort key.

---

## What is now solid

The following corrections are substantially closed:

**Authority DTO separation.** The external decision request no longer carries `PrincipalID`; the internal authenticated command does. 

**PackRef is no longer client-selected.** The plan explicitly removes it from the external decision request. The remaining problem is that the supposed server-side mapping is still hardcoded. 

**UNKNOWN ≠ EMPTY.** The per-section availability model is materially better than the previous shared availability field and preserves independent source failures.

**Migration ownership.** Moving migrations behind a public Solvent migration package and applying both migration sets from `argus migrate` closes the previous ownership gap. 

**Agent evidence authority.** Rejecting `operator_asserted` from `submit_packet` is the correct authority boundary; only the human-side creation path is still unfinished. 

**Developer lifecycle direction.** The important semantic correction is now present: `task dev` applies idempotent migrations rather than destructive reset. 

## Final disposition

I would **not freeze Plan 10r5 yet**.

The final revision should close exactly these five functional gaps before implementation:

```text
1. Persist authoritative scenario → pack binding
2. Make ARGUS retry actually allocate a new port
3. Eliminate probe-then-close listener race
4. Make dev locking atomic and safe around stale child processes
5. Implement the actual human operator_asserted attestation path
```

Then make the three small specification cleanups:

```text
6. One migration statement-splitting implementation
7. Explicit task fresh lifecycle
8. Strict canonical sorting / semver parsing
```

After those, the plan is at the point where implementation can proceed without another architecture redesign.
