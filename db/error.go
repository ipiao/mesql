// 定义错误
package medb

import (
	"errors"
)

var (
	ErrRegister   = errors.New("[medb]:同名驱动已存在")
	ErrRegisterDB = errors.New("[medb]:同名连接已存在")
	ErrDriver     = errors.New("[medb]:驱动不能为空")
	ErrCommit     = errors.New("[medb]:there's no transaction begun")
	ErrRollBack   = ErrCommit
	ErrNoNextRow  = errors.New("[medb]:it has not next row")
)
