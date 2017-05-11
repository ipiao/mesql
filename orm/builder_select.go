package meorm

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/ipiao/mesql/medb"
)

// SelectBuilder 查询
type SelectBuilder struct {
	builder  *Builder
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

// Select 选择查询列
func (s *SelectBuilder) Select(columns ...string) *SelectBuilder {
	s.columns = append(s.columns, columns...)
	return s
}

// Distinct distinct
func (s *SelectBuilder) Distinct() *SelectBuilder {
	s.distinct = true
	return s
}

// From from
func (s *SelectBuilder) From(from string) *SelectBuilder {
	s.from = from
	return s
}

// OrderBy order by
func (s *SelectBuilder) OrderBy(order string) *SelectBuilder {
	s.orderbys = append(s.orderbys, order)
	return s
}

// GroupBy group by
func (s *SelectBuilder) GroupBy(group string) *SelectBuilder {
	s.groupbys = append(s.groupbys, group)
	return s
}

// Limit limit
func (s *SelectBuilder) Limit(limit int64) *SelectBuilder {
	s.limitvalid = true
	s.limit = limit
	return s
}

// Offset offset
func (s *SelectBuilder) Offset(offset int64) *SelectBuilder {
	s.offsetvalid = true
	s.offset = offset
	return s
}

// Where where
func (s *SelectBuilder) Where(condition string, args ...interface{}) *SelectBuilder {
	s.where.where = append(s.where.where, &whereConstraint{
		condition: condition,
		values:    args,
	})
	return s
}

// Having having
func (s *SelectBuilder) Having(condition string, values ...interface{}) *SelectBuilder {
	s.having = append(s.having, &whereConstraint{
		condition: condition,
		values:    values,
	})
	return s
}

// reset
func (s *SelectBuilder) reset() *SelectBuilder {
	s.distinct = false
	s.columns = s.columns[:0]
	s.from = ""
	s.where = new(Where)
	s.orderbys = s.orderbys[:0]
	s.groupbys = s.groupbys[:0]
	s.limit = 0
	s.limitvalid = false
	s.offset = 0
	s.offsetvalid = false
	s.having = make([]*whereConstraint, 0, 0)
	s.err = nil
	s.sql = ""
	s.args = s.args[:0]
	return s
}

// ToSQL tosql
func (s *SelectBuilder) ToSQL() (string, []interface{}) {
	if len(s.sql) > 0 {
		return s.sql, s.args
	}
	return s.tosql()
}

// 把查询条件组成sql并放到查询体中
func (s *SelectBuilder) tosql() (string, []interface{}) {
	// mutex.Lock()
	// defer mutex.Unlock()
	if s.where.err != nil {
		s.err = s.where.err
		return "", nil
	}
	if len(s.columns) == 0 {
		panic("没有指定列")
	}
	if len(s.from) == 0 {
		panic("没有指定表")
	}
	buf := bufPool.Get()
	defer bufPool.Put(buf)

	var args []interface{}
	buf.WriteString("SELECT ")

	if s.distinct {
		buf.WriteString("DISTINCT ")
	}
	for i, s := range s.columns {
		if i > 0 {
			buf.WriteString(", ")
		}
		buf.WriteString(s)
	}
	buf.WriteString(" FROM ")
	buf.WriteString(s.from)

	if len(s.where.where) > 0 {
		buf.WriteString(" WHERE ")
		for i, cond := range s.where.where {
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
	if len(s.groupbys) > 0 {
		buf.WriteString(" GROUP BY ")
		for i, s := range s.groupbys {
			if i > 0 {
				buf.WriteString(", ")
			}
			buf.WriteString(s)
		}
	}
	if len(s.having) > 0 {
		buf.WriteString(" HAVING ")
		for i, cond := range s.having {
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
	if len(s.orderbys) > 0 {
		buf.WriteString(" ORDER BY ")
		for i, s := range s.orderbys {
			if i > 0 {
				buf.WriteString(", ")
			}
			buf.WriteString(s)
		}
	}
	if s.limitvalid {
		buf.WriteString(" LIMIT ?")
		args = append(args, s.limit)
	}
	if s.offsetvalid {
		buf.WriteString(" OFFSET ?")
		args = append(args, s.offset)
	}
	s.sql = buf.String()
	s.args = args
	return s.sql, s.args
}

// Exec 查询不建议使用
func (s *SelectBuilder) Exec() *medb.Result {
	var res = new(medb.Result)
	if len(s.sql) == 0 {
		s.tosql()
	}
	if s.err != nil {
		res.SetErr(s.err)
		return res
	}
	return s.builder.Exec(s.sql, s.args...)
}

// QueryTo 解析到结构体，数组。。。
func (s *SelectBuilder) QueryTo(models interface{}) (int, error) {
	if len(s.sql) == 0 {
		s.tosql()
	}
	if s.err != nil {
		return 0, s.err
	}
	return s.builder.Query(s.sql, s.args...).ScanTo(models)
}

// QueryNext 把查询组成sql并解析
func (s *SelectBuilder) QueryNext(dest ...interface{}) error {
	if len(s.sql) == 0 {
		s.tosql()
	}
	if s.err != nil {
		return s.err
	}
	return s.builder.Query(s.sql, s.args...).ScanNext(dest...)
}

// LimitPP limit和offset的复用
func (s *SelectBuilder) LimitPP(page, pagesize int64) *SelectBuilder {
	var offset = (page - 1) * pagesize
	return s.Limit(pagesize).Offset(offset)
}

// CountCond 查询符合条件的总数目
func (s *SelectBuilder) CountCond(countCond ...string) (int64, string, error) {
	var countcond string
	var count int64
	if len(countCond) == 0 {
		countcond = "count(0)"
	} else {
		countcond = countCond[0]
	}
	var sql, args = s.countsql(countcond)
	var err = s.builder.Query(sql, args...).ScanNext(&count)
	return count, sql, err
}

// CountResult 查询符合条件的总数目
func (s *SelectBuilder) CountResult(alies string, countCond ...string) (int64, string, error) {
	var countcond string
	var count int64
	if len(countCond) == 0 {
		countcond = "count(0)"
	} else {
		countcond = countCond[0]
	}
	var sql, args = s.countresult(alies, countcond)
	var err = s.builder.Query(sql, args...).ScanNext(&count)
	return count, sql, err
}

// 生成查询总条数的sql
func (s *SelectBuilder) countsql(countCond string) (string, []interface{}) {
	buf := bufPool.Get()
	defer bufPool.Put(buf)

	var args []interface{}
	buf.WriteString("SELECT ")
	buf.WriteString(countCond)
	buf.WriteString(" FROM ")
	buf.WriteString(s.from)

	if len(s.where.where) > 0 {
		buf.WriteString(" WHERE ")
		for i, cond := range s.where.where {
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
	if len(s.groupbys) > 0 {
		buf.WriteString(" GROUP BY ")
		for i, s := range s.groupbys {
			if i > 0 {
				buf.WriteString(", ")
			}
			buf.WriteString(s)
		}
	}
	if len(s.having) > 0 {
		buf.WriteString(" HAVING ")
		for i, cond := range s.having {
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
func (s *SelectBuilder) countresult(alies, countCond string) (string, []interface{}) {
	buf := bufPool.Get()
	defer bufPool.Put(buf)
	var args []interface{}
	if s.where.err != nil {
		s.err = s.where.err
		return "", nil
	}
	buf.WriteString("SELECT ")
	buf.WriteString(countCond)
	buf.WriteString(" FROM( ")
	buf.WriteString("SELECT ")

	if s.distinct {
		buf.WriteString("DISTINCT ")
	}
	for i, s := range s.columns {
		if i > 0 {
			buf.WriteString(", ")
		}
		buf.WriteString(s)
	}
	buf.WriteString(" FROM ")
	buf.WriteString(s.from)

	if len(s.where.where) > 0 {
		buf.WriteString(" WHERE ")
		for i, cond := range s.where.where {
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
	if len(s.groupbys) > 0 {
		buf.WriteString(" GROUP BY ")
		for i, s := range s.groupbys {
			if i > 0 {
				buf.WriteString(", ")
			}
			buf.WriteString(s)
		}
	}
	if len(s.having) > 0 {
		buf.WriteString(" HAVING ")
		for i, cond := range s.having {
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

// WhereIn In
func (s *SelectBuilder) WhereIn(col string, args ...interface{}) *SelectBuilder {
	s.where.wherein(col, args...)
	return s
}

// WhereNotIn In
func (s *SelectBuilder) WhereNotIn(col string, args ...interface{}) *SelectBuilder {
	s.where.wherenotin(col, args...)
	return s
}

// havingIn In
func (s *SelectBuilder) havingIn(col string, args ...interface{}) *SelectBuilder {
	return s.havingin(col, args)
}

// havingin 查询条件in的解析
func (s *SelectBuilder) havingin(col string, args interface{}) *SelectBuilder {
	var v = reflect.Indirect(reflect.ValueOf(args))
	var k = v.Kind()
	if k == reflect.Slice || k == reflect.Array {
		if v.Len() == 0 {
			return s
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
			s.err = fmt.Errorf("in不支持的类型%s", v.Index(0).Elem().Kind().String())
		}

		s.having = append(s.having, where)
	} else {
		s.err = errors.New("参数格式错误，必须为切片或数组")
	}
	return s
}

// WhereLike 全匹配
func (s *SelectBuilder) WhereLike(col string, arg interface{}) *SelectBuilder {
	s.where.whereLike(col, arg, 0)
	return s
}

// WhereLikeL 左匹配
func (s *SelectBuilder) WhereLikeL(col string, arg interface{}) *SelectBuilder {
	s.where.whereLike(col, arg, -1)
	return s
}

// WhereLikeR 右匹配
func (s *SelectBuilder) WhereLikeR(col string, arg interface{}) *SelectBuilder {
	s.where.whereLike(col, arg, 1)
	return s
}
