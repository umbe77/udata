package url

import (
	"bufio"
	"fmt"
	neturl "net/url"
	"strings"
)

type Parser struct {
	odataUrl string
	buf      struct {
		tok Token
		lit string
		n   int
	}
}

func getReader(input string) *bufio.Reader {
	return bufio.NewReader(strings.NewReader(input))
}

func NewParser(odataUrl string) *Parser {
	return &Parser{odataUrl: odataUrl}
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
		if tok != IDENT {
			return nil, fmt.Errorf("found %q, expected field", lit)
		}
		sortField := &SortField{Field: lit, Direction: "ASC"}
		sorts = append(sorts, sortField)
		//TODO: Check for sorting direction should ignore case
		if tok, lit = p.scanIgnoreWhitespace(l); tok != COMMA && (lit == "desc" || lit == "asc") {
			sortField.Direction = strings.ToUpper(lit)
		} else {
			p.unscan()
		}
		if tok, _ = p.scanIgnoreWhitespace(l); tok != COMMA {
			//TODO: Capire se rilasciare un errore sintattico
			p.unscan()
			break
		}
	}

	return sorts, nil
}

func (p *Parser) parsePrimaryExpr(l *Lexer) (Expression, error) {
	tok, lit := p.scanIgnoreWhitespace(l)
	if tok != IDENT {
		return Expression{}, fmt.Errorf("found %q, expected field", lit)
	}

	result := Expression{
		Op:   Compare,
		Args: make([]Expression, 2),
	}

	result.Args[0] = Expression{
		Op:    Literal,
		Value: lit,
	}

	tok, lit = p.scanIgnoreWhitespace(l)
	if tok == OPER_COMPARE {
		result.Value = lit
		tok, lit = p.scanIgnoreWhitespace(l)
		if tok != IDENT {
			return Expression{}, fmt.Errorf("found %q, expected field", lit)
		}
		result.Args[1] = Expression{
			Op:    Literal,
			Value: lit,
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
			op = Compare
		case "and", "or", "not":
			op = Boolean
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

	selectPart := odataParts.Get("$select")
	if selectPart == "" {
		selectPart = "*"
	}
	stmt := &StatementTree{}
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
