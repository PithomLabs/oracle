package epistemic

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
)

// BeliefView is a read-only projection of a belief row.
type BeliefView struct {
	ID         string   `json:"id"`
	Claim      string   `json:"claim"`
	ClaimType  string   `json:"claim_type"`
	Status     string   `json:"status"`
	Debt       []string `json:"debt"`
	FinalTruth bool     `json:"final_truth"`
}

// EvidenceView is a read-only projection of an evidence row.
type EvidenceView struct {
	BeliefID        string `json:"belief_id"`
	SourceURL       string `json:"source_url"`
	ProvenanceClass string `json:"provenance_class"`
	ContentSHA256   string `json:"content_sha256"`
}

// IntentView is a read-only projection of an action_intent row.
type IntentView struct {
	BeliefID string `json:"belief_id"`
	Action   string `json:"action"`
	State    string `json:"state"`
}

// Snapshot is a complete read-only view of a scenario's ledger state.
type Snapshot struct {
	Beliefs                []BeliefView   `json:"beliefs"`
	Evidence               []EvidenceView `json:"evidence,omitempty"`
	Intents                []IntentView   `json:"intents"`
	AuditLiveOnNonPromoted int            `json:"audit_live_on_nonpromoted"`
}

// SectionAvailability tracks whether a specific data source is reachable.
type SectionAvailability struct {
	Available bool   `json:"available"`
	Reason    string `json:"reason,omitempty"`
}

// RCPAvailability reports the status of each data source independently.
// An agent seeing task.available=false, snapshot.available=true knows
// that epistemic data is trustworthy but task data is not.
type RCPAvailability struct {
	Task         SectionAvailability `json:"task"`
	Dependencies SectionAvailability `json:"dependencies"`
	Snapshot     SectionAvailability `json:"snapshot"`
}

// SnapshotOpts controls which parts of the snapshot are populated.
type SnapshotOpts struct {
	BeliefID        string
	IncludeEvidence bool
}

