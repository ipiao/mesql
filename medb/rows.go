package medb

import (
	"database/sql"
	"reflect"
	"time"
)

// 单行
type Row struct {
	row *sql.Row
}

// 返回错误信息
func (this *Rows) Error() error {
	return this.err
}

// Columns,列名
func (this *Rows) Columns() []string {
	var cols, _ = this.rows.Columns()
	return cols
}

// 迭代
func (this *Rows) Next() bool {
	return this.rows.Next()
}

// 解析
func (this *Rows) Scan(dest ...interface{}) error {
	return this.rows.Scan(dest...)
}

// 组合scan和next
func (this *Rows) ScanNext(dest ...interface{}) error {
	if this.err != nil {
		return this.err
	}
	for this.rows.Next() {
		return this.rows.Scan(dest...)
	}
	return nil
}

// 关闭
func (this *Rows) Close() error {
	if this.rows.Next() {
		return this.rows.Close()
	}
	return nil
}

// parse
func (this *Rows) parse(value reflect.Value, index int, fields []interface{}) error {
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
					if err == nil {
						value.Set(reflect.ValueOf(t))
					}
				}
			} else {
				//常规结构体解析
				for i := 0; i < value.NumField(); i++ {
					var fieldValue = value.Field(i)
					var fieldType = value.Type().Field(i)
					if fieldType.Anonymous {
						//匿名字段递归解析
						this.parse(fieldValue, 0, fields)
					} else {
						//非匿名字段
						if fieldValue.CanSet() {
							var fieldName = fieldType.Tag.Get(DefaultTagName)
							if fieldName == "_" {
								continue
							}
							if fieldName == "" {
								fieldName = transFieldName(fieldType.Name)
							}
							var index, ok = this.columns[fieldName]
							if ok {
								this.parse(fieldValue, index, fields)
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
func (this *Rows) scan(v reflect.Value) error {
	if this.columns == nil {
		var cols, err = this.rows.Columns()
		if err != nil {
			return err
		}
		this.columns = make(map[string]int, len(cols))
		for i, col := range cols {
			this.columns[col] = i
		}
	}
	var fields = make([]interface{}, len(this.columns))
	for i := 0; i < len(fields); i++ {
		var pif interface{}
		fields[i] = &pif
	}
	var err = this.rows.Scan(fields...)
	if err == nil {
		err = this.parse(v, 0, fields)
	}
	return err
}

// ScanTo
func (this *Rows) ScanTo(data interface{}) (int, error) {
	if this.err == nil {
		var d, err = newData(data)
		//	类型解析
		if err != nil {
			return 0, err
		}
		//	行解析
		for this.rows.Next() && d.next() {
			var v = d.newValue()
			err = this.scan(v)
			if err != nil {
				return 0, err
			}
			d.setBack(v)
		}
		err = this.rows.Close()
		return d.length, nil
	}
	return 0, this.err
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
func (this *data) newValue() reflect.Value {
	this.length++
	if this.slice {
		var v reflect.Value
		if this.destType.Kind() == reflect.Ptr {
			v = reflect.New(this.destType.Elem()).Elem()
		} else {
			v = reflect.New(this.destType).Elem()
		}
		return v
	}
	return this.v
}

// setBack 将newValue的值设置回data
func (this *data) setBack(value reflect.Value) {
	if this.slice {
		var v = value
		if this.destType.Kind() == reflect.Ptr {
			v = v.Addr()
		}
		this.v.Set(reflect.Append(this.v, v))
	} else {
		this.v = value
	}
}

// next 能否继续获取
func (this *data) next() bool {
	if this.slice {
		return true
	}
	return this.length < 1
}

// 数据行
type Rows struct {
	rows    *sql.Rows
	err     error
	columns map[string]int
}

// 解析到结构体时使用Rows
func (this Row) Scan(dest ...interface{}) error {
	return this.row.Scan(dest...)
}
