package medb

import (
	"context"
)

// Executor 执行器
type Executor interface {
	ExecContext(ctx context.Context, sql string, args ...interface{}) *Result
	QueryContext(ctx context.Context, sql string, args ...interface{}) *Rows
	PrepareContext(ctx context.Context, sql string) Executor
	Commit() error
	Rollback() error
}
