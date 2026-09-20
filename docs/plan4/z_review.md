Your message cut off after "Please help me:" again — but this time the operative request is stated at the top, so I'll answer the design question directly, and do it the way this series works: propose the smallest schema, then attack it against the established lineage. One finding up front matters enough to lead with it: **the transcript's proposed task schema contains a field that quietly violates the matrix, and fixing it is what produces the smallest representation.**

## The adversarial finding that shapes the answer

The transcript's example task carries:

```text
Work type: consequential
Authority: Solvent
Operation: deploy:repo:workflow:ref:run_id
```

A stored `kind`/`work_type` flag re-imports the very problem the matrix spent three versions closing. Matrix §9.1: agent classification of consequence is **advisory**; the **capability declaration** defines effect classification, and Conductor "does not become a domain-specific policy router." If the task schema stores a consequential flag, then either (a) the creating Agent set it — advisory classification promoted to stored state, the ordinary-work bypass at the schema level — or (b) Conductor set it, making Conductor a classifier. Both are prohibited. The same task also shows `READY_FOR_AUTHORIZATION` as a status — that is workflow-phase-as-persisted-state, the exact thing invariant 11 forbids.

The fix *is* the subtraction: **don't store the class. Store the capability reference; derive the class from the referenced declaration.** One field instead of two, and the classification authority stays exactly where the matrix put it.

## The smallest task representation

```text
Task
  task_id           # identity + correlation (existing ID discipline)
  objective         # free text, Agent-authored: what this work is
  capability_ref    # → declared capability  (sole classification source)
  inputs            # references that must exist to start
  outputs           # references that must exist to submit
  verification_ref  # who/what accepts the outputs (domain stays external)
  deps              # task_ids that must complete first
  claim             # assignee + timestamp (atomic, exclusive)
  status            # Conductor lifecycle ONLY — kind-agnostic
```

Derived, never stored:

- `READY` = deps satisfied ∧ unclaimed ∧ open (the Beads frontier, computed not persisted)
- **consequence class = effect-classification of `capability_ref`'s declaration**
- `parent` / `discovered-from` = history edges (Beads cherry-pick), not fields — they're provenance, and the matrix already separates action lineage from scheduling

It answers the three questions by construction:

| Question | Answered by |
|---|---|
| What can I do? | `objective` + `inputs` + the capability's contract; `claim` grants **coordination ownership only, never authority** (assignment ≠ authority, per the workflow spec §5.3) |
| What must I produce? | `outputs` + `verification_ref` — the acceptance authority is referenced, never embedded (domain acceptance stays external) |
| Ordinary or consequential? | Look up `capability_ref`'s declaration. Conductor never answers this from its own judgment — it's a projection of declared fact |

Note what's *not* in the schema: no `operation` field. The operation identity definition lives in the declaration (§18); the exact operation — including `run_id` — is instantiated at proposal time per Phase 1's own identity format. Duplicating the template in the task creates a second definition that will drift from the declaration's.

## The two workflow renderings fall out unchanged

```text
Ordinary:       claim → work via declared non-effect capability
                        → submit(outputs + evidence) → verification

Consequential:  claim → construct exact operation per declaration identity
                        → Solvent authorize → Executor execute
                        → submit outcome evidence → verification
```

Conductor's lifecycle is identical in both — `OPEN → ACTIVE → SUBMITTED → DONE/REJECTED`. The authorization state is a **Solvent-owned fact referenced by evidence**, never a Conductor status. Task *creation* is always ordinary coordination; registering a task that references an effect-capable capability proposes nothing — proposing the exact operation is the consequential act.

## Protocol corrections to the transcript

Small but worth locking: the proposed Consequent-protocol ops `PROPOSE / AUTHORIZE / EXECUTE / OBSERVE` mix planes. `AUTHORIZE` and `EXECUTE` are Solvent- and Executor-owned operations — they are not Conductor protocol verbs. The Conductor work protocol is `CREATE / DISCOVER / CLAIM / RELEASE / UPDATE / BLOCK / SUBMIT / ACCEPT / REJECT` plus recording coordination-relevant results. If "OBSERVE" means *Conductor records evidence references from the consequence plane*, keep it; if it means Conductor polls or tracks authorization state, it's shadow authority bookkeeping — don't.

## What this schema quietly requires

Two dependencies the transcript doesn't surface, both already on the punch list — now promoted from "spec hygiene" to "this design cannot ship without them":

1. **H1 (identity comparison semantics).** "Arbitrary agent constructs the exact operation per declaration" is only deterministic if the declaration defines equality, not just fields. Otherwise Agent X and Agent Y construct syntactically different but "equal" identities and binding verification diverges — the Phase 5 substitution failure arriving early.
2. **H2 (default declaration owner).** The claim-any-task vision means every `capability_ref` must resolve to a real declaration — including generic capabilities (shell, file edit, web fetch) that a coding runtime uses as ordinary work. The matrix's Boundary Rule already requires this ("a tool or service used during ordinary work MUST still be covered by a capability declaration"); the task schema makes it the load-bearing path of the whole system. Until the default-owner clause exists, `capability_ref` dangles exactly where the vision is most ambitious.

## Deliberately excluded (the subtraction record)

Priorities/urgency, risk scores, assignee preconditions, embedded acceptance criteria (reference them), per-task evidence schemas (the protocol already mandates the evidence chain — §14/§16 — the task adds nothing), epics/milestones, lease/heartbeat machinery, capability-matching policy ("which tasks can *I* do?" beyond READY ∧ declaration-resolvable).

Deferred with named phases, not forgotten: **stale-claim recovery** (agent dies holding a claim — needs a timeout/release semantic, Phase 3's contention territory), claim policy beyond atomic exclusivity, priority ordering of the READY frontier, and multi-agent discovery filters.

## The one-sentence form

**A Conductor task is a claimable contract of references — objective, capability, inputs, outputs, verification, dependencies — where every semantic that could be derived is derived, and the one semantic that must never be derived from the Agent (consequence class) is derived from the one thing allowed to know it: the declaration.**

That's small enough to write code against. Settle H1 and H2 first — they're now prerequisites of this schema, not of some future phase.
