The audit confirms the invariant you were aiming for:

> **An existing belief ID cannot have its claim text rewritten through any current application path.** 

The implementation achieves this through several independent mechanisms:

```text
claim
  ↓
deterministic belief ID
  ↓
INSERT ... ON CONFLICT DO NOTHING
  ↓
no UPDATE claim path
```

The audit found no claim-mutating SQL in the kernel, no REST endpoint for changing a claim, no UI path, and no packet path that updates an existing belief.   

That means your intended lifecycle is valid:

```text
B1
original proposition
   |
   | adversarial attack
   v
B2
challenge / counterclaim
   |
   | corrective research
   v
B3
revised proposition
```

with the graph preserving the relationships and **B1 remaining historically intact**. The audit specifically confirms that the existing adversarial contradiction and reopen paths create new belief IDs rather than rewriting the old claim. 

### One correction to the audit's proposed regression test

The audit says the remaining test gap is:

> same belief ID + different claim

but the proposed test actually does this:

```text
claim C1 → UUID1
claim C2 → UUID2
```

so it demonstrates that **different claims get different IDs**, not that an explicit existing-ID mutation is rejected. 

That's still useful, but it isn't literally a same-ID mutation test.

The stronger conclusion from the audit is:

> **There is currently no legitimate application/API operation that accepts an existing belief ID and a replacement claim.**

That is actually better evidence of the invariant than a test of an operation the system does not expose.

### One architectural point I would lock now

Document this as a **hard invariant**:

```text
BELIEF IMMUTABILITY

Once a belief ID exists:

    belief.id
    belief.claim

are historical and immutable.

Research evolution occurs by:
    creating a new belief
    +
    connecting it through graph semantics
```

Status, debt, evidence relationships, and lifecycle transitions may change through their authorized paths; the proposition itself does not. The audit found that existing transitions update only debt or status, not claim. 

That is exactly the property you want for the adversarial loop.

### One thing worth checking next

The audit describes reopen as:

```text
old retracted belief
    |
    | derives
    v
new belief
```

but your conceptual model may prefer the semantic direction:

```text
new belief
    derives
old belief
```

depending on what `derives` is intended to mean in your graph grammar. The audit says the current implementation uses “new derives old.” 

I would **not change this casually**. First establish the exact intended semantics of `derives` across the packet validator, UI, and Solvent. That is separate from the immutability invariant.

Overall, though, the important result is strong:

```text
Original claim
     ↓
immutable belief ID
     ↓
adversarial challenge
     ↓
new belief(s)
     ↓
graph relationship
     ↓
historical continuity
```

That is the research lifecycle you described, and the current ARGUS implementation supports it. 
