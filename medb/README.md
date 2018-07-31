### medb

#### 功能
- 对官方sql包进行了简单的封装
- 对查询进行了重点封装，支持字符串数组，数组map，结构体等多种查询结果解析。`ScanTo`，`ScanMap`，`ScanStrings`
- 在结构体解析`ScanTo`的时候，支持字段自定义解析
- 开放相关Tag，支持自定义Tag
- 组合了sqllog，在sql出错的时候打印错误sql和错误信息，或者通过设置debug模式，在不出错的时候，也可以打印sql

#### 全局自定义字段

可以在sql执行前之前进行字段修改

``` go
var (
	MedbTag           = "db"        // medb 标签字段,不同字段间用`,`隔离
	MedbFieldName     = "col"       // 标签解析映射后，col对应字段名
	MedbFieldIgnore   = "_"         // "_" not "-"
	MedbFieldCp       = "cp"        // custome parse 字段自定义解析标签
	MedbFieldCpMethod = "MedbParse" // custome parse 字段自定义解析方法名
)
```

- 解析 `MedbTag`,若是值是`_`或者值是不带`:`的单长度映射,比如`db:"age"`，则会将`_`，或者`age`作为`MedbFieldName`映射值

#### TimeParser 时间解析器
针对`time.Time`结构解析，使用`RegisterTimeParser`可以注册自定义的时间解析函数，或者`RegisterTimeParserFunc`规避使用反射

- example

``` go
func ParseTime(s string) (time.Time, error) {
	var layout = ""
	l := len(s)
	switch l {
	case 19:
		s1 := s[4:5]
		layout = fmt.Sprintf("2006%s01%s02 15:04:05", s1, s1)
	case 10:
		s1 := s[4:5]
		layout = fmt.Sprintf("2006%s01%s02", s1, s1)
	case 8:
		layout = "15:04:05"
	case 7:
		s1 := s[4:5]
		layout = fmt.Sprintf("2006%s01", s1)
	}
	if layout == "" {
		return time.Time{}, fmt.Errorf("不支持的时间格式:%s", s)
	}
	v, err := time.ParseInLocation(layout, s, time.Local)
	return v, err
}

medb.RegisterTimeParserFunc(ParseTime)
```

#### MedbParse 自定义字段解析

要求是自定义的结构体，比如有这样一个结构

``` go
type Time time.Time

type User struct {
    Id        int
    Name      string
    Age       int  `db:"_"`
    CreatedAt Time `db:"col:created_at,cp"`
}


func (t Time) MedbParse(s string) *Time {
	tim, _ := datetool.ParseTime(s)
	t = Time(tim)
	return &t
}
```
在`CreateAt`是`time.Time`类型时，会走TimeParser解析(无法通过MedbParse进行自定义解析)。在以上结构中通过标签`cp`==`custome parse`，将会由函数`MedbParse`进行解析，因为数据库几乎所有字段都可以进行`NullString`解析，所以会将字段值以`string`类型作为入参。也就是`MedbParse(s string)`的参数类是固定的。
对于返回值，若是存在返回值，则会第一个返回值作为解析值，无论时指针还是结构体。
若是不存在返回值，务必` (t *Time)`使用指针接受器，这样在，否则讲无法对反射值进行赋值，也就是说以下形式是不行的

``` go
func (t Time) MedbParse(s string) {
	tim, _ := datetool.ParseTime(s)
	t = Time(tim)
}
```

#### Executor 执行器
这是一个接口，`db`，`tx`，`stmt`都实现了这个接口，在orm中可以方便调用
