package meorm

import "github.com/ipiao/mesql/medb"

// BareBuilder
// 最基本的sql构造，直接写sql，裸构造
type BareBuilder struct {
	builder *BaseBuilder
	sql     string
	args    []interface{}
}

// AppendSQL 继续拼接
func (b *BareBuilder) AppendSQL(sql string, args ...interface{}) *BareBuilder {
	b.sql += " " + sql
	b.args = append(b.args, args...)
	return b
}

// Exec 查询不建议使用
func (b *BareBuilder) Exec() *medb.Result {
	return b.builder.Exec(b)
}

// QueryTo 解析到结构体，数组。。。
func (b *BareBuilder) QueryTo(models interface{}) (int, error) {
	return b.builder.Query(b).ScanTo(models)
}

// QueryNext 把查询组成sql并解析
func (b *BareBuilder) QueryNext(dest ...interface{}) error {
	return b.builder.Query(b).ScanNext(dest...)
}

// ToSQL 把查询组成sql并解析
func (b *BareBuilder) ToSQL() (string, []interface{}) {
	return b.sql, b.args
}
