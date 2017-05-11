package meorm

import (
	"context"

	"github.com/ipiao/mesql/medb"
)

// SQLBuilder sql constructor
// adjust for context
type SQLBuilder interface {
	ToSQL() (string, []interface{})
}

// Exec 执行
func (c *Builder) Exec(s SQLBuilder) *medb.Result {
	sql, args := s.ToSQL()
	return c.Executor.Exec(sql, args...)
}

// ExecContext 执行
func (c *Builder) ExecContext(ctx context.Context, s SQLBuilder) *medb.Result {
	sql, args := s.ToSQL()
	return c.Executor.ExecContext(ctx, sql, args...)
}

// Query 查询
func (c *Builder) Query(s SQLBuilder) *medb.Rows {
	sql, args := s.ToSQL()
	return c.Executor.Query(sql, args...)
}

// QueryContext 查询
func (c *Builder) QueryContext(ctx context.Context, s SQLBuilder) *medb.Rows {
	sql, args := s.ToSQL()
	return c.Executor.QueryContext(ctx, sql, args...)
}

// PrepareExec 预处理执行
func (c *Builder) PrepareExec(s SQLBuilder) *medb.Result {
	sql, args := s.ToSQL()
	return c.Executor.Prepare(sql).Exec(args...)
}

// PrepareContextExec 预处理执行
func (c *Builder) PrepareContextExec(ctx context.Context, s SQLBuilder) *medb.Result {
	sql, args := s.ToSQL()
	return c.Executor.PrepareContext(ctx, sql).Exec(args...)
}

// PrepareExecContext 预处理执行
func (c *Builder) PrepareExecContext(ctx context.Context, s SQLBuilder) *medb.Result {
	sql, args := s.ToSQL()
	return c.Executor.Prepare(sql).ExecContext(ctx, args...)
}

// PrepareQuery 预处理查询
func (c *Builder) PrepareQuery(s SQLBuilder) *medb.Rows {
	sql, args := s.ToSQL()
	return c.Executor.Prepare(sql).Query(args...)
}

// PrepareContextQuery 预处理查询
func (c *Builder) PrepareContextQuery(ctx context.Context, s SQLBuilder) *medb.Rows {
	sql, args := s.ToSQL()
	return c.Executor.PrepareContext(ctx, sql).Query(args...)
}

// PrepareQueryContext 预处理查询
func (c *Builder) PrepareQueryContext(ctx context.Context, s SQLBuilder) *medb.Rows {
	sql, args := s.ToSQL()
	return c.Executor.Prepare(sql).QueryContext(ctx, args...)
}
