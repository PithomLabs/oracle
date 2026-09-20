The two-mode separation (`PLAN` vs. `WORK`) and the strict distinction between human plan approval (intent) and Solvent authorization (execution) represent major structural improvements. However, subtracting safety mechanics down to absolute zero creates critical operational deadlocks and enforcement blind spots that render Conductor incapable of coordinating untrusted or failure-prone agents.

**1. The Liveness Black Hole**
Explicitly rejecting claim leases leaves Conductor vulnerable to immediate deadlocks. In a multi-agent environment, agents regularly crash, hit token limits, or lose network connectivity. Without leases, heartbeats, or TTLs, a crashed agent leaves its claimed task permanently locked in progress. Because `READY` is dynamically derived based on dependency resolution, every downstream task blocked by that abandoned claim stalls indefinitely. Conductor cannot function as a control plane if it cannot detect or reclaim abandoned work.

**2. Solvent's Intent Context Vacuum**
Decoupling human plan approval from Solvent runtime checks creates a dangerous context vacuum. If the human-approved plan does not bound permissible operations, Solvent receives an isolated payload (e.g., `DELETE /v1/databases/prod`) at execution time with no contextual baseline to evaluate whether that mutation aligns with Plan v1. Solvent is forced to either act as an unguided policy guesser or blindly approve any syntactically valid operation, severing the chain of trust between human approval and actual execution.

**3. The Agent Self-Grading Fallacy**
Relying on the agent in `WORK` mode to decide whether a newly discovered subtask is "tactical" (autonomous) or "structural" (requires human re-approval) introduces a classic self-grading loophole. A hallucinating or runaway agent will naturally classify major scope drift—such as switching target endpoints from staging to production to resolve a failure—as a "tactical retry." Without a deterministic boundary defining structural drift, human plan approval becomes purely advisory.

**4. Cascading Corruption via Verification Blindness**
Stripping verification validation from Conductor makes task completion purely self-asserted. If Agent X marks Task B completed while generating corrupted output, Conductor instantly unlocks Task C on the `READY` frontier. Agent Y claims Task C, ingests the invalid state, and compounds the failure down the execution tree. Conductor does not need a heavy verification engine, but it must mandate a valid proof token or validator signature before updating state and unlocking dependent frontiers.

**5. The Arbitrary Execution Proxy Gap**
Shifting all enforcement to tool proxies assumes capabilities are narrowly typed. The moment an agent uses general-purpose execution tools (e.g., `bash_run` or `python_exec`), a proxy cannot distinguish between an ordinary local text transformation and an external side-effect (`curl -X DELETE`) without performing deep semantic code inspection. Pushing consequence detection strictly to tool proxies either fails on flexible tooling or forces the proxy to become a massive policy engine.

To preserve this lightweight architecture without sacrificing execution safety, three minimal primitives must be retained:

* **Claim TTLs**: A simple timeout field on claims so Conductor can automatically return stale tasks to the `READY` frontier when agents crash.
* **Bounded Plan Scopes**: Human plan approval must generate a coarse manifest of authorized capability types and target resources that Solvent can validate against runtime operation payloads.
* **Proof Tokens**: Conductor must require a valid completion or authorization token attached to `verification_ref` before closing a task and advancing dependent frontiers.