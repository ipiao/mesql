package meorm

import "github.com/ipiao/mesql/medb"

// 连接
type Conn struct {
	name string
	db   *medb.DB
}

// 返回连接名
func (this *Conn) Name() string {
	return this.name
}

// 返回连接名
func (this *Conn) DB() *medb.DB {
	return this.db
}

// 直接写sql
func (this *Conn) SQL(sql string, args ...interface{}) *commonBuilder {
	return &commonBuilder{
		sql:      sql,
		args:     args,
		connname: this.name,
	}
}

// 生成查询构造器
func (this *Conn) Select(cols ...string) *selectBuilder {
	var builder = new(selectBuilder).reset()
	builder.connname = this.name
	builder.columns = append(builder.columns, cols...)
	return builder
}

// 生成更新构造器
func (this *Conn) Update(table string) *updateBuilder {
	var builder = new(updateBuilder).reset()
	builder.connname = this.name
	builder.table = table
	return builder
}

// 生成插入构造器
func (this *Conn) InsertInto(table string) *insertBuilder {
	var builder = new(insertBuilder).reset()
	builder.connname = this.name
	builder.table = table
	return builder
}

// 生成删除构造器
func (this *Conn) DeleteFrom(table string) *deleteBuilder {
	var builder = new(deleteBuilder).reset()
	builder.connname = this.name
	builder.table = table
	return builder
}

// 生成删除构造器
func (this *Conn) Delete(column string) *deleteBuilder {
	var builder = new(deleteBuilder).reset()
	builder.connname = this.name
	builder.column = column
	return builder
}
