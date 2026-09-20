## deepseek

# Adversarial Review: ARGUS / Oracle Trust Verification POC

## Governing observation

The document is a briefing, not an implementation plan. It describes the intended architecture, the intended boundaries, and the intended workflow. It repeatedly says what the system "must not" do and what agents "should not" receive. But it does not say how any of these constraints are enforced. The gap between "must not" and "cannot" is the central vulnerability.

I will organize the review around concrete failure modes rather than speculation.

---

## 1. Authority-boundary failures

### 1.1 Coordinator is the single point of privileged write

The document says:

- Agents do not receive Solvent's MCP mutation tools.
- Coordinator must not write directly to Solvent or Conductor databases.
- Coordinator persists packet objects through Conductor/Solvent REST APIs.

So the Coordinator is the only entity that writes to Solvent on behalf of agents. It is delegated authority. But the document does not say how the Coordinator authenticates the packet's source, how it validates the packet's claims, or what prevents an agent from crafting a packet that the Coordinator will process into a privileged transition.

**Failure mode:** An agent submits a packet with a belief, evidence, and a discharge request. The Coordinator validates it mechanically, resolves references, and persists. If the mechanical validation checks only schema and reference integrity, the agent has achieved a privileged effect through the Coordinator.

**Question for the plan:** What exactly does the Coordinator's "mechanical Pack/evidence validation" check? Existence? Format? Relevance? Sufficiency? The answer determines whether the Coordinator is a conduit or an authority.

### 1.2 "Human adjudicates" is ambiguous

The document says:

> Agents produce work. Humans discharge epistemic debt. Solvent records/enforces authority.

But the human path is:

> Human → Trust UI → authenticated Coordinator → mechanical Pack/evidence validation → Solvent attributed discharge

Where in this path does the human actually adjudicate? If the mechanical validation can reject the human's request, the machine is adjudicating. If the mechanical validation is a formality, the human is adjudicating. The document does not resolve this.

**Failure mode:** The human clicks "discharge." The Coordinator's mechanical validator rejects it because the evidence does not meet some criterion. The human's adjudication is overridden by a machine rule. The system has silently inverted the intended authority chain.

**Question for the plan:** Can the mechanical validator reject a human discharge? If yes, on what grounds? If no, what is it for?

### 1.3 "Solvent is final authority" is asserted, not enforced

The document says Solvent determines whether a transition is permitted. But it does not say:

- Who configures Solvent's rules?
- What happens when Solvent rejects a Coordinator request?
- Whether Solvent's rules can be changed at runtime.
- Whether the Coordinator can retry, bypass, or appeal.

If Solvent's rules are configured by the same team that builds the Coordinator, Solvent's independence is nominal. If Solvent has no rejection path, it is a database, not an authority.

**Failure mode:** Solvent rejects a promotion because debt is open. Coordinator retries with a different packet. Solvent accepts. The "final authority" has been bypassed by retry logic.

**Question for the plan:** What does Solvent's rejection response look like, and what does the Coordinator do with it?

### 1.4 Agent capability boundary is a surface, not a perimeter

The document says agents receive only `argus.get_context` and `argus.submit_packet`. It says this is a capability boundary, not a prompt convention. Good.

But the document does not address whether the agent runtime has:

- Network access to Solvent or Conductor directly.
- Shell access to the host.
- Filesystem access to configuration or credentials.
- Other MCP servers that expose privileged tools.

Blocking the front door does not help if the back door is open. The document assumes the MCP surface is the only path to privileged actions. That is a strong assumption.

**Failure mode:** An agent uses a general-purpose HTTP tool to POST directly to Solvent's REST API. The MCP surface does not include this tool, but the runtime does. The capability boundary is bypassed without violating any documented rule.

**Question for the plan:** Is the agent runtime sandboxed? If so, how? If not, what prevents direct API access?

---

## 2. Missing prerequisites

### 2.1 ADD_DEBT loop is deferred

The document says:

> Phase 8 deliberately defers the operational ADD_DEBT / automatic "new evidence creates new debt" loop. Do not assume it is implemented.

This is a major gap. EBP v2.1's central doctrine includes "New evidence creates new debt." If this loop is not implemented, the system cannot enforce the doctrine it claims to follow. The adversarial agent may submit evidence that creates new debt, but the system will not automatically record it.

