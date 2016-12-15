package meorm

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/ipiao/mesql/medb"
)

// 查询
type selectBuilder struct {
	Executor
	connname string
	distinct bool
	columns  []string
	from     string
	//where       []*whereConstraint
	where       *Where
	orderbys    []string
	groupbys    []string
	limit       int64
	limitvalid  bool
	offset      int64
	offsetvalid bool
	having      []*whereConstraint
	err         error
	sql         string
	args        []interface{}
}

// 选择查询列
func (this *selectBuilder) Select(columns ...string) *selectBuilder {
	this.columns = append(this.columns, columns...)
	return this
}

// distinct
func (this *selectBuilder) Distinct() *selectBuilder {
	this.distinct = true
	return this
}

// from
func (this *selectBuilder) From(from string) *selectBuilder {
	this.from = from
	return this
}

// order by
func (this *selectBuilder) OrderBy(order string) *selectBuilder {
	this.orderbys = append(this.orderbys, order)
	return this
}

// group by
func (this *selectBuilder) GroupBy(group string) *selectBuilder {
	this.groupbys = append(this.groupbys, group)
	return this
}

// limit
func (this *selectBuilder) Limit(limit int64) *selectBuilder {
	this.limitvalid = true
	this.limit = limit
	return this
}

// offset
func (this *selectBuilder) Offset(offset int64) *selectBuilder {
	this.offsetvalid = true
	this.offset = offset
	return this
}

// where
func (this *selectBuilder) Where(condition string, args ...interface{}) *selectBuilder {
	this.where.where = append(this.where.where, &whereConstraint{
		condition: condition,
		values:    args,
	})
	return this
}

// having
func (this *selectBuilder) Having(condition string, values ...interface{}) *selectBuilder {
	this.having = append(this.having, &whereConstraint{
		condition: condition,
		values:    values,
	})
	return this
}

// reset
func (this *selectBuilder) reset() *selectBuilder {
	this.distinct = false
	this.columns = this.columns[:0]
	this.from = ""
	this.where = new(Where)
	this.orderbys = this.orderbys[:0]
	this.groupbys = this.groupbys[:0]
	this.limit = 0
	this.limitvalid = false
	this.offset = 0
	this.offsetvalid = false
	this.having = make([]*whereConstraint, 0, 0)
	this.err = nil
	this.sql = ""
	this.args = this.args[:0]
	return this
}

// tosql
func (this *selectBuilder) ToSQL() (string, []interface{}) {
	if len(this.sql) > 0 {
		return this.sql, this.args
	}
	return this.tosql()
}

// 把查询条件组成sql并放到查询体中
func (this *selectBuilder) tosql() (string, []interface{}) {
	mutex.Lock()
	defer mutex.Unlock()
	if this.where.err != nil {
		this.err = this.where.err
		return "", nil
	}
	if len(this.columns) == 0 {
		panic("没有指定列")
	}
	if len(this.from) == 0 {
		panic("没有指定表")
	}
	buf := bufPool.Get()
	defer bufPool.Put(buf)

	var args []interface{}
	buf.WriteString("SELECT ")

	if this.distinct {
		buf.WriteString("DISTINCT ")
	}
	for i, s := range this.columns {
		if i > 0 {
			buf.WriteString(", ")
		}
		buf.WriteString(s)
	}
	buf.WriteString(" FROM ")
	buf.WriteString(this.from)

	if len(this.where.where) > 0 {
		buf.WriteString(" WHERE ")
		for i, cond := range this.where.where {
			if i > 0 {
				buf.WriteString(" AND (")
			} else {
				buf.WriteRune('(')
			}
			buf.WriteString(cond.condition)
			buf.WriteRune(')')
			if len(cond.values) > 0 {
				args = append(args, cond.values...)
			}
		}
	}
	if len(this.groupbys) > 0 {
		buf.WriteString(" GROUP BY ")
		for i, s := range this.groupbys {
			if i > 0 {
				buf.WriteString(", ")
			}
			buf.WriteString(s)
		}
	}
	if len(this.having) > 0 {
		buf.WriteString(" HAVING ")
		for i, cond := range this.having {
			if i > 0 {
				buf.WriteString(" AND (")
			} else {
				buf.WriteRune('(')
			}
			buf.WriteString(cond.condition)
			buf.WriteRune(')')
			if len(cond.values) > 0 {
				args = append(args, cond.values...)
			}
		}
	}
	if len(this.orderbys) > 0 {
		buf.WriteString(" ORDER BY ")
		for i, s := range this.orderbys {
			if i > 0 {
				buf.WriteString(", ")
			}
			buf.WriteString(s)
		}
	}
	if this.limitvalid {
		buf.WriteString(" LIMIT ?")
		args = append(args, this.limit)
	}
	if this.offsetvalid {
		buf.WriteString(" OFFSET ?")
		args = append(args, this.offset)
	}
	this.sql = buf.String()
	this.args = args
	return this.sql, this.args
}

