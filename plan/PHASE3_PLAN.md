# PHASE 3 PLAN — EBP RESEARCH PACKET v1 CONTRACT

## Files to Create

```
oracle/packet/v1/
├── schema.json
├── types.go
├── resolve.go
├── validate.go
├── validate_test.go
└── examples/
    ├── valid_packet.json
    ├── malformed_missing_packet_id.json
    ├── malformed_empty_claim.json
    ├── malformed_unknown_claim_type.json
    ├── malformed_free_text_edge.json
    ├── malformed_bad_reference.json
    └── malformed_bad_debt.json
```

## Step 1: `types.go` — Go Types

```go
package packetv1

// Packet is the top-level EBP Research Packet v1.
type Packet struct {
    SchemaVersion string      `json:"schema_version"`
    Role          string      `json:"role"`
    PacketID      string      `json:"packet_id"`
    PackRef       PackRef     `json:"pack_ref"`
    ScenarioID    string      `json:"scenario_id,omitempty"`
    RunID         string      `json:"run_id,omitempty"`
    TaskRef       string      `json:"task_ref,omitempty"`
    ProjectRef    string      `json:"project_ref,omitempty"`
    ScenarioRef   string      `json:"scenario_ref,omitempty"`
    Scope         string      `json:"scope,omitempty"`
    CorpusRef     string      `json:"corpus_ref,omitempty"`
    Agent         Agent       `json:"agent"`
    Beliefs       []Belief    `json:"beliefs"`
    Evidence      []Evidence  `json:"evidence"`
    Debts         []Debt      `json:"debts,omitempty"`
    Edges         []Edge      `json:"edges,omitempty"`
    Tasks         []Task      `json:"tasks,omitempty"`
}

type PackRef struct {
    PackID  string `json:"pack_id"`
    Version string `json:"version"`
}

type Agent struct {
    ID    string `json:"id,omitempty"`
    Model string `json:"model,omitempty"`
    Role  string `json:"role,omitempty"`
}

type Belief struct {
    LocalID   string   `json:"local_id"`
    Claim     string   `json:"claim"`
    ClaimType string   `json:"claim_type"`
    Debt      []string `json:"debt,omitempty"`
}

type Evidence struct {
    LocalID         string `json:"local_id"`
    BeliefRef       string `json:"belief_ref"`
    ProvenanceClass string `json:"provenance_class"`
    ContentSHA256   string `json:"content_sha256"`
    SourceURL       string `json:"source_url,omitempty"`
    ArtifactRef     string `json:"artifact_ref,omitempty"`
}

type Debt struct {
    LocalID string   `json:"local_id"`
    Items   []string `json:"items"`
}

type Edge struct {
    LocalID string `json:"local_id"`
    FromRef string `json:"from_ref"`
    ToRef   string `json:"to_ref"`
    Kind    string `json:"kind"`
}

type Task struct {
    LocalID       string `json:"local_id"`
    Title         string `json:"title"`
    Description   string `json:"description,omitempty"`
    GovernanceRef string `json:"governance_ref,omitempty"`
}

// Reference prefixes
const (
    RefPrefixLocal     = "local:"
    RefPrefixCanonical = "canonical:belief:"
)

// Edge kinds
const (
    EdgeDerives     = "derives"
    EdgeContradicts = "contradicts"
)

// Roles
const (
    RoleWork        = "work"
    RoleAdversarial = "adversarial"
)

// Schema version
const SchemaVersion = "ebp-research-packet/v1"
```

## Step 2: `resolve.go` — Reference Resolution

