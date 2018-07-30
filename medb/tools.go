package medb

// SnakeName 驼峰转蛇形
func SnakeName(base string) string {
	var r = make([]rune, 0, len(base))
	var b = []rune(base)
	for i := 0; i < len(b); i++ {
		if b[i] >= 'A' && b[i] <= 'Z' {
			if i == 0 {
				r = append(r, b[i]+32)
			} else {
				r = append(r, '_', b[i]+32)
			}
		} else {
			r = append(r, b[i])
		}
	}
	return string(r)
}

func StringsPtrToInterfaces(ss []string) []interface{} {
	var ret = make([]interface{}, len(ss))
	for i := range ss {
		ret[i] = &ss[i]
	}
	return ret
}