**Failure mode:** Adversarial agent submits a counterexample. Coordinator persists it as evidence. But no new debt is filed. The belief remains unpromoted for other reasons, or is promoted because existing debt was retired. The new evidence is recorded but not acted upon.

**Question for the plan:** If ADD_DEBT is deferred, what is the manual process for creating debt from adversarial findings? Who does it? When?

### 2.2 Physics Verifier scope is undefined

The document says the Physics Verifier produces "trusted verification artifacts" but does not say:

- What claims it can verify.
- How it verifies them.
- What its output format is.
- What "trusted" means.
- Whether it is inside or outside the trust boundary.

If the verifier is a black box that produces trusted artifacts, trust is placed in the verifier, not in the human or Solvent. But the document says Solvent remains the authority layer. The relationship is unclear.

**Failure mode:** Physics Verifier returns a "verified" artifact for Gate G0. The Coordinator treats it as evidence. The human sees the evidence and discharges debt. But the verifier's scope was narrower than the claim it was used to support. The verification is invalid but appears valid.

**Question for the plan:** What is the Physics Verifier's specification, and how is its scope documented?

### 2.3 Authentication mechanism is unspecified

The document says "authenticated Coordinator" and "authenticated attributed discharge." It does not say what authentication means. Token? Certificate? Password? MFA?

For a trust-verification architecture, the authentication mechanism is critical. If authentication is weak, the entire authority chain is compromised.

**Failure mode:** Discharge is attributed to a human via a username. Anyone who knows the username can submit a discharge. The attribution is nominal, not cryptographic.

**Question for the plan:** What is the authentication mechanism, and is the attribution cryptographically bound?

### 2.4 Final-truth detection is unspecified

The document says "No final-truth claim may be promoted." It does not say how final-truth language is detected. The EBP document says the field is owner-declared and reviewer-flaggable. In ARGUS, who is the reviewer?

If the field is owner-declared, a malicious owner can declare `containsFinalTruthClaim = false` and promote a final-truth claim. If it is reviewer-flaggable, the reviewer's identity and appeal process are unspecified.

**Failure mode:** An agent declares its claim is not a final-truth claim. The claim is promoted. The human never sees the final-truth language because it was never flagged.

**Question for the plan:** Who detects final-truth language, and what is the appeal path?

---

## 3. Incorrect assumptions

### 3.1 Adversarial agent is structurally independent

The document says the adversarial agent is a "completely separate fresh OpenCode process" with "NO conversational memory of the Work Agent." This is true at the process level.

But the adversarial agent has access to the same system state (Conductor, Solvent, RCP). It can read the work agent's packets. It can see the work agent's intermediate files if they were written to disk. The independence is procedural, not structural.

**Failure mode:** The adversarial agent reads the work agent's internal notes from a shared filesystem and uses them to construct its challenge. The challenge is not independent; it is derivative of the work agent's framing.

**Question for the plan:** What is the isolation boundary between the work agent and the adversarial agent? Process? Filesystem? Network?

### 3.2 RCP "UNKNOWN ≠ EMPTY" is enforced

The document says RCP must not make the agent believe there are no beliefs when Solvent is unavailable. This is a rule, not an implementation.

How does the agent distinguish between "Solvent says there are no beliefs" and "Solvent is unavailable"? The context projection must include a marker. The document does not specify the marker or its semantics.

**Failure mode:** Solvent is down. RCP returns a context with empty belief lists. The agent treats this as "no beliefs exist" and proceeds accordingly. The rule is violated in practice.

**Question for the plan:** What is the format of the context projection, and how does it represent unknown state?

### 3.3 Coordinator does not perform adjudication

The document says the Coordinator must not perform scientific adjudication, infer truth, or become a research reasoning engine.

But the Coordinator does perform "mechanical Pack/evidence validation." The boundary between mechanical validation and adjudication is not defined. If the Coordinator checks that evidence is relevant or sufficient, it is adjudicating.

**Failure mode:** Coordinator's validation logic includes a check that evidence "supports" the belief. This is semantic, not mechanical. The Coordinator has drifted into adjudication without violating the letter of the rule.

**Question for the plan:** What is the exact list of checks the Coordinator performs, and which of them are mechanical vs. semantic?

### 3.4 One command reliably starts the environment

Question 11 asks this. The document does not answer it. The Taskfile is described as the orchestration layer, but there is no evidence it works. The services (CockroachDB, Solvent, Conductor, Coordinator, Trust UI, ARGUS MCP, Physics Verifier) have dependencies between them. If the Taskfile does not enforce the dependency order, the "one command" will fail unpredictably.

