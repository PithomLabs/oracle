package coordinator

import (
	"encoding/json"
	"fmt"
	"sort"
)

// RCPContext is the canonical RCP/v1 response.
type RCPContext struct {
	Protocol     string            `json:"protocol"`
	Task         AvailableSection  `json:"task"`
	Dependencies AvailableSection  `json:"dependencies"`
	Epistemic    EpistemicSection  `json:"epistemic"`
	Activity     AvailableSection  `json:"activity"`
}

// AvailableSection wraps data with availability metadata.
type AvailableSection struct {
	Available bool        `json:"available"`
	Reason    string      `json:"reason,omitempty"`
	Data      interface{} `json:"data"`
}

// EpistemicSection contains epistemic state with availability.
type EpistemicSection struct {
	Available bool                    `json:"available"`
	Reason    string                  `json:"reason,omitempty"`
	Beliefs   []map[string]interface{} `json:"beliefs"`
	Evidence  []map[string]interface{} `json:"evidence"`
	Edges     []map[string]interface{} `json:"edges"`
	Debt      []map[string]interface{} `json:"debt"`
	Intents   []map[string]interface{} `json:"intents"`
}

// GetContext assembles the RCP/v1 context for a given task.
func (c *Coordinator) GetContext(taskID string) (*RCPContext, error) {
	// 1. Get task from Conductor
	task, err := c.conductorClient.GetTask(taskID)
	if err != nil {
		return nil, fmt.Errorf("conductor: task not found: %w", err)
	}

	// 2. Extract scenario from governance_ref
	scenarioID := extractScenario(task)

	// 3. Read from Conductor (fail-safe: unavailable ≠ empty)
	deps, depErr := c.conductorClient.ListDependencies(taskID)
	activity, actErr := c.conductorClient.ListActivity(taskID)

	// 4. Read from Solvent (fail-safe: unavailable ≠ empty)
	beliefs, belErr := c.solventClient.ListBeliefs(scenarioID)
	evidence, evErr := c.solventClient.ListEvidence(scenarioID)
	edges, edgeErr := c.solventClient.ListEdges(scenarioID)
	intents, intErr := c.solventClient.ListIntents(scenarioID)
	solventActivity, solActErr := c.solventClient.ListActivities(scenarioID)

	// 5. Derive debt from beliefs
	debt := deriveDebt(beliefs)

	// 6. Assemble response with availability metadata
	ctx := &RCPContext{
		Protocol: "RCP/v1",
		Task: AvailableSection{
			Available: true,
			Data:      task,
		},
		Dependencies: AvailableSection{
			Available: depErr == nil,
			Reason:    depErrReason(depErr),
			Data:      deps,
		},
		Epistemic: EpistemicSection{
			Available: belErr == nil && evErr == nil && edgeErr == nil && intErr == nil,
			Reason:    epistemicReason(belErr, evErr, edgeErr, intErr),
			Beliefs:   beliefs,
			Evidence:  evidence,
			Edges:     edges,
			Debt:      debt,
			Intents:   intents,
		},
		Activity: AvailableSection{
			Available: actErr == nil && solActErr == nil,
			Reason:    activityReason(actErr, solActErr),
			Data:      mergeActivity(activity, solventActivity),
		},
	}

	return ctx, nil
}

// extractScenario extracts the scenario ID from a task's governance_ref.
func extractScenario(task map[string]interface{}) string {
	govRef, _ := task["governance_ref"].(string)
	if govRef == "" {
		return ""
	}
	var ref struct {
		ReferenceID string `json:"reference_id"`
	}
	if err := json.Unmarshal([]byte(govRef), &ref); err != nil {
		return ""
	}
	return ref.ReferenceID
}

// deriveDebt derives the debt projection from beliefs.
func deriveDebt(beliefs []map[string]interface{}) []map[string]interface{} {
	var debt []map[string]interface{}
	for _, b := range beliefs {
		beliefID, _ := b["belief_id"].(string)
		items, _ := b["debt"].([]interface{})
		if len(items) > 0 {
			debt = append(debt, map[string]interface{}{
				"belief_id": beliefID,
				"items":     items,
				"source":    "solvent",
			})
		}
	}
	if debt == nil {
		debt = []map[string]interface{}{}
	}
	return debt
}

// mergeActivity merges Conductor and Solvent activity, sorted by timestamp.
func mergeActivity(conductor, solvent []map[string]interface{}) []map[string]interface{} {
	// Tag source
	for _, a := range conductor {
		a["source"] = "conductor"
	}
	for _, a := range solvent {
		a["source"] = "solvent"
	}
	// Sort within each source group by timestamp for deterministic ordering
	sort.SliceStable(conductor, func(i, j int) bool {
		ti, _ := conductor[i]["created_at"].(string)
		tj, _ := conductor[j]["created_at"].(string)
		return ti < tj
	})
	sort.SliceStable(solvent, func(i, j int) bool {
		ti, _ := solvent[i]["created_at"].(string)
		tj, _ := solvent[j]["created_at"].(string)
		return ti < tj
	})
	// Concatenate source groups: conductor first, then solvent
	merged := make([]map[string]interface{}, 0, len(conductor)+len(solvent))
	merged = append(merged, conductor...)
	merged = append(merged, solvent...)
	return merged
}

func depErrReason(err error) string {
	if err == nil {
		return ""
	}
	return "conductor_unavailable"
}

func epistemicReason(belErr, evErr, edgeErr, intErr error) string {
	if belErr != nil {
		return "solvent_unavailable"
	}
	if evErr != nil {
		return "solvent_unavailable"
	}
	if edgeErr != nil {
		return "solvent_unavailable"
	}
	if intErr != nil {
		return "solvent_unavailable"
	}
	return ""
}

func activityReason(conductorErr, solventErr error) string {
	if conductorErr != nil && solventErr != nil {
		return "both_unavailable"
	}
	if conductorErr != nil {
		return "conductor_unavailable"
	}
	if solventErr != nil {
		return "solvent_unavailable"
	}
	return ""
}
