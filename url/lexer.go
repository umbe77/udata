package url

import (
	"bufio"
	"bytes"
)

type Token int

const (
	//Special tokens
	ILLEGAL Token = iota
	EOF
	SPACE
	//fields, function name...
	IDENT
	//Misc caracters
	COMMA
	ASTERISK
	//Operators
	OPER_COMPARE
	OPER_BOOLEAN
	//Parentesis
	OPENBRACKET
	CLOSEBRACKET
)

func isWhiteSpace(ch rune) bool {
	return ch == ' ' || ch == '\t' || ch == '\n' || ch == '\r'
}

func isIdent(ch rune) bool {
	return (ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z') || (ch >= '0' && ch <= '9') || ch == '-' || ch == '_' || ch == '\''
}

type Lexer struct {
	r *bufio.Reader
}

func NewLexer(r *bufio.Reader) *Lexer {
	return &Lexer{r: bufio.NewReader(r)}
}

var eof = rune(0)

func (l *Lexer) read() rune {
	ch, _, err := l.r.ReadRune()
	if err != nil {
		return eof
	}
	return ch
}

func (l *Lexer) unread() { _ = l.r.UnreadRune() }

func (l *Lexer) scanWhitespace() (tok Token, lit string) {
	var buf bytes.Buffer
	buf.WriteRune(l.read())

	for {
		if ch := l.read(); ch == eof {
			break
		} else if !isWhiteSpace(ch) {
			l.unread()
			break
		} else {
			buf.WriteRune(ch)
		}
	}

	return SPACE, buf.String()
}

func (l *Lexer) scanIdent() (Token, string) {
	var buf bytes.Buffer
	buf.WriteRune(l.read())

	for {
		if ch := l.read(); ch == eof {
			break
		} else if !isIdent(ch) {
			l.unread()
			break
		} else {
			buf.WriteRune(ch)
		}
	}

	lit := buf.String()
	if lit == "eq" ||
		lit == "ne" ||
		lit == "gt" ||
		lit == "gte" ||
		lit == "lt" ||
		lit == "lte" {
		return OPER_COMPARE, lit
	}

	if lit == "and" ||
		lit == "or" ||
		lit == "not" {
		return OPER_BOOLEAN, lit
	}

	return IDENT, lit
}

func (l *Lexer) Scan() (tok Token, lit string) {
	ch := l.read()

	if isWhiteSpace(ch) {
		l.unread()
		return l.scanWhitespace()
	} else if isIdent(ch) {
		l.unread()
		return l.scanIdent()
	}

	switch ch {
	case eof:
		return EOF, ""
	case ',':
		return COMMA, string(ch)
	case '*':
		return ASTERISK, string(ch)
	case '(':
		return OPENBRACKET, string(ch)
	case ')':
		return CLOSEBRACKET, string(ch)
	}

	return ILLEGAL, string(ch)
}
