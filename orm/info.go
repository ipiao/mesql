package meorm

import (
	"strings"
)

const (
	TableNameMethod = "TableName"
	OrmTag          = "meorm"
)

// 解析tag
func ParseTag(tag string) map[string]string {
	var tagMap = make(map[string]string, 0)
	var cats []string = strings.Split(tag, ";")
	for _, cat := range cats {
		var nv = strings.Split(cat, ":")
		var tagName = nv[0]
		var tagValue = ""
		if len(nv) > 1 {
			tagValue = nv[1]
		}
		tagMap[tagName] = tagValue
	}
	return tagMap
}
