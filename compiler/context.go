package compiler

import (
	"github.com/dcaiafa/xml2csv/parser"
)

type context struct {
	Names *parser.Names
	NFAs  []NFA
}
