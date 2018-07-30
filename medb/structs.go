package medb

// 专门处理查询结果与structs的关系

import (
	"database/sql"
	"errors"
	"reflect"
)

// ScanTo 解析
func (r *Rows) ScanTo(data interface{}) (int, error) {
	if r.err != nil {
		return 0, r.err
	}
	var d, err = newData(data) //	类型解析
	if err != nil {
		return 0, err
	}
	//行解析，逐行解析rows的数据
	for r.Next() && d.next() {
		var v = d.newValue()
		err = r.scan(v)
		if err != nil {
			return 0, err
		}
		d.setBack(v)
	}
	err = r.Close()
	return d.length, err
}

// scan 单行解析
func (r *Rows) scan(v reflect.Value) error {
	colM, err := r.ColumnsMap()

	if err != nil {
		return err
	}
	var fields = make([]interface{}, len(colM))
	for i := 0; i < len(fields); i++ {
		var pif interface{}
		fields[i] = &pif
	}
	err = r.Scan(fields...)
	if err == nil {
		err = r.parse(v, 0, fields)
	}
	return err
}

// parse 将fields转换成相应的类型并绑定到value上
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
			if value.Type().String() == "time.Time" { //时间结构体解析
				err := timeparse(&value, *(fields[index].(*interface{})))
				if err != nil {
					return err
				}
			} else { //常规结构体解析
				for i := 0; i < value.NumField(); i++ {
					var fieldValue = value.Field(i)
					var fieldType = value.Type().Field(i)
					if fieldType.Anonymous {
						//匿名字段递归解析
						r.parse(fieldValue, 0, fields)
					} else {
						//非匿名字段
						if fieldValue.CanSet() {
							var tagMap = ParseTag(fieldType.Tag.Get(MedbTag))
							var fieldName = tagMap[MedbFieldName]
							if fieldName == "_" {
								continue
							}
							if fieldName == "" {
								fieldName = SnakeName(fieldType.Name)
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
	case reflect.Ptr:
		indValue := reflect.New(value.Type().Elem()).Elem()
		err := r.parse(indValue, index, fields)
		if err != nil {
			return err
		}
		value.Set(indValue.Addr())
	}
	return nil
}

// metadata 解析目标的描述
type metadata struct {
	t        reflect.Type  // 原始类型
	v        reflect.Value // 原始值
	slice    bool          // 是否是数组
	destType reflect.Type  // 元素类型
	length   int           // 迭代长度
}

// newData 生成一个data的描述
//	t,v value的反射类型和反射值,如果是指针型则返回指向的类型、值
//	slice 是否是切片类型
//	destType 如果是切片类型，则指向切片元素所对应的类型
// 	data must be kind of ptr
func newData(value interface{}) (*metadata, error) {
	var d = new(metadata)
	d.t = reflect.TypeOf(value)
	d.v = reflect.ValueOf(value)

	if d.t.Kind() != reflect.Ptr {
		return nil, errors.New("destination data must be kind of ptr")
	}

	d.t = d.t.Elem()
	d.v = d.v.Elem()

	switch d.t.Kind() {
	case reflect.Slice:
		d.slice = true
		d.destType = d.t.Elem()
	default:
		d.destType = d.t
	}
	return d, nil
}

// next 能否继续获取
func (r *metadata) next() bool {
	if r.slice {
		return true
	}
	return r.length < 1
}

// newValue 获取一个要生成的目标值的反射类型
func (r *metadata) newValue() reflect.Value {
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
func (r *metadata) setBack(value reflect.Value) {
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
