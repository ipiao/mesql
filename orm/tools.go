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

// SnakeName 驼峰转蛇形
func SnakeName(base string) string {
	var r = make([]rune, 0, len(base))
	var b = []rune(base)
	for i := 0; i < len(b); i++ {
		if i > 0 && b[i] >= 'A' && b[i] <= 'Z' {
			r = append(r, '_', b[i]+32)
			continue
		}
		if i == 0 && b[i] >= 'A' && b[i] <= 'Z' {
			r = append(r, b[i]+32)
			continue
		}
		r = append(r, b[i])
	}
	return string(r)
}

// GetTableName 获取结构体对应的表的名字,v必须为结构体
func GetTableName(v reflect.Value) string {
	var tbName string
	if v.Kind() == reflect.Struct {
		var tbnameV = v.MethodByName("TableName")
		if tbnameV.IsValid() {
			var args = make([]reflect.Value, 0)
			tbName = tbnameV.Call(args)[0].String()
		} else {
			tbName = transFieldName(v.Type().Name())
		}
	} else if v.Kind() == reflect.Slice || v.Kind() == reflect.Array {
		return GetTableName(reflect.Indirect(v.Index(0)))
	} else {
		panic(fmt.Sprintf("Error kind %s", v.Kind().String()))
	}
	return tbName
}

// GetColumns 获取要插入的表列,v必须为结构体
func GetColumns(v reflect.Value) []string {
	var columns []string
	if v.Kind() == reflect.Struct {
		var l = v.NumField()
		for i := 0; i < l; i++ {
			var f = v.Type().Field(i)
			if f.Anonymous {
				columns = append(columns, GetColumns(v.Field(i))...)
			} else {
				var tagMap = medb.ParseTag(f.Tag.Get(ormTag))
				var colName = tagMap[ormFieldSelectTag]
				if colName == "" {
					colName = transFieldName(f.Name)
				}
				if colName == "_" {
					continue
				}
				columns = append(columns, colName)
			}
		}
	} else if v.Kind() == reflect.Slice || v.Kind() == reflect.Array {
		return GetColumns(reflect.Indirect(v.Index(0)))
	} else {
		panic(fmt.Sprintf("Error kind %s", v.Kind().String()))
	}
	return columns
}

// GetValues 获取值
func GetValues(v reflect.Value) [][]interface{} {
	if v.Kind() == reflect.Struct {
		var values = make([][]interface{}, 1)
		var l = v.NumField()
		for i := 0; i < l; i++ {
			var f = v.Type().Field(i)
			if f.Anonymous {
				values[0] = append(values[0], GetValues(v.Field(i))[0]...)
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
			values = append(values, GetValues(reflect.Indirect(v.Index(i)))...)
		}
		return values
	} else {
		panic(fmt.Sprintf("Error kind %s", v.Kind().String()))
	}
}
