package mesql

import (
	"errors"
)

var (
	ErrRegister  = errors.New("mesql:driver has the same name exists")
	ErrDriver    = errors.New("mesql:driver can not be null")
	ErrCommit    = errors.New("mesql:there's no transaction begun")
	ErrRollBack  = ErrCommit
	ErrNoNextRow = errors.New("mesql:it has not next row")
)
