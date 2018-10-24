package meorm

import (
	"github.com/ipiao/mesql/medb"
)

// DeleteBuilder 删除
type DeleteBuilder struct {
	*where
	builder    *BaseBuilder
	sql        string
	column     string
	args       []interface{}
	err        error
	table      string
	limit      int64
	limitvalid bool
	orderbys   []string
}

func (d *DeleteBuilder) reset() *DeleteBuilder {
	d.table = ""
	d.column = ""
	d.where = new(where)
	d.where.dialect = d.dialect
	d.orderbys = d.orderbys[:0]
	d.limit = 0
	d.limitvalid = false
	d.err = nil
	d.sql = ""
	d.args = d.args[:0]
	return d
}

// From from
// delete ... join 语法
// delete T1 from T1 join T2
func (d *DeleteBuilder) From(table string) *DeleteBuilder {
	d.table = table
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

	if d.where.err != nil {
		d.err = d.where.err
		return "", nil
	}

	holder := d.builder.dialect.Holder()
	buf := bufPool.Get()
	defer bufPool.Put(buf)

	var args []interface{}
	buf.WriteString("DELETE ")
	if len(d.column) > 0 {
		buf.WriteString(d.column)
		buf.WriteByte(' ')
	}
	buf.WriteString("FROM ")
	buf.WriteString(d.table)

	if len(d.where.conds) > 0 {
		buf.WriteString(" WHERE ")
		for i, cond := range d.where.conds {
			if i > 0 {
				buf.WriteString(" AND (")
			} else {
				buf.WriteByte('(')
			}
			buf.WriteString(cond.condition)
			buf.WriteByte(')')
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
		buf.WriteString(" LIMIT ")
		buf.WriteByte(holder)
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
	if d.err != nil {
		res.SetErr(d.err)
		return res
	}
	return d.builder.Exec(d)
}

// PrepareExec 预处理后执行
func (d *DeleteBuilder) PrepareExec() *medb.Result {
	var res = new(medb.Result)
	if d.err != nil {
		res.SetErr(d.err)
		return res
	}
	return d.builder.Prepare(d.sql).Exec(d.args...)
}
