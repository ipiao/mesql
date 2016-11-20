package meorm

import "ipiao/mesql/medb"

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

// 生成查询构造器
func (this *Conn) Select(cols ...string) *selectBuilder {
	var builder = new(selectBuilder).reset()
	builder.connname = this.name
	builder.columns = append(builder.columns, cols...)
	return builder
}

func (this *Conn) SQL(sql string, args ...interface{}) *commonBuilder {
	return &commonBuilder{
		sql:      sql,
		args:     args,
		connname: this.name,
	}
}
