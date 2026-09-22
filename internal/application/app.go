package application

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"fmt"
	"sort"
	"strconv"
	"strings"

	"github.com/google/uuid"
	"github.com/PithomLabs/oracle/internal/epistemic"
	"github.com/PithomLabs/oracle/internal/work"
	"github.com/PithomLabs/oracle/verifier"
	domainpack "github.com/PithomLabs/oracle/domain-pack"
	packetv1 "github.com/PithomLabs/oracle/packet/v1"
)

// TxExecutor is the minimal transaction interface satisfied by *sql.Tx
// and test doubles. It allows Persist to be tested with fault injection
// without requiring the wrapper to implement *sql.Tx's full API.
type TxExecutor interface {
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
	QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row
	Commit() error
	Rollback() error
}

// LocalOperatorPrincipalID is the canonical UUID for the local human operator.
// Used by both the seed and the Trust UI discharge/promote paths.
// Must match the principal row created by seed.SeedIfEmpty.
const LocalOperatorPrincipalID = "00000000-0000-0000-0000-000000000001"

// App is the application layer — the single entry point for MCP and HTTP.
// It orchestrates epistemic (Solvent) and work (Conductor) operations
// without being an authority source itself.
type App struct {
	db            *sql.DB
	kern          *epistemic.Store
	audit         *epistemic.AuditService
	workStore     *work.Store
	artifactReg   *verifier.ArtifactRegistry
	verifierSpecs []VerifierSpec
	packRegistry  *domainpack.PackRegistry
}

// VerifierSpec mirrors the pack's verifier authorization.
type VerifierSpec struct {
	VerifierID string `json:"verifier_id"`
	MinVersion string `json:"min_version"`
}

// New creates a new App.
func New(db *sql.DB) *App {
	return &App{
		db:        db,
		kern:      epistemic.NewStore(db),
		audit:     epistemic.NewAudit(db),
		workStore: work.NewStore(db),
	}
}

// NewWithRegistry creates a new App with an artifact registry for verifier enforcement.
func NewWithRegistry(db *sql.DB, registry *verifier.ArtifactRegistry, specs []VerifierSpec) *App {
	return &App{
		db:            db,
		kern:          epistemic.NewStore(db),
		audit:         epistemic.NewAudit(db),
		workStore:     work.NewStore(db),
		artifactReg:   registry,
		verifierSpecs: specs,
	}
}

// SetPackRegistry configures the pack registry for retirement-rule enforcement.
func (a *App) SetPackRegistry(reg *domainpack.PackRegistry) {
	a.packRegistry = reg
}

// DB returns the underlying database handle for use by Persist callers.
func (a *App) DB() *sql.DB { return a.db }

// entityID computes a deterministic UUID from scenario + type + content.
// Same input always produces the same ID, making ON CONFLICT (id) DO NOTHING
// genuinely idempotent across retries.
func EntityID(scenarioID, entityType, content string) string {
	h := sha256.Sum256([]byte(scenarioID + ":" + entityType + ":" + content))
	return hex.EncodeToString(h[:16])
}

