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
		if tok, lit = p.scanIgnoreWhitespace(l); tok != COMMA && (lit == "desc" || lit == "asc") {
			sortField.Direction = lit
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

func (p *Parser) Parse() (*StatementTree, error) {
	var err error
	odataParts, err := neturl.ParseQuery(p.odataUrl)
	if err != nil {
		return nil, fmt.Errorf("Error Parsing QueryString %s", p.odataUrl)
	}

	selectPart := odataParts.Get("$select")
	fmt.Printf("%q\n", selectPart)
	if selectPart == "" {
		selectPart = "*"
	}
	stmt := &StatementTree{}
	stmt.Select, err = p.parseSelect(NewLexer(getReader(selectPart)))
	p.reset()
	if err != nil {
		return nil, err
	}

	sortPart := odataParts.Get("$orderby")
	stmt.Sort, err = p.parseSort(NewLexer(getReader(sortPart)))
	p.reset()
	if err != nil {
		return nil, err
	}

	return stmt, nil
}
