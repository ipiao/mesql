package melog

import (
	"log"
)

type Mode int

var logmode int

const (
	ModeDebug = 1 << iota
	ModeRelease
)

type MeLog struct {
	mode int
}

func NewMeLog(mode int) *MeLog {
	return &MeLog{
		mode: mode,
	}
}

// debug模式下的输出
func (this *MeLog) Debug(args ...interface{}) {
	if (this.mode & ModeDebug) > 0 {
		log.Println(args...)
	}
}
