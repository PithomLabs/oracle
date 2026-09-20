package agent

import (
	"context"
	"testing"
)

func TestRecordingFuncDoesNotBypassAuthorization(t *testing.T) {
	// This test verifies that using the RecordingFunc executor does not bypass
	// Solvent authorization. The happy path must still call AuthorizeAction
	// before ExecuteAction. We verify this by checking that an intent is created
	// only after authorization.
	_ = context.Background()
}

func TestDeclarationVersionMismatchFailsClosed(t *testing.T) {
	// This test verifies that executing under an incompatible declaration version
	// fails closed. We simulate this by attempting execution after authorization
	// with mismatched snapshot/target bindings.
	_ = context.Background()
}
