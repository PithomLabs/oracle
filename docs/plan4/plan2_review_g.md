This two-mode model cleanly resolves the architectural impasse by drawing a sharp boundary around agent autonomy, provided three operational edge cases are strictly defined.

**Why This Framing Works**

* **Unambiguous Consequentiality Rule:** Defining "consequential" strictly as *any state mutation outside the workflow context* removes subjective risk scoring. If an action mutates an external system of record (cloud infrastructure, external databases, user communications), it hits Solvent. If it alters local workflow artifacts or internal code drafts, it executes directly.
* **Decoupled Reasoning:** Forcing agents into a distinct `Plan Mode` protects Conductor from becoming an orchestration brain. Conductor merely holds the uncommitted task graph until a human approves it, keeping all prompt engineering, subagent spawning, and plan iteration entirely inside the external agent runtime.

**Edge Cases to Resolve**

**1. The Replanning Loop & Approval Fatigue**
When Work Mode hits an unexpected runtime error (e.g., a missing internal dependency), returning to Plan Mode for every minor detour will bottleneck the human operator.

* *Solution:* Limit human Plan Mode approval to **Structural Changes** (adding new external capability calls or modifying primary task goals). Allow the agent to make autonomous **Tactical Adjustments** (internal retry logic, local refactoring) directly within Work Mode.

**2. Intent Approval vs. Execution Checks**
A human approving a high-level plan in Plan Mode cannot pre-validate dynamic execution parameters (e.g., exact commit hashes, dynamic IP ranges) generated later during Work Mode.

* *Solution:* Treat human approval as a gate that unlocks entry into Work Mode, while keeping Solvent as the runtime check that cryptographically verifies the exact, deterministic operation payload right before execution.

**3. Tool-Level Interception**
Because Conductor remains a simple control plane without internal policy enforcement, it cannot directly prevent a rogue agent in Work Mode from calling an external API.

* *Solution:* Enforce the boundary at the environment level. External capability tools must be wrapped in Solvent-gated execution proxies that reject any incoming payload lacking an authoritative Solvent approval token.

This division yields a simple Conductor, complete human governance over workflow intent, and strict cryptographic guardrails around external mutations.