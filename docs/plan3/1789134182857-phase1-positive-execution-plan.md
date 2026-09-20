# Phase 1 Positive Execution Experiment

## Objective
Prove the positive execution boundary with a fake/RecordingFunc executor:
valid authorization → valid intent_id → ExecuteAction → executor invoked exactly once.

Do NOT prove real external effect. RecordingFunc execution remains classified as simulated/non-external-effect proof.

## Current State
- Solvent frozen at HEAD=7602699, working tree clean
- Negative path tests pass (mismatched actor denied, missing authorization denied)
- Happy path test exists but does not verify executor invocation count
- Agent RunHappyPath records evidence but does not verify exactly-once invocation
- Existing EFFECT_CONFIRMED evidence in happy path is already labeled non_external_effect_proof=true

## Implementation Plan

### 1. Instrument fake executor in main.go
Add thread-safe invocation counter to main.go package:
- sync.Mutex-protected counter
- resetExecutorCallCount() — set to 0
- incrementExecutorCallCount() — increment on each call
- getExecutorCallCount() — read current count

Modify the registered `github_trigger_workflow` fake executor to call incrementExecutorCallCount().

### 2. Expose counter for test access
Either:
- Return counter accessor from startSolventServer, OR
- Export getExecutorCallCount() and resetExecutorCallCount() from main.go test package

### 3. Write positive_execution_test.go
Test flow:
1. Reset executor call counter to 0
2. Setup: conductor DB, Solvent DB, Solvent server, project, principal
3. Create executor registry via executor.NewRegistry("recording", ...)
4. Create evidence collector
5. Create agent and call RunHappyPath(ctx)
6. Assertions:
   - AUTHORIZED event count == 1
   - EXECUTION_ATTEMPTED event count == 1
   - Simulated effect event count == 1 (EFFECT_CONFIRMED with non_external_effect_proof=true)
   - Executor call count == 1
   - Canonical operation_id identical across proposal, authorization, intent, execution request, executor invocation
   - Canonical form: deploy:pithomlabs/reference-loop-test:ref-loop.yml:main:<run_id>
   - non_external_effect_proof == true on effect event
   - No claim of real external EFFECT_CONFIRMED

Test isolation:
- Reset counter before test
- Assert counter == 0 before execution
- Assert counter == 1 after execution
- Fail if counter > 1

### 4. Validation
```bash
cd /home/chaschel/Documents/go/oracle/reference-loop
go test ./...
go vet ./...
```

### 5. Verify Solvent freeze
```bash
cd /home/chaschel/Documents/go/solvent-main
git rev-parse HEAD    # expect 7602699
git status --short    # expect clean
git diff --name-only  # expect empty
```

### 6. Documentation
Create PHASE1_POSITIVE_EXECUTION.md:
- Scenario
- Setup
- Operation identity (canonical form)
- Authorization result
- Intent ID
- ExecuteAction result
- Executor invocation count
- Executor operation identity
- Effect classification (simulated / non-external-effect proof)
- Evidence source
- Verdict

## Failure Handling
If any assertion fails:
1. Inspect evidence records to determine actual behavior
2. Classify using existing taxonomy:
   - IMPLEMENTATION_DEFECT
   - INTEGRATION_DEFECT
   - EXECUTOR_DEFECT
   - DEPLOYMENT_ENVIRONMENT_DEFECT
   - SPECIFICATION_DEFECT
   - NEW_SECURITY_PROPERTY
Do not pre-classify. Classify after inspection.

## Stop Condition
After this experiment STOP. Do not implement:
- adversarial suite
- concurrency
- cancellation
- replay
- revocation
- BM-IST
- Agent Skill
- Conductor Skill
- protocol formalization
- Conformance Matrix
- new workflow primitives

## Expected Result
PASS or PASS WITH LIMITATIONS. Explicitly state whether valid authorization → exactly-one execution has been empirically proven.