**Failure mode:** `task up` starts Solvent before CockroachDB is ready. Solvent fails to connect. Taskfile reports success because processes started. Developer sees "all services up" and proceeds to a broken environment.

**Question for the plan:** Does the Taskfile verify readiness, not just process start? How?

---

## 4. Unenforced claims

### 4.1 "Agents cannot directly retire debt"

The document says agents cannot directly retire debt, promote a belief, or mutate Solvent. This is enforced by the MCP surface — agents only have `get_context` and `submit_packet`. Good.

But the agent can submit a packet that the Coordinator will process into a debt retirement. The agent cannot directly retire debt, but it can indirectly cause debt to be retired if the Coordinator's validation passes.

**Failure mode:** Agent submits a packet with evidence and a discharge request. Coordinator validates and persists. Debt is retired. The agent has achieved a privileged effect through the Coordinator.

**Question for the plan:** Does the Coordinator distinguish agent-submitted discharges from human-submitted discharges? If not, what prevents agent-caused retirements?

### 4.2 "Solvent records/enforces authority"

The document says Solvent is the system of record for beliefs, evidence, debt, edges, intents, promotion, retraction, authorization, and audit history. It says Solvent determines whether a transition is permitted.

But the document does not describe Solvent's promotion gate logic. It says "Solvent promotion gate" without specifying what the gate checks. If the gate is trivial (e.g., always accepts), Solvent is not enforcing authority.

**Failure mode:** Solvent's promotion gate checks only that the request is well-formed. It does not check debt or final-truth claims. Promotion is granted for any well-formed request. Solvent is a database, not an authority.

**Question for the plan:** What exactly does Solvent's promotion gate check, and how is its logic verified?

### 4.3 "Coordinator remains orchestration"

The document says the Coordinator must not become a new authority layer. But the Coordinator is the only writer to Solvent and Conductor on behalf of agents. It validates packets, resolves references, and persists objects. This is a lot of power for an "orchestration" layer.

The boundary between orchestration and authority is not defined. If the Coordinator's validation logic grows, it becomes an authority layer without anyone deciding it should.

**Failure mode:** Coordinator's validation logic includes reference resolution, which requires querying Solvent. Over time, additional checks are added. The Coordinator now enforces a subset of Solvent's rules. Solvent's authority is partially duplicated.

**Question for the plan:** What is the maximum scope of the Coordinator's validation, and how is scope creep prevented?

---

## 5. Hidden coupling

### 5.1 Conductor and Solvent share state via RCP

RCP projects state from Conductor and Solvent into a single context. The document says RCP "does NOT invent cross-system happens-before semantics between independently clocked Conductor and Solvent events."

But the agent sees a single context. If the context presents Conductor events and Solvent events as if they are ordered, the agent will infer causality that does not exist.

**Failure mode:** Context shows a Solvent belief and a Conductor task. The agent assumes the task was created because of the belief. In fact, they are unrelated. The agent's reasoning is based on a false premise.

**Question for the plan:** How does the context projection present cross-system events, and is the presentation neutral?

### 5.2 Trust UI depends on Conductor priority and Solvent state

The document says the Insights view uses "Conductor task priority/state plus Solvent epistemic state." If either source is inconsistent or lagging, the UI presents misleading information.

**Failure mode:** Conductor marks a task as high priority. Solvent has retracted the associated belief. The UI shows the task as important. The human investigates a retracted claim.

**Question for the plan:** How does the UI handle inconsistency between Conductor and Solvent?

### 5.3 Physics Verifier outputs couple to Solvent evidence

The document says the Physics Verifier's output becomes evidence/artifact material. If the verifier is tightly coupled to the BM-IST domain, its outputs may not generalize. If it is generic, its outputs may not be sufficient for BM-IST.

**Failure mode:** Verifier produces an artifact for Gate G0. The artifact is accepted as evidence. But the verifier's logic was designed for a different domain and does not actually verify G0. The evidence is invalid but appears valid.

**Question for the plan:** What is the verifier's domain scope, and how is it validated?

---

## 6. Testing gaps

### 6.1 Dry run does not test ambiguous cases

The Phase 8 dry run has two branches: dead-end (adversarial + retraction) and promotion (human discharge). It does not test:

