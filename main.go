package main

import (
	"bufio"
	"fmt"
	"strings"

	"github.com/umbe77/udata/url"
)

func main() {
	selectUrl := "field1,field2, field3"
	parser := url.NewParser(bufio.NewReader(strings.NewReader(selectUrl)))
	if ast, err := parser.Parse(); err == nil {
		fmt.Printf("%+v\n", ast.Select)
	}
}
