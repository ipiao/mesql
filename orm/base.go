package meorm

import (
	"github.com/ipiao/mesql/medb"
	"github.com/ipiao/mesql/orm/dialect"
)

// Conn 连接
type Conn struct {
	*medb.DB
	dialect dialect.Dialect
}

// SQL 直接写sql
func (c *Conn) SQL(sql string, args ...interface{}) *CommonBuilder {
	return &CommonBuilder{
		sql:  sql,
		args: args,
	}
}

// Select 生成查询构造器
func (c *Conn) Select(cols ...string) *SelectBuilder {
	var builder = new(SelectBuilder).reset()
	builder.Conn = c
	builder.columns = append(builder.columns, cols...)
	return builder
}

// Update 生成更新构造器
func (c *Conn) Update(table string) *UpdateBuilder {
	var builder = new(UpdateBuilder).reset()
	builder.Conn = c
	builder.table = table
	return builder
}

// InsertOrUpdate 生成插入或更新构造器
func (c *Conn) InsertOrUpdate(table string) *InsupBuilder {
	var builder = new(InsupBuilder).reset()
	builder.Conn = c
	builder.table = table
	return builder
}

// InsertInto 生成插入构造器
func (c *Conn) InsertInto(table string) *InsertBuilder {
	var builder = new(InsertBuilder).reset()
	builder.Conn = c
	builder.table = table
	return builder
}

// ReplaceInto 生成插入构造器
func (c *Conn) ReplaceInto(table string) *InsertBuilder {
	var builder = new(InsertBuilder).reset()
	builder.Conn = c
	builder.table = table
	builder.replace = true
	return builder
}

// DeleteFrom 生成删除构造器
func (c *Conn) DeleteFrom(table string) *DeleteBuilder {
	var builder = new(DeleteBuilder).reset()
	builder.Conn = c
	builder.table = table
	return builder
}

// Delete 生成删除构造器
func (c *Conn) Delete(column string) *DeleteBuilder {
	var builder = new(DeleteBuilder).reset()
	builder.Conn = c
	builder.column = column
	return builder
}
