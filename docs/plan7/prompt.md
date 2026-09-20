You are the senior implementation-planning agent for the ARGUS / Oracle
Physics Trust Verification POC.

DO NOT IMPLEMENT ANYTHING YET.

Your task is to draft an implementation plan for making the entire ARGUS
and Physics Trust Verifier developer experience seamless and reproducible.

The goal is:

    git clone
        ↓
    one command
        ↓
    all required services/binaries/databases running
        ↓
    developer can execute:
        - unit/integration tests
        - Physics verifier
        - ARGUS Phase 8 dry run
        - Work OpenCode agent
        - Adversarial OpenCode agent
        - Trust UI

The intended developer experience should be approximately:

    task setup
    task up
    task test
    task verify
    task dry-run

Determine the exact commands only after inspecting the actual repositories.

DO NOT implement code.
Produce only the implementation plan.

==================================================
BACKGROUND
==================================================

ARGUS currently has these frozen boundaries:

CAPABILITY ≠ WORK ≠ AUTHORITY ≠ EXECUTION

OpenCode
  = agency / research work

Conductor
  = operational workflow/task state

Solvent
  = epistemic authority:
    beliefs, evidence, debt, edges,
    promotion, retraction, authorization, audit

Coordinator
  = orchestration boundary / packet compilation /
    HTTP interface over Conductor + Solvent

Trust UI
  = human observation + adjudication

EBP v2.1
  = epistemic operating discipline

RCP
  = thin read-only research context projection

MCP
  = agent-facing transport

Physics Verifier
  = bounded trusted verification artifact producer

BM-IST
  = current domain program / Domain Pack

Phase 8 is implemented and ready for its first actual dry run.

The current README states that:
- Coordinator exists as a library + HTTP handler
- no standalone Coordinator binary currently exists
- Trust UI expects a Coordinator HTTP service
- Solvent requires CockroachDB
- Conductor is a separate service
- actual Work-Agent / Adversarial-Agent execution is still outstanding

The goal of this plan is to close the developer-experience gap.

==================================================
PRIMARY OBJECTIVE
==================================================

Design a reproducible local stack where the developer does not have to
manually discover, build, configure, or launch each dependency.

The desired architecture is:

                  task up
                     │
        ┌────────────┼────────────┐
        ▼            ▼            ▼
   CockroachDB    Solvent      Conductor
                     │            │
                     └─────┬──────┘
                           ▼
                     Coordinator
                           │
                    ┌──────┴──────┐
                    ▼             ▼
                 Trust UI      ARGUS MCP
                                  │
                           ┌──────┴──────┐
                           ▼             ▼
                       Work Agent   Adversarial Agent

The plan must determine which processes actually exist today and which
ones must be packaged into runnable binaries.

==================================================
REPOSITORIES TO INSPECT
==================================================

Before writing the plan, inspect the ACTUAL repositories and current
working tree.

At minimum inspect:

1. Oracle / ARGUS
   - go.mod
   - coordinator/
   - coordinator/http/
   - mcp/
   - trust-ui/
   - verifier/
   - domain-pack/
   - packet/
   - corpus/
   - reference-loop/
   - current Taskfile, if any
   - all cmd/ directories
   - existing scripts
   - README.md
   - test infrastructure

2. Solvent
   - go.mod
   - main binaries / cmd/
   - API server entry points
   - migration/bootstrap commands
   - MCP server
   - CockroachDB requirements
   - configuration/environment variables
   - health/readiness endpoints
   - test setup
   - existing dev/run scripts
   - existing Taskfile, if any

3. Conductor
   - go.mod
   - main binaries / cmd/
   - HTTP server
   - database initialization
   - SQLite/in-memory mode
   - configuration/environment variables
   - health/readiness endpoint
   - test setup
   - existing dev/run scripts
   - existing Taskfile, if any

4. Existing reference-loop
   - determine which services it assumes
   - determine whether it already contains useful orchestration logic
   - do NOT modify it unless the implementation plan explicitly justifies
     a reusable component

5. Physics verifier
   - verifier/
   - physics/v1/
   - artifact registry
   - required source/data artifacts
   - how trusted artifacts are produced
   - how those artifacts are surfaced to ARGUS
   - whether any external runtime is required

==================================================
IMPORTANT: DO NOT ASSUME SERVICE STARTUP
==================================================

Do not assume:
- Solvent has a binary named "solvent"
- Conductor has a binary named "conductor"
- Coordinator has a binary
- CockroachDB is installed locally
- services listen on the README's historical ports
- existing scripts still work

