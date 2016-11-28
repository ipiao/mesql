package meorm

import (
	"errors"
	"fmt"

	"github.com/ipiao/mesql/medb"
)

// insert构造器
type insertBuilder struct {
	Executor
	connname string
	table    string
	columns  []string
	values   [][]interface{}
	sql      string
	args     []interface{}
	err      error
}

// reset
func (this *insertBuilder) reset() *insertBuilder {
	this.table = ""
	this.columns = this.columns[:0]
	this.values = make([][]interface{}, 0, 0)
	this.err = nil
	this.sql = ""
	this.args = this.args[:0]
	return this
}

// 插入列
func (this *insertBuilder) Columns(columns ...string) *insertBuilder {
	// 支持重复构造
	this.columns = append(this.columns, columns...)
	return this
}

// 值
func (this *insertBuilder) Values(values ...interface{}) *insertBuilder {
	// 支持多条插入
	//	mutex.Lock()
	//	defer mutex.Unlock()
	if len(this.columns) != len(values) && len(this.columns) > 0 {
		this.err = errors.New(fmt.Sprintf("values %v 的长度 %d 不匹配 columns 的长度 %d",
			values, len(values), len(this.columns)))
	}
	this.values = append(this.values, values)
	return this
}

// tosql
func (this *insertBuilder) tosql() (string, []interface{}) {
	mutex.Lock()
	defer mutex.Unlock()
	buf := bufPool.Get()
	defer bufPool.Put(buf)

	buf.WriteString("INSERT INTO ")
	buf.WriteString(this.table)
	buf.WriteString(" (")

	for i, col := range this.columns {
		if i > 0 {
			buf.WriteString(" ,")
		}
		buf.WriteString(col)
	}

	buf.WriteString(") VALUES")

	var args []interface{}
	for i, value := range this.values {
		if i > 0 {
			buf.WriteRune(',')
		}
		buf.WriteString(" (")
		for j, val := range value {
			if j > 0 {
				buf.WriteString(" ,?")
			} else {
				buf.WriteString("?")
			}
			args = append(args, val)
		}
		buf.WriteString(")")
	}
	this.sql = buf.String()
	this.args = args
	return this.sql, this.args
}

// tosql
func (this *insertBuilder) ToSQL() (string, []interface{}) {
	if len(this.sql) > 0 {
		return this.sql, this.args
	}
	return this.tosql()
}

// 执行
func (this *insertBuilder) Exec() *medb.Result {
	if len(this.sql) == 0 {
		this.tosql()
	}
	if this.err != nil {
		var res = &medb.Result{
			Err: this.err,
		}
		return res
	}
	return connections[this.connname].db.Exec(this.sql, this.args...)
}

// 解析到结构体，数组。。。
func (this *insertBuilder) QueryTo(models interface{}) (int, error) {
	//	if len(this.sql) == 0 {
	//		this.tosql()
	//	}
	//	if this.err != nil {
	//		return 0, this.err
	//	}
	//	return connections[this.connname].db.Query(this.sql, this.args...).ScanTo(models)
	return 0, errors.New("[meorm]:Insert 不能使用该方法")
}

// 把查询组成sql并解析
func (this *insertBuilder) QueryNext(dest ...interface{}) error {
	//	if len(this.sql) == 0 {
	//		this.tosql()
	//	}
	//	if this.err != nil {
	//		return this.err
	//	}
	//	return connections[this.connname].db.Query(this.sql, this.args...).ScanNext(dest...)
	return errors.New("[meorm]:Insert 不能使用该方法")
}