// contentHash computes the canonical hash for idempotency tracking.
// Hashes the complete canonical packet: headers, beliefs, evidence, edges, tasks.
// Entity lists are sorted deterministically before hashing.
func contentHash(pkt *packetv1.Packet) string {
	h := sha256.New()

	// Headers
	h.Write([]byte(pkt.ScenarioID))
	h.Write([]byte{0})
	h.Write([]byte(pkt.PacketID))
	h.Write([]byte{0})
	h.Write([]byte(pkt.Role))
	h.Write([]byte{0})
	h.Write([]byte(pkt.PackRef))
	h.Write([]byte{1}) // entity separator

	// Beliefs sorted by LocalID
	sortedBeliefs := make([]packetv1.Belief, len(pkt.Beliefs))
	copy(sortedBeliefs, pkt.Beliefs)
	sort.Slice(sortedBeliefs, func(i, j int) bool {
		return sortedBeliefs[i].LocalID < sortedBeliefs[j].LocalID
	})
	for _, b := range sortedBeliefs {
		h.Write([]byte(b.LocalID))
		h.Write([]byte{0})
		h.Write([]byte(b.Claim))
		h.Write([]byte{0})
		h.Write([]byte(b.ClaimType))
		h.Write([]byte{0})
		debt := b.Debt
		if debt == nil {
			debt = []string{}
		}
		sortedDebt := make([]string, len(debt))
		copy(sortedDebt, debt)
		sort.Strings(sortedDebt)
		for _, d := range sortedDebt {
			h.Write([]byte(d))
			h.Write([]byte{0})
		}
		h.Write([]byte(b.InputSpec))
		h.Write([]byte{1})
	}

	// Evidence sorted by LocalID
	sortedEvidence := make([]packetv1.Evidence, len(pkt.Evidence))
	copy(sortedEvidence, pkt.Evidence)
	sort.Slice(sortedEvidence, func(i, j int) bool {
		return sortedEvidence[i].LocalID < sortedEvidence[j].LocalID
	})
	for _, e := range sortedEvidence {
		h.Write([]byte(e.LocalID))
		h.Write([]byte{0})
		h.Write([]byte(e.BeliefRef))
		h.Write([]byte{0})
		h.Write([]byte(e.ProvenanceClass))
		h.Write([]byte{0})
		h.Write([]byte(e.ContentSHA256))
		h.Write([]byte{0})
		h.Write([]byte(e.SourceURL))
		h.Write([]byte{0})
		h.Write([]byte(e.ArtifactRef))
		h.Write([]byte{1})
	}

	// Edges sorted by LocalID
	sortedEdges := make([]packetv1.Edge, len(pkt.Edges))
	copy(sortedEdges, pkt.Edges)
	sort.Slice(sortedEdges, func(i, j int) bool {
		return sortedEdges[i].LocalID < sortedEdges[j].LocalID
	})
	for _, edge := range sortedEdges {
		h.Write([]byte(edge.LocalID))
		h.Write([]byte{0})
		h.Write([]byte(edge.FromRef))
		h.Write([]byte{0})
		h.Write([]byte(edge.ToRef))
		h.Write([]byte{0})
		h.Write([]byte(edge.Kind))
		h.Write([]byte{1})
	}

	// Tasks sorted by Title+Description+GovernanceRef
	sortedTasks := make([]packetv1.Task, len(pkt.Tasks))
	copy(sortedTasks, pkt.Tasks)
	sort.Slice(sortedTasks, func(i, j int) bool {
		a, b := sortedTasks[i], sortedTasks[j]
		if a.Title != b.Title {
			return a.Title < b.Title
		}
		if a.Description != b.Description {
			return a.Description < b.Description
		}
		return a.GovernanceRef < b.GovernanceRef
	})
	for _, t := range sortedTasks {
		h.Write([]byte(t.Title))
		h.Write([]byte{0})
		h.Write([]byte(t.Description))
		h.Write([]byte{0})
		h.Write([]byte(t.GovernanceRef))
		h.Write([]byte{1})
	}

	return hex.EncodeToString(h.Sum(nil))
}

// ---- Compile ----

// Compile validates a packet structurally.
func (a *App) Compile(ctx context.Context, pkt *packetv1.Packet) error {
	if pkt.SchemaVersion != packetv1.SchemaVersion {
		return fmt.Errorf("unsupported schema version: %s", pkt.SchemaVersion)
	}
	if pkt.PacketID == "" {
		return fmt.Errorf("packet_id is required")
	}
	if pkt.ScenarioID == "" {
		return fmt.Errorf("scenario_id is required")
	}
	if pkt.PackRef == "" {
		return fmt.Errorf("pack_ref is required")
	}
	if len(pkt.Beliefs) == 0 && len(pkt.Evidence) == 0 && len(pkt.Tasks) == 0 {
		return fmt.Errorf("packet is empty")
	}
	return nil
}

// ---- Validate ----

