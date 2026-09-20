chaschel@linux:~/Documents/go/oracle/reference-loop$ go test -v -run TestPhase2FirstRealEffect -timeout 180s
=== RUN   TestPhase2FirstRealEffect
    phase2_test.go:266: [PHASE 2] scenario_id=0a9f1bc9-4835-4646-a7d5-15c9d9d5a53c run_id=phase2-20260912115022
    phase2_test.go:268: [PHASE 2] Step 1: Verify declaration file
    phase2_test.go:269: [DECLARATION] capability_ref=github:ibmendoza/reference-loop-test:deploy
    phase2_test.go:269: [DECLARATION] version=v1.0.0
    phase2_test.go:269: [DECLARATION] content_hash=sha256:88cfe622207fab3dd06bd7ee00f9c1882b1ba3aaf1fa6bd4623e3d4925d6f5c5
    phase2_test.go:269: [DECLARATION] effective_reference=.github/workflows/ref-loop.yml@main
    phase2_test.go:269: [DECLARATION] VERIFIED: all integrity checks passed
    phase2_test.go:271: [PHASE 2] Step 2: Record human-approved plan
    phase2_test.go:280: [PHASE 2] plan task created: ff1f195d-44e8-41e9-90cc-0320e3f7e9ad
    phase2_test.go:282: [PHASE 2] Step 3: Agent claims task
    phase2_test.go:287: [PHASE 2] task claimed: ff1f195d-44e8-41e9-90cc-0320e3f7e9ad status=active
    phase2_test.go:289: [PHASE 2] Step 4: Construct exact operation identity
    phase2_test.go:292: [PHASE 2] operation_id=deploy:ibmendoza/reference-loop-test:ref-loop.yml:main:phase2-20260912115022
    phase2_test.go:294: [PHASE 2] Step 5: Set up Solvent authorization
    phase2_test.go:301: [PHASE 2] belief_id=d9d03591-28db-4f54-a550-bea318295944
    phase2_test.go:338: [PHASE 2] target_id=549bc30b-4cb5-4185-a806-5581aefe3b6b
    phase2_test.go:354: [PHASE 2] Solvent authorization chain complete
    phase2_test.go:356: [PHASE 2] Step 6: AuthorizeAction
    phase2_test.go:371: [PHASE 2] AUTHORIZED: target_id=549bc30b-4cb5-4185-a806-5581aefe3b6b allowed=true
    phase2_test.go:377: [PHASE 2] intent_id=0e3af5df-7277-49b9-948e-07b50af5e603 intent_state=live
    phase2_test.go:379: [PHASE 2] Step 7: ExecuteAction (fresh authorization + real GitHub executor)
    phase2_test.go:393: [PHASE 2] EXECUTE RESULT: allowed=true success=true output="" error="" reason=
    phase2_test.go:403: [PHASE 2] Step 7.5: Poll GitHub SOR for EFFECT_CONFIRMED
    phase2_test.go:404: [PHASE 2] GitHub run 34671440409 status=in_progress conclusion= created=2026-09-12T03:50:23Z
    phase2_test.go:404: [PHASE 2] GitHub run 34671440409 status=completed conclusion=success created=2026-09-12T03:50:23Z
    phase2_test.go:414: [PHASE 2] GitHub SOR: run_id=34671440409 status=completed conclusion=success
    phase2_test.go:421: [PHASE 2] Step 8: Record result in Conductor
    phase2_test.go:433: [PHASE 2] task submitted
    phase2_test.go:435: [PHASE 2] Step 9: Evidence summary
    phase2_test.go:436: [PHASE 2] EVIDENCE:
    phase2_test.go:437: [PHASE 2]   scenario_id: 0a9f1bc9-4835-4646-a7d5-15c9d9d5a53c
    phase2_test.go:438: [PHASE 2]   run_id: phase2-20260912115022
    phase2_test.go:439: [PHASE 2]   execution_run_id: phase2-20260912115022
    phase2_test.go:440: [PHASE 2]   operation_id: deploy:ibmendoza/reference-loop-test:ref-loop.yml:main:phase2-20260912115022
    phase2_test.go:441: [PHASE 2]   belief_id: d9d03591-28db-4f54-a550-bea318295944
    phase2_test.go:442: [PHASE 2]   target_id: 549bc30b-4cb5-4185-a806-5581aefe3b6b
    phase2_test.go:443: [PHASE 2]   intent_id: 0e3af5df-7277-49b9-948e-07b50af5e603
    phase2_test.go:444: [PHASE 2]   task_id: ff1f195d-44e8-41e9-90cc-0320e3f7e9ad
    phase2_test.go:445: [PHASE 2]   project_id: phase2-project-phase2-20260912115022
    phase2_test.go:446: [PHASE 2]   github_run_id: 34671440409
    phase2_test.go:447: [PHASE 2]   declaration_version: v1.0.0
    phase2_test.go:448: [PHASE 2]   declaration_hash: sha256:32e3bd246faaeb0f291ee76bbc2b2b5397a8484f2dabf728f154a60f8986bf9f
    phase2_test.go:449: [PHASE 2]   executor_invocation: 0 (real github_trigger_workflow)
    phase2_test.go:450: [PHASE 2]
    phase2_test.go:451: [PHASE 2] PASS CRITERIA CHECK:
    phase2_test.go:452: [PHASE 2]   PASS-1: Human-approved plan recorded: YES (plan task ff1f195d-44e8-41e9-90cc-0320e3f7e9ad)
    phase2_test.go:453: [PHASE 2]   PASS-2: Exact operation constructed: YES (deploy:ibmendoza/reference-loop-test:ref-loop.yml:main:phase2-20260912115022)
    phase2_test.go:454: [PHASE 2]   PASS-3: Solvent authorized exact operation: YES (Allowed=true)
    phase2_test.go:455: [PHASE 2]   PASS-4: Authorization evidence reached boundary: YES (intent_id=0e3af5df-7277-49b9-948e-07b50af5e603)
    phase2_test.go:456: [PHASE 2]   PASS-5: Fresh authorization at execution boundary: YES (PrepareForAction + ClaimIntent CAS)
    phase2_test.go:457: [PHASE 2]   PASS-6: Authorization binding exact: YES (ClaimIntent CAS succeeded)
    phase2_test.go:458: [PHASE 2]   PASS-7: Real GitHub executor invoked: YES (github_trigger_workflow via GITHUB_TOKEN)
    phase2_test.go:460: [PHASE 2]   PASS-8: GitHub SOR confirmed EFFECT_CONFIRMED: YES (run_id=34671440409 conclusion=success)
    phase2_test.go:466: [PHASE 2]   PASS-9: Operation identity correlated: YES (same ID end-to-end)
    phase2_test.go:467: [PHASE 2]   PASS-10: Conductor records coordination only: YES
    phase2_test.go:468: [PHASE 2]   PASS-11: Solvent unmodified: YES (verified in Phase 1.5)
    phase2_test.go:469: [PHASE 2]
    phase2_test.go:472: [PHASE 2] VERDICT: PASS
--- PASS: TestPhase2FirstRealEffect (12.90s)
PASS
ok  	github.com/PithomLabs/oracle/reference-loop	12.909s
