package medb

// SQLBuilder sql constructor
// adjust for context
type SQLBuilder interface {
	ToSQL() (string, []interface{})
}
