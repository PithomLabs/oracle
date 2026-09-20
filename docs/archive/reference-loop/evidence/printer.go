package evidence

import (
	"fmt"
	"strings"
	"time"
)

// PrintTrace prints the evidence trace in the required format
func PrintTrace(pkg *EvidencePackage) {
	fmt.Println("=== REFERENCE LOOP EVIDENCE TRACE ===")
	fmt.Printf("Run ID:      %s\n", pkg.RunID)
	fmt.Printf("Scenario:    %s\n", pkg.ScenarioID)
	fmt.Printf("Duration:    %v\n", pkg.EndTime.Sub(pkg.StartTime))
	fmt.Println()

	// Print in the required sequence order
	sequence := []string{
		"TASK_CREATED",
		"TASK_CLAIMED",
		"PROPOSAL_CREATED",
		"AUTHORIZED",
		"EXECUTION_ATTEMPTED",
		"EFFECT_CONFIRMED",
		"RESULT_OBSERVED",
		"CONDUCTOR_UPDATED",
		"AGENT_CONTINUED",
	}

	for _, event := range sequence {
		records := pkg.FindByEvent(event)
		if len(records) > 0 {
			r := records[0] // Take first matching record
			fmt.Printf("%s\n", event)
			fmt.Printf("    Owner:           %s\n", r.Owner)
			fmt.Printf("    Source:          %s\n", r.Source)
			fmt.Printf("    Operation ID:    %s\n", r.OperationID)
			if r.DeclarationVer != "" {
				fmt.Printf("    Declaration Ver: %s\n", r.DeclarationVer)
			}
			fmt.Printf("    Timestamp:       %s\n", r.Timestamp.Format(time.RFC3339))
			fmt.Printf("    Run ID:          %s\n", r.RunID)
			if r.TaskID != "" {
				fmt.Printf("    Task ID:         %s\n", r.TaskID)
			}
			if r.IntentID != "" {
				fmt.Printf("    Intent ID:       %s\n", r.IntentID)
			}
			fmt.Println()
		} else {
			fmt.Printf("%s\n", event)
			fmt.Printf("    [NO EVIDENCE FOUND]\n\n")
		}
	}

	// Print all records for completeness
	fmt.Println("=== ALL EVIDENCE RECORDS ===")
	for _, r := range pkg.Records {
		fmt.Printf("[%s] %s | owner=%s | source=%s | task=%s | intent=%s\n",
			r.Timestamp.Format("15:04:05"), r.Event, r.Owner, r.Source, r.TaskID, r.IntentID)
	}
}

// PrintSummary prints a summary of evidence completeness
func PrintSummary(pkg *EvidencePackage) {
	fmt.Println("=== EVIDENCE COMPLETENESS ===")
	sequence := []string{
		"TASK_CREATED",
		"TASK_CLAIMED",
		"PROPOSAL_CREATED",
		"AUTHORIZED",
		"EXECUTION_ATTEMPTED",
		"EFFECT_CONFIRMED",
		"RESULT_OBSERVED",
		"CONDUCTOR_UPDATED",
		"AGENT_CONTINUED",
	}

	complete := true
	for _, event := range sequence {
		records := pkg.FindByEvent(event)
		status := "✓"
		if len(records) == 0 {
			status = "✗"
			complete = false
		}
		owner := ""
		if len(records) > 0 {
			owner = records[0].Owner
		}
		fmt.Printf("  %s %s (owner: %s)\n", status, event, owner)
	}

	fmt.Println()
	if complete {
		fmt.Println("RESULT: ALL EVIDENCE PRESENT - Happy path verified")
	} else {
		fmt.Println("RESULT: MISSING EVIDENCE - Happy path incomplete")
	}
}

// FindByEvent finds records matching an event
func (p *EvidencePackage) FindByEvent(event string) []EvidenceRecord {
	var result []EvidenceRecord
	for _, r := range p.Records {
		if strings.EqualFold(r.Event, event) || strings.HasPrefix(r.Event, event) {
			result = append(result, r)
		}
	}
	return result
}

// VerifyOperationBinding verifies the exact operation identity survived the full path
func VerifyOperationBinding(pkg *EvidencePackage, expectedOpID string) (bool, string) {
	events := []string{"PROPOSAL_CREATED", "AUTHORIZED", "EXECUTION_ATTEMPTED", "EFFECT_CONFIRMED"}
	var ops []string

	for _, event := range events {
		records := pkg.FindByEvent(event)
		if len(records) > 0 && records[0].OperationID != "" {
			ops = append(ops, records[0].OperationID)
		}
	}

	if len(ops) == 0 {
		return false, "no operation IDs found in evidence"
	}

	// Check all operation IDs match
	first := ops[0]
	for _, op := range ops[1:] {
		if op != first {
			return false, fmt.Sprintf("operation ID mismatch: %s != %s", first, op)
		}
	}

	if expectedOpID != "" && first != expectedOpID {
		return false, fmt.Sprintf("expected operation ID %s, got %s", expectedOpID, first)
	}

	return true, fmt.Sprintf("operation ID %s survived complete path", first)
}
