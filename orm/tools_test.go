package meorm

import "testing"

type Table struct {
}

// its reciver must be Table not *Table
func (Table) TableName() string {
	return "tab"
}

func TestGetTableName(t *testing.T) {
	tab := Table{}
	var name = GetTableName(tab)
	t.Log(name)
}
