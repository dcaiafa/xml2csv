package main

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"encoding/xml"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"runtime"
	"runtime/pprof"
	"strings"
	"sync"
	"time"

	"github.com/robertkrimen/otto"
	"gopkg.in/xmlpath.v2"
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

func run(in, out, transformPath string, append bool) error {
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

	t := &transform{}
	err = json.NewDecoder(transformFile).Decode(t)
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

	xmlFile, err := os.Open(in)
	if err != nil {
		return fmt.Errorf("failed to open %v: %v", in, err)
	}
	defer xmlFile.Close()

	var writer io.Writer
	if out != "" {
		outFileFlags := os.O_CREATE
		if append {
			outFileFlags |= os.O_APPEND
		} else {
			outFileFlags |= os.O_TRUNC
		}
		outFile, err := os.OpenFile(out, outFileFlags, 644)
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

	csvWriter := csv.NewWriter(writer)

	if !append {
		csvWriter.Write(headers)
	}

	wg := &sync.WaitGroup{}

	processorCount := runtime.NumCPU()
	xmlChan := make(chan []byte, processorCount)
	csvChan := make(chan []string)
	errChan := make(chan error)

	cancelChan := make(chan struct{})

	wg.Add(processorCount)
	for i := 0; i < processorCount; i++ {
		go process(t, xmlChan, csvChan, errChan, wg)
	}

	wg.Add(1)
	go func() {
		defer wg.Done()
		err = parse(xmlFile, t.path, func(xml []byte) error {
			xmlCopy := make([]byte, len(xml))
			copy(xmlCopy, xml)

			xmlChan <- xmlCopy

			select {
			case <-cancelChan:
				return fmt.Errorf("cancelled")
			default:
				return nil
			}
		})

		if err != nil {
			errChan <- err
		}
		close(xmlChan)
	}()

	go func() {
		wg.Wait()
		close(csvChan)
	}()

loop:
	for {
		select {
		case err = <-errChan:
			close(cancelChan)
			break loop

		case csvLine, ok := <-csvChan:
			if !ok {
				break loop
			}
			err = csvWriter.Write(csvLine)
			if err != nil {
				break loop
			}
		}
	}

	if err != nil {
		return err
	}

	csvWriter.Flush()

	return nil
}

func process(t *transform, xmlChan <-chan []byte, csvChan chan<- []string, errChan chan error, wg *sync.WaitGroup) {
	defer wg.Done()

	jsvm := otto.New()

loop:
	for xml := range xmlChan {
		node, err := xmlpath.Parse(bytes.NewReader(xml))
		if err != nil {
			errChan <- fmt.Errorf("failed to parse xml: %v", err)
			return
		}

		include := len(t.Include) == 0
		for i := 0; i < len(t.Include); i++ {
			if include, err = evalBool(jsvm, node, &t.Include[i]); err != nil {
				errChan <- fmt.Errorf("failed to eval include expression: %v", err)
				return
			}
			if include {
				break
			}
		}
		if !include {
			continue loop
		}

		for i := 0; i < len(t.Exclude); i++ {
			exclude, err := evalBool(jsvm, node, &t.Exclude[i])
			if err != nil {
				errChan <- fmt.Errorf("failed to eval exclude expression: %v", err)
				return
			}
			if exclude {
				continue loop
			}
		}

		values := make([]string, len(t.Columns))
		for i := 0; i < len(t.Columns); i++ {
			value, err := evalString(jsvm, node, &t.Columns[i].transformExpression)
			if err != nil {
				errChan <- fmt.Errorf("failed to eval column expression: %v", err)
				return
			}
			values[i] = strings.TrimSpace(value)
		}

		csvChan <- values
	}
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
		in         = flag.String("i", "", "")
		out        = flag.String("o", "", "")
		transform  = flag.String("t", "", "")
		append     = flag.Bool("a", false, "")
		cpuprofile = flag.String("cpuprofile", "", "write cpu profile to file")
		memprofile = flag.String("memprofile", "", "write memory profile to file")
	)

	flag.Parse()

	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			log.Fatal(err)
		}
		defer f.Close()
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}

	if *memprofile != "" {
		go func() {
			time.Sleep(5 * time.Second)
			f, err := os.Create(*memprofile)
			if err != nil {
				log.Fatal(err)
			}
			defer f.Close()
			pprof.Lookup("heap").WriteTo(f, 0)
		}()
	}

	err := run(*in, *out, *transform, *append)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}