- Partial debt retirement.
- Contested faithfulness review.
- Multi-round adversarial exchanges.
- Agent-submitted discharge that should be rejected.
- Solvent rejection of Coordinator request.
- RCP with one backend down.
- Concurrent submissions.

**Failure mode:** The dry run passes. The architecture is declared ready. The first real use case involves partial debt retirement, which the system does not handle. The failure is discovered in production.

**Question for the plan:** What is the test coverage beyond the happy path and the single failure path?

### 6.2 Tests may mock privileged paths

The document does not describe the test strategy. If tests mock Solvent and Conductor, they do not prove that the Coordinator's validation works against the real services. If tests use real services, they may be slow and flaky.

**Failure mode:** Tests pass with mocks. Production fails because the real Solvent API has different error semantics than the mock.

**Question for the plan:** Are integration tests run against real services, and are they part of the developer workflow?

### 6.3 No adversarial test of the capability boundary

The document says agents cannot access privileged tools. But there is no test that verifies this. An agent could be given a new tool by mistake, and no test would catch it.

**Failure mode:** A new MCP tool is added for a different purpose. It has side effects that touch Solvent. Agents now have privileged access. No test detects the regression.

**Question for the plan:** Is there a test that enumerates the agent's tools and asserts they are exactly `get_context` and `submit_packet`?

---

## 7. Developer experience risks

### 7.1 Taskfile may become an orchestration framework

The document says the Taskfile should "only orchestrate existing services/binaries, database readiness, configuration and tests." But this is a lot of responsibility for a Taskfile. Dependency management, readiness checks, and test orchestration are the responsibilities of an orchestration framework.

**Failure mode:** Taskfile grows to include retry logic, health checks, and service discovery. It becomes a framework. The document's claim that it is "only orchestration" is false.

**Question for the plan:** What is the Taskfile's maximum scope, and how is scope creep prevented?

### 7.2 One-command startup may hide failures

If `task up` starts all services and reports success when processes start (rather than when services are ready), developers may proceed with a broken environment.

**Failure mode:** CockroachDB takes 30 seconds to become ready. Solvent starts immediately, fails to connect, and exits. Taskfile reports success. Developer sees "all services up" and proceeds. The environment is broken.

**Question for the plan:** Does `task up` verify readiness, and what does it do if a service fails?

### 7.3 Developer experience work may precede architecture freeze

The document says the remaining work is developer experience. But the architecture has unresolved questions (authority boundary, capability boundary, final-truth detection, etc.). If the tooling is built before the architecture is settled, it may need to be rewritten.

**Failure mode:** Taskfile is built for the current architecture. Architecture changes. Taskfile is obsolete. Developer experience work must be redone.

**Question for the plan:** Is the architecture frozen, or is it still open? What is the freeze criterion?

---

## 8. Governance and process gaps

### 8.1 No process for changing Solvent's rules

The document says Solvent determines whether a transition is permitted. It does not say who changes Solvent's rules or how changes are audited.

**Failure mode:** A developer changes Solvent's promotion gate to accept a new class of requests. The change is not reviewed. Solvent's authority is silently weakened.

**Question for the plan:** What is the change-control process for Solvent's rules?

### 8.2 No process for handling Solvent- Coordinator disagreement

If Solvent rejects a Coordinator request, what happens? Does the Coordinator retry? Escalate to a human? Log the rejection?

**Failure mode:** Solvent rejects a promotion because debt is open. Coordinator retries with a different packet. Solvent accepts. The rejection is lost.

**Question for the plan:** What is the disagreement-handling process?

### 8.3 No process for RCP failure

If RCP cannot reach Solvent, it must return UNKNOWN, not EMPTY. But how does the agent respond to UNKNOWN? Does it retry? Escalate? Proceed with partial context?

**Failure mode:** Agent receives UNKNOWN context. It proceeds as if the context is empty. The rule is violated in practice.

**Question for the plan:** What is the agent's prescribed behavior on UNKNOWN context?

---

## 9. What the document gets right

To be fair, the document correctly identifies:

- The central invariant (CAPABILITY ≠ WORK ≠ AUTHORITY ≠ EXECUTION).
- The need for capability boundaries enforced by tools, not prompts.
- The separation of orchestration (Coordinator) from authority (Solvent).
- The need for human adjudication of epistemic debt.
- The need for adversarial review as a first-class role.
- The risk of developer tooling becoming a framework.
- The need for UNKNOWN ≠ EMPTY semantics.
- The deferral of ADD_DEBT as a deliberate scope decision.

