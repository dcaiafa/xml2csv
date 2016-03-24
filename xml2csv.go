package main

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"encoding/xml"
	"flag"
	"fmt"
	"github.com/robertkrimen/otto"
	"gopkg.in/xmlpath.v2"
	"io"
	"os"
	"strings"
)

type transformExpression struct {
	XPath string `json:"xpath"`
	Eval  string `json:"eval"`

	xpath *xmlpath.Path
}

type transformColumn struct {
	Name string `json:"name"`
	transformExpression
}

type transform struct {
	Path    string                `json:"path"`
	Include []transformExpression `json:"include"`
	Exclude []transformExpression `json:"exclude"`
	Columns []transformColumn     `json:"columns"`

	path []string
}

func run(in, out, transformPath string) error {
	var err error

	switch {
	case in == "":
		return fmt.Errorf("missing required parameter -i")
	case out == "":
		return fmt.Errorf("missing required parameter -o")
	case transformPath == "":
		return fmt.Errorf("missing required parameter -t")
	}

	transformFile, err := os.Open(transformPath)
	if err != nil {
		return fmt.Errorf("failed to open %v: %v", transformPath, err)
	}
	defer transformFile.Close()

	var t transform
	err = json.NewDecoder(transformFile).Decode(&t)
	if err != nil {
		return fmt.Errorf("failed to parse json in %v: %v", transformPath, err)
	}

	t.path = strings.Split(t.Path, "/")
	if len(t.path) < 2 && t.path[0] != "" {
		return fmt.Errorf("Path field is invalid")
	}
	t.path = t.path[1:]

	if t.path[len(t.path)-1] == "" {
		t.path = t.path[:len(t.path)-1]
		if len(t.path) == 0 {
			return fmt.Errorf("parameter -path is invalid")
		}
	}

	for i := 0; i < len(t.Include); i++ {
		if err = compile(&t.Include[i]); err != nil {
			return err
		}
	}
	for i := 0; i < len(t.Exclude); i++ {
		if err = compile(&t.Exclude[i]); err != nil {
			return err
		}
	}
	for i := 0; i < len(t.Columns); i++ {
		if err = compile(&t.Columns[i].transformExpression); err != nil {
			return err
		}
	}

	fmt.Println(t)

	xmlFile, err := os.Open(in)
	if err != nil {
		return fmt.Errorf("failed to open %v: %v", in, err)
	}
	defer xmlFile.Close()

	var writer io.Writer
	if out != "" {
		outFile, err := os.Create(out)
		if err != nil {
			return fmt.Errorf("failed to create %v: %v", out, err)
		}
		defer outFile.Close()
		writer = outFile
	} else {
		writer = os.Stdout
	}

	headers := make([]string, len(t.Columns))
	for i := 0; i < len(headers); i++ {
		headers[i] = t.Columns[i].Name
	}

	jsvm := otto.New()

	csvWriter := csv.NewWriter(writer)
	csvWriter.Write(headers)

	err = parse(xmlFile, t.path, func(xml []byte) error {
		node, err := xmlpath.Parse(bytes.NewReader(xml))
		if err != nil {
			return fmt.Errorf("failed to parse xml: %v", err)
		}

		include := len(t.Include) == 0
		for i := 0; i < len(t.Include); i++ {
			if include, err = evalBool(jsvm, node, &t.Include[i]); err != nil {
				return fmt.Errorf("failed to eval include expression: %v", err)
			}
			if include {
				break
			}
		}
		if !include {
			return nil
		}

		for i := 0; i < len(t.Exclude); i++ {
			exclude, err := evalBool(jsvm, node, &t.Exclude[i])
			if err != nil {
				return fmt.Errorf("failed to eval exclude expression: %v", err)
			}
			if exclude {
				return nil
			}
		}

		values := make([]string, len(t.Columns))
		for i := 0; i < len(t.Columns); i++ {
			value, err := evalString(jsvm, node, &t.Columns[i].transformExpression)
			if err != nil {
				return fmt.Errorf("failed to eval column expression: %v", err)
			}
			values[i] = strings.TrimSpace(value)
		}

		err = csvWriter.Write(values)
		if err != nil {
			return fmt.Errorf("failed to write to csv file: %v", err)
		}
		return nil
	})

	csvWriter.Flush()

	if err != nil {
		return err
	}

	return nil
}

func compile(expr *transformExpression) error {
	var err error
	expr.xpath, err = xmlpath.Compile(expr.XPath)
	if err != nil {
		return fmt.Errorf("Failed to compile xpath expression '%v': %v", expr.XPath, err)
	}
	return nil
}

func eval(jsvm *otto.Otto, node *xmlpath.Node, expr *transformExpression) (otto.Value, error) {
	emptyString, _ := otto.ToValue("")

	nodeValue, ok := expr.xpath.String(node)
	if !ok {
		return emptyString, nil
	}

	if expr.Eval == "" {
		stringVal, _ := otto.ToValue(nodeValue)
		return stringVal, nil
	}

	jsvm.Set("v", nodeValue)
	exprValue, err := jsvm.Run(expr.Eval)
	if err != nil {
		return emptyString,
			fmt.Errorf("failed to evaluate javascript expr '%v': %v", expr.Eval, err)
	}
	return exprValue, nil
}

func evalBool(jsvm *otto.Otto, node *xmlpath.Node, expr *transformExpression) (bool, error) {
	if expr.Eval == "" {
		return expr.xpath.Exists(node), nil
	}

	exprValue, err := eval(jsvm, node, expr)
	if err != nil {
		return false, err
	}

	boolValue, err := exprValue.ToBoolean()
	if err != nil {
		return false,
			fmt.Errorf("result of javascript expr '%v' doesn't evaluate to bool", expr.Eval)
	}

	return boolValue, nil
}

func evalString(jsvm *otto.Otto, node *xmlpath.Node, expr *transformExpression) (string, error) {
	exprValue, err := eval(jsvm, node, expr)
	if err != nil {
		return "", err
	}

	stringValue, err := exprValue.ToString()
	if err != nil {
		return "",
			fmt.Errorf("result of javascript expr '%v' doesn't evaluate to string", expr.Eval)
	}

	return stringValue, nil
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
		in        = flag.String("i", "", "")
		out       = flag.String("o", "", "")
		transform = flag.String("t", "", "")
	)

	flag.Parse()

	err := run(*in, *out, *transform)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}