Inspect the actual code.

For every required component determine:

    binary
    source path
    build command
    default port
    config
    environment variables
    dependencies
    startup command
    readiness check
    shutdown behavior
    logs
    persistent data location

==================================================
COORDINATOR DECISION
==================================================

The current Oracle repository has Coordinator as a library + HTTP handler.

Determine the smallest clean way to make it runnable.

Preferred outcome:

oracle/cmd/coordinator/main.go

or equivalent.

The plan must define:

- how Coordinator is constructed
- how Solvent client is configured
- how Conductor client is configured
- PackRegistry initialization
- ArtifactReader / verifier initialization
- operator configuration
- auth configuration
- HTTP listen address
- graceful shutdown
- health endpoint if needed
- structured logging if already supported

Do NOT turn Coordinator into a new architectural layer.

The standalone binary is merely a launcher around the existing Coordinator
library.

==================================================
SOLVENT
==================================================

Determine the exact developer startup sequence.

The plan must cover:

1. Start CockroachDB.
2. Wait for database readiness.
3. Run/apply Solvent migrations.
4. Seed required data if the existing code requires it.
5. Start Solvent API.
6. Verify Solvent health/readiness.
7. Verify required Phase 8 endpoints exist.
8. Ensure the Phase 8 edge endpoint is included:
       POST /v1/beliefs/{parent_id}/edges
9. Ensure /v1/discharge is available.
10. Ensure required evidence/belief/ledger APIs are available.
11. Document Solvent MCP separately if needed, but DO NOT expose it to
    Work/Adversarial OpenCode agents.

Determine whether Solvent itself already has a working Taskfile or
startup target that can be reused.

Prefer reusing existing mechanisms over duplicating them in Oracle.

==================================================
CONDUCTOR
==================================================

Determine the exact startup sequence.

The plan must cover:

- build/start binary
- initialize database if necessary
- configure port
- configure project
- health/readiness
- cancellation API
- task/dependency APIs
- governance_ref
- priority

Use existing SQLite/in-memory capabilities if appropriate.

Do not create a new task database for ARGUS.

The Conductor remains the operational system of record.

==================================================
TRUST UI
==================================================

Make Trust UI runnable as part of the stack.

Plan configuration for:

COORDINATOR_URL
ARGUS_OPERATOR_PRINCIPAL_ID
operator authentication/token configuration
listen port

Trust UI remains:
- separate Go module
- Go templates
- vanilla HTML/CSS/JS
- no direct DB
- no direct Solvent
- no direct Conductor

Determine whether the current Trust UI requires any other environment
variables.

==================================================
ARGUS MCP
==================================================

Determine how the MCP adapter is built and launched.

The agent-facing surface must remain exactly:

    argus.get_context
    argus.submit_packet

Do NOT expose:
- solvent-mcp
- Conductor write tools
- edge mutation tools
- database tools

The plan must explain how Taskfile starts or prepares the MCP adapter
for OpenCode.

Determine whether the adapter is:
- standalone binary
- stdio process
- HTTP client wrapper
- embedded server

Do not guess.

==================================================
OPEN CODE WORKFLOW
==================================================

The developer experience must support TWO independent OpenCode processes.

### Work Agent

Fresh OpenCode process:

    argus.get_context
        ↓
    RCP/v1
        ↓
    research
        ↓
    argus.submit_packet

### Adversarial Agent

Different fresh OpenCode process:

    argus.get_context
        ↓
    reconstruct existing state
        ↓
    attack
        ↓
    argus.submit_packet

The plan must define how a developer gets the correct MCP configuration
without accidentally exposing Solvent/Conductor mutation tools.

Prefer generated local configuration or a checked-in template.

Do not place credentials in git.

==================================================
PHYSICS VERIFIER
==================================================

The developer experience must also make the Physics Verifier runnable.

Determine:

- exact verifier entry point
- whether a CLI binary is needed
- whether `cmd/verifier` or equivalent should be created
- required input artifacts
- required corpus files
- artifact registry behavior
- trusted registration path
- expected output artifact
- hashes
- where artifacts are written
- how ARGUS can reference them
- how a developer verifies the result
- how verifier output participates in EBP evidence

The Physics Verifier MUST remain bounded.

Do not turn it into a full CAS.

Do not move scientific authority into the verifier.

The verifier produces trusted verification artifacts;
Solvent remains epistemic authority.

