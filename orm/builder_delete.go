package meorm

import (
	"github.com/ipiao/mesql/medb"
)

// DeleteBuilder 删除
type DeleteBuilder struct {
	connname   string
	sql        string
	column     string
	args       []interface{}
	err        error
	table      string
	where      *Where
	limit      int64
	limitvalid bool
	orderbys   []string
}

func (d *DeleteBuilder) reset() *DeleteBuilder {
	d.table = ""
	d.column = ""
	d.where = new(Where)
	d.orderbys = d.orderbys[:0]
	d.limit = 0
	d.limitvalid = false
	d.err = nil
	d.sql = ""
	d.args = d.args[:0]
	return d
}

// From from
func (d *DeleteBuilder) From(table string) *DeleteBuilder {
	d.table = table
	return d
}

// Where 条件
func (d *DeleteBuilder) Where(condition string, args ...interface{}) *DeleteBuilder {
	d.where.where = append(d.where.where, &whereConstraint{
		condition: condition,
		values:    args,
	})
	return d
}

// WhereIn in条件
func (d *DeleteBuilder) WhereIn(col string, args ...interface{}) *DeleteBuilder {
	d.where.wherein(col, args...)
	return d
}

// OrderBy 条件
func (d *DeleteBuilder) OrderBy(order string) *DeleteBuilder {
	d.orderbys = append(d.orderbys, order)
	return d
}

// Limit 条件
func (d *DeleteBuilder) Limit(limit int64) *DeleteBuilder {
	d.limitvalid = true
	d.limit = limit
	return d
}

// 生成sql
func (d *DeleteBuilder) tosql() (string, []interface{}) {
	// mutex.Lock()
	// defer mutex.Unlock()
	if d.where.err != nil {
		d.err = d.where.err
		return "", nil
	}
	buf := bufPool.Get()
	defer bufPool.Put(buf)

	var args []interface{}
	buf.WriteString("DELETE ")
	buf.WriteString(d.column)
	buf.WriteString(" FROM ")
	buf.WriteString(d.table)

	if len(d.where.where) > 0 {
		buf.WriteString(" WHERE ")
		for i, cond := range d.where.where {
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

	if len(d.orderbys) > 0 {
		buf.WriteString(" ORDER BY ")
		for i, s := range d.orderbys {
			if i > 0 {
				buf.WriteString(", ")
			}
			buf.WriteString(s)
		}
	}
	if d.limitvalid {
		buf.WriteString(" LIMIT ?")
		args = append(args, d.limit)
	}

	d.sql = buf.String()
	d.args = args
	return d.sql, d.args
}

// ToSQL 对外
func (d *DeleteBuilder) ToSQL() (string, []interface{}) {
	if len(d.sql) > 0 {
		return d.sql, d.args
	}
	return d.tosql()
}

// Exec 执行
func (d *DeleteBuilder) Exec() *medb.Result {
	var res = new(medb.Result)
	if len(d.sql) == 0 {
		d.tosql()
	}
	if d.err != nil {
		res.SetErr(d.err)
		return res
	}
	return connections[d.connname].Exec(d.sql, d.args...)
}
