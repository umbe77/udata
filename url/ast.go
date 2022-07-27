package url

type SelectStatement struct {
	Fields []string
}

type SortField struct {
	Field     string
	Direction string
}

type StatementTree struct {
	Select *SelectStatement
	Sort   []*SortField
}
