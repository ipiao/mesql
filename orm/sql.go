package meorm

import "github.com/ipiao/mesql/medb"

// CommonBuilder 最基本的sql构造，直接写sql
type CommonBuilder struct {
	Executor
	connname string
	sql      string
	args     []interface{}
}

// AppendSQL 继续拼接
func (comm *CommonBuilder) AppendSQL(sql string, args ...interface{}) *CommonBuilder {
	comm.sql += " " + sql
	comm.args = append(comm.args, args...)
	return comm
}

// Exec 查询不建议使用
func (comm *CommonBuilder) Exec() *medb.Result {
	return connections[comm.connname].Exec(comm.sql, comm.args...)
}

// QueryTo 解析到结构体，数组。。。
func (comm *CommonBuilder) QueryTo(models interface{}) (int, error) {
	return connections[comm.connname].Query(comm.sql, comm.args...).ScanTo(models)
}

// QueryNext 把查询组成sql并解析
func (comm *CommonBuilder) QueryNext(dest ...interface{}) error {
	return connections[comm.connname].Query(comm.sql, comm.args...).ScanNext(dest...)
}

// ToSQL 把查询组成sql并解析
func (comm *CommonBuilder) ToSQL() (string, []interface{}) {
	return comm.sql, comm.args
}
