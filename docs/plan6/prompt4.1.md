PHASE 8 ADVERSARIAL REVIEW — DIFFERENTIAL CORRECTIONS

Apply ONLY these additional corrections to the previous fix instructions.
Do not undo any already-approved fixes.

==================================================
D1 — F9: DO NOT SORT INDEPENDENT LEDGER CLOCKS
==================================================

The previous fix instruction said to merge Conductor + Solvent activity
and sort globally by created_at.

REPLACE that instruction.

Conductor and Solvent are independent stores with independent clocks.
A global timestamp sort creates a deterministic order that may falsely
claim chronology / happens-before across systems.

Do NOT manufacture cross-ledger chronology.

Use an honest deterministic representation instead:

Preferred:
- preserve each source's own authoritative ordering
- keep Conductor activity in Conductor order
- keep Solvent activity in Solvent/audit order
- expose source explicitly
- if a single merged array is required for transport, use a deterministic
  source grouping/order and DO NOT imply that cross-source adjacency is
  temporal causality

Alternatively, use a source-local sequence/order plus source identifier
as the deterministic key.

Document:

> RCP guarantees deterministic presentation, not cross-ledger
> happens-before semantics.

Add a test proving repeated reads are deterministic without claiming that
Conductor event A happened before Solvent event B merely because of
timestamps.

==================================================
D2 — F12: ADD_DEBT / LIVING-LEDGER LOOP
==================================================

The adversarial review identified a real EBP gap:

EBP v2.1 says new evidence can create new debt.

Current decision types do NOT include ADD_DEBT.

Do NOT silently claim this part of EBP is implemented.

First inspect existing Solvent APIs for an existing safe debt-add mechanism.

If an existing minimal human-gated mechanism exists and can be reused
without expanding the architecture, document and use it.

Otherwise:

DO NOT invent a broad new subsystem.

For Phase 8, record explicitly in the Freeze/Reconciliation section:

> ADD_DEBT is not exercised by the Phase 8 dry run and is deferred.
> Adversarial contradictions are recorded structurally, but reinstating
> debt from new evidence remains a future human-gated extension.

Then change EBP compliance reporting:

"New evidence creates new debt"
= DEFERRED / NOT EXERCISED

NOT PASS.

Also ensure no implementation or test claims that a contradicts edge
automatically reinstates debt.

If implementation of ADD_DEBT is chosen, it MUST be:
- human-gated
- Coordinator-mediated
- Solvent-authoritative
- auditable
- domain-pack-aware
- absent from agent MCP capability surface

Prefer deferral for this thin Phase 8 unless existing infrastructure
makes the extension genuinely trivial.

==================================================
D3 — F2: BIND DISCHARGE JUSTIFICATION TO VERIFIED EVIDENCE
==================================================

The previous fix requires actual qualifying evidence to exist.

Strengthen it further.

When retirement requires a reproducible_artifact or other evidence class,
the attributed Solvent discharge MUST carry an InstrumentRef that binds
the discharge to the verified evidence object(s) that justified it.

Do NOT use InstrumentRef as an arbitrary opaque decision ID detached from
the evidence.

Required conceptual flow:

human selects debt
→ Coordinator resolves Pack rule
→ Coordinator queries actual persisted evidence
→ Coordinator verifies evidence belongs to belief/scenario
→ Coordinator verifies evidence class
→ Coordinator obtains verified evidence IDs
→ Coordinator constructs InstrumentRef from those evidence IDs
→ authenticated operator
→ Solvent /v1/discharge
→ audit records attributable discharge + justification reference

The exact InstrumentRef format can remain POC-minimal, for example:

evidence:<uuid>[,<uuid>...]

or another deterministic representation grounded in persisted evidence IDs.

Do not duplicate evidence into Solvent.

Add tests:
- discharge references verified evidence
- fake/nonexistent evidence reference cannot justify discharge
- evidence from another belief cannot justify discharge
- evidence from another scenario cannot justify discharge
- operator_asserted discharge uses the human attribution path rather than
  pretending an agent evidence row is operator evidence

