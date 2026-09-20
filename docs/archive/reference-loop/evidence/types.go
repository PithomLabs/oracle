package evidence

import (
	"encoding/json"
	"time"
)

type EvidenceRecord struct {
	Event                string          `json:"event"`               // TASK_CREATED, AUTHORIZED, etc.
	Owner                string          `json:"owner"`               // conductor, solvent, executor, external, agent
	Source               string          `json:"source"`              // conductor.activity, solvent.audit, github.api, etc.
	OperationID          string          `json:"operation_id"`        // deploy:repo:workflow:ref:run_id
	DeclarationVer       string          `json:"declaration_ver"`     // snapshot_id from Solvent
	Timestamp            time.Time       `json:"timestamp"`
	RunID                string          `json:"run_id"`              // immutable correlation ID
	ScenarioID           string          `json:"scenario_id"`         // happy_path, wrong_operation, etc.
	TaskID               string          `json:"task_id,omitempty"`
	IntentID             string          `json:"intent_id,omitempty"`
	NonExternalEffectProof *bool          `json:"non_external_effect_proof,omitempty"`
	Raw                  json.RawMessage `json:"raw"`                 // participant's own record
}

type EvidencePackage struct {
	RunID      string           `json:"run_id"`
	ScenarioID string           `json:"scenario_id"`
	StartTime  time.Time        `json:"start_time"`
	EndTime    time.Time        `json:"end_time"`
	Records    []EvidenceRecord `json:"records"`
}

type ExpectedEvent struct {
	Event       string
	Owner       string
	Source      string
	Description string
}

var HappyPathEvents = []ExpectedEvent{
	{Event: "TASK_CREATED", Owner: "conductor", Source: "conductor.activity", Description: "Task created in Conductor"},
	{Event: "TASK_CLAIMED", Owner: "conductor", Source: "conductor.activity", Description: "Task claimed by agent"},
	{Event: "PROPOSAL_CREATED", Owner: "agent", Source: "agent.solvent_call", Description: "Agent called Solvent authorize_action"},
	{Event: "AUTHORIZED", Owner: "solvent", Source: "solvent.audit", Description: "Solvent authorized and created intent"},
	{Event: "EXECUTION_ATTEMPTED", Owner: "executor", Source: "solvent.audit", Description: "Executor claimed intent and invoked adapter"},
	{Event: "EFFECT_CONFIRMED", Owner: "external", Source: "github.api", Description: "GitHub workflow executed (authoritative SOR)"},
	{Event: "RESULT_OBSERVED", Owner: "agent", Source: "agent.solvent_response", Description: "Agent received ExecuteAction response"},
	{Event: "CONDUCTOR_UPDATED", Owner: "conductor", Source: "conductor.activity", Description: "Conductor recorded activity/status"},
	{Event: "AGENT_CONTINUED", Owner: "agent", Source: "agent.next_task", Description: "Agent discovered next task"},
}