```go
package packetv1

import (
    "fmt"
    "strings"
)

type ReferenceKind string

const (
    RefKindLocal     ReferenceKind = "local"
    RefKindCanonical ReferenceKind = "canonical"
)

type ResolvedReference struct {
    Kind ReferenceKind
    ID   string // for local: the local_id; for canonical: the UUID
    Raw  string // original reference string
}

// ParseReference parses "local:b1" or "canonical:belief:<uuid>"
func ParseReference(ref string) (ResolvedReference, error) {
    if strings.HasPrefix(ref, RefPrefixLocal) {
        id := strings.TrimPrefix(ref, RefPrefixLocal)
        if id == "" {
            return ResolvedReference{}, fmt.Errorf("empty local reference")
        }
        return ResolvedReference{Kind: RefKindLocal, ID: id, Raw: ref}, nil
    }
    if strings.HasPrefix(ref, RefPrefixCanonical) {
        id := strings.TrimPrefix(ref, RefPrefixCanonical)
        if !isValidUUID(id) {
            return ResolvedReference{}, fmt.Errorf("invalid canonical UUID: %s", id)
        }
        return ResolvedReference{Kind: RefKindCanonical, ID: id, Raw: ref}, nil
    }
    return ResolvedReference{}, fmt.Errorf("invalid reference format: %s", ref)
}

func isValidUUID(s string) bool {
    if len(s) != 36 {
        return false
    }
    positions := []int{8, 13, 18, 23}
    for _, p := range positions {
        if s[p] != '-' {
            return false
        }
    }
    for i, c := range s {
        if i == 8 || i == 13 || i == 18 || i == 23 {
            continue
        }
        if !((c >= '0' && c <= '9') || (c >= 'a' && c <= 'f') || (c >= 'A' && c <= 'F')) {
            return false
        }
    }
    return true
}
```

## Step 3: `validate.go` — Packet Validation

Key functions:
- `Validate(pkt *Packet, registry *domainpack.PackRegistry) error` — full validation
- `validateReferences(pkt *Packet) error` — all local/canonical refs
- `validateBeliefs(pkt *Packet, pack Pack) error`
- `validateEvidence(pkt *Packet, pack Pack) error`
- `validateEdges(pkt *Packet) error`
- `validateTasks(pkt *Packet) error`

Validation order (fail-closed):
1. Schema version check
2. packet_id non-empty
3. role ∈ {work, adversarial}
4. pack_ref resolves through registry
5. Belief structure (local_id unique, claim non-empty, claim_type from pack)
6. Evidence structure (local_id unique, belief_ref resolves, content_sha256 valid hex, provenance_class from pack)
7. Debt vocabulary membership
8. All references resolve
9. Edge structure (from_ref/to_ref resolve, kind ∈ {derives, contradicts})
10. Task structure (local_id unique, title non-empty, no forbidden status)

## Step 4: `schema.json` — JSON Schema

Structural constraints only. Cross-object semantic checks in Go code.

## Step 5: Test Fixtures

**Valid:** `valid_packet.json` — one belief, one evidence, one edge, one task, valid pack_ref

**Invalid (6 fixtures):**
- `malformed_missing_packet_id.json` — no packet_id
- `malformed_empty_claim.json` — belief with empty claim
- `malformed_unknown_claim_type.json` — claim_type not in pack
- `malformed_free_text_edge.json` — edge from_ref is free text
- `malformed_bad_reference.json` — evidence belief_ref is malformed
- `malformed_bad_debt.json` — debt item not in pack vocabulary

## Step 6: `validate_test.go` — ~20 tests

All negative tests from the frozen plan + positive tests.

## Key Design Decisions

1. **PackRef as object** — `{pack_id, version}` not a string. Cleaner parsing.
2. **Reference format** — `local:id` and `canonical:belief:uuid`. Simple prefix parsing.
3. **Debt as separate object** — `Debt{LocalID, Items}` allows packet to declare debt separately from beliefs.
4. **No external dependencies** — JSON Schema validation done in Go code, not a library.
5. **Domain neutral** — packet code imports `domainpack.Pack` interface, not `bmistv1.Pack` directly.

## Verification

```bash
cd /home/chaschel/Documents/go/oracle
go test ./...
go test -race ./...
go vet ./...
```