These are the right priorities. The gaps are in enforcement, not in intent.

---

## 10. Concrete questions the plan must answer before implementation

1. What exactly does the Coordinator's "mechanical Pack/evidence validation" check?
2. Can the mechanical validator reject a human discharge? If so, on what grounds?
3. What does Solvent's promotion gate check, and how is its logic verified?
4. How does the agent runtime prevent direct API access to Solvent and Conductor?
5. How is final-truth language detected, and who is the reviewer?
6. How does RCP represent UNKNOWN vs. EMPTY, and what does the agent do with UNKNOWN?
7. How is the Physics Verifier's scope defined and validated?
8. What is the authentication mechanism for human discharge?
9. How does the Taskfile verify readiness, not just process start?
10. Is the architecture frozen, and what is the freeze criterion?

---

## 11. Bottom line

The document is a well-intentioned briefing on a trust-verification architecture. It correctly identifies the boundaries that must be enforced and the roles that must be separated. But it does not say how the boundaries are enforced or the roles are prevented from collapsing.

The central vulnerability is the gap between "must not" and "cannot." The document says agents must not access privileged tools, but does not say how this is enforced. It says the Coordinator must not become an authority layer, but does not say what prevents scope creep. It says Solvent is the final authority, but does not say what its rules are or who configures them.

The dry run is a useful first test, but it does not exercise the failure modes that matter most: partial debt retirement, contested faithfulness, Solvent rejection, RCP failure, and concurrent submissions. Until these are tested, the architecture's claims are aspirational.

The developer experience work is important, but it may be premature if the architecture is not frozen. The Taskfile may become a framework, and the "one command" may hide failures rather than reveal them.

None of these are fatal. All are fixable. But they must be addressed before the architecture is declared ready for the BM–IST research work. The infrastructure should not become the complexity.


## Z

# ARGUS — INDEPENDENT ADVERSARIAL ARCHITECTURE REVIEW

Noted as "for your info only" — but the briefing is written as a review request with an explicit checklist (§14–15), so I'll deliver exactly that: concrete failures over speculative expansion, in the format the document itself demands. Overall assessment first, then findings, then the mapped checklist, the minimum honest test suite, one end-to-end red-team trace, and the BM-IST alignment notes.

**Overall:** the authority spine is correct and the separation doctrine (CAPABILITY ≠ WORK ≠ AUTHORITY ≠ EXECUTION) is genuinely enforced at the MCP layer rather than by convention — that is rarer than it should be. The failures below cluster in three places: **the edge model leaks authority through agent-authored content** (the worst one), **the Trust UI human-adjudication path has the standard localhost attack surface**, and **the Physics Verifier's "trusted artifacts" have no defined trust chain**. All are fixable without architectural growth.

---

## 1. WHAT IS ARCHITECTURALLY SOUND

- Tool-boundary as capability boundary ("if an agent cannot see a tool, it cannot invoke it") — correct mechanism, not prompt hygiene.
- UNKNOWN ≠ EMPTY in RCP — the single most important honesty rule in the read path; most systems fail exactly here.
- Solvent as final authority with Coordinator explicitly forbidden from becoming a shadow store; Conductor/Solvent state explicitly non-equivalent ("task accepted" ≠ "claim promoted").
- Retraction cascade + "retracted belief is NOT later promoted" as a hard state-machine rule.
- EBP's debt model is strictly stronger than the verdict-only ledger schema this program sketched earlier — adopting it is an upgrade, provided the deferred ADD_DEBT loop stays visibly deferred (it does, per §2 — good).

---

## 2. FINDINGS

Severity: **P0** = would invalidate the dry run's claim to demonstrate the invariant; **P1** = real defect, fix before live use; **P2** = before scale/external users.

### F1 — Agent-authored edges are live authority-mutating content. **[P0]**

`POST /v1/beliefs/{id}/edges` creates edges that Solvent's retraction logic already consumes. A `contradicts` edge authored by an agent is therefore an agent-initiated trigger for a human-driven retraction cascade. The human "adjudicates" — but under queue pressure the human is a rubber stamp on agent-manufactured epistemic content.

