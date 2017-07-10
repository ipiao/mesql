package medb

import (
	"database/sql"
	"reflect"
	"time"
)

// TimeParser parse time data to time.Time
type TimeParser func(*reflect.Value, interface{}) error

// defaultTimeParser Default time parse
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
	}
	return err
}

// set time parser to default-timeparser
var timeparse TimeParser = defaultTimeParser

// RegisterTimeParser 注册时间解析器
func RegisterTimeParser(timefun TimeParser) {
	timeparse = timefun
}
