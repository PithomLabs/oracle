# Adversarial Review: Role / Boundary Matrix v0.4

This document has clearly absorbed the majority of the v0.3 findings. The "succeeded on handoff receipt" prohibition (Section 14), the "terminated/revoked" outcome (Section 14), the fail-closed-on-network-partition rule (Section 9.2), and the Review Disposition table (Section 27) are all strong additions. The document is approaching structural maturity.

However, four v0.3 findings were **not incorporated** (two of them silently), and the tightening of the macro-rules has exposed **five new mechanical gaps** that will surface during conformance harness coding or the BM-IST POC.

Here is the full adversarial teardown.

---

## Part I: v0.3 Findings Not Incorporated

### 1. The Two-Row Merge Was Rejected Without Rationale (Section 5 vs. Section 7)

**The Trap:** The v0.3 review asked to merge "Effect-Capable Operation Invocation" and "Consequential Proposal" into a single row, because Section 7 explicitly states: *"These are sequential aspects of one boundary, not alternative paths."* In v0.4, they remain two separate rows in the Boundary Matrix.

**The Critique:** A Boundary Matrix is a normative reference. If the matrix lists two rows, implementers will build two distinct API contracts, two state transitions, and two handoff protocols. Then Section 7 says they're "one boundary." This is a structural self-contradiction. You cannot have a matrix that says "two boundaries" and a prose section that says "one boundary" in the same normative document.

If the intent is to keep them as two rows because they represent two distinct *mechanical steps*, then Section 7 must be rewritten to say: *"These are two sequential mechanical steps within one logical consequential boundary."* But if they are truly one boundary, the matrix must merge them.

**The Fix:** Either merge the rows, or rewrite Section 7 to justify the separation. Do not leave the contradiction standing.

---

### 2. Declaration Version Still Missing from Authorization Evidence (Section 9.4 & Section 16)

**The Trap:** The v0.3 review asked to mandate the declaration/identity version in the authorization evidence. In v0.4:
- Section 9.4 says authorization evidence must carry *"the authorized operation identity/context."*
- The Evidence Matrix (Section 16) Authorization row lists: *"Decision + authorization reference + bound operation."* No version.
- But the Execution Attempt row lists: *"Execution reference + requested operation + **declaration/identity version**."*

**The Critique:** The Executor must validate the execution against the specific declaration version that was active at authorization time. But the authorization evidence doesn't carry the version. How does the Executor obtain it? If the Executor looks it up independently, it might get a newer version (declaration drift). If the Agent passes it alongside the authorization evidence, the Agent is supplying an unverified claim about which version was authorized. The version must flow *through* the authorization evidence itself to be trustworthy.

**The Fix:** 
- Revise Section 9.4: *"Authorization evidence MUST carry or directly expose the authorized operation identity/context **and the declaration/identity version against which the authorization was evaluated**."*
- Revise Section 16 Authorization row: Add *"+ declaration/identity version"* to the minimum evidence column.

---

### 3. "Omission is a Conformance Failure" Still Untestable by Generic Harness (Section 9.3)

**The Trap:** Section 9.3 still states: *"Omitting an effect-relevant parameter is a conformance failure."* Section 27 still rejects the universal effect-detection oracle as impossible.

**The Critique:** The generic protocol conformance harness can only verify that *declared* parameters are correctly bound. It cannot detect *undeclared* parameters because it has no oracle for discovering hidden side effects. An unscoped mandate creates a conformance criterion that the primary conformance mechanism cannot evaluate.

**The Fix:** Add one sentence to Section 9.3: *"Detection of omitted effect-relevant parameters is a design-time and integration-conformance responsibility, not a generic protocol-harness test. The protocol harness verifies correct binding of declared parameters only."*

---

### 4. Autonomous Mode + Indefinite Wait Still Not Explicitly Banned (Section 12)

**The Trap:** Section 12 says an operator-controlled unresolved hold *"is a human intervention and MUST produce intervention evidence."* The v0.3 review asked to explicitly ban this in Autonomous Mode.

**The Critique:** Calling it "a human intervention" is suggestive but not prohibitive. An implementer could argue: *"Our autonomous system logs a 'human intervention evidence' record automatically and continues waiting."* The spec doesn't explicitly say this is invalid. In Autonomous Mode, by definition, no human is present to act on the hold. An indefinite wait with no human is a deadlock.

**The Fix:** Add to Section 12: *"In Autonomous Mode, the temporal policy MUST specify a terminal/remediation action. An operator-controlled indefinite wait is only valid when the operating mode permits human participation."*

---

## Part II: New Findings in v0.4

### 5. Identity Comparison Determinism Is Asserted but Not Constrained Across Implementations (Section 9.3 & Section 26)

**The Trap:** Section 9.3 says: *"The comparison rule MUST be deterministic."* Section 26 asks: *"Can two conforming implementations apply different operation-identity equality rules to the same declared operation?"* But the spec never **answers** its own question.

**The Critique:** "Deterministic" means "same input → same output within one implementation." It does NOT mean "same input → same output *across* implementations." Two conforming Executors could use different JSON canonicalization, different parameter-ordering rules, or different case-sensitivity conventions and reach opposite conclusions about whether X = Y. One rejects the execution, the other allows it. Both are "deterministic." Both are "conforming." The security guarantee collapses.

**The Fix:** Add to Section 9.3: *"The operation-identity definition MUST include or reference the specific comparison/canonicalization rule used to evaluate equality. All participants evaluating operation binding for the same declared operation MUST apply the same declared comparison rule. The declaration owner is responsible for specifying this rule unambiguously."*

---

