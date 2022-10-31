package compiler

import (
	"github.com/umbe77/udata/url"
)

type Compiler interface {
	Compile(stmTree *url.StatementTree) (string, error)
	CompileExpression(expr url.Expression) (string, error)
}
