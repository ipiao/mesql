package meorm

import (
	"errors"
	"reflect"

	"fmt"

	"github.com/ipiao/mesql/medb"
)

// InsertBuilder insert构造器
type InsertBuilder struct {
	connname string
	table    string
	columns  []string
	values   [][]interface{}
	sql      string
	args     []interface{}
	err      error
}

// reset
func (b *InsertBuilder) reset() *InsertBuilder {
	b.table = ""
	b.columns = b.columns[:0]
	b.values = make([][]interface{}, 0, 0)
	b.err = nil
	b.sql = ""
	b.args = b.args[:0]
	return b
}

// Columns 插入列
func (b *InsertBuilder) Columns(columns ...string) *InsertBuilder {
	// 支持重复构造
	b.columns = append(b.columns, columns...)
	return b
}

// Values 值
func (b *InsertBuilder) Values(values ...interface{}) *InsertBuilder {
	b.values = append(b.values, values)
	return b
}

// Models 插入结构体
// models必须为结构体、结构体数组，或者相应的指针
func (b *InsertBuilder) Models(models interface{}) *InsertBuilder {
	var t = reflect.TypeOf(models)
	var v = reflect.ValueOf(models)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
		v = v.Elem()
	}
	var cols = GetColumns(v)
	var vals = GetValues(v)
	if len(b.columns) == 0 || len(cols) == 0 {
		b.err = errors.New("columns can not be null")
		return b
	}
	//
	if len(b.columns) == 1 {
		if b.columns[0] == "*" {
			b.columns = cols
		}
	}
	// 获取列名和结构体字段列的映射
	var tempMap = make(map[int]int, len(b.columns))
	for i, column := range b.columns {
		flag := 0
		for j, col := range cols {
			if column == col {
				tempMap[i] = j
				flag++
				break
			}
		}
		if flag == 0 {
			b.err = fmt.Errorf("can not find column %s in models", column)
		}
	}
	// 拼接值
	for _, val := range vals {
		var value = make([]interface{}, len(b.columns))
		for i, v := range tempMap {
			value[i] = val[v]
		}
		b.values = append(b.values, value)
	}
	return b
}

// tosql
func (b *InsertBuilder) tosql() (string, []interface{}) {
	mutex.Lock()
	defer mutex.Unlock()
	buf := bufPool.Get()
	defer bufPool.Put(buf)

	buf.WriteString("INSERT INTO ")
	buf.WriteString(b.table)
	buf.WriteString(" (")

	for i, col := range b.columns {
		if i > 0 {
			buf.WriteString(" ,")
		}
		buf.WriteString(col)
	}
	buf.WriteString(") VALUES")
	var args []interface{}
	for i, value := range b.values {
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
	b.sql = buf.String()
	b.args = args
	return b.sql, b.args
}

// ToSQL tosql
func (b *InsertBuilder) ToSQL() (string, []interface{}) {
	if len(b.sql) > 0 {
		return b.sql, b.args
	}
	return b.tosql()
}

// Exec 执行
func (b *InsertBuilder) Exec() *medb.Result {
	var res = new(medb.Result)
	if len(b.sql) == 0 {
		b.tosql()
	}
	if b.err != nil {
		res.SetErr(b.err)
		return res
	}
	return connections[b.connname].Exec(b.sql, b.args...)
}
