package application

import (
	"context"
	"database/sql"
	"fmt"
	"sync"
)

// faultTx wraps *sql.Tx and intercepts ExecContext calls for deterministic
// fault injection. It satisfies the TxExecutor interface.
type faultTx struct {
	*sql.Tx
	mu         sync.Mutex
	execCount  int
	failAfterN int
	failErr    error
}

func (f *faultTx) ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.execCount++
	if f.execCount > f.failAfterN {
		return nil, f.failErr
	}
	return f.Tx.ExecContext(ctx, query, args...)
}

// newFaultTx wraps a real *sql.Tx with deterministic fault injection.
// It will fail ExecContext after failAfterN successful calls.
func newFaultTx(tx *sql.Tx, failAfterN int) *faultTx {
	return &faultTx{
		Tx:         tx,
		failAfterN: failAfterN,
		failErr:    fmt.Errorf("injected fault after %d execs", failAfterN),
	}
}

// execCount returns the number of ExecContext calls made so far.
func (f *faultTx) execCountTotal() int {
	f.mu.Lock()
	defer f.mu.Unlock()
	return f.execCount
}
