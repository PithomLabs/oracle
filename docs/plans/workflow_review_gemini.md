While Specification v0.1 successfully avoids creating new runtime infrastructure on paper, it achieves this by offloading massive operational complexity, uncoordinated race conditions, and integration debt onto external execution boundaries.

**1. The Integration Debt of Zero-Trust Executors**
The specification mandates that every effect-capable execution path must fail closed unless valid Solvent authorization evidence is present. However, by strictly forbidding a workflow proxy, gateway, or sidecar runtime, the spec forces every single tool, API wrapper, or script used by an Executor to natively parse, validate, and check Solvent tokens. If an agent uses 30 heterogeneous tools, you must write custom authorization verification wrappers for all 30. The moment an off-the-shelf tool (e.g., a standard HTTP client or SQL driver) is exposed without custom wrapping, an agent misclassifying a consequential action as "ordinary research" will execute it un-authorized.

**2. In-Flight Race Conditions and Split-Brain State**
Section 7.5 permits human intervention during execution, yet Section 6 explicitly states that Conductor task cancellation cannot revoke Solvent authority or undo executed effects. This creates an uncoordinated race condition: if a human cancels a task in Conductor while an Executor is midway through a multi-step operation, Conductor marks the task as canceled, but the Executor continues running under valid Solvent authorization. Without a dedicated coordination bus or cross-revocation contract, Conductor and the Executor enter a permanent split-brain state where infrastructure records conflict with real-world reality.

**3. The "Ambiguous Outcome" Deadlock**
Section 5.5 introduces the `ambiguous / unknown` outcome category for unconfirmed side effects. However, because the spec bans persistent workflow state machines and retry ledgers, there is no protocol-level mechanism to track retry counts, backoff intervals, or reconciliation state. A stateless agent receiving an `ambiguous` result is forced to either spin in an unmonitored retry loop or maintain internal state across turns, turning the agent into a hidden, fragile state machine that violates the spec's own invariants.

**4. "Event Store Denialism" and Correlation Fragility**
Section 10 claims cross-role correlation is achieved purely by stitching together native system IDs (`task_id`, `intent_id`, `execution_id`) across separate logs without a workflow event store. In production, querying three distinct, independently owned system ledgers to reconstruct causal lineage during an audit or incident response is notoriously brittle. If Conductor truncates metadata or an Executor logs asynchronously, causality breaks. Dispersing log correlation across three systems without a unified schema isn't avoiding an event store—it's building an unindexed distributed tracing nightmare.

**Hardening Requirements Before Baseline Sign-Off**

* **Define Egress Containment Standard:** Mandate that un-gated generic tools (e.g., raw shell, general HTTP) *must* be sandboxed behind an authorization-enforcing proxy, explicitly clarifying that proxying network boundaries does not constitute "new workflow infrastructure."


* **Specify Cancel-and-Revoke Semantics:** Require that when Conductor receives a human cancellation signal, it must emit an explicit revocation event to Solvent before terminating the coordination context.
* **Formalize Ambiguous Outcome Recovery:** Establish a hard protocol rule: `ambiguous` outcomes automatically trigger a mandatory human escalation or a read-only state query, explicitly forbidding automated agent retries.