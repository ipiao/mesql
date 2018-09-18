package meorm

// 按结构体更新
func UpdateModels(models interface{}) {
	//	var value = reflect.Indirect(reflect.ValueOf(models))
	//	var k = value.Kind()
	//	switch k {
	//	case reflect.Struct:
	//		return this.updateStruct(&value)
	//		//	case reflect.Slice, reflect.Array:
	//		//		return this.insertSlice(&value)
	//	default:
	//		panic(fmt.Sprintf("Error kind %s", k.String()))
	//	}
}

//  TODO 主键和唯一索引的排除问题

//func (this *Conn) updateStruct(v *reflect.Value) *medb.Result {
//	buf := bufPool.Get()
//	defer bufPool.Put(buf)

//	var tbName = GetTableName(*v)
//	var cols = GetColumns(*v)
//	var vals = GetValues(*v)

//	if tbName == "" || len(cols) == 0 {
//		return &medb.Result{
//			Err: errors.New("Error struct"),
//		}
//	}

//	buf.WriteString("UPDATE ")
//	buf.WriteString(tbName)
//	buf.WriteString(" SET ")

//	for i, col := range cols {
//		if col
//		if i > 0 {
//			buf.WriteString(" ,")
//		}
//		buf.WriteString(col)
//		buf.WriteString(" = ? ")
//	}

//	return this.db.Exec(buf.String(), args...)
//}
