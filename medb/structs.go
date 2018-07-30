package medb

// 专门处理查询结果与structs的关系

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"reflect"
)

type field struct {
	sf   *reflect.StructField
	kind reflect.Kind
}

// 获取结构的元素
func getDataFieldsMap(t reflect.Type) (map[string]*field, error) {
	return nil, nil
}

// 创建数据类型解析的数组
func makeDataScanSlice(t reflect.Type, cols []string) ([]interface{}, error) {
	dfMap, err := getDataFieldsMap(t)
	if err != nil {
		return nil, err
	}
	ret := make([]interface{}, len(cols))
	for _, col := range cols {
		if df, ok := dfMap[col]; ok {
			switch df.kind {
			case reflect.Bool:
				ret = append(ret, &sql.NullBool{})
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				ret = append(ret, &sql.NullInt64{})
			case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
				ret = append(ret, &sql.NullInt64{})
			case reflect.Float32, reflect.Float64:
				ret = append(ret, &sql.NullFloat64{})
			case reflect.String:
				ret = append(ret, &sql.NullString{})
			default:
				return nil, fmt.Errorf("unsupported field kind %s", df.kind.String())
			}
		} else {
			return nil, fmt.Errorf("can not find col %s", col)
		}
	}
	return ret, nil
}

// // 返回去创建反射值
// func makeDataValues(t reflect.Type, vals []interface{}, cols []string) (*reflect.Value, error) {
// 	dfMap, err := getDataFieldsMap(t)
// 	if err != nil {
// 		return nil, err
// 	}
// 	for _, val := range vals {
// 		for _, col := range cols {
// 			if df, ok := dfMap[col]; ok {
// 				switch df.kind {
// 				case reflect.Bool:

// 				case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
// 					ret = append(ret, &sql.NullInt64{})
// 				case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
// 					ret = append(ret, &sql.NullInt64{})
// 				case reflect.Float32, reflect.Float64:
// 					ret = append(ret, &sql.NullFloat64{})
// 				case reflect.String:
// 					ret = append(ret, &sql.NullString{})
// 				default:
// 					return nil, fmt.Errorf("unsupported field kind %s", df.kind.String())
// 				}
// 			} else {
// 				return nil, fmt.Errorf("can not find col %s", col)
// 			}
// 		}
// 	}
// }

// ScanSlice 解析结构数组
func (r *Rows) ScanSlice(data interface{}) (int, error) {
	if r.err != nil {
		return 0, r.err
	}

	v := reflect.Indirect(reflect.ValueOf(data))
	if v.Kind() != reflect.Slice {
		return 0, errors.New("not slice or ptr to slice")
	}
	// t := v.Elem().Type()  // 数组元素类型
	et := v.Elem().Type() // 数组元素底层类型
	if et.Kind() == reflect.Ptr {
		et = et.Elem()
	}

	// // 获取到结构体的解析映射
	// fieldsMap, err := getDataFieldsMap(et)
	// if err != nil {
	// 	return 0, err
	// }

	var cols []string
	if r.columns == nil {
		cols, err := r.Columns()
		if err != nil {
			return 0, err
		}
		r.columns = make(map[string]int, len(cols))
		for i, col := range cols {
			r.columns[col] = i
		}
	}

	defer r.Close()
	for r.Next() {
		rd, err := makeDataScanSlice(et, cols)
		if err != nil {
			return 0, err
		}
		err = r.Scan(rd...)
		if err != nil {
			return 0, err
		}
	}

	return 0, nil
}

// ScanTo 解析
func (r *Rows) ScanTo(data interface{}) (int, error) {
	if r.err == nil {
		var d, err = newData(data) //	类型解析
		if err != nil {
			return 0, err
		}
		//	行解析，逐行解析rows的数据
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
			log.Println(value.Type().String(), value.Type().Name())
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

// data 解析目标的描述
type data struct {
	t        reflect.Type
	v        reflect.Value
	slice    bool
	destType reflect.Type
	length   int
}

// newData 生成一个data的描述
//	t,v value的反射类型和反射值,如果是指针型则返回指向的类型、值
//	slice 是否是切片类型
//	destType 如果是切片类型，则指向切片元素所对应的类型
// 	data must be kind of ptr
func newData(value interface{}) (*data, error) {
	var d = new(data)
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
func (r *data) next() bool {
	if r.slice {
		return true
	}
	return r.length < 1
}

// newValue 获取一个要生成的目标值的反射类型
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

//=============
