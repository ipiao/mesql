package meorm

import (
	"github.com/ipiao/mesql/medb"
	"github.com/ipiao/mesql/orm/dialect"
)

// Builder 构造器
type Builder struct {
	medb.Executor
	dialect dialect.Dialect
}

// Commit 提交事务
func (c *Builder) Commit() error {
	return c.Executor.Commit()
}

// Rollback 回滚
func (c *Builder) Rollback() error {
	return c.Executor.Rollback()
}

// SQL 直接写sql
func (c *Builder) SQL(sql string, args ...interface{}) *CommonBuilder {
	return &CommonBuilder{
		builder: c,
		sql:     sql,
		args:    args,
	}
}

// Select 生成查询构造器
func (c *Builder) Select(cols ...string) *SelectBuilder {
	var builder = new(SelectBuilder).reset()
	builder.builder = c
	builder.columns = append(builder.columns, cols...)
	return builder
}

// Update 生成更新构造器
func (c *Builder) Update(table string) *UpdateBuilder {
	var builder = new(UpdateBuilder).reset()
	builder.builder = c
	builder.table = table
	return builder
}

// InsertOrUpdate 生成插入或更新构造器
func (c *Builder) InsertOrUpdate(table string) *InsupBuilder {
	var builder = new(InsupBuilder).reset()
	builder.builder = c
	builder.table = table
	return builder
}

// InsertInto 生成插入构造器
func (c *Builder) InsertInto(table string) *InsertBuilder {
	var builder = new(InsertBuilder).reset()
	builder.builder = c
	builder.table = table
	return builder
}

// ReplaceInto 生成插入构造器
func (c *Builder) ReplaceInto(table string) *InsertBuilder {
	var builder = new(InsertBuilder).reset()
	builder.builder = c
	builder.table = table
	builder.replace = true
	return builder
}

// DeleteFrom 生成删除构造器
func (c *Builder) DeleteFrom(table string) *DeleteBuilder {
	var builder = new(DeleteBuilder).reset()
	builder.builder = c
	builder.table = table
	return builder
}

// Delete 生成删除构造器
func (c *Builder) Delete(column string) *DeleteBuilder {
	var builder = new(DeleteBuilder).reset()
	builder.builder = c
	builder.column = column
	return builder
}

// NewBuilder 创建无连接构造器
func NewBuilder() *Builder {
	return &Builder{}
}
