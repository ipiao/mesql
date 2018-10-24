package meorm

import (
	"fmt"

	"github.com/ipiao/mesql/meorm/dialect"
)

type condValues struct {
	condition string
	values    []interface{}
}

// where where
type where struct {
	conds   []*condValues
	dialect dialect.Dialect
	err     error
}

// Where 条件
func (w *where) Where(condition string, args ...interface{}) *where {
	w.conds = append(w.conds, &condValues{
		condition: condition,
		values:    args,
	})
	return w
}

// WhereLike 全匹配
func (w *where) WhereLike(col string, arg interface{}) *where {
	return w.whereLike(col, arg, 0)
}

// WhereLikeL 左匹配
func (w *where) WhereLikeL(col string, arg interface{}) *where {
	return w.whereLike(col, arg, -1)
}

// WhereLikeR 右匹配
func (w *where) WhereLikeR(col string, arg interface{}) *where {
	return w.whereLike(col, arg, 1)
}

// like -1表示左边匹配，0全匹配，1.右边匹配
func (w *where) whereLike(col string, arg interface{}, likekind int8) *where {
	holder := w.dialect.Holder()
	var condValues = &condValues{
		condition: fmt.Sprintf("%s LIKE %c", col, holder),
	}
	var value interface{}
	if likekind == 0 {
		value = "%" + fmt.Sprintf("%v", arg) + "%"
	} else if likekind == -1 {
		value = "%" + fmt.Sprintf("%v", arg)
	} else if likekind == 1 {
		value = fmt.Sprintf("%v", arg) + "%"
	}
	condValues.values = append(condValues.values, value)
	w.conds = append(w.conds, condValues)
	return w
}

// WhereIn In
func (w *where) WhereIn(col string, args ...interface{}) *where {
	return w.wherein(col, args...)
}

// 查询条件in的解析
func (w *where) wherein(col string, args ...interface{}) *where {
	// if len(args) == 0 {
	// 	w.err = fmt.Errorf("length of args in method wherein %s can not be 0", col)
	// 	return w
	// }
	holder := w.dialect.Holder()
	var buf = bufPool.Get()
	defer bufPool.Put(buf)
	var condValues = new(condValues)
	buf.WriteString(fmt.Sprintf("%s IN(", col))
	for i := 0; i < len(args); i++ {
		if i > 0 {
			buf.WriteString(" ,")
		}
		buf.WriteByte(holder)
	}
	buf.WriteByte(')')
	condValues.condition = buf.String()
	condValues.values = args

	// switch v.Index(0).Elem().Kind() {
	// case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
	// 	for i := 0; i < v.Len(); i++ {
	// 		where.values = append(where.values, v.Index(i).Elem().Int())
	// 	}
	// case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
	// 	for i := 0; i < v.Len(); i++ {
	// 		where.values = append(where.values, v.Index(i).Elem().Uint())
	// 	}

	// case reflect.Float32, reflect.Float64:
	// 	for i := 0; i < v.Len(); i++ {
	// 		where.values = append(where.values, v.Index(i).Elem().Float())
	// 	}
	// case reflect.Bool:
	// 	for i := 0; i < v.Len(); i++ {
	// 		where.values = append(where.values, v.Index(i).Bool())
	// 	}
	// case reflect.String:
	// 	for i := 0; i < v.Len(); i++ {
	// 		where.values = append(where.values, v.Index(i).Elem().String())
	// 	}
	// case reflect.Slice, reflect.Array:
	// 	if v.Len() > 1 {
	// 		w.err = errors.New("WhereIn 参数错误")
	// 		return w
	// 	}
	// 	v = v.Index(0).Elem()
	// 	var params []interface{}
	// 	for i := 0; i < v.Len(); i++ {
	// 		params = append(params, v.Index(i).Interface())
	// 	}
	// 	return w.wherein(col, params...)
	// default:
	// 	w.err = fmt.Errorf("in不支持的类型%s", v.Index(0).Elem().Kind().String())
	// }
	w.conds = append(w.conds, condValues)
	return w
}

// 查询条件in的解析
func (w *where) wherenotin(col string, args ...interface{}) *where {
	if len(args) == 0 {
		// w.err = errors.New("length of args in method wherenotin can not be 0")
		return w
	}

	holder := w.dialect.Holder()
	var buf = bufPool.Get()
	defer bufPool.Put(buf)
	var condValues = new(condValues)
	buf.WriteString(fmt.Sprintf("%s NOT IN(", col))
	for i := 0; i < len(args); i++ {
		if i > 0 {
			buf.WriteString(" ,")
		}
		buf.WriteByte(holder)
	}
	buf.WriteByte(')')
	condValues.condition = buf.String()
	w.conds = append(w.conds, condValues)
	return w
}
