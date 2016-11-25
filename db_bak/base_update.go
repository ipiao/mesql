package mesql

//生成修改语句sql
import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"
)

// Indirect 初始化model
func Indirect(models interface{}) reflect.Value {
	var value = reflect.ValueOf(models)
	if value.Kind() == reflect.Interface || value.Kind() == reflect.Ptr {
		value = value.Elem()
	}
	return value
}

// Update 单表修改 0值跳过
func (this *BaseService) Update(model interface{}) bool {
	var value = Indirect(model)
	var sql = "update `" + SnakeName(value.Type().Name()) + "` set "
	this.parseUpdate(value)
	var columns string
	if len(this.condition) > 0 {
		for _, v := range this.condition {
			columns += v[1] + ","
		}
	}
	sql += columns[:len(columns)-1]
	var id = value.FieldByName("Id")
	if id.IsValid() {
		sql += " where id= " + strconv.Itoa(int(id.Int()))
	}
	var _, e = this.DB.Exec(sql, this.GetValues()...)
	if e != nil {
		fmt.Println("sql错误:", e)
		return false
	}
	return true
}
func (this *BaseService) parseUpdate(t reflect.Value) {
	if t.Kind() != reflect.Struct {
		return
	}
	for i := 0; i < t.NumField(); i++ {
		if t.Type().Field(i).Anonymous {
			this.parseUpdate(t.Field(i))
			continue
		}
		var key = SnakeName(t.Type().Field(i).Name)
		if strings.ToLower(key) == "id" {
			continue
		}
		var field = t.Field(i)
		switch field.Kind() {
		case reflect.Bool:
		case reflect.Struct:
			if field.Type().Name() == "Time" {
				var tStr = fmt.Sprintf("%v", field.Interface())
				var tm, err = time.Parse("2006-01-02", strings.Split(tStr, " ")[0])
				if err != nil || tm.IsZero() {
					continue
				}
				this.setCondtion(key, "`"+key+"`=?")
				this.setValue(field.Interface())
			}
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			if field.Int() != 0 {
				this.setCondtion(key, "`"+key+"`=?")
				this.setValue(field.Int())
			}
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			if field.Uint() != 0 {
				this.setCondtion(key, "`"+key+"`=?")
				this.setValue(field.Uint())
			}
		case reflect.Float32, reflect.Float64:
			if field.Float() != 0 {
				this.setCondtion(key, "`"+key+"`=?")
				this.setValue(field.Float())
			}
		case reflect.Interface:
			this.parseUpdate(field.Field(i))
		case reflect.String:
			if field.String() != "" {
				this.setCondtion(key, "`"+key+"`=?")
				this.setValue(field.String())
			}
		}
	}
}

// UpdateModels 批量更新(0值字段跳过)
//func (this *BaseService) UpdateModels(models interface{}) {
//	//	var value = Indirect(models)
//	//	var columns = GetColumns(value)
//	//	_ = columns

//}

// UpdateModels 批量更新(指定字段)
func (this *BaseService) UpdateColumn(models interface{}, columns ...string) (int, error) {
	if len(columns) < 1 {
		return 0, errors.New("无指定字段")
	}
	var value = Indirect(models)
	var tableName string
	if value.Kind() == reflect.Slice {
		tableName = "`" + SnakeName(value.Index(0).Type().Name()) + "`"
	} else if value.Kind() == reflect.Struct {
		//		tableName = "`" + SnakeName(value.Type().Name()) + "`"
	} else {
		return 0, errors.New("类型错误")
	}
	var sql = "update " + tableName + " SET "
	var length = value.Len()
	if length < 1 {
		return 0, errors.New("数组不能为空")
	}
	var setStr = " CASE id " + strings.Repeat(" WHEN ? THEN ? ", length) + " END "
	var colStr []string
	for i := 0; i < len(columns); i++ {
		var cc = "`" + columns[i] + "`=" + setStr
		colStr = append(colStr, cc)
	}
	sql += strings.Join(colStr, ",") + " WHERE id in (" +
		strings.Join(strings.Split(strings.Repeat("?", length), ""), ",") + ")"
	var vals = make(map[string][]interface{}, len(columns))
	for i := 0; i < len(columns); i++ {
		vals[columns[i]] = nil
	}
	GetValueMap(value, vals)
	var sqlVal []interface{}
	for _, v := range columns {
		for i := 0; i < length; i++ {
			sqlVal = append(sqlVal, vals["id"][i], vals[v][i])
		}
	}
	//将where条件的id拼上
	sqlVal = append(sqlVal, vals["id"]...)
	var res, err = this.DB.Exec(sql, sqlVal...)
	if err != nil {
		return 0, err
	}
	var n, _ = res.RowsAffected()
	return int(n), err
}

// UpdateModels 批量更新(指定字段)
func (this *BaseService) UpdateColumnWithTable(models interface{}, tableName string, columns ...string) (int, error) {
	if len(columns) < 1 {
		return 0, errors.New("无指定字段")
	}
	var value = Indirect(models)
	if value.Kind() == reflect.Slice {
		//tableName = "`" + SnakeName(value.Index(0).Type().Name()) + "`"
	} else if value.Kind() == reflect.Struct {
		//		tableName = "`" + SnakeName(value.Type().Name()) + "`"
	} else {
		return 0, errors.New("类型错误")
	}

	var sql = "update " + tableName + " SET "
	var length = value.Len()
	if length < 1 {
		return 0, errors.New("数组不能为空")
	}
	var setStr = " CASE id " + strings.Repeat(" WHEN ? THEN ? ", length) + " END "
	var colStr []string
	for i := 0; i < len(columns); i++ {
		var cc = "`" + columns[i] + "`=" + setStr
		colStr = append(colStr, cc)
	}
	sql += strings.Join(colStr, ",") + " WHERE id in (" +
		strings.Join(strings.Split(strings.Repeat("?", length), ""), ",") + ")"
	var vals = make(map[string][]interface{}, len(columns))
	for i := 0; i < len(columns); i++ {
		vals[columns[i]] = nil
	}
	GetValueMap(value, vals)
	var sqlVal []interface{}
	for _, v := range columns {
		for i := 0; i < length; i++ {
			sqlVal = append(sqlVal, vals["id"][i], vals[v][i])
		}
	}
	//将where条件的id拼上
	sqlVal = append(sqlVal, vals["id"]...)
	var res, err = this.DB.Exec(sql, sqlVal...)
	if err != nil {
		return 0, err
	}
	var n, _ = res.RowsAffected()
	return int(n), err
}

func GetValueMap(t reflect.Value, vals map[string][]interface{}) {
	if t.Kind() == reflect.Slice {
		for i := 0; i < t.Len(); i++ {
			GetValueMap(t.Index(i), vals)
		}
	}
	if t.Kind() == reflect.Struct {
		var n = t.NumField()
		for i := 0; i < n; i++ {
			var key = SnakeName(t.Type().Field(i).Name)
			if key != "id" {
				//检查key是否为指定字段
				if _, ok := vals[key]; !ok {
					continue
				}
			}
			if t.Type().Field(i).Anonymous {
				GetValueMap(t.Field(i), vals)
			} else {
				vals[key] = append(vals[key], t.Field(i).Interface())
			}
		}
	}
}

// SetColumn 如果字段的'0'值是必要的，用此方法加入
func (this *BaseService) SetColumn(key string, value interface{}) *BaseService {
	this.setCondtion(key, "`"+key+"`=?")
	this.setValue(value)
	return this
}
