package meorm

import (
	"context"
	"errors"

	"github.com/ipiao/mesql/medb"
)

// SQLBuilder sql constructor
// adjust for context
type SQLBuilder interface {
	ToSQL() (string, []interface{})
}

// Exec 执行
func (c *BaseBuilder) Exec(s SQLBuilder) *medb.Result {
	return c.ExecContext(context.TODO(), s)
}

// ExecContext 执行
func (c *BaseBuilder) ExecContext(ctx context.Context, s SQLBuilder) *medb.Result {
	if s == nil {
		return new(medb.Result).SetErr(errors.New("builder can not be nil"))
	}
	sql, args := s.ToSQL()
	return c.Executor.ExecContext(ctx, sql, args...)
}

// Query 查询
func (c *BaseBuilder) Query(s SQLBuilder) *medb.Rows {
	return c.QueryContext(context.TODO(), s)
}

// QueryContext 查询
func (c *BaseBuilder) QueryContext(ctx context.Context, s SQLBuilder) *medb.Rows {
	if s == nil {
		return &medb.Rows{}
	}
	sql, args := s.ToSQL()
	return c.Executor.QueryContext(ctx, sql, args...)
}

// PrepareExec 预处理执行
func (c *BaseBuilder) PrepareExec(s SQLBuilder) *medb.Result {
	return c.PrepareExecContext(context.TODO(), s)
}

// PrepareContextExec 预处理执行
func (c *BaseBuilder) PrepareContextExec(ctx context.Context, s SQLBuilder) *medb.Result {
	if s == nil {
		return new(medb.Result).SetErr(errors.New("builder can not be nil"))
	}
	sql, args := s.ToSQL()
	return c.Executor.PrepareContext(ctx, sql).Exec(args...)
}

// PrepareExecContext 预处理执行
func (c *BaseBuilder) PrepareExecContext(ctx context.Context, s SQLBuilder) *medb.Result {
	if s == nil {
		return new(medb.Result).SetErr(errors.New("builder can not be nil"))
	}
	sql, args := s.ToSQL()
	return c.Executor.Prepare(sql).ExecContext(ctx, args...)
}

// PrepareQuery 预处理查询
func (c *BaseBuilder) PrepareQuery(s SQLBuilder) *medb.Rows {
	return c.PrepareQueryContext(context.TODO(), s)
}

// PrepareContextQuery 预处理查询
func (c *BaseBuilder) PrepareContextQuery(ctx context.Context, s SQLBuilder) *medb.Rows {
	if s == nil {
		return &medb.Rows{}
	}
	sql, args := s.ToSQL()
	return c.Executor.PrepareContext(ctx, sql).Query(args...)
}

// PrepareQueryContext 预处理查询
func (c *BaseBuilder) PrepareQueryContext(ctx context.Context, s SQLBuilder) *medb.Rows {
	if s == nil {
		return &medb.Rows{}
	}
	sql, args := s.ToSQL()
	return c.Executor.Prepare(sql).QueryContext(ctx, args...)
}
