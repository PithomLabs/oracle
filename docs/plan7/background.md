## prompt

how to run argus and test the workflow per README

## answer

Per the current README, **ARGUS Phase 8 is not yet a single-command runnable stack**. The repository deliberately has the Coordinator as a library + HTTP handler, but no standalone Coordinator binary. The README explicitly calls that out. 

## What you need running

The intended runtime is:

```text
OpenCode Work Agent
        │
        │ MCP
        ▼
ARGUS MCP adapter
        │
        ▼
Coordinator HTTP
        │
   ┌────┴────┐
   ▼         ▼
Solvent   Conductor
   │
CockroachDB
```

The Trust UI sits beside that:

```text
Browser
   │
Trust UI :8081
   │
Coordinator :8080
```

The README identifies CockroachDB as required for Solvent integration/dry-run execution, and Conductor as a separate service.  

## 1. Start CockroachDB

The README expects:

```text
localhost:26260
```

This is required for the Solvent database backend and integration tests. 

## 2. Start Solvent

The README says Solvent is an external REST service, with its address configured rather than giving a concrete startup command. 

So you need the Solvent service running and connected to the CockroachDB instance.

## 3. Start Conductor

Likewise, Conductor is a separate REST service:

```text
Conductor REST API
```

Its address is configured in the Coordinator/client configuration. 

## 4. Start Coordinator

This is the current gap.

The README says:

```text
oracle/coordinator
    = Go library + HTTP handler

no standalone Coordinator binary
```

and:

```text
Trust UI expects:
http://localhost:8080
```



Therefore **there is no repository command like `go run ./cmd/coordinator` yet**.

For the complete Phase 8 test, you need some separate process that instantiates the Coordinator HTTP server and listens on `:8080`.

So at this point the workflow cannot be reproduced from the repository alone.

## 5. Start Trust UI

The README does give a runnable command:

```bash
cd trust-ui
go build -o trust-ui .

COORDINATOR_URL=http://localhost:8080 \
ARGUS_OPERATOR_PRINCIPAL_ID=<operator-uuid> \
./trust-ui
```

It listens on:

```text
http://localhost:8081
```

and provides:

```text
/insights
/debts
/api/retire
```



Then open:

```text
http://localhost:8081
```

## 6. Configure OpenCode

The architecture requires two separate OpenCode processes.

Each gets exactly:

```text
argus.get_context
argus.submit_packet
```

and does **not** get Solvent MCP, Conductor write tools, or database access. 

The README explains the tool boundary, but **does not currently give the exact OpenCode MCP configuration command/path**. So that part still needs to be supplied from the actual `argus` MCP adapter configuration.

## 7. Execute the actual Phase 8 workflow

### Work Agent

Start a fresh OpenCode process and give it the initial task for G0.

Its expected sequence is:

```text
argus.get_context
        ↓
RCP/v1
        ↓
inspect current task / claims / evidence / debt
        ↓
perform bounded G0 research
        ↓
argus.submit_packet
```

The packet should result in persisted beliefs/evidence/edges/tasks.

The README defines the work-agent workflow this way. 

### Inspect Trust UI

At:

```text
http://localhost:8081/insights
```

you should see the evolving research state.

Then:

```text
http://localhost:8081/debts
```

should show the debt, evidence, adversarial challenges and actual retirement rules. 

### Adversarial Agent

Start a **new** OpenCode process with no conversation history.

It should call:

```text
argus.get_context
```

and reconstruct the work from ARGUS state rather than inheriting the first agent's conversation.

Then it performs the adversarial review and submits:

```text
argus.submit_packet
```

The adversarial packet should create a `contradicts` edge.

That edge should become visible as an adversarial challenge in Insights.



## 8. Exercise Branch A

The first branch is:

```text
Work
  ↓
Adversarial challenge
  ↓
contradicts edge
  ↓
Human RETRACT
  ↓
Solvent RetractCascade
  ↓
Conductor task cancelled
  ↓
Dead End
```



Verify the task becomes terminal `cancelled`, not merely rejected.

## 9. Exercise Branch B

Use a surviving belief or successor:

```text
Debt
  ↓
human reviews evidence
  ↓
qualifying evidence
  ↓
Coordinator validates Pack rule
  ↓
Solvent /v1/discharge
  ↓
operator attribution
  ↓
PROMOTE request
  ↓
Solvent gate
```

A retracted belief is **not** subsequently promoted. 

## The practical problem right now

You can run:

```text
CockroachDB
Solvent
Conductor
Trust UI
```

but **the README does not provide the missing Coordinator launcher**, and it also doesn't give the exact OpenCode MCP startup configuration. The README itself identifies the Coordinator binary as the outstanding operational gap. 

So the most accurate current runbook is:

```text
CockroachDB
    ↓
Solvent
    ↓
Conductor
    ↓
[missing standalone Coordinator launcher]
    ↓
Trust UI :8081
    ↓
OpenCode Work Agent
    ↓
OpenCode Adversarial Agent
```

**Before the actual dry run, I would add a tiny `cmd/coordinator` launcher to Oracle.** That does not change the Coordinator architecture; it simply makes the already-existing HTTP server actually runnable. It would turn the current POC from "library + manually assembled runtime" into a reproducible Phase 8 stack.
