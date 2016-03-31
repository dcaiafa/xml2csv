package parser

import (
	"bytes"
	"math"
	"strings"
	"testing"
)

func newLexer(in string) *Lexer {
	r := bytes.NewReader([]byte(in))
	n := &Names{}
	l := NewLexer("test", r, n)
	return l
}

func expectToken(t *testing.T, l *Lexer, tok int) (*yySymType, bool) {
	var val yySymType
	read := l.Lex(&val)
	if tok != read {
		t.Errorf("expected token %v, actual %v", tok, read)
		return &val, false
	}
	return &val, true
}

func expectError(t *testing.T, l *Lexer) {
	var val yySymType
	read := l.Lex(&val)
	if read != eof {
		t.Errorf("expected error, actual token %v", read)
		return
	}
}

func expectID(t *testing.T, l *Lexer, id string) {
	val, ok := expectToken(t, l, ID)
	if !ok {
		return
	}
	name := l.Names.NameFromID(val.val.ID)
	if name != id {
		t.Errorf("expected ID %v, actual %v", id, name)
	}
}

func expectNUM(t *testing.T, l *Lexer, num int64) {
	val, ok := expectToken(t, l, NUM)
	if !ok {
		return
	}
	if num != val.val.Num {
		t.Errorf("expected NUM %v, actual %v", num, val.val.Num)
	}
}

func expectSTR(t *testing.T, l *Lexer, str string) {
	val, ok := expectToken(t, l, STR)
	if !ok {
		return
	}
	if str != val.val.Str {
		t.Errorf("expected STR %v, actual %v", str, val.val.Str)
	}
}

func expectPATH(t *testing.T, l *Lexer, path string) {
	val, ok := expectToken(t, l, PATH)
	if !ok {
		return
	}
	pathStr := "/" + strings.Join(val.val.Path, "/")
	if path != pathStr {
		t.Errorf("expected PATH '%v', actual '%v'", path, pathStr)
	}
}

func TestLexerID(t *testing.T) {
	l := newLexer(" a Z az ZA _b")
	expectID(t, l, "a")
	expectID(t, l, "Z")
	expectID(t, l, "az")
	expectID(t, l, "ZA")
	expectID(t, l, "_b")
	expectToken(t, l, eof)
}

func TestLexerNUM(t *testing.T) {
	l := newLexer(" 1234 5 9223372036854775807")
	expectNUM(t, l, 1234)
	expectNUM(t, l, 5)
	expectNUM(t, l, math.MaxInt64)
}

func TestLexerSTR(t *testing.T) {
	l := newLexer(`"a b c" "" "a\nb" "\\\"\n\t"`)
	expectSTR(t, l, "a b c")
	expectSTR(t, l, "")
	expectSTR(t, l, "a\nb")
	expectSTR(t, l, "\\\"\n\t")

	l = newLexer(`"abc`)
	expectError(t, l)
}

func TestLexerPATH(t *testing.T) {
	l := newLexer(" /abc/123/@attrib / /a")
	expectPATH(t, l, "/abc/123/@attrib")
	expectPATH(t, l, "/")
	expectPATH(t, l, "/a")

	l = newLexer("/a/")
	expectError(t, l)
}

func TestLexerTokens(t *testing.T) {
	l := newLexer(":(),+-*//%==&&||<> <= >=")
	expectToken(t, l, ':')
	expectToken(t, l, '(')
	expectToken(t, l, ')')
	expectToken(t, l, ',')
	expectToken(t, l, '+')
	expectToken(t, l, '-')
	expectToken(t, l, '*')
	expectToken(t, l, DIVOP)
	expectToken(t, l, '%')
	expectToken(t, l, EQ)
	expectToken(t, l, AND)
	expectToken(t, l, OR)
	expectToken(t, l, '<')
	expectToken(t, l, '>')
	expectToken(t, l, LE)
	expectToken(t, l, GE)
}

func TestLexerKeywords(t *testing.T) {
	l := newLexer("foreach where select")
	expectToken(t, l, FOREACH)
	expectToken(t, l, WHERE)
	expectToken(t, l, SELECT)
}
