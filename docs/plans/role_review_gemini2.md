Matrix v0.2 makes real structural progress by explicitly decoupling semantic ownership from enforcement responsibility and establishing clear rules for transport metadata versus semantic payloads. However, by continuing to forbid a centralized runtime or state engine, it delegates complex distributed-system enforcement to scattered adapter wrappers, creating dangerous audit gaps and dangling state traps.

**Core Architectural Traps & Structural Vulnerabilities**

* **The Integration Implementation Drift Risk (Section 8):** Delegating authorization enforcement to individual integration adapters while keeping semantic ownership in Solvent introduces major security risks. Requiring 20+ disparate tool wrappers to independently implement token parsing, age-policy evaluation, and operation matching guarantees logic drift and creates a massive, fragmented attack surface for authorization bypasses.
* **The Phantom Authorization Hanging Trap (Sections 10 & 13):** Section 10 mandates an actionable policy for `Authorization UNKNOWN` without introducing a dedicated `Timeout` or `Abandoned` workflow state. If an operator cancels a stuck task via Conductor, Section 13 dictates that Conductor cancellation *cannot* mutate Solvent authority. As a result, the authorization request remains hanging indefinitely in Solvent while Conductor moves on, creating orphaned authority contexts.
* **The "Immaterial Parameter" Exploitation Vector (Section 7.3):** Allowing individual integration declarations to designate parameters as "immaterial transport metadata" creates a subtle parameter-tampering vulnerability. If a tool developer mistakenly classifies fields like `redirect_uri`, `callback_url`, or `user_context` as transport metadata, an attacker can manipulate downstream execution without breaking the exact operation binding match.

**Matrix v0.2 Claims vs. Operational Reality**

| Spec Provision | Architectural Reality | Operational Failure |
| --- | --- | --- |
| **Integrations Are Merely Adapters** | Wrappers must evaluate tokens, parse canonical payloads, and handle retries locally. | Adapters act as unmonitored mini-gateways, scattering security enforcement. |
| **Strict Conductor/Solvent Isolation** | Canceling a task in Conductor does not alter Solvent state. | Generates zombie authorization requests in Solvent whenever Conductor tasks abort. |
| **Local Immaterial Declarations** | Each integration decides which payload fields do not affect identity. | Inconsistent parameter stripping across tools causes false-positive rejections ($X \neq Y$). |

**Mandatory Hardening Directives**

* **Provide a Standardized Verification Module:** Instead of letting every integration write custom evidence-verification logic, supply a single, imported SDK module that standardizes RFC 8785 JSON canonicalization and token verification across all adapters.
* **Enable Express Authorization Invalidation:** Allow Conductor task cancellations to emit an explicit revocation event to Solvent so aborted workflows never leave dangling, un-consumed authority tokens.
* **Centralize Immaterial Parameter Approvals:** Require Solvent—as the ultimate authority owner—to validate and approve an integration's declared "immaterial" parameter list during registration to ensure sensitive payload variables cannot be stripped.