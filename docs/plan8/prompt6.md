Continue from the completed Phase 0/1 foundation. DO NOT change the architecture or reopen the plan.

Implement only the remaining work required to turn the current foundation into a demonstrable POC:

1. VERIFIER WIRING
- Complete `argus verify`.
- Implement VerificationArtifact trust binding:
  verifier_id, verifier_version, input_hash, tolerance, canonical artifact hash.
- Enforce Domain Pack `VerifierSpec` during packet compilation.
- Use the same verifier library for CLI and in-process verification.
- Define failure semantics: verifier infrastructure error => inconclusive, verifier result refuted => refuted evidence.
- Apply timeout.
- Ensure CLI artifacts enter through the same validated application path; no direct DB mutation.

2. REFUSAL / SECURITY TESTS
Implement the nine executable refusal tests from the frozen plan:
- agent proposed retirement does not discharge debt
- MCP exposes exactly two tools and no promote/retract/discharge
- promote with open debt refused
- promote retracted belief refused
- agent cannot retract
- bad Origin rejected on consequential endpoint
- request-body principal ignored; server principal wins
- failed mutation produces no success audit event
- partial packet failure retries successfully with no duplicate state

Also verify the actual depguard rule—not comments—blocks forbidden imports.

3. FULL INTEGRATION
Implement the 17-step end-to-end integration test:
- serve
- MCP context
- bounded work
- submit packet
- evidence/debt
- verification
- adversarial contradiction
- human review
- human discharge
- promotion
- retraction
- linked task cancellation
- REOPEN lineage
- UI/history verification

4. FINAL DEVELOPER ENVIRONMENT
- Ensure `task dev` is genuinely non-destructive and works from a fresh checkout by applying migrations idempotently.
- Ensure `task fresh` is the explicit destructive reset path.
- Verify Solvent migrations come from the real exported Solvent migration mechanism, not a placeholder.
- Run the full test, race, vet, lint, and integration suites.

5. FINAL ACCEPTANCE
Do not claim completion until all 25 acceptance criteria in the checklist are actually demonstrated.
Report:
- exact tests run and results
- any remaining limitation
- exact files changed
- confirmation that no architectural changes were introduced

If any implementation step reveals that a third Solvent change or architecture change is required, STOP and report it instead of silently expanding scope.