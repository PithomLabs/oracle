## prompt

Are we risking of over bookkeeping of the context graph instead of focusing on the actual computation/logic/derivation about the problem domain, not too much accounting per EBP 2.1


## answer

Yes. **There is a real risk of over-bookkeeping now.** In fact, I think this is the most important architectural risk to watch at this stage because EBP 2.1 explicitly says:

> **Accounting must never become the work.**

The good news is that the work so far has mostly been **load-bearing accounting**, not vanity accounting. But we are approaching the boundary.

### The distinction I would use

The context graph should maintain only information that changes one of these:

```text
1. Can an agent reconstruct the current research state?
2. Can a human determine why a claim is trusted/challenged?
3. Can the system prevent unauthorized epistemic transitions?
4. Can we reproduce what happened when a dispute occurs?
```

Anything else should be deferred.

So:

```text
belief immutability
        → YES, essential

agent → packet → belief provenance
        → YES, essential

contradicts / derives edges
        → YES, essential

atomic packet persistence
        → YES, essential

agent identity
        → YES, useful provenance

packet validation
        → YES, essential trust boundary
```

But things like:

```text
persist a disposition record for every belief
context snapshot versioning
extensive evidence lineage UI
epistemic-kind taxonomy
agent attestation
graph visualization
automatic review coverage tables
```

are **not yet proven necessary**.

### The strongest warning sign

Look at where we've gone recently:

```text
Why can't the agent see edges?
→ add edges to context.

Who created this belief?
→ add packet provenance.

Who created the edge?
→ add edge provenance.

Who created the task?
→ add task provenance.

Did the adversarial agent review every belief?
→ maybe add coverage.

Which exact context snapshot did it see?
→ maybe add snapshot ID.

What epistemic stratum is this statement?
→ maybe add kind.
```

Every individual question is reasonable.

But taken together, you can easily end up building:

```text
              RESEARCH
                 |
         ┌───────┴────────┐
         ↓                ↓
      actual work      accounting
                          |
                     accounting
                          |
                     accounting
                          |
                       ...
```

and the system starts optimizing for **describing research rather than doing research**.

That would violate the point of ARGUS.

### Your recent adversarial run gives us the answer

The adversarial agent produced genuinely useful domain criticism:

```text
G0 may be attacking a straw-man formulation
the actual program architecture is different
the flow is not yet constructed
...
```

That is the value.

Notice what actually advanced the science:

```text
reading v6.3
+
reconstructing claims
+
understanding the architecture
+
finding a mismatch
+
producing counterclaims
```

Not:

```text
packet count
origin metadata
submission timestamps
UI columns
```

Those are support mechanisms.

So the correct priority is:

```text
DOMAIN REASONING
    >
EPISTEMIC CONTROL
    >
PROVENANCE
    >
OBSERVABILITY
    >
CONVENIENCE
```

### I would actually freeze the bookkeeping after Plan 13

Once the MCP validation bypass is fixed and the provenance spine is verified, I would establish a temporary rule:

> **No new research-accounting feature without a demonstrated dry-run failure.**

That means don't proactively implement:

```text
review coverage table
context snapshots
epistemic kind
agent attestation
more provenance fields
richer UI lineage
new graph queries
```

unless the next actual Work → Adversarial → Work cycle demonstrates the need.

### A very useful design test

For every proposed piece of metadata, ask:

> **What research decision becomes possible that was impossible without this field?**

If the answer is:

> "The UI becomes nicer."

defer it.

If the answer is:

> "We can prove which claim was challenged and by whom."

keep it.

If the answer is:

> "We can prevent an unauthorized state transition."

keep it.

If the answer is:

> "We might someday want to analyze agent behavior."

defer it.

### And I would change the agent mindset too

The agent should not think:

```text
First: account for everything.
Then: do research.
```

It should think:

```text
First: reconstruct enough state to avoid being wrong about what has already
been done.

Then: do the actual reasoning.

Finally: leave enough structured residue that another agent/human can
continue or challenge it.
```

That ordering matters enormously.

### The sweet spot

The ARGUS ledger should be **lossless enough for trust, but too small to become the work**.

In practice:

```text
                MINIMUM SUFFICIENT ACCOUNTING
                           |
        ┌──────────────────┼──────────────────┐
        ↓                  ↓                  ↓
    provenance          relations         authority
        |                  |                  |
   who/what made it    derives/contradicts   who may act
```

Then get out of the way and let the agent think.

So yes: **your concern is justified, and I would actively put a brake on the bookkeeping now.** The current provenance/validation work is justified because it closes concrete integrity failures. After Plan 13, the next milestone should be a substantive research cycle, not another schema expansion.

The best test of whether we've crossed the line is simple:

> **Does the next Work Agent spend most of its effort deriving, checking, and attacking the BM–IST–AS program—or reconstructing ARGUS metadata?**

It should be overwhelmingly the former.


## prompt

while waiting, does review have net valid that we need to consider per over accounting guidelines? Endorsed in full — and this is the correct application of the project's own doctrine at exactly the moment it's most needed. A few things worth adding, because this answer is about to become policy and the policy should be complete.