// Validate checks domain pack rules (debt membership, evidence classes, verifier specs).
func (a *App) Validate(ctx context.Context, pkt *packetv1.Packet) error {
	// Resolve pack from packet's pack_ref for vocabulary validation.
	var debtVocab map[string]bool
	var evidenceClasses map[string]bool
	if a.packRegistry != nil {
		packID, packVersion := parsePackRef(pkt.PackRef)
		if pack, err := a.packRegistry.Get(packID, packVersion); err == nil {
			debtVocab = make(map[string]bool)
			for _, d := range pack.GetDebtVocabulary() {
				debtVocab[d] = true
			}
			evidenceClasses = make(map[string]bool)
			for _, ec := range pack.GetEvidenceClasses() {
				evidenceClasses[ec] = true
			}
		}
	}
	// Fallback to hardcoded vocabulary if pack registry unavailable (test compat).
	if debtVocab == nil {
		debtVocab = map[string]bool{
			"needMap": true, "needInvariant": true, "needToyCheck": true,
			"needNullModel": true, "needObstruction": true, "needFaithfulnessReview": true,
			"needInitialCondition": true, "needRegularity": true,
		}
	}
	if evidenceClasses == nil {
		evidenceClasses = map[string]bool{
			"reproducible_artifact": true, "operator_asserted": true,
		}
	}
	// Validate debt membership against pack vocabulary.
	for i, b := range pkt.Beliefs {
		for _, debt := range b.Debt {
			if !debtVocab[debt] {
				return fmt.Errorf("belief[%d]: unknown debt item %q", i, debt)
			}
		}
	}
	// Validate evidence classes.
	for i, e := range pkt.Evidence {
		if !evidenceClasses[e.ProvenanceClass] {
			return fmt.Errorf("evidence[%d]: unsupported provenance_class %q", i, e.ProvenanceClass)
		}
	}
	// Runtime VerifierSpec enforcement: for each artifact reference, verify
	// the verifier is authorized by the pack.
	if a.artifactReg != nil && len(a.verifierSpecs) > 0 {
		for i, e := range pkt.Evidence {
			if e.ArtifactRef == "" {
				continue
			}
			artifact, err := a.artifactReg.Resolve(ctx, e.ArtifactRef)
			if err != nil {
				return fmt.Errorf("evidence[%d]: artifact not found: %w", i, err)
			}
			spec := a.findVerifierSpec(artifact.VerifierID)
			if spec == nil {
				return fmt.Errorf("evidence[%d]: verifier %s not authorized by pack", i, artifact.VerifierID)
			}
			if !versionGTE(artifact.VerifierVersion, spec.MinVersion) {
				return fmt.Errorf("evidence[%d]: verifier %s version %s below minimum %s",
					i, artifact.VerifierID, artifact.VerifierVersion, spec.MinVersion)
			}
		}
	}
	return nil
}

// ValidatePacket runs the canonical packet structural validator (fail-closed).
// If the pack registry is unavailable, validation is rejected — never silently skipped.
func (a *App) ValidatePacket(ctx context.Context, pkt *packetv1.Packet) error {
	if a.packRegistry == nil {
		return fmt.Errorf("packet validation unavailable: pack registry is not configured")
	}
	return packetv1.Validate(pkt, a.packRegistry)
}

func (a *App) findVerifierSpec(verifierID string) *VerifierSpec {
	for i := range a.verifierSpecs {
		if a.verifierSpecs[i].VerifierID == verifierID {
			return &a.verifierSpecs[i]
		}
	}
	return nil
}

// versionGTE returns true if version >= minVersion using strict integer semver.
// Rejects prerelease versions. Requires exactly 3 components (major.minor.patch).
func versionGTE(version, minVersion string) bool {
	v, err1 := parseSemver(version)
	m, err2 := parseSemver(minVersion)
	if err1 != nil || err2 != nil {
		return false
	}
	if v[0] != m[0] {
		return v[0] > m[0]
	}
	if v[1] != m[1] {
		return v[1] > m[1]
	}
	return v[2] >= m[2]
}

func parseSemver(v string) ([3]int, error) {
	var parts [3]int
	v = strings.TrimPrefix(v, "v")
	if strings.Contains(v, "-") {
		return parts, fmt.Errorf("prerelease versions not supported: %s", v)
	}
	segs := strings.SplitN(v, ".", 3)
	if len(segs) != 3 {
		return parts, fmt.Errorf("invalid semver: %s", v)
	}
	for i, s := range segs {
		n, err := strconv.Atoi(s)
		if err != nil {
			return parts, fmt.Errorf("invalid semver component: %s", err)
		}
		parts[i] = n
	}
	return parts, nil
}

// ---- Persist ----

