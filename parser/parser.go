//line parser.y:2
package parser

import __yyfmt__ "fmt"

//line parser.y:2
import (
	"io"
)

//line parser.y:11
type yySymType struct {
	yys int
	val *TokenValue
	res interface{}
}

const FROM = 57346
const WHERE = 57347
const SELECT = 57348
const ID = 57349
const NUM = 57350
const STR = 57351
const PATH = 57352
const DIVOP = 57353

var yyToknames = [...]string{
	"$end",
	"error",
	"$unk",
	"FROM",
	"WHERE",
	"SELECT",
	"ID",
	"NUM",
	"STR",
	"PATH",
	"':'",
	"'('",
	"')'",
	"','",
	"'+'",
	"'-'",
	"'*'",
	"DIVOP",
	"'%'",
}
var yyStatenames = [...]string{}

const yyEofCode = 1
const yyErrCode = 2
const yyInitialStackSize = 16

//line parser.y:208

func Parse(fileName string, reader io.Reader, names *Names) {
	yyErrorVerbose = true
	l := NewLexer(fileName, reader, names)
	yyParse(l)
}

//line yacctab:1
var yyExca = [...]int{
	-1, 1,
	1, -1,
	-2, 0,
}

const yyNprod = 33
const yyPrivate = 57344

var yyTokenNames []string
var yyStates []string

const yyLast = 56

var yyAct = [...]int{

	13, 24, 44, 49, 28, 29, 30, 31, 32, 28,
	29, 30, 31, 32, 50, 19, 21, 22, 20, 33,
	18, 35, 34, 38, 37, 27, 36, 26, 12, 39,
	40, 41, 42, 43, 10, 47, 9, 48, 5, 3,
	14, 46, 6, 45, 15, 17, 16, 25, 23, 11,
	7, 51, 8, 4, 2, 1,
}
var yyPact = [...]int{

	34, -1000, 34, -1000, 31, 24, -1000, 22, -1000, 8,
	-1000, -1000, 18, -6, -1000, -1000, -1000, -1000, 8, 10,
	-1000, -1000, -1000, 7, -1000, 8, 13, 12, 8, 8,
	8, 8, 8, -11, 8, 18, -6, -1000, -1000, -1000,
	-1000, -1000, -1000, -1000, -1000, -10, 0, -6, -1000, -1000,
	8, -6,
}
var yyPgo = [...]int{

	0, 55, 54, 39, 53, 52, 50, 49, 48, 1,
	47, 0, 46, 45, 44, 43, 41, 40,
}
var yyR1 = [...]int{

	0, 1, 2, 2, 3, 4, 5, 6, 6, 7,
	8, 8, 9, 10, 10, 11, 11, 11, 11, 11,
	12, 13, 13, 14, 15, 15, 16, 16, 17, 17,
	17, 17, 17,
}
var yyR2 = [...]int{

	0, 1, 2, 1, 3, 2, 2, 1, 0, 2,
	3, 1, 2, 2, 2, 1, 1, 1, 1, 3,
	1, 1, 1, 4, 1, 0, 3, 1, 3, 3,
	3, 3, 3,
}
var yyChk = [...]int{

	-1000, -1, -2, -3, -4, 4, -3, -6, -5, 5,
	10, -7, 6, -11, -17, -14, -12, -13, 12, 7,
	10, 8, 9, -8, -9, -10, 9, 7, 15, 16,
	17, 18, 19, -11, 12, 14, -11, 11, 11, -11,
	-11, -11, -11, -11, 13, -15, -16, -11, -9, 13,
	14, -11,
}
var yyDef = [...]int{

	0, -2, 1, 3, 8, 0, 2, 0, 7, 0,
	5, 4, 0, 6, 15, 16, 17, 18, 0, 0,
	20, 21, 22, 9, 11, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 25, 0, 12, 13, 14, 28,
	29, 30, 31, 32, 19, 0, 24, 27, 10, 23,
	0, 26,
}
var yyTok1 = [...]int{

	1, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 19, 3, 3,
	12, 13, 17, 15, 14, 16, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 11,
}
var yyTok2 = [...]int{

	2, 3, 4, 5, 6, 7, 8, 9, 10, 18,
}
var yyTok3 = [...]int{
	0,
}

