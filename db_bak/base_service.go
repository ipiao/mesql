package mesql

import (
	"database/sql"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"
)

type BaseService struct {
	DB        *DB
	condition [][2]string
	values    []interface{}
}

func NewConnDB(name string) *BaseService {
	return &BaseService{
		DB: OpenDB(name),
	}
}

// Query 解析查询条件，结果赋予models
func (this *BaseService) QueryTo(models interface{}) (int, error) {
	var n, err = this.DB.Query(this.GetCondition(), this.GetValues()...).ScanTo(models)
	this.Reset()
	return n, err
}

// Query 解析查询条件，结果赋予models
func (this *BaseService) QueryNext(res ...interface{}) error {
	var err = this.DB.Query(this.GetCondition(), this.GetValues()...).ScanNext(res)
	this.Reset()
	return err
}

// Query 解析查询条件，结果赋予models
func (this *BaseService) Exec() (sql.Result, error) {
	var r, err = this.DB.Exec(this.GetCondition(), this.GetValues()...)
	this.Reset()
	return r, err
}

//
func (this *BaseService) SQL(sql string, values ...interface{}) *BaseService {
	return this.AddCondition(sql).setValue(values...)
}

//
func (this *BaseService) InitSQL(sql string, values ...interface{}) *BaseService {
	return this.AddCondition(sql).setValue(values...)
}

//
func (this *BaseService) AppendSQL(sql string, values ...interface{}) *BaseService {
	return this.AddCondition(sql).setValue(values...)
}

// Limit limit条件
func (this *BaseService) Limit(num ...int) *BaseService {
	var limit string
	if len(num) == 1 {
		limit = ` limit ` + strconv.Itoa(num[0])
	} else if len(num) == 2 {
		limit = ` limit ` + strconv.Itoa(num[0]) + `,` + strconv.Itoa(num[1])
	}
	return this.AddCondition(limit)
}

// LimitPP limit分页，按page和pagesize
func (this *BaseService) LimitPP(page, pagesize int) *BaseService {
	var limit = ` limit ` + strconv.Itoa((page-1)*page) + `,` + strconv.Itoa(pagesize)
	return this.AddCondition(limit)
}

// AddCondition 添加条件 无需key值,无value值
func (this *BaseService) AddCondition(condition string) *BaseService {
	//检查是否存在此key,有则添加
	condition = strings.TrimSpace(condition)
	for k, v := range this.condition {
		if v[0] == "condition" {
			this.condition[k][1] += " " + condition
			return this
		}
	}
	this.condition = append(this.condition, [2]string{"condition", condition})
	return this
}

// Set 默认等于
func (this *BaseService) Set(key string, value interface{}) *BaseService {
	return this.SetEq(key, value)
}

// SetEq 等于
func (this *BaseService) SetEq(key string, value interface{}) *BaseService {
	this.setCondtion(key, " and `"+key+"` = ? ")
	this.setValue(value)
	return this
}

// SetNotEq 不等于
func (this *BaseService) SetNotEq(key string, value interface{}) *BaseService {
	this.setCondtion(key, " and `"+key+"` != ? ")
	this.setValue(value)
	return this
}

// SetNull 为空
func (this *BaseService) SetNull(key string) *BaseService {
	this.setCondtion(key, " and `"+key+"` is null ")
	return this
}

// SetNotNull 不为空
func (this *BaseService) SetNotNull(key string) *BaseService {
	this.setCondtion(key, " and `"+key+"` is not null ")
	return this
}

// SetGt 大于
func (this *BaseService) SetGt(key string, value interface{}) *BaseService {
	this.setCondtion(key, " and `"+key+"` > ? ")
	this.setValue(value)
	return this
}

// SetGe 大于等于
func (this *BaseService) SetGe(key string, value interface{}) *BaseService {
	this.setCondtion(key, " and `"+key+"` >=? ")
	this.setValue(value)
	return this
}

// SetLt 小于
func (this *BaseService) SetLt(key string, value interface{}) *BaseService {
	this.setCondtion(key, " and `"+key+"` < ? ")
	this.setValue(value)
	return this
}

// SetLe 小于等于
func (this *BaseService) SetLe(key string, value interface{}) *BaseService {
	this.setCondtion(key, " and `"+key+"` <= ? ")
	this.setValue(value)
	return this
}

// SetBetween 两者之间
func (this *BaseService) SetBetween(key string, value1 interface{}, value2 interface{}) *BaseService {
	this.setCondtion(key, " and `"+key+"` between ? and ？")
	this.values = append(this.values, value1, value2)
	return this
}

// SetOr 或者等于
func (this *BaseService) SetOr(key string, value interface{}) *BaseService {
	this.setCondtion(key, " or `"+key+"` = ? ")
	this.setValue(value)
	return this
}

// SetOrLike 或者相似
func (this *BaseService) SetOrLike(key string, value interface{}) *BaseService {
	this.setCondtion(key, " or `"+key+"` like ? ")
	this.setValue(value)
	return this
}

// SetLike 相似
func (this *BaseService) SetLike(key string, value interface{}) *BaseService {
	this.setCondtion(key, " and `"+key+"` like ?")
	this.setValue(value)
	return this
}

// SetDate 设置日期
func (this *BaseService) SetDate(key string, value interface{}) *BaseService {
	this.setCondtion(key, " and to_days(`"+key+"`) = to_days(?) ")
	this.setValue(value)
	return this
}

// SetDateGe 设置起始日期
func (this *BaseService) SetDateGe(key string, value interface{}) *BaseService {
	this.setCondtion(key, " and to_days(`"+key+"`) >= to_days(?) ")
	this.setValue(value)
	return this
}

