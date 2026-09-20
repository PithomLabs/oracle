Proceed now. Exit planning mode and implement the developer-experience fix.

The architecture and scope are frozen. Do not create another plan, perform another architecture review, or add new infrastructure.

Implement Parts 1–10 exactly as already specified:

1. Embed BM-IST pack so `argus serve` does not depend on the current working directory.
2. Add deterministic local operator-token default `argus-local-operator`, while preserving `ARGUS_OPERATOR_TOKEN` override and authentication enforcement.
3. Complete the credential/config audit and remove obsolete manual prerequisites without weakening real authority boundaries.
4. Verify/fix DB URL resolution so sane local defaults work automatically and missing CockroachDB produces one clear actionable error.
5. Make startup ordering deterministic.
6. Verify clean-machine startup has no hidden cwd/config prerequisites.
7. Update minimal developer documentation with local defaults.
8. Add focused tests for all new defaults/startup behavior.
9. Run:
   go test ./...
   go test -race ./...
   go vet ./...
10. Perform an actual clean-start smoke test of `argus serve`.

Critical acceptance test:

    git clone
    go build ./cmd/argus
    argus serve

must reach a usable local Trust UI without requiring the developer to manually locate the domain-pack, invent an operator token, or supply configuration that the POC can safely default.

Preserve:

- one process
- one database
- two MCP tools
- embedded Trust UI
- fixed ports
- atomic task lock
- task dev = non-destructive
- task fresh = destructive
- server-derived PrincipalID
- server-selected PackRef
- all existing authority/security invariants

Do not optimize further after the acceptance test passes.

At the end, report:
- files changed
- defaults introduced
- obsolete configuration removed
- test results
- exact clean-start command/result
- final configuration/credential matrix