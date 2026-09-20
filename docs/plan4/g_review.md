The smallest representation of a Conductor task is a four-block contract that decouples task management from execution while giving an arbitrary Agent the exact boundaries, proof obligations, and security classification it needs to act autonomously.

**Minimal Task Schema**

```yaml
id: "task-1042"
kind: "consequential"            # ordinary | consequential

interface:
  capability: "DeployWorkflow"   # Semantic interface name
  inputs:
    repo: "org/repo"
    ref: "main"
  operation: "deploy:workflow"   # Exact operation identity (required if consequential)

proof:
  outputs: ["run_id"]           # Expected artifacts/data
  evidence: ["e2e_tests_pass"]  # Verification gates needed for completion

dependencies:
  blocks_on: ["task-1038"]      # Hard blockers (used to compute READY frontier)
  discovered_from: "task-1000"  # Historical provenance tracking

```

**Field Breakdown & Semantic Ownership**

**1. Work Classification (`kind` & `operation`)**

* **`kind: ordinary`**: The task produces no external side-effects (e.g., refactoring code, analyzing docs, generating unit tests). The Agent executes it directly and reports the result to Conductor.


* **`kind: consequential`**: The task alters external authoritative state (e.g., executing a database migration or cloud deployment). This flag signals to any Agent that it **cannot** attempt execution directly. It must extract the exact `operation` identity and request authorization from Solvent before passing the authorization evidence to the Executor.



**2. What the Agent Can Do (`interface`)**

* **`capability` & `inputs**`: Defines the semantic interface and parameters. Arbitrary Agents (Agent X, Y, or Z) query Conductor for tasks matching the capabilities they are equipped to handle.


* **`operation`**: Binds the exact operation payload required for exact-operation binding, ensuring the authorized operation cannot be mutated during execution.



**3. What the Agent Must Produce (`proof`)**

* **`outputs`**: Lists the required data keys or artifact references the Agent must return to satisfy the task.


* **`evidence`**: Specifies the mechanical verification gates (e.g., passing unit tests, static checks, or telemetry verification) that must be satisfied before Conductor closes the task. The Agent satisfies proof obligations rather than merely asserting completion.



**4. Frontier Readiness & Provenance (`dependencies`)**

* **`blocks_on`**: Hard prerequisites evaluated by Conductor to compute the `READY` claimable frontier. If this array is empty, the task is immediately claimable.


* **`discovered_from`**: Tracks origin provenance without acting as an execution blocker. This preserves project history so a subsequent Agent can understand how the current state was reached without inheriting context windows.
