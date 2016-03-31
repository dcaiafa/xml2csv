package compiler

import (
	"fmt"

	"github.com/dcaiafa/xml2csv/parser"
)

type ErrorReporter func(msg string)

type Compiler struct {
	rootAST   parser.AST
	names     *parser.Names
	exprs     map[parser.AST]*expr
	errorFunc ErrorReporter

	globalNFA    NFA
	globalDFA    *DFA
	transforms   []*transform
	curTransform *transform
	failed       bool
}

func NewCompiler(
	names *parser.Names,
	prog *parser.Program,
	reportErrorFunc ErrorReporter) *Compiler {

	c := &Compiler{
		rootAST:   prog,
		names:     names,
		exprs:     make(map[parser.AST]*expr),
		errorFunc: reportErrorFunc,
	}

	return c
}

func (c *Compiler) Compile() bool {
	c.rootAST.Visit(c.check)
	if c.failed {
		return false
	}
	c.rootAST.Visit(c.emit)
	return !c.failed
}

func (c *Compiler) fail(fmsg string, args ...interface{}) {
	c.failed = true
	if c.errorFunc != nil {
		c.errorFunc(fmt.Sprintf(fmsg, args...))
	}
}
