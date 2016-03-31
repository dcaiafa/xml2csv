package compiler

type transform struct {
	GlobalPath *path
	NFA        NFA
	DFA        *DFA
}
