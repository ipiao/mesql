package meorm

import (
	"github.com/ipiao/mesql/medb"
	"github.com/ipiao/mesql/meorm/dialect"
)

// BaseBuilder 基础构造器构造器
// 有执行器和上下文构成
// 相当于一次性的执行任务
// 构建上下文环境，dialect
type BaseBuilder struct {
	medb.Executor
	dialect dialect.Dialect
}

// 默认mysql
func NewBuilder(e medb.Executor) *BaseBuilder {
	return NewDialectBuilder(dialect.Mysql, e)
}

func NewDialectBuilder(dialect dialect.Dialect, e medb.Executor) *BaseBuilder {
	return &BaseBuilder{
		Executor: e,
		dialect:  dialect,
	}
}

// SQL 直接写sql
func (c *BaseBuilder) SQL(sql string, args ...interface{}) *BareBuilder {
	return &BareBuilder{
		builder: c,
		sql:     sql,
		args:    args,
	}
}

// Select 生成查询构造器
func (c *BaseBuilder) Select(cols ...string) *SelectBuilder {
	var builder = new(SelectBuilder).reset()
	builder.builder = c
	builder.columns = append(builder.columns, cols...)
	return builder
}

// Update 生成更新构造器
func (c *BaseBuilder) Update(table string) *UpdateBuilder {
	var builder = new(UpdateBuilder).reset()
	builder.builder = c
	builder.table = table
	return builder
}

// InsertOrUpdate 生成插入或更新构造器
func (c *BaseBuilder) InsertOrUpdate(table string) *InsupBuilder {
	var builder = new(InsupBuilder).reset()
	builder.builder = c
	builder.table = table
	return builder
}

// InsertInto 生成插入构造器
func (c *BaseBuilder) InsertInto(table string) *InsertBuilder {
	var builder = new(InsertBuilder).reset()
	builder.builder = c
	builder.table = table
	return builder
}

// ReplaceInto 生成插入构造器
func (c *BaseBuilder) ReplaceInto(table string) *InsertBuilder {
	var builder = new(InsertBuilder).reset()
	builder.builder = c
	builder.table = table
	builder.replace = true
	return builder
}

// DeleteFrom 生成删除构造器
func (c *BaseBuilder) DeleteFrom(table string) *DeleteBuilder {
	var builder = new(DeleteBuilder).reset()
	builder.builder = c
	builder.table = table
	return builder
}

// Delete 生成删除构造器
func (c *BaseBuilder) Delete(column string) *DeleteBuilder {
	var builder = new(DeleteBuilder).reset()
	builder.builder = c
	builder.column = column
	return builder
}
