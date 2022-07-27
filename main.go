package main

import (
	"fmt"

	"github.com/umbe77/udata/url"
)

func main() {
	odataUrl := "$select=field1,field2 , field3&$filter=Name eq 'Pippo'&$orderby=fild1, fild2 desc, field3 asc"
	parser := url.NewParser(odataUrl)
	ast, err := parser.Parse()
	if err != nil {
		fmt.Printf("%s\n", err)
	}
	fmt.Printf("SELECT %+v\n", ast.Select)
	for _, s := range ast.Sort {
		fmt.Printf("SORTField %s --> %s\n", s.Field, s.Direction)
	}
}