If a standalone verifier binary materially improves the DX,
include it in the plan.

==================================================
TASKFILE
==================================================

Design a root Taskfile as the single developer orchestration surface.

Inspect whether a Taskfile already exists in Oracle, Solvent, or
Conductor and decide what should be reused versus consolidated.

Preferred user experience:

    task setup

    task up

    task status

    task logs

    task test

    task test:integration

    task verify

    task dry-run

    task down

    task clean

Exact names may differ if the repositories already have conventions.

The final plan must define:

### setup
- check prerequisites
- create local directories
- create local config
- initialize databases
- build required binaries
- prepare fixtures/corpus

### up
Start, in dependency order:

1. CockroachDB
2. Solvent
3. Conductor
4. Coordinator
5. ARGUS MCP adapter if required as a persistent process
6. Trust UI

Wait for readiness between dependencies.

Do NOT use arbitrary sleep durations as the only readiness mechanism.

Prefer:
- HTTP health checks
- database readiness commands
- process checks
- bounded retries/timeouts

### status
Show:
- process
- PID
- port
- readiness
- health
- log location

### logs
Provide convenient grouped logs.

### test
Run all normal repository tests.

### test:integration
Run DB-backed integration tests when dependencies are running.

### verify
Run:
- Physics Verifier
- trusted artifact checks
- corpus integrity checks
- relevant ARGUS verification tests

### dry-run
Run or prepare the actual Phase 8 workflow.

The plan must decide how much can safely be automated.

At minimum it should verify:

    Work Agent context
    → work packet
    → persistence
    → fresh Adversarial Agent context
    → contradiction
    → human decision
    → retraction/dead-end branch
    → surviving/successor debt discharge
    → promotion gate

Do not fake human actions merely to make the Taskfile green.

Where human intervention is required, the task should stop with clear
instructions.

### down
Graceful shutdown in reverse dependency order.

### clean
Only remove POC-local state, never destroy source code.

Provide a safer explicit target for destructive cleanup if needed.

==================================================
HEALTH / READINESS
==================================================

Define a single readiness strategy.

Potential examples:

    task wait:cockroach
    task wait:solvent
    task wait:conductor
    task wait:coordinator

But reuse real health endpoints if they already exist.

The plan must specify:
- endpoint/command
- timeout
- retry interval
- failure output

Do NOT rely solely on "process exists."

==================================================
PORTS
==================================================

Inspect actual defaults.

Do NOT blindly reuse:
localhost:8080
for multiple services.

Produce one final local development port map.

Example format:

    CockroachDB   localhost:26260
    Solvent       localhost:XXXX
    Conductor     localhost:YYYY
    Coordinator   localhost:ZZZZ
    Trust UI      localhost:8081

Use actual verified ports.

==================================================
CONFIGURATION
==================================================

Prefer a single local configuration strategy.

Examples:

.env.example
Taskfile.yaml
config/dev/
or existing repository mechanism.

Do not create parallel configuration systems unnecessarily.

Document:
- Solvent URL
- Conductor URL
- Coordinator URL
- Trust UI URL
- database URL
- operator principal
- auth token/key
- MCP configuration
- artifact directory
- corpus directory

Secrets:
- never hardcode credentials
- provide local development defaults only when safe
- clearly distinguish POC credentials from production credentials

==================================================
DATA / STATE MANAGEMENT
==================================================

Define local persistent directories, e.g.:

.tmp/
  cockroach/
  solvent/
  conductor/
  artifacts/
  logs/
  pids/

Use actual existing conventions if they exist.

Taskfile must make state ownership clear.

Do not introduce another persistent application database.

==================================================
TEST MATRIX
==================================================

The implementation plan MUST include tests for the developer environment.

At minimum:

### Startup
- all required processes start
- all readiness checks pass
- startup ordering works
- repeated `task up` is safe or reports a clear state

### Shutdown
- graceful shutdown
- no orphan processes
- ports released

### Coordinator
- standalone binary builds
- Coordinator connects to Solvent
- Coordinator connects to Conductor
- `/v1/context/{task_id}`
- `/decisions`
- authentication

### Solvent
- DB migration
- API startup
- required Phase 8 endpoints
- edge creation
- discharge
- evidence reads

### Conductor
- startup
- task lifecycle
- cancellation
- governance_ref
- priority

### MCP
- exactly two agent tools
- no authority tools
- OpenCode can connect

