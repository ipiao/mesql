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

func (this *Conn) Select(cols ...string) *selectBuilder {
	var selector = new(selectBuilder).reset()
	selector.connname = this.name
	selector.columns = append(selector.columns, cols...)
	return selector
}
