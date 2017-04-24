package meorm

import "github.com/ipiao/mesql/medb"

// Builder sql构造器
type Builder interface {
	Tosql() (string, []interface{})
	Exec() *medb.Result
}

// Executor 解析器
type Executor interface {
	Exec() *medb.Result
	QueryTo(dest interface{}) (int, error)
	QueryNext(dest ...interface{}) error
}
