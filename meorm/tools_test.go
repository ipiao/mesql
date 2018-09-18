package meorm

import (
	"fmt"
	"testing"

	"github.com/bmizerany/assert"
)

type Table struct {
	Id int
}

// its reciver must be Table not *Table
func (t Table) TableName() string {
	return fmt.Sprintf("table_%d", t.Id)
}

func TestGetTableName(t *testing.T) {
	assert.Equal(t, GetTableName(Table{1}), "table_1")
	assert.Equal(t, GetTableName(&Table{2}), "table_2")
	assert.Equal(t, GetTableName([]*Table{{1}}), "table_1")
	assert.Equal(t, GetTableName([]*Table{{1}, {2}}), "table_1")
	assert.Equal(t, GetTableName([]Table{{2}, {2}}), "table_2")
	assert.Equal(t, GetTableName(&[]Table{{1}, {2}}), "table_1")
}
