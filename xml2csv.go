package main

import (
	"bytes"
	"encoding/csv"
	"encoding/xml"
	"flag"
	"fmt"
	"gopkg.in/xmlpath.v2"
	"io"
	"os"
	"strings"
)

func run(xml, out, path, include, exclude string, cols []string) error {
	var err error

	switch {
	case xml == "":
		return fmt.Errorf("missing required parameter -xml")
	case out == "":
		return fmt.Errorf("missing required parameter -out")
	case path == "":
		return fmt.Errorf("missing required parameter -path")
	case len(cols) == 0:
		return fmt.Errorf("missing column xpath parameters")
	}

	pathElements := strings.Split(path, "/")
	if len(pathElements) < 2 && pathElements[0] != "" {
		return fmt.Errorf("parameter -path is invalid")
	}
	pathElements = pathElements[1:]

	if pathElements[len(pathElements)-1] == "" {
		pathElements = pathElements[:len(pathElements)-1]
		if len(pathElements) == 0 {
			return fmt.Errorf("parameter -path is invalid")
		}
	}

	var includePath *xmlpath.Path
	if include != "" {
		includePath, err = xmlpath.Compile(include)
		if err != nil {
			return fmt.Errorf(
				"parameter -include is not a valid XPath expression: %v", err)
		}
	}

	var excludePath *xmlpath.Path
	if exclude != "" {
		excludePath, err = xmlpath.Compile(exclude)
		if err != nil {
			return fmt.Errorf(
				"parameter -exclude is not a valid XPath expression: %v", err)
		}
	}

	colHeaders := make([]string, len(cols))
	colPaths := make([]*xmlpath.Path, len(cols))
	for i, col := range cols {
		elems := strings.SplitN(col, "=", 2)
		if len(elems) != 2 {
			return fmt.Errorf(
				"positional parameter '%v' should be formatted as <header>=<xpath>",
				col)
		}

		colHeaders[i] = strings.TrimSpace(elems[0])

		colPath, err := xmlpath.Compile(elems[1])
		if err != nil {
			return fmt.Errorf(
				"positional parameter '%v' is not a valid XPath expression: %v",
				col, err)
		}
		colPaths[i] = colPath
	}

	xmlFile, err := os.Open(xml)
	if err != nil {
		return fmt.Errorf("failed to open %v: %v", xml, err)
	}

	writer := csv.NewWriter(os.Stdout)
	writer.Write(colHeaders)

	err = parse(xmlFile, pathElements, func(xml []byte) error {
		node, err := xmlpath.Parse(bytes.NewReader(xml))
		if err != nil {
			return fmt.Errorf("failed to parse xml: %v", err)
		}

		if includePath != nil && !includePath.Exists(node) {
			return nil
		}

		if excludePath != nil && excludePath.Exists(node) {
			return nil
		}

		values := make([]string, len(colPaths))
		for i, colPath := range colPaths {
			value, ok := colPath.String(node)
			if ok {
				values[i] = strings.TrimSpace(value)
			}
		}

		err = writer.Write(values)
		if err != nil {
			return fmt.Errorf("failed to write to csv file: %v", err)
		}
		return nil
	})

	writer.Flush()

	if err != nil {
		return err
	}

	return nil
}

func parse(xmlFile io.Reader, path []string, emit func(xml []byte) error) error {
	buffer := &bytes.Buffer{}
	decoder := xml.NewDecoder(xmlFile)
	encoder := xml.NewEncoder(buffer)
	stack := make([]string, 0, 10)
	matching := false

	for {
		t, err := decoder.Token()
		if err == io.EOF {
			break
		} else if err != nil {
			return fmt.Errorf("failed to parse xml: %v", err)
		}

		switch t := t.(type) {
		case xml.StartElement:
			stack = append(stack, t.Name.Local)
			if !matching {
				matching = sliceEq(stack, path)
			}
			if matching {
				encoder.EncodeToken(t)
			}

		case xml.EndElement:
			stack = stack[:len(stack)-1]
			if matching {
				encoder.EncodeToken(t)
				if len(stack) < len(path) {
					matching = false
					encoder.Flush()
					err := emit(buffer.Bytes())
					if err != nil {
						return fmt.Errorf("failed to process xml node: %v", err)
					}
					buffer.Reset()
				}
			}

		default:
			if matching {
				encoder.EncodeToken(t)
			}
		}
	}
	return nil
}

func sliceEq(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := 0; i < len(a); i++ {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func main() {
	var (
		xml     = flag.String("xml", "", "")
		out     = flag.String("out", "", "")
		path    = flag.String("path", "", "")
		include = flag.String("include", "", "")
		exclude = flag.String("exclude", "", "")
	)

	flag.Parse()

	err := run(*xml, *out, *path, *include, *exclude, flag.Args())
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}
