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
	Operand Operation = iota
	Comparison
	Boolean
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
}
