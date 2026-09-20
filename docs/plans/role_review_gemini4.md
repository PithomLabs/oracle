Loop Engineering v0.4 patch-fixes operation-identity comparison and default capability ownership, but doubles down on a no-runtime dogma that opens severe offline revocation loopholes and creates operational state-machine sprawl.

**Core Architectural Exploits & Vulnerabilities**

* **The Offline Authorization Revocation Loophole (§9.2):** Section 9.2 permits independently verifiable authorization evidence to remain valid even when the authority service is unreachable. If Solvent issues an emergency revocation during a network partition, an Executor will still execute high-consequence operations using un-revokable offline evidence.


* **The Asynchronous Outcome Deadlock (§14):** Section 14 explicitly forbids reporting async handoff receipts or queue acceptances as `Succeeded` unless explicitly declared in the capability model. Because Loop Engineering forbids a centralized workflow engine, every individual integration must build its own background state machine to poll external Systems of Record (SOR) until completion.


* **The Orphaned Token Window (§23, §27):** Strictly prohibiting Conductor task cancellations from mutating Solvent authority leaves active authorization tokens valid until their natural TTL expires. Cancelling a task in Conductor leaves active authorization evidence in the wild, enabling delayed or rogue executions.


* **The "Deployment Operator" Accountability Trap (§8):** Defaulting unowned capability declarations to the deployment operator creates an administrative dumping ground. Operations teams become implicitly responsible for declaring semantic effect surfaces and operation-identity comparison rules for third-party tools they did not build.



**v0.4 Protocol Specs vs. Operational Realities**

| Spec Refinement | Stated Protocol Goal | Operational Failure Mode |
| --- | --- | --- |
| **Offline Auth Evidence (§9.2)** | Prevents false-closed failures when Solvent is unreachable.

 | Bypasses emergency revocations during authority outages, executing revoked tasks.

 |
| **Strict Success Proof (§14)** | Prevents premature completion reporting on async queue acceptances.

 | Offloads distributed state-tracking onto individual custom adapters.

 |
| **Operator Ownership Default (§8)** | Eliminates ownerless capability declarations.

 | Shifts software design responsibility to ops teams, leading to stale capability declarations.

 |

**Mandatory Hardening Directives**

* **Require Live Revocation Checks for High-Consequence Operations:** Override offline evidence validity (§9.2) for high-impact operations by requiring real-time verification against Solvent.


* **Standardize SOR Event Callbacks:** Provide a standard event-driven callback contract for external SORs to push execution state updates, avoiding custom polling logic across adapters.


* **Add Cross-Boundary Revocation Signals:** Allow Conductor task cancellations to issue an explicit revocation signal to Solvent so cancelled workflows immediately invalidate downstream tokens.