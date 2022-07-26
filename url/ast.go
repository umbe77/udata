package url

type SelectStatement struct {
	Fields []string
}

type StatementTree struct {
	Select *SelectStatement
}