### Physics Verifier
- verifier executes
- artifact generated
- artifact hash verified
- trusted registration works
- ARGUS can reference artifact

### End-to-end
- Work Agent can reconstruct context
- packet persists
- Adversarial Agent can reconstruct context
- adversarial packet persists
- contradicts edge exists
- human retraction path works
- linked task cancellation works
- dead-end appears
- debt-discharge path works
- Solvent promotion gate works

==================================================
CI / AUTOMATION
==================================================

Determine what should run in CI.

At minimum propose:

    task test

For DB-dependent tests:
- decide whether CI should launch CockroachDB automatically
- prefer the same Taskfile path developers use
- avoid a developer-only setup that CI cannot reproduce

For the Physics Verifier:
- deterministic fixtures
- artifact hashes
- no network dependency unless explicitly required

==================================================
PHASE 8 DRY RUN SCRIPT
==================================================

The plan should provide a concrete operator playbook.

Example:

    task setup
    task up
    task status
    task verify

Then:

    Terminal A:
      OpenCode Work Agent

    Terminal B:
      OpenCode Adversarial Agent

    Browser:
      http://localhost:<trust-ui-port>

Document exactly what the developer should observe at each stage.

The playbook must make the two-agent separation explicit.

==================================================
IMPORTANT ARCHITECTURAL CONSTRAINTS
==================================================

DO NOT:

- merge Solvent into Oracle
- merge Conductor into Oracle
- embed Solvent database access into Coordinator
- give agents Solvent MCP tools
- give agents Conductor mutation tools
- create a second research database
- add direct database writes to the Taskfile
- make Taskfile bypass service APIs
- make Trust UI bypass Coordinator
- turn Coordinator into a scientific reasoning engine
- turn Physics Verifier into a general CAS
- make the dry run entirely autonomous if human adjudication is required

Taskfile orchestrates processes.
It does not become an authority layer.

==================================================
EXPECTED IMPLEMENTATION ARTIFACTS
==================================================

The plan should determine whether the following are needed:

    Taskfile.yaml
    oracle/cmd/coordinator/main.go
    oracle/cmd/verifier/main.go
    oracle/cmd/argus-mcp/main.go

Only add binaries that are genuinely required after inspecting the
existing repository.

If equivalent binaries already exist, reuse them.

Also determine whether:
- health endpoints need to be added
- readiness helpers are required
- local config files are required
- Docker/Podman is necessary
- CockroachDB can be launched through an existing mechanism
- process supervision should be Taskfile-only

Prefer the smallest solution.

==================================================
DEVELOPER EXPERIENCE ACCEPTANCE CRITERION
==================================================

The final plan must target this experience:

Fresh checkout:

    task setup
    task up
    task status

produces:

    CockroachDB READY
    Solvent READY
    Conductor READY
    Coordinator READY
    MCP READY
    Trust UI READY

Then:

    task test
    task verify

works.

Then the developer can execute the actual Phase 8 dry run with two
fresh OpenCode processes and the Trust UI.

The developer should NOT need to:
- discover undocumented ports
- manually construct service commands
- manually run database migrations
- manually build Coordinator
- manually discover MCP configuration
- manually wire Solvent and Conductor URLs
- manually create runtime directories
- manually guess readiness
- manually clean up orphan processes

==================================================
DELIVERABLE
==================================================

Create:

PHASE8_DEV_ENV_IMPLEMENTATION_PLAN.md

Structure:

1. Objective
2. Actual current runtime architecture
3. Repository/service inventory
4. Required binaries
5. Required databases
6. Coordinator launcher plan
7. Physics Verifier launcher plan
8. Solvent startup plan
9. Conductor startup plan
10. Trust UI startup plan
11. MCP/OpenCode startup plan
12. Port/configuration map
13. Taskfile design
14. Readiness/health strategy
15. Local state/log management
16. Testing strategy
17. CI strategy
18. End-to-end Phase 8 dry-run procedure
19. Developer experience acceptance criteria
20. Risks / known limitations
21. Explicit non-goals

For every proposed change include:
- repository
- file/directory
- reason
- dependencies
- exact behavior
- tests

Do not implement.

Do not invent commands or binaries.

Inspect the actual codebases first and make the plan fit the existing
architecture.

The goal is not "more tooling."

The goal is:

> One reproducible developer command starts the entire ARGUS trust
> verification environment so the developer can immediately test the
> Physics Verifier and execute the Phase 8 Work-Agent / Adversarial-Agent
> workflow.