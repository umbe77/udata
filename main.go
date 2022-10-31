package main

import (
	"fmt"

	"github.com/umbe77/udata/compiler"
	"github.com/umbe77/udata/url"
)

func main() {
	//odataUrl := "$select=field1,field2 , field3&$filter=Name eq 'Pippo'&$orderby=fild1, fild2 DESC, field3 asc"
	//odataUrl := "$filter=Name eq 'Pippo' and LastName ne 'Pluto' or Age gte 20"
	//odataUrl := "$filter=(Name eq 'Pippo' or (LastName ne 'Pluto' and Age gte 20 and (a eq 8 or b gt 4)))"
	//odataUrl := "$filter=contains(Name, 'Pippo') or Name eq 'Pluto'"
	odataUrl := "$select=field1,field2 , field3&$orderby=fild1, fild2 DESC, field3 asc&$filter=(Name eq 'Pippo' and contains(Role, 'test') or (LastName ne 'Pluto' and Age gte 20 and (a eq 8 or b gt 4)))"

	parser := url.NewParser("sample_table", odataUrl)
	ast, err := parser.Parse()
	if err != nil {
		fmt.Printf("%s\n", err)
	}

	c := compiler.Sql{}
	result, err := c.Compile(ast)
	if err != nil {
		fmt.Printf("%s\n", err)
	}
	fmt.Println(result)

	//filterJson, err := json.Marshal(ast)
	//fmt.Printf("%s\n", string(filterJson))
}