var yyErrorMessages = [...]struct {
	state int
	token int
	msg   string
}{}

//line yaccpar:1

/*	parser for yacc output	*/

var (
	yyDebug        = 0
	yyErrorVerbose = false
)

type yyLexer interface {
	Lex(lval *yySymType) int
	Error(s string)
}

type yyParser interface {
	Parse(yyLexer) int
	Lookahead() int
}

type yyParserImpl struct {
	lval  yySymType
	stack [yyInitialStackSize]yySymType
	char  int
}

func (p *yyParserImpl) Lookahead() int {
	return p.char
}

func yyNewParser() yyParser {
	return &yyParserImpl{}
}

const yyFlag = -1000

func yyTokname(c int) string {
	if c >= 1 && c-1 < len(yyToknames) {
		if yyToknames[c-1] != "" {
			return yyToknames[c-1]
		}
	}
	return __yyfmt__.Sprintf("tok-%v", c)
}

func yyStatname(s int) string {
	if s >= 0 && s < len(yyStatenames) {
		if yyStatenames[s] != "" {
			return yyStatenames[s]
		}
	}
	return __yyfmt__.Sprintf("state-%v", s)
}

func yyErrorMessage(state, lookAhead int) string {
	const TOKSTART = 4

	if !yyErrorVerbose {
		return "syntax error"
	}

	for _, e := range yyErrorMessages {
		if e.state == state && e.token == lookAhead {
			return "syntax error: " + e.msg
		}
	}

	res := "syntax error: unexpected " + yyTokname(lookAhead)

	// To match Bison, suggest at most four expected tokens.
	expected := make([]int, 0, 4)

	// Look for shiftable tokens.
	base := yyPact[state]
	for tok := TOKSTART; tok-1 < len(yyToknames); tok++ {
		if n := base + tok; n >= 0 && n < yyLast && yyChk[yyAct[n]] == tok {
			if len(expected) == cap(expected) {
				return res
			}
			expected = append(expected, tok)
		}
	}

	if yyDef[state] == -2 {
		i := 0
		for yyExca[i] != -1 || yyExca[i+1] != state {
			i += 2
		}

		// Look for tokens that we accept or reduce.
		for i += 2; yyExca[i] >= 0; i += 2 {
			tok := yyExca[i]
			if tok < TOKSTART || yyExca[i+1] == 0 {
				continue
			}
			if len(expected) == cap(expected) {
				return res
			}
			expected = append(expected, tok)
		}

		// If the default action is to accept or reduce, give up.
		if yyExca[i+1] != 0 {
			return res
		}
	}

	for i, tok := range expected {
		if i == 0 {
			res += ", expecting "
		} else {
			res += " or "
		}
		res += yyTokname(tok)
	}
	return res
}

func yylex1(lex yyLexer, lval *yySymType) (char, token int) {
	token = 0
	char = lex.Lex(lval)
	if char <= 0 {
		token = yyTok1[0]
		goto out
	}
	if char < len(yyTok1) {
		token = yyTok1[char]
		goto out
	}
	if char >= yyPrivate {
		if char < yyPrivate+len(yyTok2) {
			token = yyTok2[char-yyPrivate]
			goto out
		}
	}
	for i := 0; i < len(yyTok3); i += 2 {
		token = yyTok3[i+0]
		if token == char {
			token = yyTok3[i+1]
			goto out
		}
	}

out:
	if token == 0 {
		token = yyTok2[1] /* unknown char */
	}
	if yyDebug >= 3 {
		__yyfmt__.Printf("lex %s(%d)\n", yyTokname(token), uint(char))
	}
	return char, token
}

func yyParse(yylex yyLexer) int {
	return yyNewParser().Parse(yylex)
}