// Persist writes packet entities to the provided transaction with entity-level
// idempotency. The caller must begin and commit/rollback the transaction.
// Write order: beliefs → evidence → edges → tasks → proposed_retirement → idempotency.
func (a *App) Persist(ctx context.Context, tx TxExecutor, pkt *packetv1.Packet) (*Result, error) {
	hash := contentHash(pkt)
	result := &Result{
		PacketID:  pkt.PacketID,
		BeliefIDs: make(map[string]string),
	}

	// Defense-in-depth: global local_id uniqueness across all entity types.
	// The local:<id> reference grammar means all local_ids share one namespace.
	{
		seen := make(map[string]bool)
		type entry struct{ kind, id string }
		var collisions []entry
		for _, b := range pkt.Beliefs {
			if seen[b.LocalID] {
				collisions = append(collisions, entry{"belief", b.LocalID})
			}
			seen[b.LocalID] = true
		}
		for _, e := range pkt.Evidence {
			if seen[e.LocalID] {
				collisions = append(collisions, entry{"evidence", e.LocalID})
			}
			seen[e.LocalID] = true
		}
		for _, edge := range pkt.Edges {
			if seen[edge.LocalID] {
				collisions = append(collisions, entry{"edge", edge.LocalID})
			}
			seen[edge.LocalID] = true
		}
		for _, t := range pkt.Tasks {
			if seen[t.LocalID] {
				collisions = append(collisions, entry{"task", t.LocalID})
			}
			seen[t.LocalID] = true
		}
		if len(collisions) > 0 {
			return nil, fmt.Errorf("duplicate local_id across packet entities: %v", collisions)
		}
	}

	// Defense-in-depth: validate claim_type enum values before INSERT.
	for i, b := range pkt.Beliefs {
		switch b.ClaimType {
		case "derived", "accommodated", "postulated":
			// valid
		default:
			return nil, fmt.Errorf("belief[%d]: invalid claim_type %q", i, b.ClaimType)
		}
	}

	// Defense-in-depth: validate edge kind values before INSERT.
	for i, edge := range pkt.Edges {
		if edge.Kind != packetv1.EdgeDerives && edge.Kind != packetv1.EdgeContradicts {
			return nil, fmt.Errorf("edge[%d]: invalid kind %q", i, edge.Kind)
		}
	}

	// 0. Persist packet submission provenance FIRST (FK dependency).
	// origin_packet_id on belief/evidence/task/edge_provenance references this row.
	// TaskRef handling: empty → NULL, valid UUID → store, invalid UUID → reject.
	var taskRef sql.NullString
	if pkt.TaskRef != "" {
		if _, err := uuid.Parse(pkt.TaskRef); err != nil {
			return nil, fmt.Errorf("invalid task_ref: %w", err)
		}
		taskRef = sql.NullString{String: pkt.TaskRef, Valid: true}
	}
	if _, err := tx.ExecContext(ctx,
		`INSERT INTO packet_submission
			 (packet_id, scenario_id, task_id, agent_id, role, harness, model, content_sha256)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		 ON CONFLICT (packet_id) DO NOTHING`,
		pkt.PacketID, pkt.ScenarioID, taskRef,
		pkt.Agent.ID, pkt.Agent.Role, pkt.Agent.Harness, pkt.Agent.Model, hash); err != nil {
		return nil, fmt.Errorf("insert packet submission: %w", err)
	}

	// 1. Persist beliefs with deterministic IDs.
	for _, b := range pkt.Beliefs {
		debt := b.Debt
		if debt == nil {
			debt = []string{}
		}
		beliefID := EntityID(pkt.ScenarioID, "belief", b.Claim)
		_, err := tx.ExecContext(ctx,
			`INSERT INTO belief (id, scenario_id, claim, claim_type, debt, origin_packet_id)
			 VALUES ($1::UUID, $2::UUID, $3, $4, $5, $6)
			 ON CONFLICT (id) DO NOTHING`,
			beliefID, pkt.ScenarioID, b.Claim, b.ClaimType, debt, pkt.PacketID)
		if err != nil {
			return nil, fmt.Errorf("insert belief: %w", err)
		}
		result.BeliefIDs[b.LocalID] = beliefID
	}

	// 2. Persist evidence with deterministic IDs.
	for i, e := range pkt.Evidence {
		// operator_asserted requires human attestation — agents cannot create it
		if e.ProvenanceClass == "operator_asserted" {
			return nil, fmt.Errorf("evidence[%d]: operator_asserted evidence requires human attestation; agent-submitted packets cannot create it", i)
		}
		// reproducible_artifact requires a trusted artifact reference
		if e.ProvenanceClass == "reproducible_artifact" && e.ArtifactRef == "" {
			return nil, fmt.Errorf("evidence[%d]: reproducible_artifact requires artifact_ref", i)
		}
		beliefID := resolveRef(e.BeliefRef, result.BeliefIDs)
		if beliefID == "" {
			return nil, fmt.Errorf("evidence references unresolved belief: %s", e.BeliefRef)
		}
		evidenceID := EntityID(pkt.ScenarioID, "evidence", e.ContentSHA256)
		_, err := tx.ExecContext(ctx,
			`INSERT INTO evidence (id, scenario_id, belief_id, provenance_class, source_url, content_sha256, origin_packet_id)
			 VALUES ($1::UUID, $2::UUID, $3::UUID, $4, $5, $6, $7)
			 ON CONFLICT (id) DO NOTHING`,
			evidenceID, pkt.ScenarioID, beliefID, e.ProvenanceClass, e.SourceURL, e.ContentSHA256, pkt.PacketID)
		if err != nil {
			return nil, fmt.Errorf("insert evidence: %w", err)
		}
		result.EvidenceIDs = append(result.EvidenceIDs, evidenceID)
	}

	// 2.5. Validate input-to-claim binding for verification evidence.
	if a.artifactReg != nil {
		for i, e := range pkt.Evidence {
			if e.ProvenanceClass != "reproducible_artifact" {
				continue
			}
			beliefID := resolveRef(e.BeliefRef, result.BeliefIDs)
			if beliefID == "" {
				continue // already caught above
			}
			// Find the belief's declared InputSpec.
			var inputSpec string
			for _, b := range pkt.Beliefs {
				if EntityID(pkt.ScenarioID, "belief", b.Claim) == beliefID {
					inputSpec = b.InputSpec
					break
				}
			}
			if inputSpec == "" {
				return nil, fmt.Errorf("evidence[%d]: belief %s has no input_spec for verification binding", i, beliefID)
			}
			// Load artifact from registry and check input hash match.
			artifact, err := a.artifactReg.Resolve(ctx, e.ArtifactRef)
			if err != nil {
				return nil, fmt.Errorf("evidence[%d]: artifact not found: %w", i, err)
			}
			if artifact.InputHash != inputSpec {
				return nil, fmt.Errorf(
					"evidence[%d]: input hash mismatch: artifact has %s, belief declares %s",
					i, artifact.InputHash, inputSpec)
			}
		}
	}

	// 3. Persist edges with deterministic IDs.
	for _, edge := range pkt.Edges {
		fromID := resolveRef(edge.FromRef, result.BeliefIDs)
		toID := resolveRef(edge.ToRef, result.BeliefIDs)
		if fromID == "" || toID == "" {
			return nil, fmt.Errorf("edge references unresolved: from=%s to=%s", edge.FromRef, edge.ToRef)
		}
		// Edge validation: no self-edges.
		if fromID == toID {
			return nil, fmt.Errorf("edge[%s]: self-edge not allowed (parent=child=%s)", edge.LocalID, fromID)
		}
		// Edge validation: contradicts must target canonical existing belief, not local.
		if edge.Kind == packetv1.EdgeContradicts && strings.HasPrefix(edge.ToRef, packetv1.RefPrefixLocal) {
			return nil, fmt.Errorf("edge[%s]: contradicts target must be a canonical existing belief, not local reference %s", edge.LocalID, edge.ToRef)
		}
		// Edge validation: all endpoints must belong to same scenario.
		if err := a.validateEdgeScenario(ctx, tx, fromID, toID, pkt.ScenarioID); err != nil {
			return nil, fmt.Errorf("edge[%s]: %w", edge.LocalID, err)
		}
		result2, err := tx.ExecContext(ctx,
			`INSERT INTO belief_edge (parent_id, child_id, kind)
			 VALUES ($1::UUID, $2::UUID, $3)
			 ON CONFLICT (parent_id, child_id) DO NOTHING`,
			fromID, toID, edge.Kind)
		if err != nil {
			return nil, fmt.Errorf("insert edge: %w", err)
		}
		rowsAffected, err := result2.RowsAffected()
		if err != nil {
			return nil, fmt.Errorf("edge rows affected: %w", err)
		}
		if rowsAffected == 0 {
			var existingKind string
			err := tx.QueryRowContext(ctx,
				`SELECT kind FROM belief_edge WHERE parent_id = $1::UUID AND child_id = $2::UUID`,
				fromID, toID).Scan(&existingKind)
			if err != nil {
				return nil, fmt.Errorf("edge conflict re-read failed: %w", err)
			}
			if existingKind != edge.Kind {
				return nil, fmt.Errorf("edge conflict: %s->%s already has kind=%s, cannot add kind=%s",
					fromID, toID, existingKind, edge.Kind)
			}
		}
		// Record edge provenance.
		_, err = tx.ExecContext(ctx,
			`INSERT INTO edge_provenance (parent_id, child_id, origin_packet_id)
			 VALUES ($1::UUID, $2::UUID, $3)
			 ON CONFLICT (parent_id, child_id) DO NOTHING`,
			fromID, toID, pkt.PacketID)
		if err != nil {
			return nil, fmt.Errorf("insert edge provenance: %w", err)
		}
		result.EdgeCount++
	}

	// 4. Persist tasks.
	// Task integrity preflight: verify conductor_project exists for this scenario.
	if len(pkt.Tasks) > 0 {
		var projectExists bool
		err := tx.QueryRowContext(ctx,
			`SELECT EXISTS(SELECT 1 FROM conductor_project WHERE id = $1::UUID)`,
			pkt.ScenarioID,
		).Scan(&projectExists)
		if err != nil {
			return nil, fmt.Errorf("verify project for scenario: %w", err)
		}
		if !projectExists {
			return nil, fmt.Errorf("cannot create tasks: project %s does not exist", pkt.ScenarioID)
		}
	}
	for _, t := range pkt.Tasks {
		taskID := EntityID(pkt.ScenarioID, "task", t.Title)
		governanceRef := t.GovernanceRef
		_, err := tx.ExecContext(ctx,
			`INSERT INTO conductor_task (id, project_id, title, description, status, priority, governance_ref, origin_packet_id, created_at, updated_at)
			 VALUES ($1, $2, $3, $4, 'proposed', 'medium', $5, $6, now(), now())
			 ON CONFLICT (id) DO NOTHING`,
			taskID, pkt.ScenarioID, t.Title, t.Description, governanceRef, pkt.PacketID)
		if err != nil {
			return nil, fmt.Errorf("insert task: %w", err)
		}
		result.TaskIDs = append(result.TaskIDs, taskID)
	}

	// 5. Write idempotency row LAST (row existence = completed).
	_, err := tx.ExecContext(ctx,
		`INSERT INTO submission_idempotency (content_hash, scenario_id, packet_id)
		 VALUES ($1, $2, $3)
		 ON CONFLICT (content_hash, scenario_id) DO NOTHING`,
		hash, pkt.ScenarioID, pkt.PacketID)
	if err != nil {
		return nil, fmt.Errorf("insert idempotency: %w", err)
	}

	return result, nil
}

