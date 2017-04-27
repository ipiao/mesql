package meorm

// // InsertModels 插入结构体或结构体数组
// func (c *Conn) InsertModels(models interface{}) *medb.Result {
// 	var r = new(medb.Result)
// 	var value = reflect.Indirect(reflect.ValueOf(models))
// 	var k = value.Kind()
// 	switch k {
// 	case reflect.Struct:
// 		return c.insertStruct(&value)
// 	case reflect.Slice, reflect.Array:
// 		if value.Type().Elem().Kind() == reflect.Struct {
// 			return c.insertSlice(&value)
// 		}
// 		r.SetErr(fmt.Errorf("Error kind of models []%s", value.Type().Elem().Kind().String()))
// 	default:
// 		r.SetErr(fmt.Errorf("Error kind of models %s", k.String()))
// 	}
// 	return r
// }

// // 插入结构体
// // mysql 的主键遇 0 值可以自动忽略
// func (c *Conn) insertStruct(v *reflect.Value) *medb.Result {
// 	buf := bufPool.Get()
// 	defer bufPool.Put(buf)

// 	var tbName = GetTableName(*v)
// 	var cols = GetColumns(*v)
// 	var values = GetValues(*v)

// 	if tbName == "" || len(cols) == 0 {
// 		var res = new(medb.Result)
// 		res.SetErr(errors.New("Error struct"))
// 		return res
// 	}

// 	buf.WriteString("INSERT INTO ")
// 	buf.WriteString(tbName)
// 	buf.WriteString(" (")

// 	var valueStr = "("
// 	for i, col := range cols {
// 		if i > 0 {
// 			buf.WriteString(" ,")
// 			valueStr += " ,"
// 		}
// 		buf.WriteString(col)
// 		valueStr += "?"
// 	}
// 	buf.WriteString(") VALUES ")
// 	valueStr += ")"
// 	buf.WriteString(valueStr)
// 	var args = values[0]
// 	return c.Exec(buf.String(), args...)
// }

// // 插入数组
// func (c *Conn) insertSlice(v *reflect.Value) *medb.Result {
// 	buf := bufPool.Get()
// 	defer bufPool.Put(buf)

// 	var tbName = GetTableName(*v)
// 	var cols = GetColumns(*v)
// 	var values = GetValues(*v)

// 	if tbName == "" || len(cols) == 0 {
// 		var res = new(medb.Result)
// 		res.SetErr(errors.New("Error slice"))
// 		return res
// 	}

// 	buf.WriteString("INSERT INTO ")
// 	buf.WriteString(tbName)
// 	buf.WriteString(" (")

// 	var valueStr = "("
// 	for i, col := range cols {
// 		if i > 0 {
// 			buf.WriteString(" ,")
// 			valueStr += " ,"
// 		}
// 		buf.WriteString(col)
// 		valueStr += "?"
// 	}
// 	buf.WriteString(") VALUES ")
// 	valueStr += "),"

// 	valueStr = strings.Repeat(valueStr, len(values))
// 	buf.WriteString(valueStr[:len(valueStr)-1])
// 	var args []interface{}
// 	for i := range values {
// 		args = append(args, values[i]...)
// 	}
// 	return c.Exec(buf.String(), args...)
// }
