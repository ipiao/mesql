package meorm

import "testing"

func TestDelete(t *testing.T) {
	b := DeleteFrom("table")
	b.Where("id<>1")
	b.WhereIn("col", 1, 2, 3)
	t.Log(b.ToSQL())
}

func TestSelect(t *testing.T) {
	b := Select("id,name").From("table")
	b.WhereIn("id", 1, 2, 3)
	b.WhereLike("name", "a")
	t.Log(b.ToSQL())
}
