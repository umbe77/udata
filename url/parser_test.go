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

func TestParseSort(t *testing.T) {
	tbl := "sample"
	odataUrl := "$orderby=Field1, Field2 DESC,Field3 desc,Field4 ASC,Field5 asc,"

	p := NewParser(tbl, odataUrl)

	stmTree, err := p.Parse()
	if err != nil {
		t.Error(err)
	}
	if len(stmTree.Sort) != 5 {
		t.Errorf("Expect 5 Sorting elements got: %d", len(stmTree.Sort))
	}
	if stmTree.Sort[0].Field != "Field1" || stmTree.Sort[0].Direction != "ASC" {
		t.Errorf("Expected 'Field1 ASC', got: '%s %s'", stmTree.Sort[0].Field, stmTree.Sort[0].Direction)
	}
	if stmTree.Sort[1].Field != "Field2" || stmTree.Sort[1].Direction != "DESC" {
		t.Errorf("Expected 'Field2 DESC', got: '%s %s'", stmTree.Sort[1].Field, stmTree.Sort[1].Direction)
	}
	if stmTree.Sort[2].Field != "Field3" || stmTree.Sort[2].Direction != "DESC" {
		t.Errorf("Expected 'Field3 DESC', got: '%s %s'", stmTree.Sort[2].Field, stmTree.Sort[2].Direction)
	}
	if stmTree.Sort[3].Field != "Field4" || stmTree.Sort[3].Direction != "ASC" {
		t.Errorf("Expected 'Field4 ASC', got: '%s %s'", stmTree.Sort[3].Field, stmTree.Sort[3].Direction)
	}
	if stmTree.Sort[4].Field != "Field5" || stmTree.Sort[4].Direction != "ASC" {
		t.Errorf("Expected 'Field5 ASC', got: '%s %s'", stmTree.Sort[4].Field, stmTree.Sort[4].Direction)
	}
}
