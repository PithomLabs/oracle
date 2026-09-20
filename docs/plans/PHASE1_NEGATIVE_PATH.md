# PHASE 1 — NEGATIVE-PATH VALIDATION

## TEST A — MISMATCHED ACTOR AUTHORIZATION DENIAL

### Scenario
Authorize a consequential operation as actor A, then attempt authorization as actor B against the same approved target.

### Setup
- Created belief and promoted it.
- Created authority target owned by `principalA` (`testPrincipalID`).
- Attached justification, requested authorization, approved target.
- Successfully authorized as `principalA` to establish a valid baseline.
- Attempted authorization as `principalB` (`00000000-0000-0000-0000-000000000002`).

### Operation identity
`deploy:pithomlabs/reference-loop-test:ref-loop.yml:main:<run_id>`

### Actor
- Authorized actor: `principalA`
- Attempted actor: `principalB`

### Authorization state before second attempt
Target approved under `principalA`; `action_intent` exists for `principalA`; operation is otherwise valid.

### Execution request
Second `AuthorizeAction` call with `ActorID = principalB`.

### Expected result
Denial with `actor_id_mismatch`.

### Actual result
Denial returned by Solvent REST handler: `actor_id does not match authenticated principal`.

### Executor invoked
No executor invoked. Authorization never reached execution.

### Effect observed
None.

### Evidence source
Solvent audit record / direct REST response error.

### Verdict
PASS — mismatched actor authorization fails closed through the real Solvent authorization boundary.

---

## TEST B — EXECUTION WITHOUT AUTHORIZATION DENIAL

### Scenario
Attempt to execute a consequential operation without a valid authorization intent.

### Setup
- Created belief, promoted it.
- Created and approved target under `principalA`.
- Did NOT authorize an intent for execution.
- Called `ExecuteAction` with a zero UUID `intent_id`.

### Operation identity
`deploy:pithomlabs/reference-loop-test:ref-loop.yml:main:<run_id>`

### Actor
`principalA`

### Authorization state
No live intent exists for the target/action combination at execution time.

### Execution request
`ExecuteAction` with `IntentID = 00000000-0000-0000-0000-000000000000`.

### Expected result
Execution denied; executor not invoked.

### Actual result
Execution denied by Solvent REST handler.

### Executor invoked
No executor invoked.

### Effect observed
None.

### Evidence source
Solvent REST response error.

### Verdict
PASS — execution without valid authorization fails closed through the real Solvent execution boundary.

---

## Classification
- Test A finding type: None — expected behavior confirmed.
- Test B finding type: None — expected behavior confirmed.
- Genuine defects observed: None in this experiment.
- Specification defects observed: None.
- New security properties discovered: None beyond the existing `AUTHORIZE ≠ EXECUTE` invariant, which is now empirically confirmed.

---

## Solvent freeze check during this phase
- HEAD: `7602699`
- Working tree: clean
- Source modifications: none

---

## Final assessment for this experiment

1. Does mismatched actor authorization fail closed? YES
2. Does execution without authorization fail closed? YES
3. Is the executor prevented from executing in both cases? YES
4. Is any false effect confirmation produced? NO
5. Does Conductor remain a coordination component rather than an authority? YES
6. Does Solvent remain the sole authority? YES
7. Did these tests reveal a specification defect? NO
8. Did these tests reveal a genuinely new security property? NO
9. Is a new cross-role protocol primitive now justified? NO

---

## Final verdict
PASS WITH LIMITATIONS

Recommended next experiment (do not implement yet):
Add one real failure-path test where authorization succeeds, execution is attempted with a valid `intent_id`, and the fake executor is instrumented to confirm it is invoked exactly once. This would be the next smallest evidence step before any adversarial suite.