// ---- Context Assembly ----

// Dashboard is the Trust UI dashboard view — all tasks, all beliefs, and recent submissions.
type Dashboard struct {
	Tasks       []*work.Task                    `json:"tasks"`
	Beliefs     []epistemic.BeliefView          `json:"beliefs"`
	Submissions []epistemic.PacketSubmissionView `json:"submissions"`
}

// GetDashboard returns all tasks, beliefs, and recent submissions for the Trust UI dashboard.
func (a *App) GetDashboard(ctx context.Context, scenarioID string) (*Dashboard, error) {
	tasks, err := a.workStore.ListAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("list tasks: %w", err)
	}
	beliefs, err := epistemic.GetAllBeliefs(ctx, a.db)
	if err != nil {
		return nil, fmt.Errorf("get beliefs: %w", err)
	}
	var submissions []epistemic.PacketSubmissionView
	if scenarioID != "" {
		subs, err := epistemic.GetSubmissionsForScenario(ctx, a.db, scenarioID)
		if err == nil {
			submissions = subs
		}
	}
	if submissions == nil {
		submissions = []epistemic.PacketSubmissionView{}
	}
	return &Dashboard{
		Tasks:       tasks,
		Beliefs:     beliefs,
		Submissions: submissions,
	}, nil
}

// GetContext assembles the RCP context for a task.
// Never returns error for availability issues — always returns Context
// with per-section availability metadata.
func (a *App) GetContext(ctx context.Context, taskID string) (*Context, error) {
	result := &Context{
		Availability: epistemic.RCPAvailability{
			Task:         epistemic.SectionAvailability{Available: true},
			Dependencies: epistemic.SectionAvailability{Available: true},
			Snapshot:     epistemic.SectionAvailability{Available: true},
		},
		Snapshot: &epistemic.Snapshot{},
	}

	task, err := a.workStore.GetByID(ctx, taskID)
	if err != nil {
		result.Availability.Task.Available = false
		result.Availability.Task.Reason = "work_unavailable"
	} else {
		result.Task = task
	}

	deps, err := a.listDependencies(ctx, taskID)
	if err != nil {
		result.Availability.Dependencies.Available = false
		result.Availability.Dependencies.Reason = "dependencies_unavailable"
	} else {
		result.Dependencies = deps
	}

	if result.Task != nil {
		snap, err := epistemic.GetSnapshot(ctx, a.db, result.Task.ProjectID, epistemic.SnapshotOpts{
			IncludeEvidence: true,
		})
		if err != nil {
			result.Availability.Snapshot.Available = false
			result.Availability.Snapshot.Reason = "solvent_unavailable"
		} else {
			result.Snapshot = snap
		}
	} else {
		result.Availability.Snapshot.Available = false
		result.Availability.Snapshot.Reason = "skipped: task unavailable"
	}

	return result, nil
}

