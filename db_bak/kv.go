package mesql

// KV 键值结构
type KV struct {
	key   string
	value interface{}
}

type KVs struct {
	key    string
	values []interface{}
}

func newKV(k string, v interface{}) KV {
	return KV{key: k, value: v}
}

func newKVs(k string, v ...interface{}) KVs {
	return KVs{key: k, values: v}
}
