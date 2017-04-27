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

// ParseTag 解析标签
func ParseTag(tag string) map[string]string {
	var res = make(map[string]string)
	var arr = strings.Split(tag, ";")
	for _, a := range arr {
		if strings.Contains(a, ":") {
			brr := strings.Split(a, ":")
			res[brr[0]] = brr[1]
		} else {
			res[MedbFieldName] = a
		}
	}
	return res
}
