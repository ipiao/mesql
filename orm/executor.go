package meorm

// 解析器
type Executor interface {
	Exec()
	QueryTo()
	QueryNext()
}

type Scop interface {
	ToSQL()
}
