package meorm

import (
	"github.com/ipiao/mesql/medb"
)

// UpdateBuilder 更新构造器
// 只支持单个更新
type UpdateBuilder struct {
	builder    *Builder
	connname   string
	table      string
	columns    []string
	values     []interface{}
	where      *Where
	orderbys   []string
	limit      int64
	limitvalid bool
	err        error
	sql        string
	args       []interface{}
}

// set项
type setClause struct {
	column string
	value  interface{}
}

// reset
func (u *UpdateBuilder) reset() *UpdateBuilder {
	u.table = ""
	u.columns = u.columns[:0]
	u.values = u.values[:0]
	//u.setClause = make([]*setClause, 0, 0)
	u.where = new(Where)
	u.orderbys = u.orderbys[:0]
	u.limit = 0
	u.limitvalid = false
	u.err = nil
	u.sql = ""
	u.args = u.args[:0]
	return u
}

// Set 设置值
func (u *UpdateBuilder) Set(column string, value interface{}) *UpdateBuilder {
	u.columns = append(u.columns, column+"=?")
	u.values = append(u.values, value)
	return u
}

// SetS 设置值
func (u *UpdateBuilder) SetS(column string, values ...interface{}) *UpdateBuilder {
	u.columns = append(u.columns, column)
	u.values = append(u.values, values...)
	return u
}

// Where where 条件
func (u *UpdateBuilder) Where(condition string, args ...interface{}) *UpdateBuilder {
	u.where.where = append(u.where.where, &whereConstraint{
		condition: condition,
		values:    args,
	})
	return u
}

// WhereIn 条件
func (u *UpdateBuilder) WhereIn(col string, args ...interface{}) *UpdateBuilder {
	u.where.wherein(col, args...)
	return u
}

// OrderBy orderby 条件
func (u *UpdateBuilder) OrderBy(order string) *UpdateBuilder {
	u.orderbys = append(u.orderbys, order)
	return u
}

// Limit limit
func (u *UpdateBuilder) Limit(limit int64) *UpdateBuilder {
	u.limitvalid = true
	u.limit = limit
	return u
}

// 生成sql
func (u *UpdateBuilder) tosql() (string, []interface{}) {
	// mutex.Lock()
	// defer mutex.Unlock()
	if u.where.err != nil {
		u.err = u.where.err
		return "", nil
	}
	buf := bufPool.Get()
	defer bufPool.Put(buf)

	var args []interface{}
	buf.WriteString("UPDATE ")
	buf.WriteString(u.table)
	buf.WriteString(" SET ")
	for i, s := range u.columns {
		if i > 0 {
			buf.WriteString(", ")
		}
		buf.WriteString(s)
	}
	args = append(args, u.values...)

	if len(u.where.where) > 0 {
		buf.WriteString(" WHERE ")
		for i, cond := range u.where.where {
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

	if len(u.orderbys) > 0 {
		buf.WriteString(" ORDER BY ")
		for i, s := range u.orderbys {
			if i > 0 {
				buf.WriteString(", ")
			}
			buf.WriteString(s)
		}
	}
	if u.limitvalid {
		buf.WriteString(" LIMIT ?")
		args = append(args, u.limit)
	}

	u.sql = buf.String()
	u.args = args
	return u.sql, u.args
}

// ToSQL tosql
func (u *UpdateBuilder) ToSQL() (string, []interface{}) {
	if len(u.sql) > 0 {
		return u.sql, u.args
	}
	return u.tosql()
}

// Exec 执行
func (u *UpdateBuilder) Exec() *medb.Result {
	if len(u.sql) == 0 {
		u.tosql()
	}
	if u.err != nil {
		var res = new(medb.Result).SetErr(u.err)
		return res
	}
	return u.builder.Exec(u.sql, u.args...)
}