// ---- Human Decisions ----

// SubmitDecision handles human decisions (discharge, promote, retract).
// Takes AuthenticatedDecisionCommand with server-derived principal.
func (a *App) SubmitDecision(ctx context.Context, cmd *AuthenticatedDecisionCommand) error {
	switch cmd.Type {
	case "discharge":
		if a.packRegistry == nil {
			return fmt.Errorf("pack registry unavailable: cannot validate retirement")
		}
		pack, err := a.resolvePackFromScenario(ctx, cmd.ScenarioID)
		if err != nil {
			return fmt.Errorf("pack resolution failed: %w", err)
		}
		rules := pack.GetRetirementRules()
		if rules == nil {
			return fmt.Errorf("pack does not support retirement rules")
		}
		rule, exists := rules[cmd.ObligationKey]
		if !exists {
			return fmt.Errorf("no retirement rule for debt item: %s", cmd.ObligationKey)
		}
		if cmd.EvidenceClass == "" {
			return fmt.Errorf("evidence_class is required for debt discharge of: %s", cmd.ObligationKey)
		}
		if cmd.EvidenceClass != rule.EvidenceClass {
			return fmt.Errorf("evidence class mismatch: debt %q requires %q, got %q",
				cmd.ObligationKey, rule.EvidenceClass, cmd.EvidenceClass)
		}
		evidenceIDs, err := a.verifyPersistedEvidence(ctx, cmd.BeliefID, cmd.ScenarioID, cmd.EvidenceClass)
		if err != nil {
			return fmt.Errorf("evidence verification failed: %w", err)
		}
		if len(evidenceIDs) == 0 {
			return fmt.Errorf("no qualifying evidence of class %q for belief %s in scenario %s",
				cmd.EvidenceClass, cmd.BeliefID, cmd.ScenarioID)
		}
		instrumentRef := buildInstrumentRef(evidenceIDs)
		return a.kern.Discharge(ctx, cmd.ScenarioID, cmd.BeliefID, cmd.ObligationKey, instrumentRef, cmd.PrincipalID)
	case "promote":
		return a.kern.Promote(ctx, cmd.ScenarioID, cmd.BeliefID)
	case "retract":
		retracted, err := a.kern.RetractCascade(ctx, cmd.ScenarioID, cmd.BeliefID)
		if err != nil {
			return err
		}
		_ = retracted
		return a.cancelLinkedTasks(ctx, cmd.BeliefID)
	default:
		return fmt.Errorf("unknown decision type: %s", cmd.Type)
	}
}

