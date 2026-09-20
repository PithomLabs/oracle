The **Reference Loop** is our deliberately small, real-world experiment for answering one question:

> **Can an agent actually do consequential work through Conductor, Solvent, and an Executor without those responsibilities collapsing into one system?**

It is not another product. It is not a workflow engine. It is not BM-IST. It is not the final Loop Engineering implementation.

It is the **proof**.

## Think of it like a factory

Imagine an AI agent is an employee.

It decides:

> “I need to perform this action.”

That alone doesn't mean the action should happen.

So we give different jobs to different parts of the system:

```text
Agent
  = the employee who decides what to do

Conductor
  = the foreman who organizes the work

Solvent
  = the security officer who decides whether the exact action is permitted

Executor
  = the worker/tool that actually performs it

External system
  = the real world where the change happens
```

Then the result comes back through the system so the work can continue.

That's the Reference Loop:

```text
Agent
  ↓
Conductor
  ↓
Solvent
  ↓
Executor
  ↓
External system
  ↓
Result
  ↓
Conductor
  ↓
Agent
```

## Why do we need it?

Because so far, we have mostly described this architecture on paper.

We've said:

> Conductor should coordinate, but not authorize.

> Solvent should authorize, but not execute.

> Executor should execute, but not create authority.

Those are good principles.

But they are still principles until we make them work in an actual end-to-end scenario.

The Reference Loop asks:

> **Can we actually make these components cooperate while keeping their boundaries intact?**

That is the whole purpose.

---

# A concrete example

Suppose the agent needs to perform a tiny consequential action.

It might be something as simple as:

> Change X to Y in some safe test system.

The exact domain doesn't matter at first.

### Step 1 — Agent discovers work

The agent asks Conductor:

> “What work should I do?”

Conductor says:

> “Here is Task 123.”

Now we have coordination.

Conductor knows:

```text
Task 123 exists
Agent A owns it
Task is active
```

But Conductor has **not** decided that the eventual action is allowed.

---

### Step 2 — Agent decides what it wants to do

The agent reasons:

> “To complete Task 123, I want to perform operation X against target T.”

That is agency.

The important thing is that the agent is **proposing** an operation, not granting itself permission.

---

### Step 3 — Solvent decides whether X is authorized

The exact operation goes to Solvent.

Solvent asks its authoritative state:

> “Is this exact operation, against this exact target/state, authorized?”

Suppose the answer is yes.

Now we have:

```text
AUTHORIZED(X)
```

Notice what just happened.

Conductor didn't authorize it.

Agent didn't authorize it.

Executor didn't authorize it.

**Solvent did.**

That's one of the main things the Reference Loop is proving.

---

### Step 4 — Executor actually does X

Now the Executor gets the operation and its authorization evidence.

The Executor performs the real action.

For example:

```text
change test record 123
```

The Executor is allowed to act **because the authority already exists**.

It does not create that authority itself.

---

### Step 5 — Something real happens

The external system changes.

This matters a lot.

We don't want:

```text
Conductor says "done"
→ therefore we assume the world changed
```

We want evidence from the actual effect boundary.

Depending on the system:

```text
External SOR confirms the effect
```

or:

```text
Executor itself is the authoritative source
```

Now we have actual execution evidence.

---

### Step 6 — The result comes back

The result is made available to:

```text
Agent       → so it can interpret what happened
Conductor   → so it can continue coordinating the work
```

The Agent can now say:

> “The operation succeeded. I should do the next thing.”

And Conductor can say:

> “Task 123 has progressed.”

Again, neither of them gets to rewrite what actually happened.

---

# Why this is more important than it first appears

The Reference Loop is really testing whether we can keep five different things separate:

```text
What the agent wants to do
        ≠
What work is assigned
        ≠
What is authorized
        ≠
What was executed
        ≠
What actually happened
```

That's the heart of the architecture.

Without this separation, systems tend to collapse everything into something like:

```text
task.status = approved
```

and then everyone starts treating "approved" as:

> assigned
> authorized
> executed
> successful

Those are not the same thing.

The Reference Loop forces us to prove that they aren't accidentally becoming the same thing in implementation.

---

# The happy path is only half the point

The really valuable part begins when we deliberately break the loop.

For example:

### Agent asks to do X, but Executor tries Y

```text
Authorized: X
Attempted:  Y
```

That must fail.

Otherwise our authority boundary is fake.

### There is no authorization

Executor tries anyway.

That must fail closed.

### Authorization is stale

Executor tries to use it.

That must fail.

### Solvent cannot be reached

The integration has to behave according to the declared verification model.

It cannot simply say:

> “Solvent is unavailable, so I'll assume it's okay.”

### Executor fails

We must record:

```text
FAILED
```

not:

```text
SUCCESS
```

### We can't determine whether the external effect happened

We must record:

```text
AMBIGUOUS
```

not blindly choose success or failure.

### Agent misunderstands the result

The Agent might think:

> “It succeeded!”

when the authoritative result says:

> “Failed.”

That's okay from an architecture perspective.

The Agent is allowed to reason incorrectly.

What must **not** happen is that its interpretation rewrites the authoritative result.

---

# Why we are doing this now

This is the most important part of the pivot.

We spent a lot of time defining Loop Engineering because we wanted to be certain we weren't accidentally creating another workflow system or mixing authority into coordination.

We've now reached the point where more specification work has diminishing returns.

The Reference Loop lets reality answer the question.

Instead of asking:

> “What should happen if X?”

we can ask:

> “Let's actually run X and see what happens.”

And when something breaks, we classify the problem:

```text
implementation problem
integration problem
Executor problem
deployment problem
new security property
```

That last category is the only one that potentially tells us:

> “Our architecture is missing something fundamental.”

Everything else is an engineering fix.

---

# What the Reference Loop is NOT

This distinction is important.

It is **not**:

```text
a new workflow engine
a Loop Engineering server
a replacement for Conductor
a replacement for Solvent
a new orchestration framework
a generic execution gateway
a BM-IST runtime
```

It is basically a **thin proof harness and reference scenario** around the existing systems.

We are trying to prove the architecture, not build another architecture.

---

# What success looks like

At the end of the first experiment, I want to be able to point at one piece of work and reconstruct:

```text
Agent:
"I proposed X."

Conductor:
"Task X was assigned to this agent."

Solvent:
"Exact operation X was authorized."

Executor:
"I attempted exact operation X."

External system / SOR:
"X actually happened."

Conductor:
"The work progressed."

Agent:
"I received the result and continued."
```

And then we should be able to deliberately make one part of that sequence wrong and show that the system catches the violation.

That's much stronger evidence than another hundred pages of architecture.

## In one sentence

**The Reference Loop is our smallest real experiment for proving that Agent, Conductor, Solvent, and Executor can work together on a consequential action while keeping agency, coordination, authority, and execution genuinely separate.**

And once that works, we can start substituting agents, Executors, and orchestration mechanisms to prove the broader Loop Engineering claim.

