package meorm

import "github.com/ipiao/mesql/medb"

// InsupBuilder insert or update构造器
type InsupBuilder struct {
	builder    *BaseBuilder
	table      string
	inscolumns []string
	upcolumns  []string
	insvalues  []interface{}
	upvalues   []interface{}
	sql        string
	args       []interface{}
	err        error
	dupkeys    []string
	dup        bool
}

// reset
func (b *InsupBuilder) reset() *InsupBuilder {
	b.table = ""
	b.inscolumns = b.inscolumns[:0]
	b.upcolumns = b.upcolumns[:0]
	b.dupkeys = b.dupkeys[:0]
	b.insvalues = make([]interface{}, 0)
	b.err = nil
	b.sql = ""
	b.args = b.args[:0]
	return b
}

// Columns 列
func (b *InsupBuilder) Columns(cols ...string) *InsupBuilder {
	if b.dup {
		b.upcolumns = append(b.upcolumns, cols...)
	} else {
		b.inscolumns = append(b.inscolumns, cols...)
	}
	return b
}

// Values 值
func (b *InsupBuilder) Values(values ...interface{}) *InsupBuilder {
	if b.dup {
		b.upvalues = append(b.upvalues, values...)
	} else {
		b.insvalues = append(b.insvalues, values...)
	}
	return b
}

// DupKeys 冲突列
func (b *InsupBuilder) DupKeys(cols ...string) *InsupBuilder {
	b.dup = true
	b.dupkeys = append(b.dupkeys, cols...)
	return b
}

// tosql
func (b *InsupBuilder) tosql() (string, []interface{}) {
	buf := bufPool.Get()
	defer bufPool.Put(buf)

	buf.WriteString("INSERT INTO ")
	buf.WriteString(b.table)
	buf.WriteString(" (")

	for i, col := range b.inscolumns {
		if i > 0 {
			buf.WriteString(" ,")
		}
		buf.WriteString(col)
	}
	buf.WriteString(") VALUES")
	var args []interface{}
	buf.WriteString(" (")
	for i, value := range b.insvalues {
		if i > 0 {
			buf.WriteString(" ,?")
		} else {
			buf.WriteString("?")
		}
		args = append(args, value)
	}
	buf.WriteString(")")

	// 默认更新字段
	if len(b.upcolumns) == 0 {
		b.upcolumns = b.inscolumns
	}
	if len(b.upvalues) == 0 {
		b.upvalues = b.insvalues
	}
	// 去除更新字段中的冲突字段
	for _, dup := range b.dupkeys {
		for i := range b.upcolumns {
			if b.upcolumns[i] == dup {
				b.upcolumns = append(b.upcolumns[:i], b.upcolumns[i+1:]...)
				b.insvalues = append(b.insvalues[:i], b.insvalues[i+1:]...)
				break
			}
		}
	}
	buf.WriteString(" ON DUPLICATE KEY UPDATE ")
	for i := range b.upcolumns {
		if i > 0 {
			buf.WriteString(" ,")
		}
		buf.WriteString(b.upcolumns[i] + "=?")
		args = append(args, b.upvalues[i])
	}

	b.sql = buf.String()
	b.args = args
	return b.sql, b.args
}

// ToSQL tosql
func (b *InsupBuilder) ToSQL() (string, []interface{}) {
	if len(b.sql) > 0 {
		return b.sql, b.args
	}
	return b.tosql()
}

// Exec 执行
func (b *InsupBuilder) Exec() *medb.Result {
	var res = new(medb.Result)
	if b.err != nil {
		res.SetErr(b.err)
		return res
	}
	return b.builder.Exec(b)
}
