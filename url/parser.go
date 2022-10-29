package url

import (
	"bufio"
	"fmt"
	neturl "net/url"
	"strings"
)

type Parser struct {
	odataUrl string
	entity   string
	buf      struct {
		tok Token
		lit string
		n   int
	}
}

func getReader(input string) *bufio.Reader {
	return bufio.NewReader(strings.NewReader(input))
}

func NewParser(entity string, odataUrl string) *Parser {
	return &Parser{entity: entity, odataUrl: odataUrl}
}

func (p *Parser) scan(l *Lexer) (tok Token, lit string) {
	if p.buf.n != 0 {
		p.buf.n = 0
		return p.buf.tok, p.buf.lit
	}

	tok, lit = l.Scan()

	p.buf.tok, p.buf.lit = tok, lit
	return
}

func (p *Parser) scanIgnoreWhitespace(l *Lexer) (tok Token, lit string) {
	tok, lit = p.scan(l)
	if tok == SPACE {
		tok, lit = p.scan(l)
	}
	return
}

func (p *Parser) unscan() { p.buf.n = 1 }
func (p *Parser) reset()  { p.buf.n = 0 }

func (p *Parser) parseSelect(l *Lexer) (*SelectStatement, error) {
	stmt := &SelectStatement{}

	for {
		tok, lit := p.scanIgnoreWhitespace(l)
		if tok == EOF {
			break
		}
		if tok != IDENT && tok != ASTERISK {
			return nil, fmt.Errorf("found %q, expected field", lit)
		}
		stmt.Fields = append(stmt.Fields, lit)
		if tok, _ = p.scanIgnoreWhitespace(l); tok != COMMA {
			p.unscan()
			break
		}
	}

	return stmt, nil
}

func (p *Parser) parseSort(l *Lexer) ([]*SortField, error) {
	sorts := make([]*SortField, 0)

	for {
		tok, lit := p.scanIgnoreWhitespace(l)
		if tok == EOF {
			break
		}
		if tok != IDENT {
			return nil, fmt.Errorf("found %q, expected field", lit)
		}
		sortField := &SortField{Field: lit, Direction: "ASC"}
		sorts = append(sorts, sortField)
		if tok, lit = p.scanIgnoreWhitespace(l); tok != COMMA && (lit == "DESC" || lit == "ASC" || lit == "desc" || lit == "asc") {
			sortField.Direction = strings.ToUpper(lit)
		} else {
			p.unscan()
		}
		if tok, _ = p.scanIgnoreWhitespace(l); tok != COMMA {
			p.unscan()
			break
		}
	}

	return sorts, nil
}

func (p *Parser) parseFunction(functionName string, l *Lexer) (Expression, error) {
	fc := Expression{
		Op:    FUNCTION,
		Args:  make([]Expression, 0),
		Value: functionName,
	}

	for {
		tok, _ := p.scanIgnoreWhitespace(l)
		if tok == CLOSEBRACKET {
			break
		}
		if tok != COMMA {
			p.unscan()
			fcArgs, err := p.parsePrimaryExpr(l)
			if err != nil {
				return Expression{}, err
			}
			fc.Args = append(fc.Args, fcArgs)
		}
	}

	return fc, nil
}
func (p *Parser) parsePrimaryExpr(l *Lexer) (Expression, error) {
	tok, lit := p.scanIgnoreWhitespace(l)
	if tok == OPENBRACKET {
		expr, err := p.parseBinaryOp(l)
		if err != nil {
			return Expression{}, err
		}

		return Expression{
			Op:    GROUP,
			Args:  []Expression{expr},
			Value: "()",
		}, nil

	}
	if tok != IDENT {
		return Expression{}, fmt.Errorf("found %q, expected field", lit)
	}

	lhs := Expression{
		Op:    LITERAL,
		Value: lit,
	}

	tmpFunctionName := lit
	tok, _ = p.scanIgnoreWhitespace(l)
	switch tok {
	case OPENBRACKET:
		var err error
		lhs, err = p.parseFunction(tmpFunctionName, l)
		if err != nil {
			return Expression{}, err
		}
		return lhs, nil
	case COMMA:
		return lhs, nil
	case CLOSEBRACKET:
		p.unscan()
		return lhs, nil
	default:
		p.unscan()

	}

	result := Expression{
		Op:   COMPARE,
		Args: make([]Expression, 2),
	}
	result.Args[0] = lhs

	tok, lit = p.scanIgnoreWhitespace(l)
	if tok == OPER_COMPARE {
		result.Value = lit
		tok, lit = p.scanIgnoreWhitespace(l)
		if tok != IDENT {
			return Expression{}, fmt.Errorf("found %q, expected field", lit)
		}
		result.Args[1] = Expression{
			Op:    LITERAL,
			Value: lit,
		}
		tok, lit = p.scanIgnoreWhitespace(l)
		if tok != CLOSEBRACKET {
			p.unscan()
		}
		return result, nil
	}

	p.unscan()
	return Expression{}, fmt.Errorf("Syntax Error %q", lit)

}

func (p *Parser) parseBinaryOp(l *Lexer) (Expression, error) {
	lhs, err := p.parsePrimaryExpr(l)
	if err != nil {
		return Expression{}, err
	}
	tok, lit := p.scanIgnoreWhitespace(l)
	if tok == OPER_BOOLEAN {
		rhs, err := p.parseBinaryOp(l)
		if err != nil {
			return Expression{}, err
		}
		var op Operation
		switch lit {
		case "eq", "ne", "gt", "gte", "lt", "lte":
			op = COMPARE
		case "and", "or", "not":
			op = BOOLEAN
		}

		return Expression{
			Op:    op,
			Args:  []Expression{lhs, rhs},
			Value: lit,
		}, nil
	}
	return lhs, nil
}

func (p *Parser) parseFilter(l *Lexer) (FilterStatement, error) {
	var stmt = FilterStatement{}
	expr, err := p.parseBinaryOp(l)
	if err != nil {
		return stmt, nil
	}
	stmt.Expressions = append(stmt.Expressions, expr)
	return stmt, nil
}

func (p *Parser) Parse() (*StatementTree, error) {
	var err error
	odataParts, err := neturl.ParseQuery(p.odataUrl)
	if err != nil {
		return nil, fmt.Errorf("Error Parsing QueryString %s", p.odataUrl)
	}

	stmt := &StatementTree{}
	stmt.From = p.entity

	selectPart := odataParts.Get("$select")
	if selectPart == "" {
		selectPart = "*"
	}
	stmt.Select, err = p.parseSelect(NewLexer(getReader(selectPart)))
	p.reset()
	if err != nil {
		return nil, err
	}

	filterPart := odataParts.Get("$filter")
	if filterPart != "" {
		stmt.Filter, err = p.parseFilter(NewLexer(getReader(filterPart)))
		p.reset()
		if err != nil {
			return nil, err
		}
	}
	sortPart := odataParts.Get("$orderby")
	if sortPart != "" {
		stmt.Sort, err = p.parseSort(NewLexer(getReader(sortPart)))
		p.reset()
		if err != nil {
			return nil, err
		}
	}

	return stmt, nil
}
