package meorm

import (
	"errors"

	"github.com/ipiao/mesql/medb"
)

// 更新构造器
type updateBuilder struct {
	Executor
	connname   string
	table      string
	setClause  []*setClause
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
func (this *updateBuilder) reset() *updateBuilder {
	this.table = ""
	this.setClause = make([]*setClause, 0, 0)
	this.where = new(Where)
	this.orderbys = this.orderbys[:0]
	this.limit = 0
	this.limitvalid = false
	this.err = nil
	this.sql = ""
	this.args = this.args[:0]
	return this
}

// 设置值
func (this *updateBuilder) Set(column string, value interface{}) *updateBuilder {
	this.setClause = append(this.setClause, &setClause{
		column: column,
		value:  value,
	})
	return this
}

// where 条件
func (this *updateBuilder) Where(condition string, args ...interface{}) *updateBuilder {
	this.where.where = append(this.where.where, &whereConstraint{
		condition: condition,
		values:    args,
	})
	return this
}

// WhereIn 条件
func (this *updateBuilder) WhereIn(col string, args ...interface{}) *updateBuilder {
	this.where.wherein(col, args...)
	return this
}

// orderby 条件
func (this *updateBuilder) OrderBy(order string) *updateBuilder {
	this.orderbys = append(this.orderbys, order)
	return this
}

// limit
func (this *updateBuilder) Limit(limit int64) *updateBuilder {
	this.limitvalid = true
	this.limit = limit
	return this
}

// 生成sql
func (this *updateBuilder) tosql() (string, []interface{}) {
	mutex.Lock()
	defer mutex.Unlock()
	if this.where.err != nil {
		this.err = this.where.err
		return "", nil
	}
	buf := bufPool.Get()
	defer bufPool.Put(buf)

	var args []interface{}
	buf.WriteString("UPDATE ")
	buf.WriteString(this.table)
	buf.WriteString(" SET ")
	for i, s := range this.setClause {
		if i > 0 {
			buf.WriteString(", ")
		}
		buf.WriteString(s.column + " = ?")
		args = append(args, s.value)
	}

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

	this.sql = buf.String()
	this.args = args
	return this.sql, this.args
}

// tosql
func (this *updateBuilder) ToSQL() (string, []interface{}) {
	if len(this.sql) > 0 {
		return this.sql, this.args
	}
	return this.tosql()
}

// 执行
func (this *updateBuilder) Exec() *medb.Result {
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
func (this *updateBuilder) QueryTo(models interface{}) (int, error) {
	//	if len(this.sql) == 0 {
	//		this.tosql()
	//	}
	//	if this.err != nil {
	//		return 0, this.err
	//	}
	//	return connections[this.connname].db.Query(this.sql, this.args...).ScanTo(models)
	return 0, errors.New("[meorm]:Update 不能使用该方法")
}

// 把查询组成sql并解析
func (this *updateBuilder) QueryNext(dest ...interface{}) error {
	//	if len(this.sql) == 0 {
	//		this.tosql()
	//	}
	//	if this.err != nil {
	//		return this.err
	//	}
	//	return connections[this.connname].db.Query(this.sql, this.args...).ScanNext(dest...)
	return errors.New("[meorm]:Update 不能使用该方法")
}
