package medb

import (
	"strings"
)

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
