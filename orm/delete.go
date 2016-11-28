package meorm

import (
	"errors"

	"github.com/ipiao/mesql/medb"
)

type deleteBuilder struct {
	Executor
	connname   string
	sql        string
	args       []interface{}
	err        error
	table      string
	where      *Where
	limit      int64
	limitvalid bool
	orderbys   []string
}

// reset
func (this *deleteBuilder) reset() *deleteBuilder {
	this.table = ""
	this.where = new(Where)
	this.orderbys = this.orderbys[:0]
	this.limit = 0
	this.limitvalid = false
	this.err = nil
	this.sql = ""
	this.args = this.args[:0]
	return this
}

// where 条件
func (this *deleteBuilder) Where(condition string, args ...interface{}) *deleteBuilder {
	this.where.where = append(this.where.where, &whereConstraint{
		condition: condition,
		values:    args,
	})
	return this
}

// orderby 条件
func (this *deleteBuilder) OrderBy(order string) *deleteBuilder {
	this.orderbys = append(this.orderbys, order)
	return this
}

// limit
func (this *deleteBuilder) Limit(limit int64) *deleteBuilder {
	this.limitvalid = true
	this.limit = limit
	return this
}

// 生成sql
func (this *deleteBuilder) tosql() (string, []interface{}) {
	mutex.Lock()
	defer mutex.Unlock()
	if this.where.err != nil {
		this.err = this.where.err
		return "", nil
	}
	buf := bufPool.Get()
	defer bufPool.Put(buf)

	var args []interface{}
	buf.WriteString("DELETE FROM ")
	buf.WriteString(this.table)

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
func (this *deleteBuilder) ToSQL() (string, []interface{}) {
	if len(this.sql) > 0 {
		return this.sql, this.args
	}
	return this.tosql()
}

// 执行
func (this *deleteBuilder) Exec() *medb.Result {
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
func (this *deleteBuilder) QueryTo(models interface{}) (int, error) {
	//	if len(this.sql) == 0 {
	//		this.tosql()
	//	}
	//	if this.err != nil {
	//		return 0, this.err
	//	}
	//	return connections[this.connname].db.Query(this.sql, this.args...).ScanTo(models)
	return 0, errors.New("[meorm]:Delete 不能使用该方法")
}

// 把查询组成sql并解析
func (this *deleteBuilder) QueryNext(dest ...interface{}) error {
	//	if len(this.sql) == 0 {
	//		this.tosql()
	//	}
	//	if this.err != nil {
	//		return this.err
	//	}
	//	return connections[this.connname].db.Query(this.sql, this.args...).ScanNext(dest...)
	return errors.New("[meorm]:Delete 不能使用该方法")
}