// 查询不建议使用
func (this *selectBuilder) Exec() *medb.Result {
	if len(this.sql) == 0 {
		this.tosql()
	}
	if this.err != nil {
		var res = &medb.Result{
			Err: this.err,
		}
		return res
	}
	return connections[this.connname].db.Exec(this.sql, this.args...)
}

// 解析到结构体，数组。。。
func (this *selectBuilder) QueryTo(models interface{}) (int, error) {
	if len(this.sql) == 0 {
		this.tosql()
	}
	if this.err != nil {
		return 0, this.err
	}
	return connections[this.connname].db.Query(this.sql, this.args...).ScanTo(models)
}

// 把查询组成sql并解析
func (this *selectBuilder) QueryNext(dest ...interface{}) error {
	if len(this.sql) == 0 {
		this.tosql()
	}
	if this.err != nil {
		return this.err
	}
	return connections[this.connname].db.Query(this.sql, this.args...).ScanNext(dest...)
}

// limit和offset的复用
func (this *selectBuilder) LimitPP(page, pagesize int64) *selectBuilder {
	var offset = (page - 1) * pagesize
	return this.Limit(pagesize).Offset(offset)
}

// 查询符合条件的总数目
func (this *selectBuilder) CountCond(countCond ...string) (int64, string, error) {
	var countcond string
	var count int64
	if len(countCond) == 0 {
		countcond = "count(0)"
	} else {
		countcond = countCond[0]
	}
	//	if len(this.having) > 0 {
	//		var sql, args = this.countresult("alie", countcond)
	//		var err = connections[this.connname].db.Query(sql, args...).ScanNext(&count)
	//		return count, sql, err
	//	}
	var sql, args = this.countsql(countcond)
	var err = connections[this.connname].db.Query(sql, args...).ScanNext(&count)
	return count, sql, err
}

// 查询符合条件的总数目
func (this *selectBuilder) CountResult(alies string, countCond ...string) (int64, string, error) {
	var countcond string
	var count int64
	if len(countCond) == 0 {
		countcond = "count(0)"
	} else {
		countcond = countCond[0]
	}
	var sql, args = this.countresult(alies, countcond)
	var err = connections[this.connname].db.Query(sql, args...).ScanNext(&count)
	return count, sql, err
}

// 生成查询总条数的sql
func (this *selectBuilder) countsql(countCond string) (string, []interface{}) {
	buf := bufPool.Get()
	defer bufPool.Put(buf)

	var args []interface{}
	buf.WriteString("SELECT ")
	buf.WriteString(countCond)
	buf.WriteString(" FROM ")
	buf.WriteString(this.from)

	if len(this.where.where) > 0 {
		buf.WriteString(" WHERE ")
		for i, cond := range this.where.where {
			if i > 0 {
				buf.WriteString(" AND (")
			} else {
				buf.WriteRune('(')
			}
			buf.WriteString(cond.condition)
			buf.WriteRune(')')
			if len(cond.values) > 0 {
				args = append(args, cond.values...)
			}
		}
	}
	if len(this.groupbys) > 0 {
		buf.WriteString(" GROUP BY ")
		for i, s := range this.groupbys {
			if i > 0 {
				buf.WriteString(", ")
			}
			buf.WriteString(s)
		}
	}
	if len(this.having) > 0 {
		buf.WriteString(" HAVING ")
		for i, cond := range this.having {
			if i > 0 {
				buf.WriteString(" AND (")
			} else {
				buf.WriteRune('(')
			}
			buf.WriteString(cond.condition)
			buf.WriteRune(')')
			if len(cond.values) > 0 {
				args = append(args, cond.values...)
			}
		}
	}
	return buf.String(), args
}

