package meorm

import (
	"fmt"
	"reflect"
	"regexp"
	"strings"

	"github.com/ipiao/mesql/medb"
)

var reg = regexp.MustCompile(`\B[A-Z]`)

// transFieldName 转换字段名称
func transFieldName(name string) string {
	return strings.ToLower(reg.ReplaceAllString(name, "_$0"))
}

// 获取结构体对应的表的名字,v必须为结构体
func GetTableName(v reflect.Value) string {
	var tbName string
	if v.Kind() == reflect.Struct {
		var tbnameV = v.MethodByName(TableNameMethod)
		if tbnameV.IsValid() {
			var args = make([]reflect.Value, 0)
			tbName = tbnameV.Call(args)[0].String()
		} else {
			tbName = transFieldName(v.Type().Name())
		}
	} else if v.Kind() == reflect.Slice || v.Kind() == reflect.Array {
		return GetTableName(v.Index(0))
	} else {
		panic(fmt.Sprintf("Error kind %s", v.Kind().String()))
	}
	return tbName
}

// 获取要插入的表列,v必须为结构体
func GetColumns(v reflect.Value) []string {
	var columns []string
	if v.Kind() == reflect.Struct {
		var l = v.NumField()
		for i := 0; i < l; i++ {
			var f = v.Type().Field(i)
			if f.Anonymous {
				columns = append(columns, GetColumns(v.Field(i))...)
			} else {
				var colName = f.Tag.Get(medb.MeTag)
				if colName == "" || colName == "_" {
					colName = transFieldName(f.Name)
				}
				var ormtag = f.Tag.Get(OrmTag)
				if ormtag == "_" {
					continue
				}
				columns = append(columns, colName)
			}
		}
	} else if v.Kind() == reflect.Slice || v.Kind() == reflect.Array {
		return GetColumns(v.Index(0))
	} else {
		panic(fmt.Sprintf("Error kind %s", v.Kind().String()))
	}
	return columns
}

// 获取值
func GetValues(v reflect.Value) [][]interface{} {
	if v.Kind() == reflect.Struct {
		var values [][]interface{} = make([][]interface{}, 1)
		var l = v.NumField()
		for i := 0; i < l; i++ {
			var f = v.Type().Field(i)
			if f.Anonymous {
				values[0] = append(values[0], GetValues(v.Field(i))[0]...)
			} else {
				var ormtag = f.Tag.Get(OrmTag)
				if ormtag == "_" {
					continue
				}
				values[0] = append(values[0], v.Field(i).Interface())
			}
		}
		return values
	} else if v.Kind() == reflect.Slice || v.Kind() == reflect.Array {
		var values = make([][]interface{}, 0)
		for i := 0; i < v.Len(); i++ {
			values = append(values, GetValues(v.Index(i))...)
		}
		return values
	} else {
		panic(fmt.Sprintf("Error kind %s", v.Kind().String()))
	}
}
