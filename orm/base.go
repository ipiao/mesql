package meorm

import "ipiao/mesql/medb"

// 连接
type Conn struct {
	db       *medb.DB
	selector *selectBuilder
}

func (this *Conn) Select(cols ...string) *selectBuilder {
	if this.selector == nil {
		this.selector = new(selectBuilder)
	}
	this.selector.reset()
	this.selector.columns = append(this.selector.columns, cols...)
	return this.selector
}
