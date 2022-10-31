package compiler

import (
	"fmt"
	"strings"

	"github.com/umbe77/udata/url"
)

type Sql struct{}

var compareOperatorMap = map[string]string{
	"eq":  "=",
	"ne":  "<>",
	"gt":  ">",
	"gte": ">=",
	"lt":  "<",
	"lte": "<=",
}

var booleanOperatorMap = map[string]string{
	"and": "AND",
	"or":  "OR",
	"NOT": "NOT",
}

func clearSqlInjection(value string) string {
	if strings.HasPrefix(value, "'") && strings.HasSuffix(value, "'") {

		return fmt.Sprintf("'%s'", strings.ReplaceAll(strings.Trim(value, "'"), "'", "''"))
	}
	return value
}

func (c Sql) formatBinaryOp(lhs, rhs url.Expression, op string) (string, error) {
	lhsStr, err := c.CompileExpression(lhs)
	if err != nil {
		return "", err
	}
	rhsStr, err := c.CompileExpression(rhs)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%s %s %s", lhsStr, op, rhsStr), nil

}

func (c Sql) formatKnownFunctions(functionName string, arguments []url.Expression) (string, error) {
	var result string
	switch functionName {
	case "contains":
		lhs, err := c.CompileExpression(arguments[0])
		if err != nil {
			return "", err
		}
		rhs, err := c.CompileExpression(arguments[1])
		if err != nil {
			return "", err
		}
		result = fmt.Sprintf("%s LIKE '%%%s%%'", lhs, strings.Trim(rhs, "'"))
		break
	default:
		return "", fmt.Errorf("%s function, not valid", functionName)
	}
	return result, nil
}

func (c Sql) Compile(stmTree *url.StatementTree) (string, error) {
	selectStmt := fmt.Sprintf("SELECT %s", strings.Join(stmTree.Select.Fields, ", "))

	orderbyFields := make([]string, len(stmTree.Sort))
	for i, sortField := range stmTree.Sort {
		orderbyFields[i] = fmt.Sprintf("%s %s", sortField.Field, sortField.Direction)
	}
	orderbyStmt := fmt.Sprintf("GROUP BY %s", strings.Join(orderbyFields, ", "))

	whereStr, err := c.CompileExpression(stmTree.Filter.Expressions[0])
	if err != nil {
		return "", nil
	}
	whereStmt := fmt.Sprintf("WHERE %s", whereStr)

	return fmt.Sprintf("%s\n%s\n%s", selectStmt, whereStmt, orderbyStmt), nil

}
func (c Sql) CompileExpression(expr url.Expression) (string, error) {
	var result string
	var err error
	switch expr.Op {
	case url.LITERAL:
		result = fmt.Sprintf("%s", clearSqlInjection(expr.Value))
		break
	case url.COMPARE:
		result, err = c.formatBinaryOp(expr.Args[0], expr.Args[1], compareOperatorMap[expr.Value])
		if err != nil {
			return "", err
		}
		break
	case url.BOOLEAN:
		result, err = c.formatBinaryOp(expr.Args[0], expr.Args[1], booleanOperatorMap[expr.Value])
		if err != nil {
			return "", err
		}
		break
	case url.GROUP:
		groupStr, err := c.CompileExpression(expr.Args[0])
		if err != nil {
			return "", err
		}
		result = fmt.Sprintf("(%s)", groupStr)
		break
	case url.FUNCTION:
		result, err = c.formatKnownFunctions(expr.Value, expr.Args)
		if err != nil {
			return "", nil
		}
		break
	default:
		result = ""
	}
	return result, nil
}