## The warning sign you cited has a name and a precedent

The list — edges→context, packet provenance, edge provenance, task provenance, coverage, snapshots, kind, attestation — each individually justified, collectively a schema industry. That pattern already occurred once in this program, at larger scale: the specification era itself, seven rounds of boundary refinement that ended only when the project declared "the architecture is done; the remaining findings are wiring." The brake you're proposing now is the same pivot at the meta level. Worth saying in the policy text so the pattern is named: **the last time this question came up, the answer was "prove it empirically" — and that was correct then too.**

## One strengthening: the freeze needs scope, not just timing

"Freeze after Plan 13" is the right instinct; make the freeze condition-shaped rather than date-shaped, or it erodes the first time someone has a genuinely good idea:

```text
BOOKKEEPING FREEZE (effective after Plan 13 merges)

No new research-accounting capability (fields, tables, projections,
columns, coverage mechanisms, snapshot machinery) unless ALL of:

  1. A dry-run/work-cycle FAILURE is cited — a concrete decision that
     was impossible, a dispute that couldn't be resolved, an
     unauthorized transition, an unreconstructible state. Not a
     "would be nicer," not an agent suggestion, not a review insight.
  2. The failure is classified: trust/integrity (rarely deferred),
     research-blocking (fix now), or convenience (deferred to the
     post-POC ledger).
  3. The smallest mechanism that fixes the failure is specified
     before any schema is touched.

Standing deferrals carried with triggers (unchanged): coverage
persistence, context snapshots, epistemic-kind field, attestation,
priority derivation.
```

This is the disposition-log discipline applied forward: every past finding already has a disposition; now every *future* accounting proposal gets one at proposal time, with evidence. It also settles the currently-open items cleanly — the origin-anchor fix and Plan 13 proceed (they close demonstrated failures: provenance spine, MCP bypass), everything else on the deferred list stays deferred.

## Two additions to the decision test

Your four-question test is good; two more that catch cases it misses:

- **The token question:** *"If this field were deleted tomorrow, would any past or future research conclusion change?"* Accounting that changes nothing when deleted is decoration. (Belief immutability: catastrophic. Origin anchor: disputes become unresolvable. Coverage table: nothing. Epistemic kind: one prose line lost.)
- **The agent-effort question, run in both directions:** you asked whether the agent spends its effort deriving vs. reconstructing. Also ask whether the *bookkeeping itself* consumes agent effort — disposition records, per-belief coverage lines, snapshot references all cost the agent tokens and attention mid-reasoning. The cheapest accounting is the kind the Coordinator computes (provenance anchors, hashes, validation) and the most expensive is the kind agents must author. Weight deferred proposals accordingly: **prefer machinery-computed provenance over agent-authored provenance** wherever both would do.

## The measurement you proposed, operationalized

"Does the next Work Agent spend most of its effort on the domain?" is testable, and given the dry run's pattern, worth testing cheaply: the role card already instructs reconstruct-then-reason. At the next run, check where the agent's packet mass went — claims/evidence/edges vs. process/governance text. The category error found in the last run was exactly domain reasoning leaking into governance stratum; the inverse leak (process accounting crowding out domain work) is the failure mode your freeze guards. Two numbers, one line in the run report, evidence instead of vibes.

## The one counterweight, stated so the brake doesn't over-correct

The freeze governs *accounting*. It must not quietly extend to the research program itself — v6.3's debt classes, gates, and pre-registrations are not bookkeeping; they're the actual epistemology of the domain, and they keep moving (substrate-class budget, Valentini fork pre-registration, null models are all still owed). The distinction to write down: **EBP obligations on claims = the work. Infrastructure obligations on the ledger = accounting.** The freeze binds only the second. Otherwise "less is more" becomes an excuse to under-populate the ledger with real scientific debts, which would be the opposite failure.

And the closing calibration, because EBP itself supplies it: the same document that says "accounting must never become the work" also says the accounting exists *so that* promotion never happens silently. The current provenance push existed because a trust loop without an origin spine can't tell earned promotion from accumulated assertion. That load-bearing accounting is now nearly complete — which is exactly why the brake arrives on time rather than late. After Plan 13: **the ledger is frozen, the agents get to work, and the next version number that matters is the research program's, not the schema's.**



## answer

Yes. **Most of that review is net-valid and worth incorporating**, but I would tighten it in a few places so the bookkeeping freeze does not become another abstraction layer.

### What I would adopt

**1. Make the freeze condition-shaped, not date-shaped.**

This is the strongest addition.

Use:

```text
BOOKKEEPING FREEZE

Do not add new research-accounting capability unless:

1. A concrete dry-run failure is identified.
2. The failure is classified as:
   - trust/integrity
   - research-blocking
   - convenience
3. The smallest mechanism that fixes that failure is specified before
   schema or infrastructure changes begin.
```

That is better than "freeze after Plan 13," because the actual sequence has already moved beyond Plan 13 into the provenance/P0 remediation. I would make the trigger **after the current remediation and its full DB-backed verification pass**, not a particular plan number.

