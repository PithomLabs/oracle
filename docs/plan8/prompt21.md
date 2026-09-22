````text
# ARGUS Invariant Audit — Belief Claim Immutability

Audit the current ARGUS implementation for one specific architectural invariant:

> Once a belief has been persisted under a belief ID, can any application
> path mutate the claim text associated with that existing belief ID?

Do NOT change code initially.

The purpose of this audit is to determine whether belief identity + original
claim content are effectively append-only, as required for the research
ledger lifecycle.

## Desired invariant

For an existing belief:

    belief ID
        +
    original claim text

must be immutable.

Research evolution must occur by creating a NEW belief ID and connecting it
to the previous belief through the existing graph semantics (`derives`,
`contradicts`, etc.), rather than overwriting the original claim.

Authorized lifecycle changes such as status, debt, evidence relationships,
retraction, and promotion may change state where the existing architecture
permits them. The question is specifically whether the CLAIM TEXT itself can
be mutated for the SAME belief ID.

## Audit the complete write path

Trace every possible path that can reach the authoritative belief record.

Inspect at minimum:

1. `argus.submit_packet`
2. packet compilation
3. packet persistence
4. application-layer belief persistence
5. Solvent/epistemic store methods
6. any HTTP handlers
7. Trust UI decision handlers
8. retraction/reopen paths
9. migration/schema definitions
10. tests and fixtures
11. direct SQL writes
12. any generic update/upsert helpers
13. any internal methods that accept a belief ID plus a claim
14. any import/seed/bootstrap path

Search for operations such as:

```text
UPDATE belief SET claim ...
UPDATE ... SET claim ...
UPSERT ... claim ...
INSERT ... ON CONFLICT ... DO UPDATE ...
Save(...)
Update(...)
Upsert(...)
Patch(...)
Replace(...)
````

Also search for generic persistence methods that could indirectly mutate
claim text even if they do not explicitly say `claim`.

## Distinguish these cases

Classify each discovered path as:

A. Can mutate existing claim text
B. Cannot mutate claim text
C. Creates a new belief ID instead
D. Unknown / requires deeper inspection

Pay particular attention to idempotency.

An idempotent resubmission of the same packet is acceptable if it returns the
existing belief rather than changing its claim.

A resubmission containing a DIFFERENT claim but the SAME belief identity must
not silently overwrite the original belief.

## Check schema-level enforcement

Determine whether the database schema itself provides any protection.

Inspect:

* primary key constraints
* unique constraints
* triggers
* stored procedures/functions
* immutable/audit fields
* foreign-key relationships

If there is no DB-level immutability constraint, explicitly say so.

Do NOT add one automatically. This is an audit first.

## Check semantic lifecycle

Trace what happens when an agent responds to an adversarial challenge.

Determine whether the current implementation supports the intended pattern:

```text
B1
claim = original proposition

      |
      | adversarial challenge
      v

B2
claim = contradiction / objection
B2 --contradicts--> B1

      |
      | corrective research
      v

B3
claim = revised / narrowed / strengthened proposition
B3 --derives--> B1 or otherwise connects through existing graph semantics
```

The important requirement is:

```text
B1.claim NEVER becomes B3.claim
```

B1 remains the historical proposition.

## Test specifically for same-ID mutation

Inspect existing tests for an explicit regression test.

If one does not exist, propose the smallest test needed to prove:

1. Insert belief B1 with claim C1.

2. Attempt any supported application path that would provide the same belief
   ID with claim C2.

3. Verify the operation is rejected, ignored, or otherwise prevented.

4. Re-read B1.

5. Verify:

   ```
   B1.claim == C1
   ```

6. Verify no second state silently replaces C1 under B1's ID.

If the current architecture does not expose any legitimate same-ID update
path, state that clearly and add a focused test only if necessary.

## Do NOT broaden scope

Do not:

* redesign the belief model
* add new MCP tools
* add versioning infrastructure
* introduce event sourcing
* change the schema
* modify Solvent
* change packet semantics
* implement new lifecycle features

This is an invariant audit.

## Final report

Return:

### 1. Verdict

One of:

* PASS — existing belief claim text is effectively immutable
* FAIL — an application path can mutate existing claim text
* PARTIAL — no normal path found, but a credible bypass remains

### 2. Evidence

For every relevant path, give:

```text
file
function
operation
result
```

Example:

```text
internal/epistemic/store.go
CreateBelief(...)
INSERT only
PASS — creates new ID, does not update existing claim
```

### 3. Mutation surface

Explicitly answer:

> Is there any application path that can mutate the claim text of an existing
> belief ID?

If yes, identify the exact path.

### 4. Lifecycle assessment

State whether the current implementation supports:

```text
enter → challenge → new proposition → graph link → historical continuity
```

without rewriting prior beliefs.

### 5. Test gap

State whether an explicit same-ID claim-immutability regression test exists.

If absent, provide the smallest concrete test to add.

Do not modify code unless the audit identifies a clear invariant violation and
you are explicitly asked to repair it.

```
```
