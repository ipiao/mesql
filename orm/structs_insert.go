package meorm

import (
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/ipiao/mesql/medb"
)

// 插入结构体或结构体数组
func (this *Conn) InsertModels(models interface{}) *medb.Result {

	var value = reflect.Indirect(reflect.ValueOf(models))
	var k = value.Kind()
	switch k {
	case reflect.Struct:
		return this.insertStruct(&value)
	case reflect.Slice, reflect.Array:
		return this.insertSlice(&value)
	default:
		panic(fmt.Sprintf("Error kind %s", k.String()))
	}
}

// 插入结构体
// mysql 的主键遇 0 值可以自动忽略
func (this *Conn) insertStruct(v *reflect.Value) *medb.Result {
	buf := bufPool.Get()
	defer bufPool.Put(buf)

	var tbName = GetTableName(*v)
	var cols = GetColumns(*v)
	var values = GetValues(*v)

	if tbName == "" || len(cols) == 0 {
		return &medb.Result{
			Err: errors.New("Error struct"),
		}
	}

	buf.WriteString("INSERT INTO ")
	buf.WriteString(tbName)
	buf.WriteString(" (")

	var valueStr = "("
	for i, col := range cols {
		if i > 0 {
			buf.WriteString(" ,")
			valueStr += " ,"
		}
		buf.WriteString(col)
		valueStr += "?"
	}
	buf.WriteString(") VALUES ")
	valueStr += ")"
	buf.WriteString(valueStr)
	var args = values[0]
	return this.db.Exec(buf.String(), args...)
}

// 插入数组
func (this *Conn) insertSlice(v *reflect.Value) *medb.Result {
	buf := bufPool.Get()
	defer bufPool.Put(buf)

	var tbName = GetTableName(*v)
	var cols = GetColumns(*v)
	var values = GetValues(*v)

	if tbName == "" || len(cols) == 0 {
		return &medb.Result{
			Err: errors.New("Error slice"),
		}
	}

	buf.WriteString("INSERT INTO ")
	buf.WriteString(tbName)
	buf.WriteString(" (")

	var valueStr = "("
	for i, col := range cols {
		if i > 0 {
			buf.WriteString(" ,")
			valueStr += " ,"
		}
		buf.WriteString(col)
		valueStr += "?"
	}
	buf.WriteString(") VALUES ")
	valueStr += "),"

	valueStr = strings.Repeat(valueStr, len(values))
	buf.WriteString(valueStr[:len(valueStr)-1])
	var args []interface{}
	for i := range values {
		args = append(args, values[i]...)
	}
	return this.db.Exec(buf.String(), args...)
}
