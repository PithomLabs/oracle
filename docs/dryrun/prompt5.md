Yes. I would make **Debts a full dashboard as well**, and I would approve Plan 3 with one refinement: **scope dashboard data to the active research scenario/project rather than literally all scenarios**.

The file correctly identifies the architectural mismatch: `GetContext` is a task-scoped RCP operation, while Insights is a human dashboard. Using `GetContext("default")` was the wrong abstraction, and the same logic applies to Debts. 

### Answer to the open question

Use:

```text
/ui/insights
    → dashboard
    → all relevant tasks
    → all relevant beliefs

/ui/debts
    → dashboard
    → all relevant beliefs with debt state
```

Do **not** make Debts require a `task_id` merely to make the current implementation work.

The reason is visible in the existing template design: Debts already presents beliefs as a flat collection, so forcing it through a single-task `GetContext` is artificial. The plan itself notes this. 

### One correction to Plan 3

I would change:

```go
GetAllBeliefs(ctx)
```

from literally:

```sql
SELECT ...
FROM belief
```

across every scenario, to whatever mechanism the application already uses to identify the **current research scenario/project**.

Otherwise today's single-scenario POC works, but a future second scenario could cause `/ui/debts` and `/ui/insights` to silently mix unrelated epistemic states.

The desired semantics should be:

```text
Human dashboard
    ↓
current ARGUS research workspace/scenario
    ↓
all tasks / all beliefs in that scope
```

while:

```text
Agent RCP
    ↓
specific task
    ↓
GetContext(taskID)
```

That preserves the distinction cleanly:

```text
GetContext      = task context for agents
GetDashboard    = workspace context for humans
```

### Plan 3 itself is otherwise sound

The proposed fix gets the important things right:

* leave `GetContext` unchanged;
* add a dashboard-specific application method;
* eliminate the fake `"default"` task;
* render into a buffer before writing the HTTP response;
* add a regression test for `/ui/insights` without `task_id`;
* apply the same dashboard pattern to Debts. 

The buffering change is particularly worth keeping because it fixes the secondary `superfluous response.WriteHeader` failure rather than merely masking the template crash. 

### I would answer the coding agent's open question this way

```text
Make /ui/debts a full dashboard, consistent with /ui/insights.

Both are human-facing dashboard projections and should not depend on
GetContext(task_id).

Use a dashboard-specific application query for both.

However, do not interpret "all beliefs" as all beliefs across every
scenario in the database. Scope the dashboard to the active/current
research scenario or workspace using the repository's existing scenario
selection mechanism.

Keep GetContext task-scoped and unchanged for MCP/RCP agents.

Apply the same buffered template rendering and HTTP error handling to both
Insights and Debts.

Add regression coverage for:
- GET /ui/insights with no task_id
- GET /ui/debts with no task_id
- both pages return 200 on a seeded environment
- seeded task/beliefs/debts are visible
- template failures do not partially write a 200 response
```

That gives you the clean architecture:

```text
                  ARGUS application
                         |
          ┌──────────────┴──────────────┐
          |                             |
      GetContext                   GetDashboard
          |                             |
       Agents                    Human Trust UI
          |                             |
      one task              current research scope
```

I would approve **Plan 3 with that scenario-scoping amendment** and then implement it.
