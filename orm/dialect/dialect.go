package dialect

// Dialect dialect
type Dialect int

// Holder holder
type Holder string

// dialects
const (
	UnKnown Dialect = iota
	Mysql
)

// holders
const (
	mysqlHolder = "?"
)

func ConvertDriverNameToDialect(s string) Dialect {
	var d Dialect
	switch s {
	case "mysql":
		d = Mysql
	default:
		d = UnKnown
	}
	return d
}
