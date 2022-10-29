package url

import "testing"

func TestParseFrom(t *testing.T) {
	tbl := "sample"
	odataUrl := ""

	p := NewParser(tbl, odataUrl)
	stmTree, err := p.Parse()
	if err != nil {
		t.Error(err)
	}
	if stmTree.From != "sample" {
		t.Errorf("Expected 'sample', got: %s", stmTree.From)
	}
}

func TestParseSelect(t *testing.T) {
	tbl := "sample"
	odataUrl := "$select=Field1, Field2 ,Field3,Field4,Field5"

	p := NewParser(tbl, odataUrl)

	stmTree, err := p.Parse()
	if err != nil {
		t.Error(err)
	}

	if stmTree.Select.Fields[0] != "Field1" {
		t.Errorf("Expected 'Field1', got: '%s'", stmTree.Select.Fields[0])
	}
	if stmTree.Select.Fields[1] != "Field2" {
		t.Errorf("Expected 'Field1', got: '%s'", stmTree.Select.Fields[1])
	}
	if stmTree.Select.Fields[2] != "Field3" {
		t.Errorf("Expected 'Field1', got: '%s'", stmTree.Select.Fields[2])
	}
	if stmTree.Select.Fields[3] != "Field4" {
		t.Errorf("Expected 'Field1', got: '%s'", stmTree.Select.Fields[3])
	}
	if stmTree.Select.Fields[4] != "Field5" {
		t.Errorf("Expected 'Field1', got: '%s'", stmTree.Select.Fields[4])
	}
}
func TestParseSelectTrailingComma(t *testing.T) {
	tbl := "sample"
	odataUrl := "$select=Field1,"

	p := NewParser(tbl, odataUrl)

	stmTree, err := p.Parse()
	if err != nil {
		t.Error(err)
	}

	if stmTree.Select.Fields[0] != "Field1" {
		t.Errorf("Expected 'Field1', got: '%s'", stmTree.Select.Fields[0])
	}
}
