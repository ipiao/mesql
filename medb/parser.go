package medb

import (
	"database/sql"
	"encoding/json"
	"log"
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

// RegisterTimeParser 注册时间解析器
func RegisterTimeParserFunc(fn func(string) (time.Time, error)) {
	timeparse = func(value *reflect.Value, field interface{}) error {
		var err error
		var s = sql.NullString{}
		err = s.Scan(field)
		if err == nil && s.Valid {
			t, err1 := fn(s.String)
			if err1 != nil {
				return err1
			}
			value.Set(reflect.ValueOf(t))
		}
		return err
	}
}

type JsonParser func(*reflect.Value, interface{}) error

func defaultJsonParser(value *reflect.Value, field interface{}) error {
	var err error
	var s = sql.NullString{}
	err = s.Scan(field)
	if err != nil {
		return err
	}
	if s.Valid {
		log.Println(s.String)
		val := reflect.New(value.Type())
		err = json.Unmarshal([]byte(s.String), val)
		if err != nil {
			log.Println(val)
			return err
		}
		value.Set(reflect.ValueOf(val))
	}
	return err
}

var jsonparse JsonParser = defaultJsonParser