func (yyrcvr *yyParserImpl) Parse(yylex yyLexer) int {
	var yyn int
	var yyVAL yySymType
	var yyDollar []yySymType
	_ = yyDollar // silence set and not used
	yyS := yyrcvr.stack[:]

	Nerrs := 0   /* number of errors */
	Errflag := 0 /* error recovery flag */
	yystate := 0
	yyrcvr.char = -1
	yytoken := -1 // yyrcvr.char translated into internal numbering
	defer func() {
		// Make sure we report no lookahead when not parsing.
		yystate = -1
		yyrcvr.char = -1
		yytoken = -1
	}()
	yyp := -1
	goto yystack

ret0:
	return 0

ret1:
	return 1

yystack:
	/* put a state and value onto the stack */
	if yyDebug >= 4 {
		__yyfmt__.Printf("char %v in %v\n", yyTokname(yytoken), yyStatname(yystate))
	}

	yyp++
	if yyp >= len(yyS) {
		nyys := make([]yySymType, len(yyS)*2)
		copy(nyys, yyS)
		yyS = nyys
	}
	yyS[yyp] = yyVAL
	yyS[yyp].yys = yystate

yynewstate:
	yyn = yyPact[yystate]
	if yyn <= yyFlag {
		goto yydefault /* simple state */
	}
	if yyrcvr.char < 0 {
		yyrcvr.char, yytoken = yylex1(yylex, &yyrcvr.lval)
	}
	yyn += yytoken
	if yyn < 0 || yyn >= yyLast {
		goto yydefault
	}
	yyn = yyAct[yyn]
	if yyChk[yyn] == yytoken { /* valid shift */
		yyrcvr.char = -1
		yytoken = -1
		yyVAL = yyrcvr.lval
		yystate = yyn
		if Errflag > 0 {
			Errflag--
		}
		goto yystack
	}

yydefault:
	/* default state action */
	yyn = yyDef[yystate]
	if yyn == -2 {
		if yyrcvr.char < 0 {
			yyrcvr.char, yytoken = yylex1(yylex, &yyrcvr.lval)
		}

		/* look through exception table */
		xi := 0
		for {
			if yyExca[xi+0] == -1 && yyExca[xi+1] == yystate {
				break
			}
			xi += 2
		}
		for xi += 2; ; xi += 2 {
			yyn = yyExca[xi+0]
			if yyn < 0 || yyn == yytoken {
				break
			}
		}
		yyn = yyExca[xi+1]
		if yyn < 0 {
			goto ret0
		}
	}
	if yyn == 0 {
		/* error ... attempt to resume parsing */
		switch Errflag {
		case 0: /* brand new error */
			yylex.Error(yyErrorMessage(yystate, yytoken))
			Nerrs++
			if yyDebug >= 1 {
				__yyfmt__.Printf("%s", yyStatname(yystate))
				__yyfmt__.Printf(" saw %s\n", yyTokname(yytoken))
			}
			fallthrough

		case 1, 2: /* incompletely recovered error ... try again */
			Errflag = 3

			/* find a state where "error" is a legal shift action */
			for yyp >= 0 {
				yyn = yyPact[yyS[yyp].yys] + yyErrCode
				if yyn >= 0 && yyn < yyLast {
					yystate = yyAct[yyn] /* simulate a shift of "error" */
					if yyChk[yystate] == yyErrCode {
						goto yystack
					}
				}

				/* the current p has no shift on "error", pop stack */
				if yyDebug >= 2 {
					__yyfmt__.Printf("error recovery pops state %d\n", yyS[yyp].yys)
				}
				yyp--
			}
			/* there is no state on the stack with an error shift ... abort */
			goto ret1

		case 3: /* no shift yet; clobber input char */
			if yyDebug >= 2 {
				__yyfmt__.Printf("error recovery discards %s\n", yyTokname(yytoken))
			}
			if yytoken == yyEofCode {
				goto ret1
			}
			yyrcvr.char = -1
			yytoken = -1
			goto yynewstate /* try again in the same state */
		}
	}

	/* reduction by production yyn */
	if yyDebug >= 2 {
		__yyfmt__.Printf("reduce %v in:\n\t%v\n", yyn, yyStatname(yystate))
	}

	yynt := yyn
	yypt := yyp
	_ = yypt // guard against "declared and not used"

	yyp -= yyR2[yyn]
	// yyp is now the index of $0. Perform the default action. Iff the
	// reduced production is ε, $1 is possibly out of range.
	if yyp+1 >= len(yyS) {
		nyys := make([]yySymType, len(yyS)*2)
		copy(nyys, yyS)
		yyS = nyys
	}
	yyVAL = yyS[yyp+1]

	/* consult goto table to find next state */
	yyn = yyR1[yyn]
	yyg := yyPgo[yyn]
	yyj := yyg + yyS[yyp].yys + 1

	if yyj >= yyLast {
		yystate = yyAct[yyg]
	} else {
		yystate = yyAct[yyj]
		if yyChk[yystate] != -yyn {
			yystate = yyAct[yyg]
		}
	}
	// dummy call; replaced with literal code
	switch yynt {

	case 1:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:32
		{
			yyVAL.res = yyDollar[1].res
		}
	case 2:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:38
		{
			yyVAL.res = append(yyDollar[1].res.([]Transform), yyDollar[2].res.(Transform))
		}
	case 3:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:42
		{
			yyVAL.res = []Transform{yyDollar[1].res.(Transform)}
		}
	case 4:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:48
		{
			yyVAL.res = &Transform{
				From:   yyDollar[1].res.(PathExpr),
				Where:  yyDollar[2].res.(Expr),
				Select: yyDollar[3].res.([]Column),
			}
		}
	case 5:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:58
		{
			yyVAL.res = &PathExpr{Path: yyDollar[2].val.Path}
		}
	case 6:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:64
		{
			yyVAL.res = yyDollar[2].res
		}
	case 7:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:70
		{
			yyVAL.res = yyDollar[1].res
		}
	case 8:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.y:73
		{
			yyVAL.res = nil
		}
	case 9:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:79
		{
			yyVAL.res = yyDollar[2].res
		}
	case 10:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:85
		{
			yyVAL.res = append(yyDollar[1].res.([]Column), yyDollar[3].res.(Column))
		}
	case 11:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:89
		{
			yyVAL.res = []Column{yyDollar[1].res.(Column)}
		}
	case 12:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:95
		{
			yyVAL.res = &Column{
				Name: yyDollar[1].res.(string),
				Expr: yyDollar[2].res.(Expr),
			}
		}
	case 13:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:104
		{
			yyVAL.res = yyDollar[1].val.Str
		}
	case 14:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:108
		{
			yyVAL.res = yylex.(*Lexer).Names.NameFromID(yyDollar[1].val.ID)
		}
	case 19:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:118
		{
			yyVAL.res = yyDollar[2].res
		}
	case 20:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:124
		{
			yyVAL.res = &PathExpr{Path: yyDollar[1].val.Path}
		}
	case 21:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:130
		{
			yyVAL.res = &LiteralExpr{Val: yyDollar[1].val.Num}
		}
	case 22:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:134
		{
			yyVAL.res = &LiteralExpr{Val: yyDollar[1].val.Str}
		}
	case 23:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parser.y:140
		{
			yyVAL.res = &CallExpr{
				FuncName: yyDollar[1].val.ID,
				Args:     yyDollar[3].res.([]Expr),
			}
		}
	case 24:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:149
		{
			yyVAL.res = yyDollar[1].res
		}
	case 25:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.y:152
		{
			yyVAL.res = nil
		}
	case 26:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:158
		{
			yyVAL.res = append(yyDollar[1].res.([]Expr), yyDollar[3].res.(Expr))
		}
	case 27:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:162
		{
			yyVAL.res = []Expr{yyDollar[1].res.(Expr)}
		}
	case 28:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:168
		{
			yyVAL.res = &BinaryExpr{
				LHS: yyDollar[1].res.(Expr),
				RHS: yyDollar[3].res.(Expr),
				Op:  OpAdd,
			}
		}
	case 29:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:176
		{
			yyVAL.res = &BinaryExpr{
				LHS: yyDollar[1].res.(Expr),
				RHS: yyDollar[3].res.(Expr),
				Op:  OpSub,
			}
		}
	case 30:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:184
		{
			yyVAL.res = &BinaryExpr{
				LHS: yyDollar[1].res.(Expr),
				RHS: yyDollar[3].res.(Expr),
				Op:  OpMul,
			}
		}
	case 31:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:192
		{
			yyVAL.res = &BinaryExpr{
				LHS: yyDollar[1].res.(Expr),
				RHS: yyDollar[3].res.(Expr),
				Op:  OpDiv,
			}
		}
	case 32:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:200
		{
			yyVAL.res = &BinaryExpr{
				LHS: yyDollar[1].res.(Expr),
				RHS: yyDollar[3].res.(Expr),
				Op:  OpMod,
			}
		}
	}
	goto yystack /* stack new state and value */
}