### 6. The "Executor IS the System of Record" Case Is Undefined (Section 3 & Section 14)

**The Trap:** Section 3 says: *"Where an external system of record exists, that system is authoritative for whether the external effect actually occurred."* Section 14 says an Executor must not report "succeeded" based on a handoff receipt when the effect is owned by an external SOR.

**The Critique:** Many Executors **are** the system of record. A direct synchronous database write, a REST API call that returns `200 OK` with the mutated resource, a file-system operation — in all these cases, there is no separate SOR. The Executor's confirmation IS the authoritative evidence. But the spec only addresses the case "where an external SOR exists." It never addresses the case where the Executor is the SOR.

This matters because the "succeeded on handoff receipt" prohibition in Section 14 is conditional on an external SOR existing. If the Executor is the SOR, then a `200 OK` IS sufficient to report "succeeded." But the spec doesn't say this explicitly, leaving implementers uncertain.

**The Fix:** Add to Section 3: *"When no separate external system of record exists, the Executor is authoritative for both execution reporting and actual effect occurrence. In this case, the Executor's own confirmation constitutes sufficient evidence for a 'succeeded' outcome."*

---

### 7. Conductor's Coordination State During the Consequential Branch Is Invisible (Section 5 & Section 7)

**The Trap:** The consequential path in the Boundary Matrix goes: Agent → Effect-capable integration → Authority boundary → Solvent → Executor → Result. Conductor does not appear in any of these rows as an initiating or receiving participant.

**The Critique:** Conductor "owns work lifecycle." But when the Agent submits a consequential proposal and enters the authorization wait, what happens to the Conductor task? The spec says workflow phases are vocabulary, not state. But the Conductor must at minimum **record** that the task is blocked on an external authorization, otherwise:
- The "long-running work doesn't block unrelated work" invariant (Scenario G) cannot be verified, because there's no observable evidence that Conductor knew the work was pending.
- A human looking at Conductor's activity view cannot tell whether a task is actively being worked on or stuck waiting for Solvent.
- The Agent has no coordination-level signal that its proposal was received and forwarded.

**The Fix:** Add a note to Section 7 or Section 5: *"Conductor does not classify or route consequential proposals, but Conductor MAY record coordination-relevant status transitions (e.g., 'awaiting authorization') using its existing lifecycle model. Such records are coordination facts, not authority facts, and do not imply Conductor participation in the authorization decision."*

---

### 8. The "Result" Boundary Has Ambiguous Routing (Section 5, Result Row)

**The Trap:** The Result row lists the receiving participant as *"Agent/client / Conductor."* The slash is ambiguous: does it mean "one or the other," "either," or "both"?

**The Critique:** If the result goes only to Conductor, the Agent cannot interpret it (violating the Agent's interpretation responsibility). If it goes only to the Agent, Conductor loses lifecycle visibility (violating the coordination ownership invariant). The spec must clarify the minimum delivery requirement.

**The Fix:** Revise the receiving participant to: *"Agent/client (for interpretation) AND Conductor (for coordination lifecycle), through applicable interfaces."* Or, if the spec intentionally leaves routing flexible, add: *"The result MUST be available to both the Agent (for interpretation) and Conductor (for coordination), regardless of which participant receives the initial notification."*

---

### 9. Dangling "Cached-Auth Evidence Marker" Reference (Section 27)

**The Trap:** The Review Disposition table includes: *"Cached-auth evidence marker/version — Adopted in boundary requirements — Needed for deterministic evidence and queued/long-running execution."*

**The Critique:** I cannot find where this was actually adopted in the boundary requirements. Section 9.4 doesn't mention cached auth. Section 16 doesn't list a cached-auth marker. Section 13 (long-running execution) doesn't reference it. If it was adopted, it should be visible in the normative sections. If it was deferred or partially adopted, the disposition table should say so.

**The Fix:** Either add the cached-auth evidence marker to Sections 9.4 and 16 explicitly, or update the disposition table to reflect its actual status.

---

### 10. Section 8 / Section 18 Duplication Persists

**The Trap:** Section 8 (Capability Declaration Lifecycle) and Section 18 (Integration / Capability Declaration) both contain nearly identical bulleted lists of declaration contents.

**The Critique:** This is a maintenance hazard. If one list is updated and the other is not, the spec contains a silent contradiction. The v0.3 review flagged this. It was not addressed.

**The Fix:** Keep the full list in Section 8 (where the lifecycle mechanics live). Replace the Section 18 list with: *"Every declared capability boundary MUST contain the elements defined in Section 8."*

---

## Summary Verdict

**Status: Not yet lockable. Four unaddressed v0.3 findings and five new findings require resolution.**

| # | Finding | Severity | Type |
|---|---------|----------|------|
| 1 | Two-row merge contradiction (Sec 5 vs 7) | High | Structural |
| 2 | Declaration version missing from auth evidence | Critical | Mechanical gap |
| 3 | Untestable omission mandate | Medium | Scope clarity |
| 4 | Autonomous mode + indefinite wait | Medium | Edge case |
| 5 | Identity comparison not constrained cross-implementation | Critical | Security |
| 6 | Executor-as-SOR case undefined | High | Edge case |
| 7 | Conductor invisible during consequential branch | High | Coordination gap |
| 8 | Result boundary routing ambiguous | High | Routing clarity |
| 9 | Dangling cached-auth reference | Low | Documentation |
| 10 | Section 8/18 duplication | Low | Maintenance |

**Critical path to lock:** Fix findings 2 and 5 first (security-critical mechanical gaps), then 1, 6, 7, and 8 (structural/coordination clarity), then the remaining medium/low items.