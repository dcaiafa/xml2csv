package compiler

import (
	"github.com/dcaiafa/xml2csv/parser"
)

type Compiler struct {
	context context
	root    parser.AST
}

func NewCompiler(names *parser.Names, ast parser.AST) *Compiler {
	return &Compiler{
		context: context{
			Names: names,
		},
		root: ast,
	}
}

func (c *Compiler) CreateSymbols() {
	c.root.Visit(func(order parser.Order, ast parser.AST) {
		createSymbolsPass(&c.context, order, ast)
	})
}

func createSymbolsPass(c *context, order parser.Order, ast parser.AST) {
	if order != parser.Prefix {
		return
	}

	switch ast := ast.(type) {
	case *parser.PathExpr:
		createNfaFromPath(c, ast.Path)
	}
}

func createNfaFromPath(c *context, path []string) {

}
