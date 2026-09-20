This is an exceptionally strong, tightly scoped document. It successfully internalizes every correction from the previous review cycles and translates them into a rigorous, implementation-independent specification. 

You have successfully navigated the "Growth Gate" by explicitly defining what this is **not** (a runtime, a UI, a database, a shadow authority) while clearly defining what it **is** (a cross-role interaction contract and conformance suite).

Here is a detailed review, highlighting what works perfectly and offering a few minor refinements to make it bulletproof before moving to the conformance test design.

---

### 1. What Works Exceptionally Well

*   **The Branched Topology (Section 5):** You perfectly captured the most critical architectural correction. By explicitly showing ordinary work bypassing Solvent and only consequential work branching into it, you prevent the "over-invoke Solvent" anti-pattern.
*   **The Vocabulary Firewall (Section 2 & 5):** Explicitly banning a persisted workflow `Intent` and restricting workflow phases to "specification-level vocabulary only" completely neutralizes the shadow-authority risk identified in the first review.
*   **The Invariants List (Section 10):** This is the intellectual core of the document. Distilling the lessons from the Solvent writeup (e.g., *Capability ≠ Authority*, *Authorization ≠ Execution*, *Coordination activity ≠ proof of external execution*) into 12 hard rules gives implementers and reviewers a definitive checklist.
*   **Operating Modes (Section 3):** Framing "Autonomous" vs. "Semi-Autonomous" as behavioral modes over the *same* contract, rather than separate architectures, is a massive architectural win. It prevents the system from fracturing into two different codebases.
*   **Long-Running Operations (Section 9):** Explicitly decoupling authorization from execution state (`auth granted ≠ execution started ≠ completed ≠ accepted`) prevents the workflow layer from accidentally absorbing orchestration logic (like a naive Temporal implementation would).

---

### 2. Targeted Refinements for v0.3

While the document is 95% there, here are a few precise adjustments to close the remaining gaps and ensure the conformance suite is fully aligned with the BM-IST proof points.

#### A. Expand Section 14 (Initial BM-IST Validation Scenarios)
Section 13 correctly lists the *requirements* for conformance, but Section 14 only provides three scenarios (A, B, C). To fully prove the thesis established in the reviews, you should add two more explicit scenarios to Section 14:

*   **Add Scenario D — Tool Substitution (The Agnosticism Proof):**
    *   *Action:* Run Scenario B (Consequential Shared Computation) using three different clients: (1) an LLM Agent (e.g., via MCP/API), (2) a deterministic scripted client, and (3) a human using `curl`/REST.
    *   *Assertion:* The Conductor, Solvent, and Executor semantics remain identical; only the client transport changes.
*   **Add Scenario E — Denial and Agent Revision (The Boundary Proof):**
    *   *Action:* Agent proposes a consequential action. Solvent denies it (e.g., due to stale snapshot or policy). Agent observes the denial, revises the `Proposal`, and resubmits.
    *   *Assertion:* The Agent never attempts to execute the denied action. The workflow loop correctly handles the "revise and retry" without the workflow layer absorbing the retry logic.

#### B. Tighten Vocabulary in Section 5 (Workflow Vocabulary)
In Section 5, the diagram uses the term `CONSEQUENTIAL PROPOSAL`. 
*   *Refinement:* Change this to just `PROPOSAL` (or `PROPOSAL / REQUEST`). 
*   *Reasoning:* The reviews specifically warned against creating a workflow-level object that sounds like an authority object. By calling it simply a `Proposal`, you reinforce that the *Agent* is just making a request. It is the *branching logic* (the consequential boundary) that elevates it to require Solvent, not the name of the object itself.

#### C. Clarify Human Authorization in Section 7 (Human Intervention)
Section 7 makes a brilliant distinction between human *review*, human *intervention*, and human *authorization*. 
*   *Refinement:* Add one sentence to explicitly state: *"If a human provides authorization, that authorization MUST still be routed through and evaluated by Solvent; a human cannot bypass the Solvent checkpoint to directly authorize an Executor."*
*   *Reasoning:* This closes a subtle loophole where a human reviewer might be assumed to have "override" powers that bypass the authority kernel.

#### D. Add a "Stale Authority" Conformance Check to Section 13
Section 13 covers denial and execution failure, but it misses one of Solvent's flagship security stories from the writeup: **Stale Authorization**.
*   *Refinement:* Add a bullet to Section 13: *"an authorization granted for a specific target/snapshot is invalidated if the target state changes before execution (stale authority rejection)."*
*   *Reasoning:* This proves that the consequential boundary isn't just a one-time "yes/no" gate, but a temporally aware checkpoint.

---

### 3. Summary Verdict

**Status: Approved with minor deltas.**

This document successfully executes the pivot from "building a workflow framework" to "formalizing a cross-system contract." It is lean, defensively architected, and perfectly scoped for the BM-IST demonstration. 

If you apply the minor additions to **Section 14 (Scenarios D & E)** and the vocabulary tweak in **Section 5**, this specification is ready to serve as the authoritative baseline for building the Conformance Test Suite.