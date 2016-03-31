package compiler

import (
	"bytes"
	"fmt"
	"strconv"
)

func nfaToDot(nfa *NFA, symToStr func(sym int) string) string {
	if symToStr == nil {
		symToStr = strconv.Itoa
	}
	buf := &bytes.Buffer{}
	fmt.Fprintf(buf, "digraph G {")
	fmt.Fprintf(buf, "  rankdir = \"LR\"; ")
	for stateID := 0; stateID < len(nfa.States); stateID++ {
		state := &nfa.States[stateID]

		shape := "circle"
		if state.Accepting {
			shape = "doublecircle"
		}
		fmt.Fprintf(buf, "  %v [shape=\"%v\"]; ", stateID, shape)

		for sym, destStateIDs := range state.Trans {
			var symName string
			if sym == Epsilon {
				symName = "-eps-"
			} else {
				symName = symToStr(sym)
			}

			for _, destStateID := range destStateIDs {
				fmt.Fprintf(buf, "  %v -> %v [label=\"%v\"]; ",
					stateID, destStateID, symName)
			}
		}
	}

	fmt.Fprintf(buf, "}")
	return buf.String()
}

func buildTestNFA() *NFA {
	nfa := &NFA{}

	nfa.AddStates(11)

	nfa.AddTransition(0, 1, Epsilon)
	nfa.AddTransition(0, 7, Epsilon)
	nfa.AddTransition(1, 2, Epsilon)
	nfa.AddTransition(2, 3, 10)
	nfa.AddTransition(3, 6, Epsilon)
	nfa.AddTransition(1, 4, Epsilon)
	nfa.AddTransition(4, 5, 20)
	nfa.AddTransition(5, 6, Epsilon)
	nfa.AddTransition(6, 1, Epsilon)
	nfa.AddTransition(6, 7, Epsilon)
	nfa.AddTransition(7, 8, 10)
	nfa.AddTransition(8, 9, 20)
	nfa.AddTransition(9, 10, 20)

	nfa.SetAccepting(10)

	return nfa
}