*Trace:* prompt-injected or buggy Work Agent submits a packet containing `contradicts` edges against the load-bearing G0 beliefs → edges live in Solvent → Adversarial-Agent/Branch-A flow leads a human to "confirm" the dead end → RetractCascade cancels tasks. Every mutation was human-clicked; every mutation was agent-authored.

*Fix:* edges arriving via packets enter as `PROPOSED`; only human-confirmed (or rule-validated + human-discharged) edges become `ACTIVE`; retraction is legal only from `ACTIVE` edges; the UI must render "proposed by agent ⟨id⟩, awaiting human confirmation." This is the same doctrine as evidence≠authority applied one level down — edges *are* evidence.

### F2 — Trust UI human path: localhost CSRF/DNS-rebinding and body-supplied attribution. **[P0]**

§14-Q4 and Q5 both concentrate here. If discharge/retraction are loopback POSTs with cookie-or-no auth: a malicious webpage can reach `http://localhost:…` via DNS rebinding (attacker domain resolves to 127.0.0.1; requests are then same-origin, cookies flow), triggering a retraction cascade or a false discharge. And if the POST body carries `discharged_by: "alice"`, attribution is forgeable — violating Solvent's own §13 (actor type is not identity).

*Fix:* principal identity exclusively from the auth token; per-action confirmation token rendered by the UI and required in the mutation; `Origin`/`Host` allowlist on all Coordinator mutation endpoints; step-up auth on the two consequential actions (discharge, retraction).

### F3 — Physics Verifier artifacts have no trust chain. **[P0 for the epistemic claim, P1 otherwise]**

"Trusted verification artifacts" is currently a label. Four concrete gaps:

1. **Claim↔script binding:** who guarantees the executed script tests the claim it is attached to? An agent that can invoke the verifier chooses what to run — evidence laundering (run a check of a weaker statement, attach output to the stronger claim).
2. **Numeric reproducibility:** the natural checks (variational finite differences, one-step Madelung comparisons) are floating-point and *not* bit-reproducible across platforms. Without a declared tolerance in the artifact schema, two runs can disagree and the ledger has no adjudication rule — the evidence is weaker than it looks.
3. **Environment pinning:** artifact must record code hash, library versions, seed, host class.
4. **Invocation path:** if the verifier is a third agent tool, agents farm it. Cleaner: agents submit artifact *references*; Coordinator invokes the verifier during packet validation and binds the result. This also keeps the agent MCP surface at exactly two tools, as §9 promises.

*Fix:* artifact schema = {claim_id, script_hash (pre-registered per obligation), raw output, exit code, tolerance, env hash}; retirement rules accept only hash-matching artifacts; verifier invoked by Coordinator, not by agents.

### F4 — `submit_packet` is a mutation gateway; the plan's "agents cannot mutate Solvent state" is overstated. **[P1]**

Agents create beliefs, evidence, edges, and tasks on every submission. The true invariant is narrower and should be stated precisely everywhere: *agents can create proposals; only humans can discharge, confirm, promote, or retract.* Without that restatement, the security model claims more than the mechanism delivers — and F1 is the proof. Add: per-agent quotas and content-hash dedup (a runaway or injected agent can otherwise flood the epistemic record and DoS the Debts queue).

### F5 — Dual-write packets with no saga. **[P1]**

One packet writes to both Solvent (beliefs/evidence/edges) and Conductor (tasks). Solvent commits; Conductor times out → partial packet. Retry without an idempotency key duplicates beliefs (inflating apparent evidential support — a corrupted record, not just a mess).

*Fix:* client-supplied idempotency keys; content-hash unique indexes in Solvent (their own §22 doctrine: structural truth); compensation or explicit orphan-marking; a reconciliation check in `task verify` (every evidence references an existing belief; every packet task exists).

### F6 — RetractCascade service coupling. **[P1]**

If Solvent calls Conductor to cancel tasks, the authority service owns workflow semantics and a Conductor outage strands a half-applied cascade (claim retracted, task still accepted — Insights now lies). Apply Solvent's own doctrine internally: retraction is recorded in Solvent; task cancellation is Coordinator-mediated follow-through with retry and a reconciliation status surfaced in Insights ("cascade pending: 1 task").

### F7 — Partial backend failure ≠ the UNKNOWN case you covered. **[P1]**

