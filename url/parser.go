package url

import (
	"bufio"
	"fmt"
)

type Parser struct {
	l   *Lexer
	buf struct {
		tok Token
		lit string
		n   int
	}
}

func NewParser(r *bufio.Reader) *Parser {
	return &Parser{l: NewLexer(r)}
}

func (p *Parser) scan() (tok Token, lit string) {
	if p.buf.n != 0 {
		p.buf.n = 0
		return p.buf.tok, p.buf.lit
	}

	tok, lit = p.l.Scan()

	p.buf.tok, p.buf.lit = tok, lit
	return
}

func (p *Parser) scanIgnoreWhitespace() (tok Token, lit string) {
	tok, lit = p.scan()
	if tok == SPACE {
		tok, lit = p.scan()
	}
	return
}

func (p *Parser) unscan() { p.buf.n = 1 }

func (p *Parser) parseSelect() (*SelectStatement, error) {
	stmt := &SelectStatement{}

	for {
		tok, lit := p.scanIgnoreWhitespace()
		if tok != IDENT {
			return nil, fmt.Errorf("found %q, expected field", lit)
		}
		stmt.Fields = append(stmt.Fields, lit)
		if tok, _ = p.scanIgnoreWhitespace(); tok != COMMA {
			p.unscan()
			break
		}
	}

	return stmt, nil
}

func (p *Parser) Parse() (*StatementTree, error) {
	stmt := &StatementTree{}
	var err error
	stmt.Select, err = p.parseSelect()
	if err != nil {
		return nil, err
	}
	return stmt, nil
}
