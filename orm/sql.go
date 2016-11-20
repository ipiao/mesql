package meorm

import "ipiao/mesql/medb"

// 最基本的sql构造，直接写sql
type commonBuilder struct {
	Executor
	connname string
	sql      string
	args     []interface{}
}

// 查询不建议使用
func (this *commonBuilder) Exec() *medb.Result {
	return connections[this.connname].db.Exec(this.sql, this.args...)
}

// 解析到结构体，数组。。。
func (this *commonBuilder) QueryTo(models interface{}) (int, error) {
	return connections[this.connname].db.Query(this.sql, this.args...).ScanTo(models)
}

// 把查询组成sql并解析
func (this *commonBuilder) QueryNext(dest ...interface{}) error {
	return connections[this.connname].db.Query(this.sql, this.args...).ScanNext(dest...)
}
