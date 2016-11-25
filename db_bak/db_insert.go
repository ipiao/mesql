package mesql

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
)

//
func (this *DB) InsertModels(model interface{}) (int, error) {
	var value = reflect.ValueOf(model)
	if value.Kind() == reflect.Ptr || value.Kind() == reflect.Interface {
		value = value.Elem()
	}
	if value.Kind() == reflect.Struct {
		return this.insertStruct("", value)
	}
	if value.Kind() == reflect.Slice {
		return 0, this.insertSlice("", value)
	}
	return -1, errors.New("you hava an error with your model kind:" + value.Kind().String())
}

//
func (this *DB) InsertModels2(table string, model interface{}) (int, error) {
	var value = reflect.ValueOf(model)
	if value.Kind() == reflect.Ptr || value.Kind() == reflect.Interface {
		value = value.Elem()
	}
	if value.Kind() == reflect.Struct {
		return this.insertStruct(table, value)
	}
	if value.Kind() == reflect.Slice {
		return 0, this.insertSlice(table, value)
	}
	return -1, errors.New("you hava an error with your model kind:" + value.Kind().String())
}

//
func (this *DB) insertStruct(table string, value reflect.Value) (int, error) {
	if table == "" {
		table = transFieldName(value.Type().Name())
	}
	var cols = GetColumns(value)
	var vals = GetValues(value)
	var tempVals = make([]string, len(cols))
	for i := 0; i < len(cols); i++ {
		tempVals[i] = "?"
	}
	var insert = "insert into `" + table + "`"
	var colsStr = "(" + strings.Join(cols, ",") + ")"
	var valsStr = "values (" + strings.Join(tempVals, ",") + ")"
	var sql = insert + colsStr + valsStr
	fmt.Println("[mesql]sql:", sql, "	values:", vals)
	var r, err = this.Exec(sql, vals...)
	if err != nil {
		return -1, err
	}
	var newId, _ = r.LastInsertId()
	return int(newId), err
}

//
func (this *DB) insertSlice(table string, value reflect.Value) error {
	//	if value.Kind() == reflect.Interface || value.Kind() == reflect.Ptr {
	//		value = value.Elem()
	//	}

	if value.Len() < 1 {
		return errors.New("the length of slice must ge 1")
	}
	var md = value.Index(0)
	if table == "" {
		table = transFieldName(md.Type().Name())
	}
	var cols = GetColumns(value)
	var vals = GetValues(value)

	var tempVal = make([]string, len(cols))
	for i := 0; i < len(cols); i++ {
		tempVal[i] = "?"
	}
	var tempVals = make([]string, value.Len())
	for i := 0; i < value.Len(); i++ {
		tempVals[i] = "(" + strings.Join(tempVal, ",") + ")"
	}
	var insert = "insert into `" + table + "`"
	var colsStr = "(" + strings.Join(cols, ",") + ")"

	var sql = insert + colsStr + " values " + strings.Join(tempVals, ",")
	fmt.Println("[mesql]sql:", sql, "	values:", vals)
	var _, err = this.Exec(sql, vals...)
	return err
}