// 把查询条件组成sql并放到查询体中
func (this *selectBuilder) countresult(alies, countCond string) (string, []interface{}) {
	buf := bufPool.Get()
	defer bufPool.Put(buf)
	var args []interface{}
	if this.where.err != nil {
		this.err = this.where.err
		return "", nil
	}
	buf.WriteString("SELECT ")
	buf.WriteString(countCond)
	buf.WriteString(" FROM( ")
	buf.WriteString("SELECT ")

	if this.distinct {
		buf.WriteString("DISTINCT ")
	}
	for i, s := range this.columns {
		if i > 0 {
			buf.WriteString(", ")
		}
		buf.WriteString(s)
	}
	buf.WriteString(" FROM ")
	buf.WriteString(this.from)

	if len(this.where.where) > 0 {
		buf.WriteString(" WHERE ")
		for i, cond := range this.where.where {
			if i > 0 {
				buf.WriteString(" AND (")
			} else {
				buf.WriteRune('(')
			}
			buf.WriteString(cond.condition)
			buf.WriteRune(')')
			if len(cond.values) > 0 {
				args = append(args, cond.values...)
			}
		}
	}
	if len(this.groupbys) > 0 {
		buf.WriteString(" GROUP BY ")
		for i, s := range this.groupbys {
			if i > 0 {
				buf.WriteString(", ")
			}
			buf.WriteString(s)
		}
	}
	if len(this.having) > 0 {
		buf.WriteString(" HAVING ")
		for i, cond := range this.having {
			if i > 0 {
				buf.WriteString(" AND (")
			} else {
				buf.WriteRune('(')
			}
			buf.WriteString(cond.condition)
			buf.WriteRune(')')
			if len(cond.values) > 0 {
				args = append(args, cond.values...)
			}
		}
	}
	buf.WriteRune(')')
	buf.WriteString(alies)
	return buf.String(), args
}

//-------------- 关于Where条件的补充
// In
func (this *selectBuilder) WhereIn(col string, args ...interface{}) *selectBuilder {
	if len(args) < 1 {
		this.err = errors.New("WhereIn 条件缺失")
		return this
	}
	this.where.wherein(col, args)
	return this
}

//// 查询条件in的解析
//func (this *selectBuilder) wherein(col string, args interface{}) *selectBuilder {
//	var v = reflect.Indirect(reflect.ValueOf(args))
//	var k = v.Kind()
//	if k == reflect.Slice || k == reflect.Array {
//		if v.Len() == 0 {
//			return this
//		}
//		var buf = bufPool.Get()
//		defer bufPool.Put(buf)
//		var where = new(whereConstraint)
//		buf.WriteString(fmt.Sprintf("%s IN(", col))
//		for i := 0; i < v.Len(); i++ {
//			if i > 0 {
//				buf.WriteString(" ,?")
//			} else {
//				buf.WriteRune('?')
//			}
//		}
//		buf.WriteRune(')')
//		where.condition = buf.String()
//		switch v.Index(0).Elem().Kind() {
//		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
//			for i := 0; i < v.Len(); i++ {
//				where.values = append(where.values, v.Index(i).Elem().Int())
//			}
//		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
//			for i := 0; i < v.Len(); i++ {
//				where.values = append(where.values, v.Index(i).Elem().Uint())
//			}

//		case reflect.Float32, reflect.Float64:
//			for i := 0; i < v.Len(); i++ {
//				where.values = append(where.values, v.Index(i).Elem().Float())
//			}
//		case reflect.Bool:
//			for i := 0; i < v.Len(); i++ {
//				where.values = append(where.values, v.Index(i).Bool())
//			}
//		case reflect.String:
//			for i := 0; i < v.Len(); i++ {
//				where.values = append(where.values, v.Index(i).Elem().String())
//			}
//		default:
//			this.err = errors.New(fmt.Sprintf("in不支持的类型%s", v.Index(0).Elem().Kind().String()))
//		}

//		this.where.where = append(this.where.where, where)
//	} else {
//		this.err = errors.New("参数格式错误，必须为切片或数组")
//	}
//	return this
//}

// In
func (this *selectBuilder) havingIn(col string, args ...interface{}) *selectBuilder {
	return this.havingin(col, args)
}

// 查询条件in的解析
func (this *selectBuilder) havingin(col string, args interface{}) *selectBuilder {
	var v = reflect.Indirect(reflect.ValueOf(args))
	var k = v.Kind()
	if k == reflect.Slice || k == reflect.Array {
		if v.Len() == 0 {
			return this
		}
		var buf = bufPool.Get()
		defer bufPool.Put(buf)
		var where = new(whereConstraint)
		buf.WriteString(fmt.Sprintf("%s in(", col))
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
		default:
			this.err = errors.New(fmt.Sprintf("in不支持的类型%s", v.Index(0).Elem().Kind().String()))
		}

		this.having = append(this.having, where)
	} else {
		this.err = errors.New("参数格式错误，必须为切片或数组")
	}
	return this
}

// 全匹配
func (this *selectBuilder) WhereLike(col string, arg interface{}) *selectBuilder {
	this.where.whereLike(col, arg, 0)
	return this
}

// 左匹配
func (this *selectBuilder) WhereLikeL(col string, arg interface{}) *selectBuilder {
	this.where.whereLike(col, arg, -1)
	return this
}

// 右匹配
func (this *selectBuilder) WhereLikeR(col string, arg interface{}) *selectBuilder {
	this.where.whereLike(col, arg, 1)
	return this
}
