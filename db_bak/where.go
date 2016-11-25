package mesql

//where体
type where struct {
	eqs       []KV
	ges       []KV
	gts       []KV
	les       []KV
	lts       []KV
	ors       []KV
	orGts     []KV
	orLts     []KV
	between   []KVs
	orBetween []KVs
	likes     []KV
	orLikes   []KV
	in        []KVs
	exists    []KV
}

//
func (this *where) Eq(k string, v interface{}) {
	this.eqs = append(this.eqs, newKV(k, v))
}
