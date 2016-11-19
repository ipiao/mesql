package meorm

import "ipiao/mesql/medb"

// 解析器
type Executor interface {
	Exec() *medb.Result
	QueryTo(dest interface{}) (int, error)
	QueryNext(dest ...interface{}) error
}

type Scop interface {
	ToSQL()
}