// ---- Helpers ----

func (a *App) listDependencies(ctx context.Context, taskID string) ([]*work.Dependency, error) {
	rows, err := a.db.QueryContext(ctx,
		`SELECT id, task_id, blocked_by_id, created_at FROM conductor_dependency WHERE task_id = $1`, taskID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var deps []*work.Dependency
	for rows.Next() {
		d := &work.Dependency{}
		if err := rows.Scan(&d.ID, &d.TaskID, &d.BlockedByID, &d.CreatedAt); err != nil {
			return nil, err
		}
		deps = append(deps, d)
	}
	return deps, rows.Err()
}

func (a *App) cancelLinkedTasks(ctx context.Context, beliefID string) error {
	tasks, err := a.workStore.GetTasksByGovernanceRef(ctx, beliefID)
	if err != nil {
		return err
	}
	for _, t := range tasks {
		if t.Status != work.TaskStatusCancelled && t.Status != work.TaskStatusAccepted {
			_ = a.workStore.Transition(ctx, t.ID, t.Status, work.TaskStatusCancelled,
				work.ActorTypeSystem, "retraction-cascade", "task.cancelled",
				fmt.Sprintf("belief %s retracted", beliefID))
		}
	}
	return nil
}

// resolvePackFromScenario derives the pack for a scenario.
// For the POC, BM-IST is the only pack, so it is server-selected by exclusion.
func (a *App) resolvePackFromScenario(ctx context.Context, scenarioID string) (domainpack.Pack, error) {
	if a.packRegistry == nil {
		return nil, fmt.Errorf("pack registry unavailable")
	}
	packID, packVersion := scenarioPackMapping(scenarioID)
	return a.packRegistry.Get(packID, packVersion)
}

func scenarioPackMapping(scenarioID string) (string, string) {
	return "bmist", "1.1.0"
}

func parsePackRef(ref string) (string, string) {
	parts := strings.SplitN(ref, "@", 2)
	if len(parts) == 2 {
		return parts[0], parts[1]
	}
	return ref, ""
}

func (a *App) verifyPersistedEvidence(ctx context.Context, beliefID, scenarioID, evidenceClass string) ([]string, error) {
	rows, err := a.db.QueryContext(ctx,
		`SELECT id FROM evidence
		 WHERE belief_id = $1::UUID AND scenario_id = $2::UUID AND provenance_class = $3`,
		beliefID, scenarioID, evidenceClass)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var ids []string
	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}
	return ids, rows.Err()
}

