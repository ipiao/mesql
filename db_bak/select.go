package mesql

// queryBuilder 查询构造器
type queryBuilder struct {
	distinct bool
	cols     []Column
	tables   []Table
}

////distinct
//type Distinct struct {
//	isDistinct bool //默认false
//}

//列，表的字段
type Column struct {
	name   string //查询列命,可以是聚合函数
	isFunc bool   //是否为函数
	table  Table
}

//查询表，有连接
type Table struct {
	name  string //表名
	alias string //别名
	join  Join   //连接方式
}

type Join string

//常量 连接方式
func (this Join) toString() string {
	if this == JOIN {
		return "join"
	} else if this == LEFT_JOIN {
		return "left join"
	} else if this == RIGHT_JOIN {
		return "right join"
	} else if this == FULL_JOIN {
		return "full join"
	} else if this == INNER_JOIN {
		return "inner join"
	} else {
		return " "
	}
}

const (
	JOIN       Join = "join"
	LEFT_JOIN  Join = "left join"
	RIGHT_JOIN Join = "right join"
	FULL_JOIN  Join = "full join"
	INNER_JOIN Join = "inner join"
)

//分组
type GroupBy struct {
	column string
}

//排序
type OrderBy struct {
	column string
	rule   string
}

//分页
type Limit struct {
	start int
	count int
}
