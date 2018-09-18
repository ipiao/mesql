package meorm

// 默认构造，只构建sql

var defaultBuilder = &BaseBuilder{}

// SQL 直接写sql
func SQL(sql string, args ...interface{}) *BareBuilder {
	return defaultBuilder.SQL(sql, args...)
}

// Select 生成查询构造器
func Select(cols ...string) *SelectBuilder {
	return defaultBuilder.Select(cols...)
}

// Update 生成更新构造器
func Update(table string) *UpdateBuilder {
	return defaultBuilder.Update(table)
}

// InsertOrUpdate 生成插入或更新构造器
func InsertOrUpdate(table string) *InsupBuilder {
	return defaultBuilder.InsertOrUpdate(table)
}

// InsertInto 生成插入构造器
func InsertInto(table string) *InsertBuilder {
	return defaultBuilder.InsertInto(table)
}

// ReplaceInto 生成插入构造器
func ReplaceInto(table string) *InsertBuilder {
	return defaultBuilder.ReplaceInto(table)
}

// DeleteFrom 生成删除构造器
func DeleteFrom(table string) *DeleteBuilder {
	return defaultBuilder.DeleteFrom(table)
}

// Delete 生成删除构造器
func Delete(column string) *DeleteBuilder {
	return defaultBuilder.Delete(column)
}
