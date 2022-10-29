package url

import "testing"

func TestScanWhiteSpaceOnlySpace(t *testing.T) {

	src := "    "
	l := NewLexer(getReader(src))
	tok, lit := l.scanWhitespace()
	if tok != SPACE || lit != "    " {
		t.Errorf("Expected ' ', got %v -- %s", tok, lit)
	}
}
func TestScanWhiteSpace(t *testing.T) {

	src := "   x "
	l := NewLexer(getReader(src))
	tok, lit := l.scanWhitespace()
	if tok != SPACE || lit != "   " {
		t.Errorf("Expected ' ', got %v -- %s", tok, lit)
	}
}

func TestScanIdentNumbers(t *testing.T) {
	src := "1234567890"
	l := NewLexer(getReader(src))
	tok, lit := l.scanIdent()
	if tok != IDENT || lit != "1234567890" {
		t.Errorf("Expected '1234567890', got %v -- %s", tok, lit)
	}
}

func TestScanIdentLettersAndSingleQuote(t *testing.T) {
	src := "ABCDEFGHIJKLMNOPQRSTUVWXYZ'abcdefghijklmnopqrstuvwxyz"
	l := NewLexer(getReader(src))
	tok, lit := l.scanIdent()
	if tok != IDENT || lit != "ABCDEFGHIJKLMNOPQRSTUVWXYZ'abcdefghijklmnopqrstuvwxyz" {
		t.Errorf("Expected 'ABCDEFGHIJKLMNOPQRSTUVWXYZ'abcdefghijklmnopqrstuvwxyz', got %v -- %s", tok, lit)
	}
}

func TestScanIdentNonIdent(t *testing.T) {
	src := ",*()"
	l := NewLexer(getReader(src))
	tok, lit := l.Scan()
	if tok != COMMA {
		t.Errorf("Expected ',', got %v -- %s", tok, lit)
	}
	tok, lit = l.Scan()
	if tok != ASTERISK {
		t.Errorf("Expected '*', got %v -- %s", tok, lit)
	}
	tok, lit = l.Scan()
	if tok != OPENBRACKET {
		t.Errorf("Expected '(', got %v -- %s", tok, lit)
	}
	tok, lit = l.Scan()
	if tok != CLOSEBRACKET {
		t.Errorf("Expected ')', got %v -- %s", tok, lit)
	}
}

func TestScanIdentIllegal(t *testing.T) {
	src := "\""
	l := NewLexer(getReader(src))
	tok, lit := l.Scan()
	if tok != ILLEGAL || lit != "\"" {
		t.Errorf("Expected '\"', got %v -- %s", tok, lit)
	}
}
