The underlying problem is **not the fixed port itself**. The real defect is that the Reference Loop test does not own the lifecycle of the Solvent test server reliably.

It currently appears to:

```text
startSolventServer()
    ↓
spawn ListenAndServe in goroutine
    ↓
continue immediately
```

while a previous failed test can leave a listener alive, and the new server can fail to bind without the test detecting it. The readiness probe can then accidentally connect to the **old** server.

So the fix should establish a proper lifecycle contract:

```text
create listener
    ↓
bind succeeds or test fails immediately
    ↓
start server on that exact listener
    ↓
wait until that exact server is accepting connections
    ↓
run test
    ↓
shutdown server
    ↓
wait for shutdown
```

Use this prompt:

```text id="9z8w4t"
The current Phase 2 failure is NOT a simple "hardcoded port" problem.

Root cause:

    startSolventServer() launches ListenAndServe asynchronously without
    reliably owning the listener lifecycle or surfacing bind/startup failure.

This permits a stale server from a previous failed run to occupy port 18080.
A new server then fails to bind, but the test's readiness probe connects to the
OLD server and the test later fails unpredictably.

Do NOT apply a band-aid such as merely replacing 18080 with a random port.

Fix the underlying server lifecycle in the Reference Loop.

======================================================================
OBJECTIVE
======================================================================

Make every Phase 2 test instance own exactly one Solvent test server:

    bind
      ↓
    start
      ↓
    verify readiness of THAT server
      ↓
    use server
      ↓
    shutdown
      ↓
    verify shutdown

A failed server startup must fail the test immediately.

A stale server from another run must never be mistaken for the current server.

All changes remain confined to:

    /home/chaschel/Documents/go/oracle/reference-loop

Never modify frozen Solvent.

======================================================================
1. BIND THE LISTENER SYNCHRONOUSLY
======================================================================

Do not let ListenAndServe choose/bind the port inside a detached goroutine.

Create the TCP listener explicitly first:

    net.Listen("tcp", "127.0.0.1:0")

The OS chooses an available ephemeral port.

If net.Listen fails:

    return the error immediately
    fail the test
    do not continue

The returned listener is the authoritative server endpoint.

Derive:

    solventAddr
    solventPort

from that exact listener.

======================================================================
2. START HTTP SERVER FROM THAT EXISTING LISTENER
======================================================================

Construct the http.Server explicitly.

Start it using the already-bound listener:

    go server.Serve(listener)

Do NOT use:

    server.ListenAndServe()

This eliminates the race between "port appears free" and actual bind.

The listener binding must already have succeeded before the test continues.

======================================================================
3. RETURN OWNERSHIP OF THE SERVER
======================================================================

Replace any API that only returns a loosely-started server with a test-owned
server handle containing at minimum:

    httpServer
    listener
    address
    shutdown function
    startup/shutdown error channel

Conceptually:

    type testSolventServer struct {
        server   *http.Server
        listener net.Listener
        addr     string
    }

The test owns the lifecycle.

Do not rely on global process state.

======================================================================
4. READINESS MUST VERIFY THE CURRENT SERVER
======================================================================

Do not merely test:

    GET http://localhost:<port>

because that can still accidentally reach the wrong process if the address is
shared.

The address must come from THIS test's listener.

After Serve starts, perform a small readiness check against the exact bound
address.

The test must fail if:

    Serve exits unexpectedly
    listener closes
    startup error occurs
    readiness timeout expires

Do not let the test continue after startup failure.

======================================================================
5. SURFACE SERVE ERRORS
======================================================================

Capture the result of:

    server.Serve(listener)

through an error channel.

Treat:

    http.ErrServerClosed

as normal shutdown.

Any other error before or during startup must be visible to the test.

The previous failure mode:

    bind fails in goroutine
    test unknowingly talks to stale server

must become impossible.

======================================================================
6. SHUTDOWN IN TEST CLEANUP
======================================================================

Every test that starts Solvent must register cleanup immediately:

    t.Cleanup(func() {
        shutdown test server
    })

Shutdown must:

    stop accepting new connections
    wait for Serve to return
    close the listener
    confirm shutdown completed

Do not rely on the OS eventually reclaiming the port.

Do not use process-kill hacks.

======================================================================
7. NO GLOBAL FIXED PORT
======================================================================

The test server must use:

    127.0.0.1:0

or an equivalent ephemeral-port mechanism.

Do not hardcode:

    18080

Do not create a shared global server across tests unless that is already an
intentional fixture owned by the test package.

Each test should be isolated.

======================================================================
8. PREVENT STALE-SERVER FALSE POSITIVES
======================================================================

Add a testable startup invariant:

    the test may only proceed after its own listener has been successfully
    bound and its own HTTP server has started.

Do not use external process discovery to decide whether Solvent is running.

Do not probe a well-known port and assume that the responding server belongs
to the current test.

======================================================================
9. APPLY TO EXISTING PHASE 2 SETUP

Refactor:

    setupPhase2Env
    startSolventServer

and any related Phase 1.5 server helpers so they use the same lifecycle model.

Do not create a second server framework.

Reuse the same helper for:

    Phase 1 happy path
    Phase 1.5
    Phase 2

where practical.

The goal is one correct test-server lifecycle abstraction, not separate fixes
for each test.

======================================================================
10. SOLVENT FREEZE
======================================================================

Before and after implementation:

    cd /home/chaschel/Documents/go/solvent-main

verify:

    git rev-parse HEAD
    git status --short
    git diff --name-only

Expected:

    HEAD = 7602699
    working tree clean
    no source modifications

If Solvent is modified:

    STOP
    report BLOCKED

======================================================================
11. VALIDATION

First run:

    go test ./...

Then specifically:

    go test -v -run TestPhase2FirstRealEffect -timeout 180s

Run the Phase 2 test repeatedly, including at least:

    two consecutive runs
    one run after an induced/previous startup failure if possible

The purpose is to prove there is no stale-server contamination.

Also run:

    go test -race ./...

if the package permits it.

======================================================================
12. REQUIRED FAILURE BEHAVIOR

Verify these cases fail cleanly:

    listener cannot bind
        → test fails immediately

    server Serve() exits unexpectedly
        → test fails with server error

    readiness timeout
        → test fails

    shutdown does not complete
        → test reports cleanup failure

Never silently reconnect to another Solvent server.

======================================================================
13. DO NOT CHANGE PHASE 2 SEMANTICS

Do not alter:

    operation identity
    declaration
    authorization chain
    Conductor semantics
    GitHub executor
    GitHub SOR logic
    PASS criteria

This task is ONLY about reliable test-server lifecycle ownership.

======================================================================
14. FINAL REPORT

After fixing the lifecycle:

Report:

    Root cause:
        server lifecycle / listener ownership

    Fixed:
        YES / NO

    Ephemeral listener:
        YES / NO

    Startup errors surfaced:
        YES / NO

    Readiness tied to current server:
        YES / NO

    Cleanup guaranteed:
        YES / NO

    Repeated Phase 2 runs stable:
        YES / NO

    Solvent modified:
        NO

Then STOP.

Do not proceed to Phase 3 or redesign the workflow.
```

The key architectural correction is **server lifecycle ownership**, not randomizing the port. Random ports are a consequence of that design, not the solution itself.
