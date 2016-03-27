package main

import (
	"flag"
	"fmt"
	"github.com/dcaiafa/xml2csv/parser"
	"os"
)

func run(in, transform string) error {
	if transform == "" {
		return fmt.Errorf("missing required parameter -t")
	}

	transformFile, err := os.Open(transform)
	if err != nil {
		return fmt.Errorf("failed to open '%v'", err)
	}

	names := &parser.Names{}
	parser.Parse(transform, transformFile, names)

	return nil
}

func main() {
	var (
		in = flag.String("i", "", "input xml file")
		t  = flag.String("t", "", "transformation definition file")
	)

	flag.Parse()

	err := run(*in, *t)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
	}
}
