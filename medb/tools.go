package medb

import (
	"regexp"
	"strings"
)

var reg = regexp.MustCompile(`\B[A-Z]`)

// transFieldName 转换字段名称,驼峰转蛇形
func transFieldName(name string) string {
	return strings.ToLower(reg.ReplaceAllString(name, "_$0"))
}
