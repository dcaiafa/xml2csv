package parser

type Op int

const (
	OpAdd Op = iota
	OpSub
	OpMul
	OpDiv
	OpMod
)

type Order int

const (
	Prefix Order = iota
	Postfix
)

type Visitor func(order Order, ast AST)

type AST interface {
	Visit(visitor Visitor)
}

type Expr interface {
	AST
}

type Program struct {
	Transforms []Transform
}

func (a *Program) Visit(visitor Visitor) {
	visitor(Prefix, a)
	for _, t := range a.Transforms {
		t.Visit(visitor)
	}
	visitor(Postfix, a)
}

type Transform struct {
	From   PathExpr
	Where  Expr
	Select []Column
}

func (a *Transform) Visit(visitor Visitor) {
	visitor(Prefix, a)
	a.From.Visit(visitor)
	a.Where.Visit(visitor)
	for _, c := range a.Select {
		c.Visit(visitor)
	}
	visitor(Postfix, a)
}

type PathExpr struct {
	Path []string
}

func (a *PathExpr) Visit(visitor Visitor) {
	visitor(Prefix, a)
}

type Column struct {
	Name string
	Expr Expr
}

func (a *Column) Visit(visitor Visitor) {
	visitor(Prefix, a)
	a.Expr.Visit(visitor)
	visitor(Postfix, a)
}

type BinaryExpr struct {
	LHS, RHS Expr
	Op       Op
}

func (a *BinaryExpr) Visit(visitor Visitor) {
	visitor(Prefix, a)
	a.LHS.Visit(visitor)
	a.RHS.Visit(visitor)
	visitor(Postfix, a)
}

type CallExpr struct {
	FuncName int
	Args     []Expr
}

func (a *CallExpr) Visit(visitor Visitor) {
	visitor(Prefix, a)
	for _, c := range a.Args {
		c.Visit(visitor)
	}
	visitor(Postfix, a)
}

type LiteralExpr struct {
	Val interface{}
}

func (a *LiteralExpr) Visit(visitor Visitor) {
	visitor(Prefix, a)
}