==================================================
D4 — F3: PRECISION OF THE AUTHORITY DEFECT
==================================================

Clarify the security model.

The defect is:

> unauthenticated invocation of a human-gated decision endpoint

NOT:

> browser can spoof operator identity

The existing configured operator identity is conceptually correct.

The fix is:

request
→ authentication boundary
→ server-configured operator principal
→ decision handling

Do NOT add browser-supplied identity as authority.

Do NOT build a full multi-user identity system for Phase 8.

Implement the smallest server-to-server authentication mechanism
compatible with the existing system.

Document explicitly:

> Phase 8 uses configuration-trusted operator identity behind an
> authenticated Coordinator boundary. Full human authentication and
> multi-operator identity management are deferred beyond the POC.

For a localhost-only deployment, note that network binding reduces
exposure but does NOT replace the authentication boundary.

Add tests distinguishing:
- unauthenticated invocation → refused
- authenticated invocation → allowed to reach configured operator
- browser actor field cannot override configured operator

==================================================
D5 — F5 IS BLOCKING, NOT MERELY HIGH
==================================================

Update the implementation report and acceptance matrix:

Human debt discharge MUST satisfy:

Trust UI
→ authenticated Coordinator
→ Pack validation
→ actual evidence verification
→ attributed Solvent /v1/discharge

Failure to use attributed Discharge means:

PHASE 8 BLOCKED

Treat this as an acceptance-criterion failure, not a soft High finding.

Specifically verify:
- no human retirement through bare /debt/retire
- discharge contains operator identity
- discharge contains evidence justification where applicable
- audit_activity / debt_discharge preserves attribution

==================================================
D6 — MAKE PACK MEMBERSHIP NORMATIVE IN THE SPEC
==================================================

Update the Phase 8 implementation specification itself.

State explicitly:

> A debt item MUST belong to the active Domain Pack debt vocabulary
> before retirement. Absence of a retirement rule or pack membership is
> fail-closed and MUST NOT be interpreted as permission.

This prevents the exact F1 implementation error from recurring.

Add this to:
- retirement-rule section
- acceptance criteria
- test matrix

The implementation must not rely on inference from the presence/absence
of a retirement rule.

==================================================
D7 — UPDATE THE EBP COMPLIANCE MATRIX
==================================================

Correct the EBP row:

"New evidence creates new debt"

Do NOT mark PASS merely because evidence ingestion does not mutate debt.

Use:

DEFERRED / NOT EXERCISED

unless ADD_DEBT is actually implemented.

If deferred, record the exact boundary:

adversarial contradiction
→ recorded as contradicts edge
→ human may later initiate debt reintroduction
→ not part of Phase 8 acceptance

==================================================
FINAL DIFFERENTIAL TESTS
==================================================

Add these tests to the existing suite:

1. Activity ordering does not fabricate cross-ledger chronology.
2. RCP activity representation is deterministic.
3. InstrumentRef binds human discharge to verified evidence IDs.
4. Evidence from another belief cannot justify discharge.
5. Evidence from another scenario cannot justify discharge.
6. Unknown debt remains fail-closed.
7. Unauthenticated /decisions invocation is refused.
8. Configured operator identity cannot be overridden by browser content.
9. ADD_DEBT is either fully tested if implemented, or explicitly marked
   deferred and absent from acceptance claims.

==================================================
FINAL REVIEW STANDARD
==================================================

After applying these differential changes, re-evaluate the implementation
against this question:

> Can an untrusted client cause an epistemic debt transition merely by
> asserting that the right evidence exists, or by invoking the decision
> endpoint without authenticated operator authority?

The answer must be NO.

Also verify:

> Does RCP present deterministic information without inventing
> cross-ledger causal chronology?

The answer must be YES.

And:

> Does the Phase 8 EBP compliance matrix distinguish implemented doctrine
> from deliberately deferred doctrine?

The answer must be YES.

Do not implement unrelated improvements.
Do not expand Phase 8 beyond the frozen architecture.
