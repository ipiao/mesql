package medb

import (
	"database/sql"
	"reflect"
	"time"
)

// Rows 数据行
type Rows struct {
	*sql.Rows
	err     error
	columns map[string]int
}

// Row 单行
type Row struct {
	row *sql.Row
}

// 返回错误信息
func (r *Rows) Error() error {
	return r.err
}

// ScanNext 组合scan和next
func (r *Rows) ScanNext(dest ...interface{}) error {
	if r.err != nil {
		return r.err
	}
	for r.Next() {
		r.Scan(dest...)
	}
	return r.Close()
}

// parse
func (r *Rows) parse(value reflect.Value, index int, fields []interface{}) error {
	switch value.Kind() {
	case reflect.Bool:
		var b = sql.NullBool{}
		var err = b.Scan(*(fields[index].(*interface{})))
		if err != nil {
			return err
		}
		if b.Valid {
			value.SetBool(b.Bool)
		}
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		var i = sql.NullInt64{}
		var err = i.Scan(*(fields[index].(*interface{})))
		if err != nil {
			return err
		}
		if i.Valid {
			value.SetInt(i.Int64)
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		var i = sql.NullInt64{}
		var err = i.Scan(*(fields[index].(*interface{})))
		if err != nil {
			return err
		}
		if i.Valid {
			value.SetUint(uint64(i.Int64))
		}
	case reflect.Float32, reflect.Float64:
		var f = sql.NullFloat64{}
		var err = f.Scan(*(fields[index].(*interface{})))
		if err != nil {
			return err
		}
		if f.Valid {
			value.SetFloat(f.Float64)
		}
	case reflect.String:
		var s = sql.NullString{}
		var err = s.Scan(*(fields[index].(*interface{})))
		if err != nil {
			return err
		}
		if s.Valid {
			value.SetString(s.String)
		}
	case reflect.Struct:
		{
			if value.Type().String() == "time.Time" {
				//时间结构体解析
				var s = sql.NullString{}
				var err = s.Scan(*(fields[index].(*interface{})))
				if err != nil {
					return err
				}
				if s.Valid {
					t, err := time.ParseInLocation("2006-01-02 15:04:05", s.String, time.Local)
					if err != nil {
						t, err = time.ParseInLocation("2006-01-02", s.String, time.Local)
					}
					if err == nil {
						value.Set(reflect.ValueOf(t))
					}
				} else {
					var i = sql.NullInt64{}
					var err = i.Scan(*(fields[index].(*interface{})))
					if err != nil {
						return err
					}
					if i.Valid {
						t := time.Unix(i.Int64, 0)
						if err == nil {
							value.Set(reflect.ValueOf(t))
						}
					}
				}
			} else {
				//常规结构体解析
				for i := 0; i < value.NumField(); i++ {
					var fieldValue = value.Field(i)
					var fieldType = value.Type().Field(i)
					if fieldType.Anonymous {
						//匿名字段递归解析
						r.parse(fieldValue, 0, fields)
					} else {
						//非匿名字段
						if fieldValue.CanSet() {
							var fieldName = fieldType.Tag.Get("db")
							if fieldName == "_" {
								continue
							}
							if fieldName == "" {
								fieldName = transFieldName(fieldType.Name)
							}
							var index, ok = r.columns[fieldName]
							if ok {
								r.parse(fieldValue, index, fields)
							}
						}
					}
				}
			}
		}
	}
	return nil
}

// scan 单行解析
func (r *Rows) scan(v reflect.Value) error {
	if r.columns == nil {
		var cols, err = r.Columns()
		if err != nil {
			return err
		}
		r.columns = make(map[string]int, len(cols))
		for i, col := range cols {
			r.columns[col] = i
		}
	}
	var fields = make([]interface{}, len(r.columns))
	for i := 0; i < len(fields); i++ {
		var pif interface{}
		fields[i] = &pif
	}
	var err = r.Scan(fields...)
	if err == nil {
		err = r.parse(v, 0, fields)
	}
	return err
}

// ScanTo 解析
func (r *Rows) ScanTo(data interface{}) (int, error) {
	if r.err == nil {
		var d, err = newData(data)
		//	类型解析
		if err != nil {
			return 0, err
		}
		//	行解析
		for r.Next() && d.next() {
			var v = d.newValue()
			err = r.scan(v)
			if err != nil {
				return 0, err
			}
			d.setBack(v)
		}
		err = r.Close()
		return d.length, nil
	}
	return 0, r.err
}

// data解析目标的描述
type data struct {
	t        reflect.Type
	v        reflect.Value
	slice    bool
	destType reflect.Type
	length   int
}

// 生成一个data的描述
func newData(value interface{}) (*data, error) {
	var d = new(data)
	d.t = reflect.TypeOf(value)
	d.v = reflect.ValueOf(value)
	if d.t.Kind() == reflect.Ptr {
		d.t = d.t.Elem()
		d.v = d.v.Elem()
	}
	switch d.t.Kind() {
	case reflect.Slice:
		{
			d.slice = true
			d.destType = d.t.Elem()
		}
	default:
		{
			d.destType = d.t
		}
	}
	return d, nil
}

// newValue 获取一个可Set的值
func (r *data) newValue() reflect.Value {
	r.length++
	if r.slice {
		var v reflect.Value
		if r.destType.Kind() == reflect.Ptr {
			v = reflect.New(r.destType.Elem()).Elem()
		} else {
			v = reflect.New(r.destType).Elem()
		}
		return v
	}
	return r.v
}

// setBack 将newValue的值设置回data
func (r *data) setBack(value reflect.Value) {
	if r.slice {
		var v = value
		if r.destType.Kind() == reflect.Ptr {
			v = v.Addr()
		}
		r.v.Set(reflect.Append(r.v, v))
	} else {
		r.v = value
	}
}

// next 能否继续获取
func (r *data) next() bool {
	if r.slice {
		return true
	}
	return r.length < 1
}
