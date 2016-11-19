package medb

import (
	"math/rand"
	"regexp"
	"strconv"
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
	var res = time.Now().Format("20060102150405")
	for i := 0; i < 6; i++ {
		res += strconv.FormatInt(seed.Int63n(1001), 16)
	}
	return res
}
