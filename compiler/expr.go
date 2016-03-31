package compiler

type exprType int

const (
	exprUnkn exprType = iota
	exprBool
	exprNum
	exprStr
	exprPath
)

type expr struct {
	typ exprType
	val interface{}
}