// SetDateLe 设置结束日期
func (this *BaseService) SetDateLe(key string, value interface{}) *BaseService {
	this.setCondtion(key, " and to_days(`"+key+"`) <= to_days(?) ")
	this.setValue(value)
	return this
}

// SetIn In条件
func (this *BaseService) SetIn(key string, value ...interface{}) *BaseService {
	var length = len(value)
	if length == 0 {
		return nil
	}
	var tmp = make([]string, length)
	for i := 0; i < length; i++ {
		tmp[i] = "?"
	}
	this.setCondtion(key, " and `"+key+"` in ("+strings.Join(tmp, ",")+")")
	this.setValue(value...)
	return this
}

// SetIn In条件
func (this *BaseService) SetIn2(key string, values interface{}) *BaseService {
	var v = reflect.ValueOf(values)
	var k = v.Kind()
	switch k {
	case reflect.String, reflect.Bool, reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Float32, reflect.Float64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		this.setCondtion(key, " and `"+key+"` = ?")
		this.setValue(values)
	case reflect.Slice, reflect.Array:
		var l = v.Len()
		if l == 0 {
			return nil
		}
		var tmp = make([]string, l)
		for i := 0; i < l; i++ {
			tmp[i] = "?"
			this.setValue(v.Index(i).Interface())
		}
		this.setCondtion(key, " and `"+key+"` in ("+strings.Join(tmp, ",")+")")
	}
	return this
}

func (this *BaseService) SetSql(sql string, value ...interface{}) *BaseService {
	var length = len(value)
	if length == 0 {
		return nil
	}
	this.setCondtion("", " and "+sql+" ")
	this.setValue(value...)
	return this
}

// SetValue 直接添加value值
func (this *BaseService) SetValue(value ...interface{}) *BaseService {
	var length = len(value)
	if length == 0 {
		return nil
	}
	this.setValue(value...)
	return this
}

// GetCondition 返回所有条件
func (this *BaseService) GetCondition() string {
	var sql string
	var condition string
	for _, v := range this.condition {
		if v[0] == "sql" {
			sql = v[1]
			continue
		}
		condition += v[1]
	}
	if sql == "" {
		return condition
	}
	return sql + " where 1=1 " + condition
}

// GetWhere 返回where 1=1条件语句
func (this *BaseService) GetWhere() string {
	var condition string
	for _, v := range this.condition {

		condition += v[1]
	}

	return " where 1=1 " + condition
}

// GetValues 返回值
func (this *BaseService) GetValues() []interface{} {
	return this.values
}

// ParseQuery 将结构体解析到查询语句
func (this *BaseService) ParseQuery(model interface{}, flag ...bool) *BaseService {
	var value = reflect.ValueOf(model)
	if value.Kind() == reflect.Ptr || value.Kind() == reflect.Interface {
		value = value.Elem()
	}
	if value.Kind() != reflect.Struct {
		//		return errors.New("参数类型错误")
	}
	//	this.setValue(GetValues(value)...)
	if len(flag) == 0 {
		this.parseQuery(value, true)
	} else {
		this.parseQuery(value, flag[0])
	}
	return this
}
func (this *BaseService) parseQuery(t reflect.Value, flag bool) {
	if t.Kind() != reflect.Struct {
		return
	}
	for i := 0; i < t.NumField(); i++ {
		if t.Type().Field(i).Anonymous {
			this.parseQuery(t.Field(i), flag)
			continue
		}
		var key = SnakeName(t.Type().Field(i).Name)
		if flag {
			if key == "page" || key == "pagesize" || key == "page_size" {
				continue
			}
		}
		var field = t.Field(i)
		switch field.Kind() {
		case reflect.Bool:
		case reflect.Struct:
			if field.Type().Name() == "Time" {
				var tStr = fmt.Sprintf("%v", field.Interface())
				var tm, err = time.Parse("2006-01-02", strings.Split(tStr, " ")[0])
				if err != nil || tm.IsZero() {
					continue
				}
				this.setCondtion(key, " and `"+key+"`=?")
				this.setValue(field.Interface())
			}
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			if field.Int() != 0 {
				this.setCondtion(key, " and `"+key+"`=?")
				this.setValue(field.Int())
			}
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			if field.Uint() != 0 {
				this.setCondtion(key, " and `"+key+"`=?")
				this.setValue(field.Uint())
			}
		case reflect.Float32, reflect.Float64:
			if field.Float() != 0 {
				this.setCondtion(key, " and `"+key+"`=?")
				this.setValue(field.Float())
			}
		case reflect.Interface:
			this.parseQuery(field.Field(i), flag)
		case reflect.String:
			if field.String() != "" {
				this.setCondtion(key, " and `"+key+"`=?")
				this.setValue(field.String())
			}
		}
	}

}

func (this *BaseService) setCondtion(key, condition string) *BaseService {
	if strings.Contains(key, ".") {
		condition = strings.Replace(condition, "`", "", -1)
	}
	//检查是否存在此key,有则替换
	for k, v := range this.condition {
		if v[0] == key {
			this.condition[k][1] = condition
			return this
		}
	}
	this.condition = append(this.condition, [2]string{key, condition})
	return this
}
func (this *BaseService) setValue(value ...interface{}) *BaseService {
	this.values = append(this.values, value...)
	return this
}

// ResetCondition 重置condition
func (this *BaseService) ResetCondition() *BaseService {
	this.condition = make([][2]string, 0)
	return this
}

// ResetValue 重置value
func (this *BaseService) ResetValue() *BaseService {
	this.values = nil
	return this
}

// Reset 重置condition和value
func (this *BaseService) Reset() *BaseService {
	return this.ResetCondition().ResetValue()
}