func buildInstrumentRef(evidenceIDs []string) string {
	if len(evidenceIDs) == 0 {
		return ""
	}
	ref := "evidence:"
	for i, id := range evidenceIDs {
		if i > 0 {
			ref += ","
		}
		ref += id
	}
	return ref
}

// validateEdgeScenario checks that both edge endpoints belong to the same scenario.
func (a *App) validateEdgeScenario(ctx context.Context, tx TxExecutor, fromID, toID, scenarioID string) error {
	var fromScenario, toScenario string
	err := tx.QueryRowContext(ctx,
		`SELECT scenario_id FROM belief WHERE id = $1::UUID`, fromID).Scan(&fromScenario)
	if err != nil {
		return fmt.Errorf("from-belief %s not found: %w", fromID, err)
	}
	err = tx.QueryRowContext(ctx,
		`SELECT scenario_id FROM belief WHERE id = $1::UUID`, toID).Scan(&toScenario)
	if err != nil {
		return fmt.Errorf("to-belief %s not found: %w", toID, err)
	}
	if !strings.EqualFold(fromScenario, scenarioID) || !strings.EqualFold(toScenario, scenarioID) {
		return fmt.Errorf("edge endpoints must belong to packet scenario %s, got from=%s to=%s", scenarioID, fromScenario, toScenario)
	}
	return nil
}

func resolveRef(ref string, idMap map[string]string) string {
	if strings.HasPrefix(ref, packetv1.RefPrefixLocal) {
		localID := strings.TrimPrefix(ref, packetv1.RefPrefixLocal)
		return idMap[localID]
	}
	if strings.HasPrefix(ref, packetv1.RefPrefixCanonical) {
		return strings.TrimPrefix(ref, packetv1.RefPrefixCanonical)
	}
	return ""
}

// ---- Types ----

// Result is the outcome of persisting a packet.
type Result struct {
	PacketID    string            `json:"packet_id"`
	BeliefIDs   map[string]string `json:"belief_ids"`
	EvidenceIDs []string          `json:"evidence_ids"`
	EdgeCount   int               `json:"edge_count"`
	TaskIDs     []string          `json:"task_ids"`
}

// Context is the RCP context for a task.
type Context struct {
	Task         *work.Task                 `json:"task"`
	Dependencies []*work.Dependency         `json:"dependencies"`
	Snapshot     *epistemic.Snapshot        `json:"snapshot"`
	Availability epistemic.RCPAvailability  `json:"availability"`
}

// ExternalDecisionRequest is the externally-supplied request from the Trust UI.
// It does NOT carry PrincipalID — the principal is derived at the auth boundary.
type ExternalDecisionRequest struct {
	Type          string `json:"type"`
	ScenarioID    string `json:"scenario_id"`
	BeliefID      string `json:"belief_id"`
	ObligationKey string `json:"obligation_key,omitempty"`
	InstrumentRef string `json:"instrument_ref,omitempty"`
	EvidenceClass string `json:"evidence_class,omitempty"`
}

// AuthenticatedDecisionCommand is the internally constructed command
// after authentication. PrincipalID is server-derived, never from request body.
type AuthenticatedDecisionCommand struct {
	Type          string
	ScenarioID    string
	BeliefID      string
	ObligationKey string
	InstrumentRef string
	PrincipalID   string // server-derived from token validation
	EvidenceClass string
}

// ---- Pack-driven vocabulary validation ----
// Debt and evidence-class validation is now derived from the loaded domain pack.
// The hardcoded fallback in Validate() handles test environments where the pack registry is not wired.
