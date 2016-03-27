//go:generate -command yacc go tool yacc
//go:generate yacc -o parser.go parser.y

package parser

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"strconv"
)

const eof = 0

type lexerState int

const (
	stateStart lexerState = iota
	stateId
	stateNum
	stateStr
	stateStrEscape
	stateComment
	stateFirstSlash
	statePath
)

type TokenValue struct {
	Num  int64
	Str  string
	Path []string
	ID   int
}

type Lexer struct {
	FileName string
	Reader   *bufio.Reader
	Names    *Names

	keywords map[string]int
	buf      bytes.Buffer
	line     int
}

func NewLexer(fileName string, reader io.Reader, names *Names) *Lexer {
	keywords := make(map[string]int)
	keywords["from"] = FROM
	keywords["where"] = WHERE
	keywords["select"] = SELECT

	return &Lexer{
		FileName: fileName,
		Reader:   bufio.NewReader(reader),
		Names:    names,
		keywords: keywords,
		line:     1,
	}
}

func isStartIdRune(r rune) bool {
	return (r >= 'A' && r <= 'Z') || (r >= 'a' && r <= 'z') || r == '_'
}

func isNumRune(r rune) bool {
	return r >= '0' && r <= '9'
}

func (l *Lexer) Lex(yyval *yySymType) int {
	var path []string
	l.buf.Reset()

	state := stateStart
	var tok int
L:
	for {
		r, _, err := l.Reader.ReadRune()
		if err != nil {
			if err != io.EOF {
				l.Error(fmt.Sprintf("failed to read rune: %v", err))
				return eof
			}
			r = 0
		}

		switch state {
		case stateStart:

			switch r {
			case ' ', '\t', '\r':
				break

			case '\n':
				l.line++

			case 0:
				tok = eof
				break L

			case '/':
				state = stateFirstSlash

			case ':', '(', ')', ',', '+', '-', '*', '%':
				tok = int(r)
				break L

			case '"':
				state = stateStr

			default:
				if isStartIdRune(r) {
					l.buf.WriteRune(r)
					state = stateId
				} else if r >= '0' && r <= '9' {
					l.buf.WriteRune(r)
					state = stateNum
				} else {
					l.Error(fmt.Sprintf("unexpected char %v", strconv.QuoteRune(r)))
				}
			}

		case stateId:
			if isStartIdRune(r) || (r >= '0' && r <= '9') {
				l.buf.WriteRune(r)
			} else {
				if r != 0 {
					l.Reader.UnreadRune()
				}
				id := l.buf.String()
				tok = l.keywords[id]
				if tok == 0 {
					yyval.val = &TokenValue{ID: l.Names.IDFromName(id)}
					tok = ID
				}
				break L
			}

		case stateNum:
			if r >= '0' && r <= '9' {
				l.buf.WriteRune(r)
			} else {
				if r != 0 {
					l.Reader.UnreadRune()
				}
				yyval.val = &TokenValue{}
				yyval.val.Num, err = strconv.ParseInt(l.buf.String(), 10, 64)
				if err != nil {
					l.Error(fmt.Sprintf("failed to parser number: %v", err))
					return eof
				}
				tok = NUM
				break L
			}

		case stateStr:
			if r == '\\' {
				state = stateStrEscape
			} else if r == '"' {
				yyval.val = &TokenValue{Str: l.buf.String()}
				tok = STR
				break L
			} else if strconv.IsPrint(r) {
				l.buf.WriteRune(r)
			} else {
				if r == eof {
					l.Error("closing quotes missing in string literal")
				} else {
					l.Error(fmt.Sprintf("unexpected character %v", strconv.QuoteRune(r)))
				}
				return eof
			}

		case stateStrEscape:
			switch r {
			case '"':
			case '\\':
			case 'n':
				r = '\n'
			case 't':
				r = '\t'
			default:
				l.Error(fmt.Sprintf("unexpected character %v in string literal",
					strconv.QuoteRune(r)))
				return eof
			}
			l.buf.WriteRune(r)
			state = stateStr

		case stateFirstSlash:
			if r == '/' {
				tok = DIVOP
				break L
			} else {
				if r != 0 {
					l.Reader.UnreadRune()
				}
				path = path[:0]
				state = statePath
			}

		case statePath:
			if isStartIdRune(r) || isNumRune(r) || r == '@' {
				l.buf.WriteRune(r)
			} else {
				if l.buf.Len() != 0 {
					path = append(path, l.buf.String())
					l.buf.Reset()
				} else if len(path) > 0 {
					// The path '/' is valid, but '/a/' is invalid.
					l.Error("invalid path")
					return eof
				}

				if r != '/' {
					if r != 0 {
						l.Reader.UnreadRune()
					}
					yyval.val = &TokenValue{Path: path}

					tok = PATH
					break L
				}
			}
		}
	}

	return tok
}

func (l *Lexer) Error(s string) {
	fmt.Fprintf(os.Stderr, "%v:%v: %v\n", l.FileName, l.line, s)
}
