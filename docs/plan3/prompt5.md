The plan is close, but I would **not approve it unchanged**. There are three important corrections.

### 1. The `EFFECT_CONFIRMED` assertion is too strong

The plan currently says:

```text
EFFECT_CONFIRMED event count == 1
```

But this experiment uses the **fake/RecordingFunc executor**. That does not prove a real external effect. Your Phase 1 evidence already established that effect confirmation was simulated/non-external-effect proof.

So the test should assert:

```text
AUTHORIZED                == 1
EXECUTION_ATTEMPTED        == 1
SIMULATED_EFFECT / equivalent == 1
executor call count        == 1
```

or, if the existing system represents the simulated result using `EFFECT_CONFIRMED`, the test must additionally assert:

```text
non_external_effect_proof == true
```

and explicitly state that this is **not** external-effect confirmation.

This preserves the distinction we have been enforcing.

### 2. The exact operation identity needs to be asserted

The plan says:

> Recording params match repo/workflow/ref/run_id

That is weaker than the invariant we established.

The test should assert that the executor received the **exact canonical `operation_id`**:

```text
deploy:pithomlabs/reference-loop-test:ref-loop.yml:main:<run_id>
```

and that the identity is identical across:

```text
proposal
authorization
intent
execution request
executor invocation
```

Do not merely compare individual parameters.

### 3. Do not pre-classify failures

The plan says:

```text
If instrumentation shows count > 1:
    IMPLEMENTATION_DEFECT in Solvent or reference-loop harness.

If ExecuteAction fails after valid authorization:
    INTEGRATION_DEFECT.
```

That is too premature.

The correct process is:

```text
failure
  ↓
inspect evidence
  ↓
classify actual owner/type
```

For example, `count > 1` could be:

* harness instrumentation error,
* retry behavior,
* duplicate client invocation,
* executor behavior,
* integration defect,
* or potentially a Solvent defect.

Likewise, `ExecuteAction` failing after authorization is not automatically an integration defect.

Use the existing taxonomy:

```text
IMPLEMENTATION_DEFECT
INTEGRATION_DEFECT
EXECUTOR_DEFECT
DEPLOYMENT_ENVIRONMENT_DEFECT
SPECIFICATION_DEFECT
NEW_SECURITY_PROPERTY
```

### One additional implementation refinement

The counter should be reset **per test/run**, not globally shared across tests without isolation. Since this is a call-count experiment, the strongest assertion is:

```text
before execution: 0
after execution: 1
```

and the test should fail if a second invocation occurs.

A mutex-protected counter is fine. An exposed global counter is acceptable for this narrow test, but test isolation must be explicit.

## Revised plan

I would change the core test to:

```text
1. Instrument fake executor with thread-safe invocation counter.

2. Reset counter = 0.

3. Establish valid authorization.

4. Obtain exact intent_id.

5. Execute exact authorized operation.

6. Assert:
   AUTHORIZED == 1
   EXECUTION_ATTEMPTED == 1
   executor call count == 1

7. Assert exact canonical operation_id at executor boundary.

8. Assert the recorded execution result is classified as:
   simulated / non-external-effect proof

9. Do NOT claim real EFFECT_CONFIRMED.

10. Re-run complete suite.

11. Verify Solvent:
      HEAD == 7602699
      working tree clean
      no source modifications.
```

### Verdict

**PASS WITH MINOR CHANGES.**

The experiment itself is exactly the right next experiment. It closes the remaining positive-path gap without expanding architecture. The main correction is conceptual: **prove exactly-once executor invocation, not external effect confirmation.**

Once this passes, the evidence ladder becomes substantially stronger:

```text
wrong actor              → DENY
no authorization         → DENY
valid authorization      → EXACTLY ONE EXECUTION
real external effect     → still unproven
```

That is the correct boundary before moving on to the next experiment.