RCP handles "Solvent down ⇒ UNKNOWN." The common failure is *partial*: Solvent up, Conductor down. The fresh agent reconstructs half a context and proceeds — misaligned with workflow state. *Fix:* per-source status in `get_context`; overall `DEGRADED` flag; agent contract: degraded context ⇒ bounded work only, no consequential submissions. Same class, opposite direction: **TRUNCATED ≠ COMPLETE** — if the broad scenario projection hits a size cap silently, the agent believes it saw all debts ("no unresolved obligations found"). Explicit truncation markers and counts, or the UNKNOWN≠EMPTY rule is honored at the top and violated at the bottom.

### F8 — Agent freshness is procedural, not enforced. **[P1]**

"Completely separate fresh OpenCode process with no conversational memory" is a methodological claim the harness doesn't yet enforce. Both agents in one working directory share files, caches, env — the adversarial agent can read the work agent's notes and the "independence" branch is theater.

*Fix:* container per agent; no shared mounts; the only I/O channel is the two MCP tools; record container/image hash in the packet provenance so independence is auditable after the fact.

### F9 — Coordinator must not accrete retirement-rule code. **[P1]**

"Mechanically validate retirement rules" — if the rules live as Go code in Coordinator, Coordinator becomes a shadow authority layer (violating §5's own rule). Rules as declarative config (policy data, validated at runtime), referenced by Solvent debt vocabulary. Phase-8 Gate-G0-specific rules are fine — as config, flagged as such.

### F10 — Evidence typing is missing. **[P1, BM-IST-specific]**

A proof, a CAS artifact, and a citation are different evidence types with different validation rules. The verifier is explicitly "not a theorem prover" — so the retirement rule for L1/L2 (proof-obligations) *cannot* be mechanically satisfiable and must require human discharge of the proof itself; the verifier can only supply numeric sub-artifacts. If the schema has one evidence type, either proofs get validated like computations (meaningless) or everything does (worse).

### F11 — RCP list ordering invites the happens-before hallucination back in. **[P2]**

Deterministic presentation + no invented cross-system ordering is stated — but agents *will* infer causality from list position. Sort by an explicit documented key, tag every item with (source, sequence), and state that ordering carries no temporal semantics.

### F12 — Cross-task/cross-scenario references. **[P2]**

Packet grammar allows `canonical:belief:<uuid>`. Without a scope check, an agent attaches evidence to another scenario's beliefs. Enforce scenario containment at packet validation (Solvent already has scenario isolation — use it structurally).

### F13 — DX/Taskfile: semantic readiness, migration idempotency, reset discipline. **[P2]**

`task up` must gate on *semantic* readiness (health endpoint that pings the DB, not TCP-listen), or first requests 500 during schema readiness. Two migration chains (Solvent, Conductor) must be idempotent on re-run. `task down` must not destroy data by default; `task reset` is the explicit destructive op. One config source; no per-service URL drift.

### F14 — Documentation drift. **[P2]**

The briefing expands BM-IST as "**B**ounded **M**athematics — **I**nvariant **S**tructure **T**heory." Our program is *Bohmian Mechanics – Invariant Set Theory*. Trivial here, exactly the genre §14's "misleading documentation" item exists to catch — fix before the name fossilizes.

### F15 — Promotion wording and the promotion→retraction path. **[P2]**

UI must never say "verified true" — promoted means debt-clear-as-of-now. State machine must allow promotion → later retraction (new evidence), and Branch-A dry run must *assert* the refusal case (promote after retract ⇒ rejected).

---

## 3. THE §14 CHECKLIST, ANSWERED COMPACTLY

| Q | Answer | Finding |
|---|---|---|
| 1 API reuse vs DB bypass | Assume yes, but verify no service shares CRDB credentials; scoped creds per service | — |
| 2 Coordinator ≠ authority | At risk via hardcoded retirement rules | F9 |
| 3 Tool-constrained capabilities | Yes for MCP surface; verifier invocation path undecided | F3, F4 |
| 4 Browser bypass of adjudication | **Yes, currently plausible** (rebinding/CSRF; body attribution) | F2 |
| 5 Caller claims vs persisted evidence | Attribution must come from token, not body; artifact binding missing | F2, F3 |
| 6 RCP ≠ new database | Holds only if derived views are projections, never materialized | — |
| 7 UNKNOWN ≠ EMPTY | Holds; **partial** and **truncated** cases uncovered | F7 |
| 8 No invented cross-ledger ordering | Stated; list-order inference still leaks it | F11 |
| 9 Implemented vs deferred EBP | Cleanly declared (ADD_DEBT deferred) — keep it visible in RCP output too | — |
| 10 Tooling reduces complexity | Taskfile-only scope is right; the risk is F9 rule-accretion | F9 |
| 11 One-command start | Achievable; gated on semantic readiness + migration idempotency | F13 |
| 12 Fresh agents against it | Yes — **if** isolation is enforced, not assumed | F8 |
| 13 Verifier output ≠ authority | Correctly positioned; trust chain undefined | F3, F10 |

---

## 4. MINIMUM HONEST TEST SUITE (`task verify` / dry-run assertions)

Contract tests against the *real* services (no mocks) must include failure injection:

1. Solvent 500 mid-packet → orphan detection/reconciliation fires.
2. Conductor down → `get_context` returns per-source UNKNOWN/DEGRADED; agent contract: stop.
3. Same idempotency key resubmitted → exactly one belief (content-hash unique index).
4. Malformed/dangling `canonical:` reference → 4xx, nothing persisted.
5. Retraction requested against `PROPOSED` edge → refused; after human confirmation → allowed; cascade cancels exactly the linked task; audit shows agent-proposed → human-adjudicated → Solvent-mutated.
6. Discharge with open debt → refused; with hash-matching pre-registered artifact → allowed; attribution from token, body field ignored/verified.
7. Promote a retracted belief → refused.
8. Cross-origin POST to localhost mutation endpoints (rebinding simulation) → rejected by Origin/Host check + missing confirmation token.
9. Artifact re-run in pinned env within declared tolerance → pass; outside → flagged, not silently accepted.

**Branch A's real value is negative-assertion:** the dry run passes only if (a) the agent's contradicts edge did *not* auto-retract, (b) the cascade hit exactly one task, (c) retracted-then-promote is rejected, (d) the audit trail is complete and attributed at every hop.

---

## 5. ONE END-TO-END RED-TEAM TRACE (the design's exam question)

A single compromised Work Agent attempts: (1) flood beliefs → **stopped by quotas + dedup**; (2) manufacture a dead end via contradiction edges → **stopped by F1 PROPOSED gate**; (3) forge a discharge → **stopped by token attribution + confirmation token**; (4) launder verifier evidence → **stopped by script-hash pre-registration (F3)**; (5) bypass via direct DB writes → **stopped by scoped credentials**. Five attack surfaces, five mechanisms — that is the demonstration Phase 8 should be able to claim, and currently can only claim (5) by default.

---

## 6. BM-IST ALIGNMENT NOTES

- **Convergence:** ARGUS beliefs = ledger claim atoms; `contradicts` edges = our `kills` edges; RetractCascade ≈ our ledger propagation rule; human discharge = our verification pass; promotion = SETTLED. The debt model supersedes my Round-7 verdict-only schema — adopt it. Carry the **claim-id ↔ Solvent-belief-id join key** in both directions so the research ledger and the authority store stay one record, not two.
- **Slice choice (G0/L1/L2) is right** — it matches Implementation-Plan P0/P1, and L1/L2 are `PROVED IN PROGRAM; write-up owed`, so the first honest debts are exactly "write the note." Their retirement rule must be typed `PROOF_OBLIGATION` (F10): human discharges the proof; the verifier supplies only numeric sub-artifacts.
- **When the slice widens:** the verifier's first pre-registered scripts should be the G-1 self-audit set (shear counterexample reproduction; G2 Madelung cancellation) — the program's flagship self-falsification checks, now with a trust chain, which is precisely what they've lacked.
- **§37.9 mapping:** external-AI output auto-tagged ASSUMED is implemented here as: agent packets are proposals, human adjudication is the audit pass. The architecture already encodes the rule; the missing piece is only the named-object/citation checker inside packet validation (Checker-CITE as a Coordinator validation step, producing typed evidence, never authority).

---

## 7. CLOSE

The smallest honest system your §15 asks for is achievable with the current component list and **four pre-live repairs**: the edge PROPOSED/ACTIVE gate (F1), the Trust UI auth/attribution hardening (F2), the verifier artifact trust chain (F3/F10), and precise restatement of the mutation invariant (F4). Everything else is P2 hardening that can ride with the dry run. The architecture's own doctrine — *prefer structural truth over labels* — is the right test for every finding above, and by that test the system is one gate away from being able to claim what it is named for: that agent work, adversarial pressure, human adjudication, and Solvent authority are not just separated on a diagram, but separated in the request path.