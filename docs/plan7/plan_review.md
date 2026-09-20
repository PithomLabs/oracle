I reviewed the implementation plan. The overall direction is right, but **I would not approve this plan as written**. There are several concrete inconsistencies that will break the promised “`git clone → task setup → task up`” experience.

The good news is that these are implementation-plan defects, not architecture defects.

## Blocking issues

### 1. MCP architecture is currently wrong

The plan says the MCP server should be:

```text
OpenCode
→ stdio MCP
→ Coordinator :8080
```

but its proposed `cmd/argus-mcp/main.go` actually constructs a **new in-process Coordinator**:

```text
cmd/argus-mcp
→ coordinator.New(...)
→ adapter.NewAdapter(coord)
```

That contradicts the frozen boundary and defeats the standalone Coordinator process. The plan explicitly says the MCP adapter connects to Coordinator over HTTP, but the implementation sketch does not. 

**Fix:** MCP must use a thin **Coordinator HTTP client/adapter**.

```text
OpenCode
  ↓ stdio MCP
cmd/argus-mcp
  ↓ HTTP
Coordinator :8080
```

It must not construct Coordinator, load packs, create an ArtifactRegistry, or talk to Solvent/Conductor directly.

This is the most important correction.

---

### 2. `task setup` tries to create the Conductor project before Conductor is running

The plan says:

```text
setup
 → setup:seed
    → conductor:create-project
```

but `up` is a later step. 

That cannot work if `conductor:create-project` uses the REST API.

**Fix:**

```text
task setup
  → build
  → DB reset
  → seed principal
  → config

task up
  → DB ready
  → Solvent start
  → wait Solvent
  → Conductor start
  → wait Conductor
  → create default project
  → Coordinator start
  → wait Coordinator
  → Trust UI start
  → wait Trust UI
```

Or make project creation a one-time idempotent step after Conductor readiness.

---

### 3. `up:solvent` is circular/wrong

This is a clear implementation bug:

```text
up:solvent
  → wait:solvent
  → start Solvent
```

The service is being waited on **before it is started**. 

It must be:

```text
up:solvent
  → start Solvent
  → wait:solvent
```

The same sequencing should be applied consistently:

```text
start → wait → next service
```

`up` currently does not wait after starting Conductor, Coordinator, or Trust UI either. 

---

### 4. Fresh checkout still assumes sibling repositories

The plan promises:

> `git clone → task setup`

but the Taskfile assumes:

```text
../solvent-main
../conductor
```



Those repositories will not exist after cloning only Oracle.

This is a direct contradiction with the stated DX goal.

Choose one explicit strategy:

```text
task setup:deps
```

that clones **pinned revisions** of Solvent and Conductor into `.tmp/src/`, or formally define a multi-repository workspace requirement.

Given the goal is seamless onboarding, I recommend:

```text
.tmp/src/solvent
.tmp/src/conductor
```

with repository URLs and pinned commits/tags in Taskfile/config.

Do not silently depend on a particular sibling-directory layout.

---

### 5. `task verify` does not actually run the standalone Physics Verifier

The plan creates:

```text
cmd/verifier/main.go
```

but `task verify` currently does:

```text
task verify:corpus
task verify:verifier

verify:verifier
→ go test -run TestPhysicsVerifier ./verifier/...
```



That tests the library. It does **not exercise the newly promised CLI binary**.

The developer experience should actually be:

```text
task verify
  → corpus integrity
  → build verifier
  → run verifier against a known fixture
  → verify artifact/hash/output
```

Otherwise the standalone verifier binary exists but is not part of the developer workflow.

---

## High-priority issues

### 6. Solvent migrations should use the canonical `db:reset`

The plan manually reproduces all ten migration filenames. 

Since Solvent already has the canonical lifecycle target, reuse it rather than copying migration knowledge into Oracle.

Otherwise the two repositories can silently diverge.

Preferred:

```text
task setup
  → task -d <solvent-repo> db:reset
```

or the equivalent canonical Solvent command.

---

### 7. `task test` does not mean “test everything”

Current target:

```text
go test ./...
go test -race ./...
go vet ./...
cd trust-ui && go test ./...
```



It does not run Solvent or Conductor tests.

For the promised DX, either:

```text
task test
```

means the whole local stack's test suite, or introduce:

```text
task test
task test:all
task test:integration
```

I would make `task test` the normal full suite and delegate to Oracle, Trust UI, Solvent and Conductor.

---

### 8. Process lifecycle is not actually idempotent

The Taskfile comments say:

> “If already running, skip”

but `up:solvent` does not implement that; it unconditionally launches `nohup`. 

The same problem exists for the other services.

Also, `down` can leave stale PID files when the process has already exited.

For a developer-oriented Taskfile, require:

```text
task up
task up
```

to be safe.

No duplicate processes, no stale PID confusion.

---

### 9. Coordinator readiness via `/packs/rules` is reasonable, but MCP must not duplicate Coordinator state

Once the MCP architecture is corrected, `argus-mcp` only needs:

```text
ARGUS_COORDINATOR_URL
```

It should not need:

```text
SOLVENT_URL
CONDUCTOR_URL
ARGUS_OPERATOR_PRINCIPAL_ID
```

because that would mean the MCP process is independently reconstructing Coordinator dependencies. The current plan exposes all three. 

After the correction:

```text
OpenCode
→ argus-mcp
→ Coordinator HTTP
```

one downstream dependency: Coordinator.

---

## Important but smaller

The standalone Coordinator launcher currently uses:

```text
PACK_DIR="./domain-pack"
```



Make this robust against working-directory changes. The Taskfile can pass an absolute repository-local path.

Also, the plan hardcodes the operator UUID:

```text
00000000-0000-0000-0000-000000000001
```

That's acceptable for a POC, but put it in one Taskfile variable and generate all dependent configuration from that single source.

---

# The corrected startup topology

After these changes, the intended flow should be:

```text
task setup
│
├── obtain pinned Solvent + Conductor sources
├── build binaries
├── start CockroachDB
├── wait DB
├── Solvent db:reset
├── seed operator
└── generate local config

task up
│
├── start Solvent :9090
├── wait Solvent
├── start Conductor :9091
├── wait Conductor
├── create default project
├── start Coordinator :8080
├── wait Coordinator
└── start Trust UI :8081

OpenCode
│
└── stdio MCP
      ↓
  argus-mcp
      ↓ HTTP
  Coordinator :8080
      ├── Solvent :9090
      └── Conductor :9091
```

Physics verification remains independently runnable:

```text
task verify
   ↓
Physics Verifier CLI
   ↓
VerificationArtifact
```

## Verdict

**The plan is conceptually correct but needs one more revision before implementation.**

The four changes I would consider mandatory are:

```text
1. MCP → Coordinator HTTP client, NOT in-process Coordinator.
2. Fix Taskfile startup ordering and wait semantics.
3. Make fresh checkout obtain pinned Solvent + Conductor dependencies.
4. Make task verify actually execute the Physics Verifier binary.
```

Then add the lifecycle/test improvements above.

Once those are patched, the plan will actually match the promised developer experience rather than merely describe it.
