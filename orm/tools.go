package meorm

import (
	"fmt"
	"reflect"

	"github.com/ipiao/mesql/medb"
	metools "github.com/ipiao/metools/utils"
)

// GetTableName get the table name of a obj
// ths obj intends to be kind of struct/slice/ptr
func GetTableName(obj interface{}) string {
	var v = reflect.Indirect(reflect.ValueOf(obj))
	return getTableName(v)
}

// getTableName 获取结构体对应的表的名字,v必须为结构体
func getTableName(v reflect.Value) string {
	var tbName string
	if v.Kind() == reflect.Struct {
		var tbnameV = v.MethodByName("TableName")
		if tbnameV.IsValid() {
			var args = make([]reflect.Value, 0)
			tbName = tbnameV.Call(args)[0].String()
		} else {
			tbName = metools.SnakeName(v.Type().Name())
		}
	} else if v.Kind() == reflect.Slice || v.Kind() == reflect.Array {
		return getTableName(reflect.Indirect(v.Index(0)))
	} else {
		return fmt.Sprintf("Invalid kind %s", v.Kind().String())
	}
	return tbName
}

// getColumns 获取要插入的表列,v必须为结构体
func getColumns(v reflect.Value) []string {
	var columns []string
	if v.Kind() == reflect.Struct {
		var l = v.NumField()
		for i := 0; i < l; i++ {
			var f = v.Type().Field(i)
			if f.Anonymous {
				columns = append(columns, getColumns(v.Field(i))...)
			} else {
				var tagMap = medb.ParseTag(f.Tag.Get(ormTag))
				var colName = tagMap[ormFieldSelectTag]
				if colName == "" {
					colName = metools.SnakeName(f.Name)
				}
				if colName == "_" {
					continue
				}
				columns = append(columns, colName)
			}
		}
	} else if v.Kind() == reflect.Slice || v.Kind() == reflect.Array {
		return getColumns(reflect.Indirect(v.Index(0)))
	} else {
		panic(fmt.Sprintf("Error kind %s", v.Kind().String()))
	}
	return columns
}

// getValues 获取值
func getValues(v reflect.Value) [][]interface{} {
	if v.Kind() == reflect.Struct {
		var values = make([][]interface{}, 1)
		var l = v.NumField()
		for i := 0; i < l; i++ {
			var f = v.Type().Field(i)
			if f.Anonymous {
				values[0] = append(values[0], getValues(v.Field(i))[0]...)
			} else {
				var tagMap = medb.ParseTag(f.Tag.Get(ormTag))
				var colName = tagMap[ormFieldSelectTag]
				if colName == "_" {
					continue
				}
				values[0] = append(values[0], v.Field(i).Interface())
			}
		}
		return values
	} else if v.Kind() == reflect.Slice || v.Kind() == reflect.Array {
		var values = make([][]interface{}, 0)
		for i := 0; i < v.Len(); i++ {
			values = append(values, getValues(reflect.Indirect(v.Index(i)))...)
		}
		return values
	} else {
		panic(fmt.Sprintf("Error kind %s", v.Kind().String()))
	}
}
