package meorm

import (
	"github.com/ipiao/mesql/medb"
)

// Conn 连接
type Conn struct {
	name string
	*medb.DB
}

// Name 返回连接名
func (c *Conn) Name() string {
	return c.name
}

// SQL 直接写sql
func (c *Conn) SQL(sql string, args ...interface{}) *CommonBuilder {
	return &CommonBuilder{
		sql:      sql,
		args:     args,
		connname: c.name,
	}
}

// Select 生成查询构造器
func (c *Conn) Select(cols ...string) *SelectBuilder {
	var builder = new(SelectBuilder).reset()
	builder.connname = c.name
	builder.columns = append(builder.columns, cols...)
	return builder
}

// Update 生成更新构造器
func (c *Conn) Update(table string) *UpdateBuilder {
	var builder = new(UpdateBuilder).reset()
	builder.connname = c.name
	builder.table = table
	return builder
}

// InsertInto 生成插入构造器
func (c *Conn) InsertInto(table string) *InsertBuilder {
	var builder = new(InsertBuilder).reset()
	builder.connname = c.name
	builder.table = table
	return builder
}

// // InsertModels 生成插入构造器
// // 表名有结构体方法 TableName 生成
// func (c *Conn) InsertModels(models interface{}) *InsertBuilder {
// 	var table = GetTableName(reflect.Indirect(reflect.ValueOf(models)))
// 	var builder = c.InsertInto(table).Models(models)
// 	builder.connname = c.name
// 	return builder
// }

// DeleteFrom 生成删除构造器
func (c *Conn) DeleteFrom(table string) *DeleteBuilder {
	var builder = new(DeleteBuilder).reset()
	builder.connname = c.name
	builder.table = table
	return builder
}
