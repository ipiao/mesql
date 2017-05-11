package medb

import "context"

// Executor 执行器
type Executor interface {
	Exec(sql string, args ...interface{}) *Result
	ExecContext(ctx context.Context, sql string, args ...interface{}) *Result
	Query(sql string, args ...interface{}) *Rows
	QueryContext(ctx context.Context, sql string, args ...interface{}) *Rows
	Rollback() error
	Commit() error
}
