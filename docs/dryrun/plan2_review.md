I reviewed `plan2`. The overall direction is sound: the audit identifies real gaps, keeps the two-tool MCP surface, adds no services or databases, and leaves the Solvent kernel untouched.  

I would **not hand this exact plan to the coding agent yet**. There are several important design corrections.

### 1. The corpus schema has a structural bug

The proposed table uses:

```sql
UNIQUE (scenario_id, issue_number)
```

but ingestion assigns `issue_number` as the index **within each document**. With seven documents, every document will have chunk `1`, chunk `2`, etc. Those rows will collide.

The plan therefore cannot actually ingest the seven documents as written. The relevant schema and ingestion design are in the file.  

Use something like:

```sql
artifact_id   TEXT NOT NULL,
chunk_index   INT NOT NULL,

UNIQUE (scenario_id, artifact_id, chunk_index)
```

and retain the artifact identity on every passage.

I would make the minimal corpus record:

```sql
corpus_issue (
    id,
    scenario_id,
    artifact_id,
    chunk_index,
    title,
    body,
    state,
    stable_locator,
    content_sha256,
    provenance_class,
    ingested_at,
    embedding
)
```

You can keep the existing `corpus_issue` name for compatibility, but its rows are really corpus passages now.

### 2. `CorpusHit.provenance` cannot work with the proposed table

The response type explicitly returns:

```go
Provenance string `json:"provenance"`
```

but the proposed `corpus_issue` schema contains no provenance column. 

This is more than cosmetic because the plan explicitly relies on provenance honesty for `adv-review*.md`. 

Store `provenance_class` on the passage row, or derive it deterministically from a registered artifact table. For this POC, storing it directly is simpler.

### 3. Do not mutate `bmist@1.0.0` in place

This is the biggest provenance problem in the plan.

The document proposes adding two debt types while retaining:

```text
bmist@1.0.0
```

because the pack is internal to the POC. 

That undermines exactly what Solvent is supposed to preserve: reproducibility of past epistemic state. A persisted belief saying `bmist@1.0.0` could mean either the old six-debt methodology or the new eight-debt methodology depending on when it was read.

For this research program, I would use:

```text
bmist-as@0.1.0
```

because the methodology now explicitly covers **BM + IST + Asymptotic Safety**, rather than silently changing the meaning of BM–IST's existing pack.

At minimum:

```text
bmist@1.1.0
```

must be used. The important point is: **never mutate the semantics of an already-referenced pack version.**

### 4. Seed the debt vocabulary from the pack

The plan duplicates the eight debt names in both `pack.json` and `seed.go`:

```text
needMap
needInvariant
...
needRegularity
```

That creates another drift point.

The seed should obtain:

```go
initialDebt := pack.GetInitialDebt()
```

and persist that.

Then the pack is the source of truth for methodology, while the seed merely instantiates it.

This also makes the planned replacement of hardcoded debt validation with registry lookup much cleaner. 

### 5. Keep `argus mcp` as the stdio process

The plan recommends:

> add `--mcp-stdio` to `argus serve` so a single process serves both HTTP and MCP

but earlier the architecture correctly establishes:

```text
argus serve
argus mcp
```

The important distinction is:

> **one binary does not require one OS process.**

MCP stdio is naturally a dedicated process because OpenCode owns its stdin/stdout stream. I would implement:

```text
argus serve
    HTTP + Trust UI + DB bootstrap

argus mcp
    stdin/stdout JSON-RPC MCP transport
```

Both are the same ARGUS binary, same code, same database, same kernel, same migrations.

That still satisfies the architectural goal of one deployable system and avoids contaminating the HTTP server lifecycle with stdio semantics.

Also, absolutely keep logs off stdout in `argus mcp`; stdout must remain protocol-only.

### 6. Resolve the seed contract before implementation

The plan correctly identifies this as an open issue:

```text
governance_ref type
```

but I would move it out of "open questions" and make it a prerequisite. 

The coding agent should inspect the actual current schema and application contract and establish one canonical representation before writing `seed.go`.

Likewise, confirm exactly how `scenario_id` is created/resolved. The seed section assumes it exists but does not make its lifecycle explicit. 

### 7. Make the "context only, never authoritative evidence" rule testable

The plan says the `agent_derived` corpus material enters as research context and never as authoritative evidence. That's correct. 

But right now that is primarily a statement, not an enforced invariant.

Add one verification:

```text
Corpus retrieval may inform an agent.
Corpus retrieval alone cannot satisfy a debt retirement rule.
Corpus provenance must never be silently converted into authoritative evidence.
```

The acceptance test should deliberately retrieve an `agent_derived` adversarial review and verify that it cannot discharge a human-gated or reproducible-artifact debt.

### 8. Hash embeddings are acceptable, but name them correctly

The POC embedder is reasonable for zero-credential startup. The plan's FNV/n-gram approach is essentially a **lexical retrieval mechanism encoded as vectors**, not semantic embedding.

That distinction matters because otherwise someone will interpret "vector search works" as "semantic retrieval works."

I would change the language to:

```text
Deterministic lexical-vector retrieval for zero-friction POC.
```

Then acceptance should test known queries against expected seed passages, not merely assert that cosine search returns rows.

### Revised shape

After those corrections, the implementation remains impressively small:

```text
OpenCode
    |
    | MCP stdio
    v
argus mcp
    |
    +--> existing application layer
    |
    +--> Solvent kernel
    |
    +--> corpus retrieval
    |
    v
CockroachDB

argus serve
    |
    +--> same application
    +--> same Solvent
    +--> same DB
    +--> Trust UI
```

with:

```text
Domain Pack
    bmist-as@0.1.0
        |
        +-- 8 debt classes
        +-- retirement rules
        +-- initial debt
```

and:

```text
Corpus
    artifact
       |
       +-- passage 0
       +-- passage 1
       +-- passage 2
       ...
```

### Verdict on `plan2`

The **architecture is approved**, but the plan needs these corrections before implementation:

1. Fix corpus identity/chunk uniqueness.
2. Persist provenance on corpus passages.
3. Version the domain pack immutably.
4. Derive seed debt from the pack.
5. Keep `argus mcp` as dedicated stdio mode rather than merging stdio into `serve`.
6. Resolve `scenario_id` and `governance_ref` contracts before coding.
7. Turn corpus provenance boundaries into executable tests.

Everything else in the plan is consistent with the desired subtraction principle: no new service, no third MCP tool, no Solvent kernel rewrite, no ADD_DEBT, and no preconstructed research graph. 

The next coding-agent prompt should incorporate those seven changes and otherwise preserve `plan2` rather than expanding its scope.
