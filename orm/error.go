package meorm

import (
	"errors"
)

var (
	ErrSelectCols = errors.New("[mesql]:不能选择空列")
)
