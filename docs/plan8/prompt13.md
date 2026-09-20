Plan 10r is much closer, but I would make **six corrections before execution**. The architecture remains right; these are implementation-consistency issues.

### 1. P0 — Do not default `ARGUS_OPERATOR_TOKEN`

The plan says:

```go
if operatorToken == "" {
    operatorToken = "dev-token-not-for-production"
}
```

That violates the earlier fail-closed requirement. An absent credential must make consequential actions unavailable, not silently install a known credential. 

Use:

```text
missing ARGUS_OPERATOR_TOKEN
→ startup error for UI
→ no consequential endpoints exposed
```

### 2. P0 — Fix the retirement-rule type inconsistency

The plan defines:

```go
GetRetirementRules() map[string]RetirementRule
```

but then `SubmitDecision` type-asserts:

```go
interface{ GetRetirementRules() map[string]interface{} }
```

and subsequently type-asserts each rule again. That will not satisfy the typed interface specified earlier.  

Use one typed contract end-to-end:

```go
rules := pack.GetRetirementRules()
rule, ok := rules[req.ObligationKey]
```

No `interface{}` detour.

### 3. P0 — `reproducible_artifact` must require a trusted artifact

The current hash validation only checks the artifact registry when `ArtifactRef != ""`; otherwise a supposedly reproducible artifact can still carry an agent-supplied hash. 

The rule should be:

```text
reproducible_artifact
    → ArtifactRef REQUIRED
    → trusted registry artifact REQUIRED
    → hash must match
    → verifier must be authorized
```

Only `operator_asserted` may rely on agent-attested hash metadata.

### 4. P1 — Migration `Apply()` must actually be idempotent

The plan says “apply twice, no errors,” but simply embedding and replaying ten SQL files does not guarantee that. 

Either verify every existing migration is inherently repeatable or add a lightweight migration-version mechanism. Do not discover on implementation that migration 004, 006, etc. fails on the second application.

### 5. P1 — Dynamic port plan still has a few inconsistencies

The core design is good: reserve ports before startup and resolve one runtime topology. 

But:

* `cmdMCP` still defaults to hardcoded `:26260`. 
* `resolve-ports` is referenced but not defined in the plan.
* `Reserve()` releases the listener before the actual service bind, so the race is reduced but not eliminated.
* The cleanest approach is to make startup retry on `EADDRINUSE` with a newly resolved topology.

The final acceptance should test the exact observed failure scenario: CRDB admin port occupied while ARGUS wants `:8080`.

### 6. P1 — Do not overclaim retirement-rule “mechanical evaluation”

The plan now mechanically verifies:

```text
debt item exists
evidence class matches
qualifying evidence exists
```

That is good, but the `Rule` field itself remains a human-readable string. 

So describe the POC honestly as **mechanically enforcing the Pack's declared evidence-class retirement gate**, not interpreting arbitrary natural-language retirement rules.

---

### Prompt to coding agent

```text
Revise Plan 10r before implementation. Keep the architecture frozen.

Fix these six items:

1. FAIL CLOSED on missing ARGUS_OPERATOR_TOKEN. Never install a default token.

2. Use the typed retirement-rule contract consistently:
   GetRetirementRules() map[string]RetirementRule.
   Remove all map[string]interface{} / secondary type assertions.

3. For reproducible_artifact evidence, require ArtifactRef + trusted registry artifact +
   matching hash + authorized VerifierSpec. Only operator_asserted may use agent-attested hash.

4. Make Solvent migration Apply genuinely idempotent. Verify existing SQL or add a minimal
   migration-version mechanism; Apply twice must succeed.

5. Finish port topology consistency:
   - eliminate hardcoded CRDB port in argus mcp
   - define resolve-ports explicitly
   - propagate resolved runtime config everywhere
   - retry startup on EADDRINUSE after reallocation
   - test CRDB-admin :8080 collision and occupied preferred ports
   - preserve PID/signal cleanup

6. Clarify retirement-rule semantics: the POC mechanically enforces declared debt-item,
   required-evidence-class, and qualifying-evidence gates; it does not interpret the
   human-readable Rule string.

Do not add services, databases, REST boundaries, or a third Solvent change.

Then freeze the plan and implement it.
```

Everything else in Plan 10r is directionally consistent with the pivot, particularly the one-process topology, official MCP SDK boundary, and runtime availability model. 
