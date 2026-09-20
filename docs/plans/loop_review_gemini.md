`Loop Engineering Requirements v0.2` reaches **GREEN-** status. It successfully resolves the major structural risks identified in earlier iterations by explicitly banning persisted workflow state machines, scoping human intervention to bounded checkpoints, abstracting provider names out of conformance criteria, and anchoring Solvent as the sole authority engine.

However, four residual edge-case traps and implementation gaps must be addressed before drafting the Workflow Contract or running the POC.

---

**Critical Residual Vulnerabilities & Implementation Traps**

| Vulnerability / Ambiguity | Root Cause in Requirements v0.2 | Consequence if Unfixed | Required Architectural Fix |
| --- | --- | --- | --- |
| **1. Consequential Boundary Enforcement Gap** | Section 8 mandates that agents cannot establish authority and Conductor cannot contain policy, but does not state *who* physically blocks an un-authorized execution call.

 | A buggy, misconfigured, or adversarial agent could bypass Solvent and invoke Executor tools directly.

 | Specify that **Executors (or their ingress gateways) must operate zero-trust**. An Executor MUST reject any action payload that lacks a valid, claimed Solvent authorization token.

 |
| **2. Zombie Authorization State on Human Cancellation** | Section 7 specifies that human intervention cannot undo an executed effect, but leaves handoffs during in-flight cancellation undefined.

 | If a human cancels a task in Conductor *after* Solvent authorizes an intent but *before* the Executor executes it, the authorization remains active in Solvent.

 | Mandate a cross-role cleanup rule: **Conductor task cancellation/redirection MUST trigger a revocation check/call against any linked Solvent Action Intent**.

 |
| **3. In-Flight Authority Revocation for Long Jobs** | Section 9 decouples long-running execution from Conductor thread blocking, but ignores authority validity during multi-hour compute runs.

 | If target snapshots or credentials change 2 hours into a 10-hour BM-IST compute job, the job continues executing on stale authority.

 | Clarify whether long-running Executors evaluate Solvent authority **once at claim-time** or require **periodic heartbeat re-authorization**.

 |
| **4. Boundary Catalog "Persistence" Leak** | Section 6 requires cataloging each boundary's "persistence location, if any".

 | Contradicts the strict Non-Goals in Section 5 and 15 banning workflow databases or event stores.

 | Reframe Section 6 to clarify that "persistence location" refers **exclusively to existing Conductor activity logs or Solvent audit tables**, forbid custom tables.

 |

---

**Recommended Amendments for Requirements v0.3**

1. **Add Ingress Enforcement Rule to Section 8 (Consequential Boundary):**
> *"Executor interfaces MUST operate as zero-trust boundaries, independently validating claimed Solvent authorization prior to triggering side effects, regardless of how the execution call was initiated."*
> 


2. **Add Revocation Handoff Invariant to Section 10 (Cross-Role Invariants):**
> *"13. Task cancellation or redirection in Conductor MUST invalidate or revoke any unexecuted Solvent Action Intent bound to that task."*
> 


3. **Clarify Long-Running Execution Semantics in Section 9:**
> *"Long-running execution authorization is evaluated atomically at the time of intent claiming unless the Executor specifically implements a periodic authority re-check interface."*
> 


4. **Tighten Section 6 (Boundary Catalog):**
> *"Boundary event recording SHALL rely strictly on Conductor task activity logs or Solvent audit ledgers. No additional persistence store may be created to track boundary transitions."*
> 



With these four clarification amendments added, the requirements document will be fully sealed against semantic creep and ready to serve as the baseline for Deliverable #2 (Workflow Contract v0.1) and the BM-IST POC.