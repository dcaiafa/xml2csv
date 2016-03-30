package compiler

type DFA struct {
	States []DFAState
}

type DFAState struct {
	NFAStates []int
	Trans     map[int]int
	Accepting bool
}
