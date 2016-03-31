package compiler

import (
	"bytes"
	"testing"

	"github.com/dcaiafa/xml2csv/parser"
)

func symbolName(n *parser.Names, sym int) string {
	nameID := sym / 10
	enter := (sym % 10) == 1
	label := n.NameFromID(nameID)
	if enter {
		label = ">" + label
	} else {
		label = "<" + label
	}
	return label
}

func TestCompiler(t *testing.T) {
	const prog = `
  	foreach /books/book
		where (/author/firstname == "Will" && /author/lastname == "Wight") ||
			/genre/primary == "fantasy"
		select a: /info/isbn, b: /info/title`

	n := &parser.Names{}
	p, err := parser.Parse("test", bytes.NewReader([]byte(prog)), n)
	if err != nil {
		t.Fatal("failed to parse program")
	}
	c := NewCompiler(n, p, func(msg string) {
		t.Fatal(msg)
	})
	if !c.Compile() {
		t.Errorf("failed to compile")
	}

	/*
		for _, t := range c.transforms {
				fmt.Fprintln(os.Stderr, nfaToDot(&t.NFA, func(sym int) string {
					return symbolName(n, sym)
				}))
			fmt.Fprintln(os.Stderr, dfaToDot(t.DFA, func(sym int) string {
				return symbolName(n, sym)
			}))
		}
	*/

}
