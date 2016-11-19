// 定义错误
package medb

import (
	"errors"
)

var (
	ErrRegister   = errors.New("[medb]:同名驱动已存在")
	ErrRegisterDB = errors.New("[medb]:同名连接已存在")
	ErrDriver     = errors.New("[medb]:驱动不能为空")
	ErrMountDB    = errors.New("[medb]:db嵌入失败")
	ErrCommit     = errors.New("[medb]:事务未开启")
	ErrNoNextRow  = errors.New("[medb]:没有下一个结果")
	ErrRollBack   = ErrCommit
)
