package compiler

const Epsilon int = 0

type NFA struct {
	States []NFAState
}

type NFAState struct {
	Trans     map[int][]int
	Accepting bool
}

func (n *NFA) AddStates(count int) int {
	firstStateID := len(n.States)
	for i := 0; i < count; i++ {
		n.States = append(n.States, NFAState{})
	}
	return firstStateID
}

func (n *NFA) SetAccepting(stateID int) {
	n.States[stateID].Accepting = true
}

func (n *NFA) TransitionsFor(state int, sym int) []int {
	return n.States[state].Trans[sym]
}

func (n *NFA) AddTransition(from, to int, symbol int) {
	fromState := &n.States[from]
	if fromState.Trans == nil {
		fromState.Trans = make(map[int][]int)
	}
	fromState.Trans[symbol] = append(fromState.Trans[symbol], to)
}

func (n *NFA) StateCount() int {
	return len(n.States)
}
