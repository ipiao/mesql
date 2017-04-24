package meorm

import (
	"errors"
	"fmt"
	"reflect"
)

type whereConstraint struct {
	condition string
	values    []interface{}
}

// Where where
type Where struct {
	where []*whereConstraint
	err   error
}

// WhereLike 全匹配
func (w *Where) WhereLike(col string, arg interface{}) *Where {
	return w.whereLike(col, arg, 0)
}

// WhereLikeL 左匹配
func (w *Where) WhereLikeL(col string, arg interface{}) *Where {
	return w.whereLike(col, arg, -1)
}

// WhereLikeR 右匹配
func (w *Where) WhereLikeR(col string, arg interface{}) *Where {
	return w.whereLike(col, arg, 1)
}

// like -1表示左边匹配，0全匹配，1.右边匹配
func (w *Where) whereLike(col string, arg interface{}, likekind int8) *Where {
	var where = &whereConstraint{
		condition: fmt.Sprintf("%s LIKE ?", col),
	}
	var value interface{}
	if likekind == 0 {
		value = fmt.Sprint("%", fmt.Sprintf("%v", arg), "%")
	} else if likekind == -1 {
		value = fmt.Sprint("%", fmt.Sprintf("%v", arg))
	} else if likekind == 1 {
		value = fmt.Sprint(fmt.Sprintf("%v", arg), "%")
	}
	where.values = append(where.values, value)
	w.where = append(w.where, where)
	return w
}

// WhereIn In
func (w *Where) WhereIn(col string, args ...interface{}) *Where {
	return w.wherein(col, args...)
}

// 查询条件in的解析
func (w *Where) wherein(col string, args ...interface{}) *Where {
	var v = reflect.Indirect(reflect.ValueOf(args))
	var k = v.Kind()
	if k == reflect.Slice || k == reflect.Array {
		if v.Len() == 0 {
			return w
		}
		var buf = bufPool.Get()
		defer bufPool.Put(buf)
		var where = new(whereConstraint)
		buf.WriteString(fmt.Sprintf("%s IN(", col))
		for i := 0; i < v.Len(); i++ {
			if i > 0 {
				buf.WriteString(" ,?")
			} else {
				buf.WriteRune('?')
			}
		}
		buf.WriteRune(')')
		where.condition = buf.String()
		switch v.Index(0).Elem().Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			for i := 0; i < v.Len(); i++ {
				where.values = append(where.values, v.Index(i).Elem().Int())
			}
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			for i := 0; i < v.Len(); i++ {
				where.values = append(where.values, v.Index(i).Elem().Uint())
			}

		case reflect.Float32, reflect.Float64:
			for i := 0; i < v.Len(); i++ {
				where.values = append(where.values, v.Index(i).Elem().Float())
			}
		case reflect.Bool:
			for i := 0; i < v.Len(); i++ {
				where.values = append(where.values, v.Index(i).Bool())
			}
		case reflect.String:
			for i := 0; i < v.Len(); i++ {
				where.values = append(where.values, v.Index(i).Elem().String())
			}
		case reflect.Slice, reflect.Array:
			if v.Len() > 1 {
				w.err = errors.New("WhereIn 参数错误")
				return w
			}
			v = v.Index(0).Elem()
			var params []interface{}
			for i := 0; i < v.Len(); i++ {
				params = append(params, v.Index(i).Interface())
			}
			return w.wherein(col, params...)
		default:
			w.err = fmt.Errorf("in不支持的类型%s", v.Index(0).Elem().Kind().String())
		}

		w.where = append(w.where, where)
	} else {
		w.err = errors.New("参数格式错误，必须为切片或数组")
	}
	return w
}

// 查询条件in的解析
func (w *Where) wherenotin(col string, args ...interface{}) *Where {
	var v = reflect.Indirect(reflect.ValueOf(args))
	var k = v.Kind()
	if k == reflect.Slice || k == reflect.Array {
		if v.Len() == 0 {
			return w
		}
		var buf = bufPool.Get()
		defer bufPool.Put(buf)
		var where = new(whereConstraint)
		buf.WriteString(fmt.Sprintf("%s NOT IN(", col))
		for i := 0; i < v.Len(); i++ {
			if i > 0 {
				buf.WriteString(" ,?")
			} else {
				buf.WriteRune('?')
			}
		}
		buf.WriteRune(')')
		where.condition = buf.String()
		switch v.Index(0).Elem().Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			for i := 0; i < v.Len(); i++ {
				where.values = append(where.values, v.Index(i).Elem().Int())
			}
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			for i := 0; i < v.Len(); i++ {
				where.values = append(where.values, v.Index(i).Elem().Uint())
			}

		case reflect.Float32, reflect.Float64:
			for i := 0; i < v.Len(); i++ {
				where.values = append(where.values, v.Index(i).Elem().Float())
			}
		case reflect.Bool:
			for i := 0; i < v.Len(); i++ {
				where.values = append(where.values, v.Index(i).Bool())
			}
		case reflect.String:
			for i := 0; i < v.Len(); i++ {
				where.values = append(where.values, v.Index(i).Elem().String())
			}
		case reflect.Slice, reflect.Array:
			if v.Len() > 1 {
				w.err = errors.New("WhereIn 参数错误")
				return w
			}
			v = v.Index(0).Elem()
			var params []interface{}
			for i := 0; i < v.Len(); i++ {
				params = append(params, v.Index(i).Interface())
			}
			return w.wherein(col, params...)
		default:
			w.err = fmt.Errorf("in不支持的类型%s", v.Index(0).Elem().Kind().String())
		}

		w.where = append(w.where, where)
	} else {
		w.err = errors.New("参数格式错误，必须为切片或数组")
	}
	return w
}
