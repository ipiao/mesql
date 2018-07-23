package medb

import (
	"fmt"
	"math/rand"
	"regexp"
	"strings"
	"time"
)

var seed = rand.New(rand.NewSource(time.Now().UnixNano()))
var reg = regexp.MustCompile(`\B[A-Z]`)

// transFieldName 转换字段名称,驼峰转蛇形
func transFieldName(name string) string {
	return strings.ToLower(reg.ReplaceAllString(name, "_$0"))
}

// 这是解决连接命名的临时方案
func RandomName() string {
	var res = fmt.Sprintf("%d%d", time.Now().UnixNano(), rand.Intn(8999)+1000)
	return res
}

func ParseTime(s string) (time.Time, error) {
	var layout = ""
	l := len(s)
	if l > 19 {
		s = s[:10] + " " + s[11:19]
		l = 19
	}
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
