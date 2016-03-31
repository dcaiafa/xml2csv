package compiler

import "github.com/dcaiafa/xml2csv/parser"

func (c *Compiler) emit(order parser.Order, ast parser.AST) {
	if order != parser.Prefix {
		return
	}

	switch ast.(type) {
	case *parser.Program:
		c.globalDFA = NFAToDFA(&c.globalNFA)
		for _, t := range c.transforms {
			t.DFA = NFAToDFA(&t.NFA)
		}
	}
}
