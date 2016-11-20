package mesql

// sql构造器
type Builder interface {
	ToSql() string
	Exec() error
}

//
type builder struct {
	db *DB
}

//
func (this *builder) reset() {

}