This is consistent with the evidence we've accumulated: the packet-validator work, provenance spine, and atomic-persistence repairs were all responses to concrete failures, not speculative improvements.

### 2. Keep the freeze strictly about infrastructure accounting

This is essential.

The freeze must **not** constrain the research methodology itself.

Keep:

```text
EBP debt on claims
research gates
null models
counterexamples
faithfulness reviews
substrate tests
pre-registered forks
```

because these are **the research**.

Freeze:

```text
new provenance tables
new metadata fields
new review ledgers
new snapshots
new UI bookkeeping
new agent telemetry
```

unless a demonstrated failure requires them.

The distinction is:

```text
research obligations = epistemology
ledger obligations    = accounting
```

That is probably the most important sentence to preserve.

### 3. Adopt the "agent effort" principle

This is also strong:

> Prefer accounting computed by the system over accounting the agent has to author.

I would make this a design preference:

```text
Prefer:
Coordinator-generated provenance
Coordinator-generated hashes
Coordinator-generated timestamps
Coordinator-generated relationships

Avoid:
Agent-authored provenance summaries
Agent-authored coverage ledgers
Agent-authored bookkeeping records
```

The agent's scarce attention should go to:

```text
reasoning
derivation
counterexample search
evidence evaluation
```

not maintaining the ledger.

This also fits the system's current direction: `origin_packet_id` is computed mechanically from the packet rather than asking the agent to separately declare "this belief was created by packet X."

### 4. Adopt the token/cost question, but soften it

The proposed:

> "If this field were deleted tomorrow, would any past or future research conclusion change?"

is a useful heuristic, but I would not make it an absolute test.

Some accounting fields do not change the scientific conclusion, but are essential for **trusting or adjudicating** the conclusion.

For example:

```text origin_packet_id
```

doesn't change the math, but deleting it could make a dispute about authorship/provenance unresolvable.

So use:

> **If this field were deleted, would we lose a necessary capability for research correctness, reproducibility, authority, or adjudication?**

That is a better test.

### 5. Measure agent effort without adding telemetry

The suggested measurement is good, but don't build a new "agent effort accounting" system.

For the next dry run, simply inspect:

```text
Work packet:
  domain claims/evidence/derivations
  vs.
  process/governance/accounting text
```

and similarly for the Adversarial packet.

That is enough to answer:

> Is the protocol helping the agent do research, or making it spend its context window talking about ARGUS?

No new schema needed.

### 6. Keep the category-error lesson

This is definitely worth preserving:

```text
open debt
    ≠
EBP violation
```

And more generally:

```text
research objection
    ≠
governance verdict
```

That is exactly what the live adversarial run exposed.

But I would **not** introduce a new `epistemic_kind` field yet. The dry run demonstrated the semantic distinction, not the need for another database field.

---

## What I would *not* carry forward

I would reject or defer these parts:

```text
"the corpus layer flattens all strata"
```

That is no longer a good description of the current architecture. You've deliberately simplified the corpus into static repository background, so don't build a new corpus ontology around that observation.

Likewise:

```text
"one field solves the stratification problem"
```

is too strong. A field can make a distinction visible; it does not establish the legal transitions between strata. The protocol and authority boundary do that.

And don't introduce a formal "type system" as a product feature. The useful lesson is **typed boundaries**, not a new type-theory subsystem.

---

# The policy I would actually freeze

I'd put this in `AGENTS.md` or a short architecture-governance section:

```text
BOOKKEEPING FREEZE

ARGUS must not accumulate research-accounting infrastructure merely because
it appears useful.

A new accounting capability requires all three:

1. A concrete failure observed in an actual work/review cycle.
2. Evidence that the failure affects trust, research continuity, authority,
   reproducibility, or adjudication.
3. The smallest mechanism capable of fixing that failure.

Prefer system-generated accounting over agent-authored bookkeeping.

Keep infrastructure accounting separate from research epistemology:

    research obligations = the work
    ledger obligations   = accounting

The freeze does not restrict the research program, EBP debt, scientific
gates, null models, counterexamples, or methodology.

When an accounting proposal is deferred, record its trigger condition rather
than implementing it preemptively.
```

Then carry a tiny deferred list:

```text
Deferred until demonstrated need:
- review-coverage persistence
- context snapshots
- epistemic-kind field
- agent attestation
- additional provenance projections
```

### One final correction

I would say the **bookkeeping freeze begins after the current P0 persistence repair is fully verified**, not "after Plan 13." Plan 14 is still part of fixing a concrete integrity failure, so it belongs on the **required** side of the line.

After that:

```text
                      ARGUS INFRASTRUCTURE
                              |
                 ┌────────────┴────────────┐
                 |                         |
           proven failure             no proven failure
                 |                         |
              fix it                    defer it
                 |
                 v
          smallest mechanism
```

That is the right balance for EBP 2.1.

The next meaningful experiment should therefore be **research**, not more ARGUS bookkeeping: get the valid packet path green, then let Work → Adversarial → Work operate on the BM–IST–AS problem and see what the system actually teaches you.


