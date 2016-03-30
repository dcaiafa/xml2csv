package compiler

type nfa2dfa struct {
	// The reference NFA.
	nfa *NFA

	// Map of subset key to subset.
	subsetMap map[string]*intset

	// All subsets in order of creation.
	subsetList []*intset

	// Maps subset key to map of symbols to destination subset key.
	trans map[string]map[int]string

	subsetToDFA    map[string]int
	startSubset    *intset
	closurePending []int
}

func NFAToDFA(nfa *NFA) *DFA {
	estDFAStateCount := nfa.StateCount() * 2

	n2d := &nfa2dfa{
		nfa:         nfa,
		subsetMap:   make(map[string]*intset, estDFAStateCount),
		trans:       make(map[string]map[int]string, estDFAStateCount),
		subsetToDFA: make(map[string]int),
	}

	n2d.startSubset = &intset{}
	n2d.startSubset.Add(0)
	n2d.closure(n2d.startSubset)

	n2d.subsetMap[n2d.startSubset.Key()] = n2d.startSubset
	n2d.subsetList = append(n2d.subsetList, n2d.startSubset)
	pending := []*intset{n2d.startSubset}

	for len(pending) > 0 {
		var subset *intset
		subset, pending = pending[len(pending)-1], pending[:len(pending)-1]

		syms := n2d.symbols(subset)
		for _, sym := range syms {
			destSubset := n2d.move(subset, sym)
			n2d.closure(destSubset)

			if _, ok := n2d.subsetMap[destSubset.Key()]; !ok {
				n2d.subsetMap[destSubset.Key()] = destSubset
				n2d.subsetList = append(n2d.subsetList, destSubset)
				pending = append(pending, destSubset)
			}
			if n2d.trans[subset.Key()] == nil {
				n2d.trans[subset.Key()] = make(map[int]string)
			}
			n2d.trans[subset.Key()][sym] = destSubset.Key()
		}
	}

	return n2d.buildDFA()
}

func (n2d *nfa2dfa) buildDFA() *DFA {
	dfa := &DFA{States: make([]DFAState, 0, len(n2d.subsetMap))}

	buildDFAState := func(subset *intset) {
		dfaS := DFAState{
			NFAStates: make([]int, 0, subset.Len()),
		}
		subset.ForEach(func(nfaState int) {
			dfaS.NFAStates = append(dfaS.NFAStates, nfaState)
			dfaS.Accepting = dfaS.Accepting || n2d.nfa.States[nfaState].Accepting
		})
		dfa.States = append(dfa.States, dfaS)
		n2d.subsetToDFA[subset.Key()] = len(dfa.States) - 1
	}

	for _, subset := range n2d.subsetList {
		buildDFAState(subset)
	}

	for subset, dfaStateIndex := range n2d.subsetToDFA {
		dfaState := &dfa.States[dfaStateIndex]

		for sym, destSubsetKey := range n2d.trans[subset] {
			if dfaState.Trans == nil {
				dfaState.Trans = make(map[int]int)
			}

			destDFAS := n2d.subsetToDFA[n2d.subsetMap[destSubsetKey].Key()]
			dfaState.Trans[sym] = destDFAS
		}
	}

	return dfa
}

func (n2d *nfa2dfa) closure(s *intset) {
	pending := n2d.closurePending
	pending = pending[:0]

	s.ForEach(func(state int) {
		pending = append(pending, state)
	})

	for len(pending) > 0 {
		var state int
		state, pending = pending[len(pending)-1], pending[:len(pending)-1]
		toStates := n2d.nfa.TransitionsFor(state, Epsilon)
		for _, toState := range toStates {
			if !s.Has(toState) {
				s.Add(toState)
				pending = append(pending, toState)
			}
		}
	}
}

func (n2d *nfa2dfa) move(stateSet *intset, sym int) *intset {
	destStates := &intset{}
	stateSet.ForEach(func(state int) {
		for _, destState := range n2d.nfa.States[state].Trans[sym] {
			destStates.Add(destState)
		}
	})
	return destStates
}

func (n2d *nfa2dfa) symbols(s *intset) []int {
	syms := &intset{}
	s.ForEach(func(state int) {
		for sym, _ := range n2d.nfa.States[state].Trans {
			if sym != Epsilon {
				syms.Add(sym)
			}
		}
	})
	return syms.Items()
}
