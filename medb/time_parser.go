package medb

import (
	"database/sql"
	"reflect"
	"time"
)

// TimeParser 时间解析
type TimeParser func(*reflect.Value, interface{}) error

func defaultTimeParser(value *reflect.Value, field interface{}) error {
	var err error
	var s = sql.NullString{}
	err = s.Scan(field)
	if err != nil {
		return err
	}
	if s.Valid {
		t, err := time.ParseInLocation("2006-01-02 15:04:05", s.String, time.Local)
		if err != nil {
			t, err = time.ParseInLocation("2006-01-02", s.String, time.Local)
		}
		if err == nil {
			value.Set(reflect.ValueOf(t))
		}
	} else {
		var i = sql.NullInt64{}
		var err = i.Scan(field)
		if err != nil {
			return err
		}
		if i.Valid {
			t := time.Unix(i.Int64, 0)
			if err == nil {
				value.Set(reflect.ValueOf(t))
			}
		}
	}
	return err
}

var timeparse TimeParser = defaultTimeParser

// RegisterTimeParser 注册时间解析器
func RegisterTimeParser(timefun TimeParser) {
	timeparse = timefun
}
