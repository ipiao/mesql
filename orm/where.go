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

type Where struct {
	where []*whereConstraint
	err   error
}

// 全匹配
func (this *Where) WhereLike(col string, arg interface{}) *Where {
	return this.whereLike(col, arg, 0)
}

// 左匹配
func (this *Where) WhereLikeL(col string, arg interface{}) *Where {
	return this.whereLike(col, arg, -1)
}

// 右匹配
func (this *Where) WhereLikeR(col string, arg interface{}) *Where {
	return this.whereLike(col, arg, 1)
}

// like -1表示左边匹配，0全匹配，1.右边匹配
func (this *Where) whereLike(col string, arg interface{}, likekind int8) *Where {
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
	this.where = append(this.where, where)
	return this
}

// In
func (this *Where) WhereIn(col string, args ...interface{}) *Where {
	return this.wherein(col, args)
}

// 查询条件in的解析
func (this *Where) wherein(col string, args interface{}) *Where {
	var v = reflect.Indirect(reflect.ValueOf(args))
	var k = v.Kind()
	if k == reflect.Slice || k == reflect.Array {
		if v.Len() == 0 {
			return this
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
			//		case reflect.Slice:
			//			if v.Len() > 1 {
			//				this.err = errors.New(fmt.Sprintf("WhereIn 参数错误"))
			//				return this
			//			}
			//			return this.wherein(col, v.Index(0).Interface())
		default:
			this.err = errors.New(fmt.Sprintf("in不支持的类型%s", v.Index(0).Elem().Kind().String()))
		}

		this.where = append(this.where, where)
	} else {
		this.err = errors.New("参数格式错误，必须为切片或数组")
	}
	return this
}