// GetSnapshot returns a read-only view of the ledger for a scenario.
func GetSnapshot(ctx context.Context, db *sql.DB, scenarioID string, opts SnapshotOpts) (*Snapshot, error) {
	snap := &Snapshot{}

	if opts.BeliefID != "" {
		row := db.QueryRowContext(ctx,
			`SELECT id, claim, claim_type, status, debt::STRING, final_truth
			 FROM belief WHERE scenario_id=$1::UUID AND id=$2::UUID`,
			scenarioID, opts.BeliefID)
		b, err := scanBelief(row)
		if err != nil {
			return nil, err
		}
		snap.Beliefs = []BeliefView{*b}
	} else {
		rows, err := db.QueryContext(ctx,
			`SELECT id, claim, claim_type, status, debt::STRING, final_truth
			 FROM belief WHERE scenario_id=$1::UUID ORDER BY claim`,
			scenarioID)
		if err != nil {
			return nil, err
		}
		defer rows.Close()
		for rows.Next() {
			b, err := scanBelief(rows)
			if err != nil {
				return nil, err
			}
			snap.Beliefs = append(snap.Beliefs, *b)
		}
		if err := rows.Err(); err != nil {
			return nil, err
		}
	}

	if opts.IncludeEvidence {
		rows, err := db.QueryContext(ctx,
			`SELECT belief_id, source_url, provenance_class, content_sha256
			 FROM evidence WHERE scenario_id=$1::UUID ORDER BY belief_id, ingested_at`,
			scenarioID)
		if err != nil {
			return nil, err
		}
		defer rows.Close()
		for rows.Next() {
			var e EvidenceView
			if err := rows.Scan(&e.BeliefID, &e.SourceURL, &e.ProvenanceClass, &e.ContentSHA256); err != nil {
				return nil, err
			}
			snap.Evidence = append(snap.Evidence, e)
		}
		if err := rows.Err(); err != nil {
			return nil, err
		}
	}

	rows, err := db.QueryContext(ctx,
		`SELECT belief_id, action, state
		 FROM action_intent WHERE scenario_id=$1::UUID ORDER BY belief_id`,
		scenarioID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var i IntentView
		if err := rows.Scan(&i.BeliefID, &i.Action, &i.State); err != nil {
			return nil, err
		}
		snap.Intents = append(snap.Intents, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return snap, nil
}

type scanner interface {
	Scan(dest ...any) error
}

func scanBelief(s scanner) (*BeliefView, error) {
	var b BeliefView
	var debtRaw string
	if err := s.Scan(&b.ID, &b.Claim, &b.ClaimType, &b.Status, &debtRaw, &b.FinalTruth); err != nil {
		return nil, err
	}
	b.Debt = parsePGArray(debtRaw)
	return &b, nil
}

func parsePGArray(s string) []string {
	s = strings.TrimSpace(s)
	if s == "" || s == "{}" {
		return []string{}
	}
	s = strings.TrimPrefix(s, "{")
	s = strings.TrimSuffix(s, "}")
	if s == "" {
		return []string{}
	}
	return strings.Split(s, ",")
}

// GetAllBeliefs returns all beliefs across all scenarios for the dashboard view.
func GetAllBeliefs(ctx context.Context, db *sql.DB) ([]BeliefView, error) {
	rows, err := db.QueryContext(ctx,
		`SELECT id, claim, claim_type, status, debt::STRING, final_truth
		 FROM belief ORDER BY claim`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var beliefs []BeliefView
	for rows.Next() {
		b, err := scanBelief(rows)
		if err != nil {
			return nil, err
		}
		beliefs = append(beliefs, *b)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	if beliefs == nil {
		beliefs = []BeliefView{}
	}
	return beliefs, nil
}

// BeliefExplain is a per-belief human-readable explanation derived from a Snapshot.
type BeliefExplain struct {
	BeliefID                     string       `json:"belief_id"`
	Claim                        string       `json:"claim"`
	ClaimType                    string       `json:"claim_type"`
	Status                       string       `json:"status"`
	RemainingDebt                []string     `json:"remaining_debt"`
	FinalTruth                   bool         `json:"final_truth"`
	IsPromoted                   bool         `json:"is_promoted"`
	IsRetracted                  bool         `json:"is_retracted"`
	CanPromote                   bool         `json:"can_promote"`
	PromotionBlockedReason       string       `json:"promotion_blocked_reason,omitempty"`
	CanAuthorize                 bool         `json:"can_authorize"`
	AuthorizationBlockedReason   string       `json:"authorization_blocked_reason,omitempty"`
	LiveIntents                  []IntentView `json:"live_intents"`
	EvidenceCount                int          `json:"evidence_count"`
	HumanSummary                 string       `json:"human_summary"`
}

// ExplainResult is the top-level response for explain.
type ExplainResult struct {
	Scenario               string          `json:"scenario"`
	ScenarioID             string          `json:"scenario_id"`
	AuditLiveOnNonPromoted int             `json:"audit_live_on_nonpromoted"`
	Beliefs                []BeliefExplain `json:"beliefs"`
	GlobalSummary          string          `json:"global_summary"`
}

// ExplainSnapshot derives a human-readable explanation from a Snapshot.
func ExplainSnapshot(scenario, scenarioID string, snap *Snapshot) *ExplainResult {
	res := &ExplainResult{
		Scenario:               scenario,
		ScenarioID:             scenarioID,
		AuditLiveOnNonPromoted: snap.AuditLiveOnNonPromoted,
		Beliefs:                []BeliefExplain{},
	}

	evidenceByBelief := make(map[string]int)
	for _, e := range snap.Evidence {
		evidenceByBelief[e.BeliefID]++
	}
	intentsByBelief := make(map[string][]IntentView)
	for _, it := range snap.Intents {
		intentsByBelief[it.BeliefID] = append(intentsByBelief[it.BeliefID], it)
	}

	for _, b := range snap.Beliefs {
		be := explainBelief(b, intentsByBelief[b.ID], evidenceByBelief[b.ID])
		res.Beliefs = append(res.Beliefs, be)
	}

	res.GlobalSummary = buildGlobalSummary(res)
	return res
}

func explainBelief(b BeliefView, intents []IntentView, evidenceCount int) BeliefExplain {
	be := BeliefExplain{
		BeliefID:      b.ID,
		Claim:         b.Claim,
		ClaimType:     b.ClaimType,
		Status:        b.Status,
		RemainingDebt: b.Debt,
		FinalTruth:    b.FinalTruth,
		EvidenceCount: evidenceCount,
		IsPromoted:    b.Status == "promoted",
		IsRetracted:   b.Status == "retracted",
	}
	if be.RemainingDebt == nil {
		be.RemainingDebt = []string{}
	}

	var live []IntentView
	for _, it := range intents {
		if it.State == "live" {
			live = append(live, it)
		}
	}
	if live == nil {
		live = []IntentView{}
	}
	be.LiveIntents = live

	switch {
	case b.Status == "retracted":
		be.CanPromote = false
		be.PromotionBlockedReason = "belief is retracted"
	case b.Status == "promoted":
		be.CanPromote = false
		be.PromotionBlockedReason = "already promoted"
	case b.FinalTruth:
		be.CanPromote = false
		be.PromotionBlockedReason = "blocked by final_truth=true"
	case len(b.Debt) > 0:
		be.CanPromote = false
		be.PromotionBlockedReason = fmt.Sprintf("blocked: %d unresolved obligation(s)", len(b.Debt))
	case b.Status == "entered":
		be.CanPromote = true
	default:
		be.CanPromote = false
		be.PromotionBlockedReason = fmt.Sprintf("unexpected status %q", b.Status)
	}

	switch {
	case b.Status == "promoted" && len(b.Debt) == 0 && !b.FinalTruth:
		be.CanAuthorize = true
	case b.Status == "retracted":
		be.CanAuthorize = false
		be.AuthorizationBlockedReason = "belief is retracted"
	case b.Status != "promoted":
		be.CanAuthorize = false
		be.AuthorizationBlockedReason = fmt.Sprintf("belief is not promoted (status=%q)", b.Status)
	case b.FinalTruth:
		be.CanAuthorize = false
		be.AuthorizationBlockedReason = "belief carries final_truth"
	default:
		be.CanAuthorize = false
		be.AuthorizationBlockedReason = "not promotable"
	}

	be.HumanSummary = buildBeliefSummary(be)
	return be
}

func buildBeliefSummary(be BeliefExplain) string {
	var parts []string
	parts = append(parts, fmt.Sprintf("Belief %q", be.Claim))

	switch be.Status {
	case "entered":
		if len(be.RemainingDebt) == 0 && !be.FinalTruth {
			parts = append(parts, "is entered and ready to promote")
		} else {
			parts = append(parts, fmt.Sprintf("is entered with %d unresolved obligation(s)", len(be.RemainingDebt)))
		}
	case "promoted":
		parts = append(parts, "is promoted")
	case "retracted":
		parts = append(parts, "is retracted")
	}

	if be.CanPromote {
		parts = append(parts, "promotion would succeed")
	} else if be.PromotionBlockedReason != "" && be.Status != "promoted" && be.Status != "retracted" {
		parts = append(parts, fmt.Sprintf("promotion blocked: %s", be.PromotionBlockedReason))
	}

	if be.CanAuthorize {
		parts = append(parts, "authorization would succeed")
	} else if be.AuthorizationBlockedReason != "" {
		parts = append(parts, fmt.Sprintf("authorization blocked: %s", be.AuthorizationBlockedReason))
	}

	return strings.Join(parts, ". ") + "."
}

func buildGlobalSummary(res *ExplainResult) string {
	if len(res.Beliefs) == 0 {
		return fmt.Sprintf("Scenario %q has no beliefs.", res.Scenario)
	}
	promoted := 0
	retracted := 0
	for _, be := range res.Beliefs {
		if be.IsPromoted {
			promoted++
		}
		if be.IsRetracted {
			retracted++
		}
	}
	return fmt.Sprintf("Scenario %q: %d belief(s), %d promoted, %d retracted.",
		res.Scenario, len(res.Beliefs), promoted, retracted)
}
