package url

type SelectStatement struct {
	Fields []string
}

type SortField struct {
	Field     string
	Direction string
}

type Operation int

const (
	LITERAL Operation = iota
	COMPARE
	BOOLEAN
	GROUP
	FUNCTION
)

type Expression struct {
	Op    Operation
	Args  []Expression
	Value string
}

type FilterStatement struct {
	Expressions []Expression
}

type StatementTree struct {
	Select *SelectStatement
	Sort   []*SortField
	Filter FilterStatement
	From   string
}
