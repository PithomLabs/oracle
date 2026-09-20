## bug


chaschel@linux:~/Documents/go/oracle/bin$ ./argus serve
2026/09/20 17:25:33 bootstrap.go:117: ARGUS: started CockroachDB (PID 8730) on :26257
*
* WARNING: ALL SECURITY CONTROLS HAVE BEEN DISABLED!
* 
* This mode is intended for non-production testing only.
* 
* In this mode:
* - Your cluster is open to any client that can access any of your IP addresses.
* - Intruders with access to your machine or network can observe client-server traffic.
* - Intruders can log in without password and read or write any data in the cluster.
* - Intruders can consume all your server's resources and cause unavailability.
*
*
* INFO: To start a secure server without mandating TLS for clients,
* consider --accept-sql-without-tls instead. For other options, see:
* 
* - https://go.crdb.dev/issue-v/53404/v26.2
* - https://www.cockroachlabs.com/docs/v26.2/secure-a-cluster.html
*
*
* WARNING: Running a server without --sql-addr, with a combined RPC/SQL listener, is deprecated.
* This feature will be removed in a later version of CockroachDB.
*
CockroachDB node starting at 2026-09-20 09:25:35.382446166 +0000 UTC m=+1.437355349 (took 0.8s)
build:               CCL v26.2.0 @ 2026/04/21 18:36:57 (go1.25.5)
malloc_conf:         prof:false,prof_active:true,background_thread:true,thp:never,metadata_thp:disabled,narenas:64,dirty_decay_ms:2000,muzzy_decay_ms:0
webui:               http://linux:8081
sql:                 postgresql://root@linux:26257/defaultdb?sslmode=disable
sql (JDBC):          jdbc:postgresql://linux:26257/defaultdb?sslmode=disable&user=root
RPC client flags:    /home/chaschel/.local/bin/cockroach <client cmd> --host=linux:26257 --insecure
logs:                /home/chaschel/Documents/go/oracle/bin/.cockroach-data/logs
temp dir:            /home/chaschel/Documents/go/oracle/bin/.cockroach-data/cockroach-temp1308587809
external I/O path:   /home/chaschel/Documents/go/oracle/bin/.cockroach-data/extern
store[0]:            path=/home/chaschel/Documents/go/oracle/bin/.cockroach-data
clusterID:           e48ff006-a563-49dd-82bf-21c4083b5764
status:              restarted pre-existing node
nodeID:              1
2026/09/20 17:25:35 bootstrap.go:163: ARGUS: database "argus" ready
2026/09/20 17:25:36 bootstrap.go:184: ARGUS: migrations applied
2026/09/20 17:25:36 seed.go:46: seed: inserting BM-IST-AS POC seed...
2026/09/20 17:25:36 seed.go:99: seed: BM-IST-AS POC seed inserted (project=00000000-0000-0000-0000-000000000100, task=00000000-0000-0000-0000-000000000101, belief=00000000-0000-0000-0000-000000000102)
2026/09/20 17:25:36 main.go:115: ARGUS_OPERATOR_TOKEN not set, using default local operator token
2026/09/20 17:25:36 main.go:145: argus serve listening on :8080
2026/09/20 17:25:48 server.go:3679: http: superfluous response.WriteHeader call from github.com/PithomLabs/oracle/internal/ui.(*Server).HandleInsights (ui.go:131)



## prompt

That changes the diagnosis materially. **The seed and bootstrap are working; the failure is in the Insights projection/rendering path.**

Your log proves:

```text
database ready
migrations applied
seed inserted
argus serve listening on :8080
```

Then immediately:

```text
http: superfluous response.WriteHeader call
from internal/ui.(*Server).HandleInsights (ui.go:131)
```

And the browser shows:

```text
template: insights.html:10:6:
executing "insights.html" at <.ID>:
nil pointer evaluating interface {}.ID
```

So I would **not touch the seed, MCP, pack, or Solvent at this point**.

### What is actually broken

The Insights handler is passing a data structure containing a **nil element/interface** into `insights.html`.

The template reaches:

```text
<ID>
```

while the current template object is nil.

Then `HandleInsights` apparently tries to write an HTTP error after template execution has already started writing the response, which explains:

```text
superfluous response.WriteHeader call
```

So there are two bugs, with one being downstream of the other:

```text
Insights query/projection
        ↓
nil object in template model
        ↓
template execution fails at .ID
        ↓
handler attempts error response
        ↓
superfluous WriteHeader
```

### The important part

Don't "fix" this by simply changing the template to:

```go
{{if .}}{{.ID}}{{end}}
```

That could hide a broken projection.

The first thing to establish is **why the Insights view model contains nil**.

Given that the seed now creates:

```text
project
task
belief
operator principal
```

I would inspect:

```text
internal/ui/ui.go
templates/insights.html
```

and trace:

```text
HandleInsights
    → query/load data
    → construct Insights view model
    → template.Execute(...)
```

Specifically determine which object corresponds to the template's `.ID` on line 10 and why it can be nil.

### Give the coding agent this exact instruction

```text
The first-run seed is working, but /ui/insights is broken.

Observed runtime:

ARGUS starts successfully:
- CockroachDB ready
- migrations applied
- BM-IST-AS seed inserted
- HTTP server listening on :8080

GET /ui/insights then fails with:

template: insights.html:10:6:
executing "insights.html" at <.ID>:
nil pointer evaluating interface {}.ID

and server logs:

http: superfluous response.WriteHeader call
from internal/ui.(*Server).HandleInsights (ui.go:131)

Investigate and fix the root cause in the Insights projection/rendering path.

Do NOT modify:
- seed logic
- MCP
- domain pack
- Solvent kernel
- research architecture

First inspect:
- internal/ui/ui.go, especially HandleInsights around line 131
- templates/insights.html, especially line 10
- the query/model construction feeding the template

Determine exactly which template value is nil and why.

Do NOT merely add template nil guards unless the nil is an intentionally valid state.
Prefer fixing the view-model/query construction so the template receives valid objects.

Also fix the HTTP error handling so a template execution error does not result in
a second WriteHeader after the response has already begun.

Add a regression test that:
1. starts from a clean seeded database,
2. renders /ui/insights,
3. verifies HTTP 200,
4. verifies no template execution error,
5. verifies the seeded task/belief actually appears.

Then run:
go test ./...

Finally run argus serve against a clean local state and manually verify:
http://localhost:8080/ui/insights

Do not proceed to additional feature work until this is fixed.
```

### One subtle point

The browser output showing the error **inside the table** strongly suggests the template has already emitted the header before failing. That is why the `WriteHeader` warning appears. So the handler should probably render into a buffer first and only write the successful response after `Execute` succeeds, or otherwise use the repository's existing safe-rendering pattern.

That is a proper robustness fix, not merely cosmetic.

**The good news:** this is now a localized UI bug. The underlying bootstrap/seed architecture has already passed the more important initialization checks.
