// Package epistemic re-exports Solvent's kernel, ledger, and audit types
// as a thin wrapper. This keeps Solvent's domain model accessible while
// maintaining a clean import boundary for the application layer.
package epistemic

import (
	"github.com/PithomLabs/solvent/kernel"
	"github.com/PithomLabs/solvent/service/audit"
	"github.com/PithomLabs/solvent/service/ledger"
)

// Re-export kernel types.
type Store = kernel.Store
type Contract = kernel.Contract
type ClaimType = kernel.ClaimType
type AuthorityTuple = kernel.AuthorityTuple
type AuthorizeResult = kernel.AuthorizeResult

const (
	Derived      = kernel.Derived
	Accommodated = kernel.Accommodated
	Postulated   = kernel.Postulated
)

// Re-export service types.
type LedgerService = ledger.Service
type AuditService = audit.Service

// Re-export kernel constructor.
var NewStore = kernel.New

// Re-export ledger constructor.
var NewLedger = ledger.New

// Re-export audit constructor.
var NewAudit = audit.New
