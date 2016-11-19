package meorm

import "ipiao/mesql/medb"

// 查询
type selectBuilder struct {
	Executor
	connname    string
	distinct    bool
	columns     []string
	from        string
	where       []*Where
	orderbys    []string
	groupbys    []string
	limit       int
	limitvalid  bool
	offset      int
	offsetvalid bool
	having      []*Where
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

// limit
func (this *selectBuilder) Limit(limit int) *selectBuilder {
	this.limitvalid = true
	this.limit = limit
	return this
}

// offset
func (this *selectBuilder) Offset(offset int) *selectBuilder {
	this.offsetvalid = true
	this.offset = offset
	return this
}

// where
func (this *selectBuilder) Where(condition string, values ...interface{}) *selectBuilder {
	this.where = append(this.where, &Where{
		condition: condition,
		values:    values,
	})
	return this
}

// having
func (this *selectBuilder) Having(condition string, values ...interface{}) *selectBuilder {
	this.having = append(this.having, &Where{
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
	this.where = make([]*Where, 0, 0)
	this.orderbys = this.orderbys[:0]
	this.groupbys = this.groupbys[:0]
	this.limit = 0
	this.limitvalid = false
	this.offset = 0
	this.offsetvalid = false
	this.having = make([]*Where, 0, 0)
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

//
func (this *selectBuilder) tosql() (string, []interface{}) {
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

	if len(this.where) > 0 {
		buf.WriteString(" WHERE ")
		for i, cond := range this.where {
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

//
func (this *selectBuilder) Exec() *medb.Result {
	if len(this.sql) == 0 {
		this.tosql()
	}
	return connections[this.connname].db.Exec(this.sql, this.args...)
}

//
func (this *selectBuilder) QueryTo(models interface{}) (int, error) {
	if len(this.sql) == 0 {
		this.tosql()
	}
	return connections[this.connname].db.Query(this.sql, this.args...).ScanTo(models)
}

//
func (this *selectBuilder) QueryNext(dest ...interface{}) error {
	if len(this.sql) == 0 {
		this.tosql()
	}
	return connections[this.connname].db.Query(this.sql, this.args...).ScanNext(dest...)
}
