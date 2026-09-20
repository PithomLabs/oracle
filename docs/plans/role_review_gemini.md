While Matrix v0.1 successfully formalizes role boundaries on paper, its absolute ban on dedicated runtime infrastructure creates dangerous operational vacuums. By pushing state tracking, operation matching, and reconciliation onto stateless endpoints, the specification introduces severe integration deadlocks and silent security failures.

**Critical Structural Vulnerabilities**

* **Stateless Reconciliation Deadlock (Sections 7 & 8):** The matrix assigns ambiguous execution outcomes strictly to the Executor. If a serverless or ephemeral Executor times out during an external write, its process dies. Without a persistent execution daemon or recovery queue (prohibited in Section 10), the action sits permanently in `Execution ambiguous` with no active component capable of running reconciliation.
* **The Parameter Normalization Security Gap (Section 5.3):** Requiring total operation identity over all "effect-relevant parameters" creates an operational trap. Strict byte-level matching fails whenever non-semantic transport metadata (e.g., `trace_id`, `timestamp`, `retry_count`) shifts. Conversely, loose structural stripping opens parameter-injection vectors where undeclared parameters alter external execution without invalidating the authorization token.
* **Circular Human Authorization Dependencies (Sections 6 & 8):** Section 6 allows a Human to "revise" a proposal before authorization, while maintaining that human action "does not create a second authority path" and must route through Solvent. If an emergency revision requires re-authorization, but Solvent only accepts programmatic Agent signatures, the human is paralyzed unless Solvent explicitly implements a human-authoring interface.
* **Self-Reporting Integration Vulnerability (Section 11):** Pushing classification down to integration self-declarations creates a critical single point of failure. If a third-party tool misclassifies a destructive `DELETE` endpoint as non-consequential "Ordinary Work," the fail-closed boundary is bypassed entirely because no central proxy inspects the payload.

**Boundary Deficit Analysis**

| Boundary | Hidden Architectural Assumption | Failure Scenario |
| --- | --- | --- |
| **Execution Eligibility** | Native tools can validate complex authorization tokens. | Off-the-shelf CLI binaries or external REST endpoints cannot evaluate custom token freshness or operation bindings without custom proxy sidecars. |
| **Consequential Proposal** | Agents reliably route actions to the authority boundary. | A compromised or misconfigured Agent calls an effect-capable tool directly via the Ordinary Work boundary. |
| **Interpretation** | Agent interpretation is bound to actual execution evidence. | An Agent hallucinates that an `Execution failure` was a success, leading Conductor to advance the workflow incorrectly. |

**Mandatory Hardening Directives**

* **Permit an Operations Gateway Exception:** Allow a lightweight, generic proxy sidecar for third-party Executors to standardize token validation and canonicalization without declaring it a prohibited "new workflow runtime."
* **Mandate Timeout Escalation to Conductor:** Require any `Execution ambiguous` outcome that exceeds a defined TTL to automatically transition the Conductor task to a human escalation state, resolving the stateless worker deadlock.
* **Centralize Effect Schemes:** Move operation effect-capability declarations out of individual integrations and into a centralized, immutable schema registry owned by Solvent.