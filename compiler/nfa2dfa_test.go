package compiler

import (
	"bytes"
	"fmt"
	"reflect"
	"sort"
	"strconv"
	"testing"
)

func dfaToDot(dfa *DFA, symName func(sym int) string) string {
	if symName == nil {
		symName = strconv.Itoa
	}

	var buf bytes.Buffer
	buf.WriteString("digraph G { rankdir=\"LR\"; ")
	for stateID, state := range dfa.States {
		shape := "circle"
		if state.Accepting {
			shape = "doublecircle"
		}
		buf.WriteString(fmt.Sprintf("%v [shape=\"%v\"]; ", stateID, shape))
		syms := make([]int, 0, len(state.Trans))
		for sym, _ := range state.Trans {
			syms = append(syms, sym)
		}
		sort.Ints(syms)
		for _, sym := range syms {
			toStateID := state.Trans[sym]
			buf.WriteString(fmt.Sprintf("%v -> %v [label=\"%v\"]; ",
				stateID, toStateID, symName(sym)))
		}
	}
	buf.WriteString("}")
	return buf.String()
}

func expectSet(t *testing.T, s *intset, set ...int) {
	if len(set) != s.Len() {
		t.Errorf("expected %v in set, actual %v", len(set), s.Len())
		return
	}
	for _, item := range set {
		if !s.Has(item) {
			t.Errorf("set doesn't match expected: %v", set)
			t.Logf("actual:")
			s.ForEach(func(item int) {
				t.Logf("  %v", item)
			})
			return
		}
	}
}

func TestNFA2DFAClosure(t *testing.T) {
	nfa := buildTestNFA()

	n2d := &nfa2dfa{nfa: nfa}

	var s intset
	s.Add(0)
	n2d.closure(&s)
	expectSet(t, &s, 0, 1, 2, 4, 7)

	s.Reset()
	s.Add(3)
	s.Add(8)
	n2d.closure(&s)
	expectSet(t, &s, 1, 2, 3, 4, 6, 7, 8)
}

func TestNFA2DFASymbols(t *testing.T) {
	nfa := buildTestNFA()

	n2d := &nfa2dfa{nfa: nfa}

	var s intset
	s.Add(2)
	s.Add(4)
	s.Add(8)

	syms := n2d.symbols(&s)
	expected := []int{10, 20}
	if !reflect.DeepEqual(syms, expected) {
		t.Errorf("expected symbols %v, actual %v", expected, syms)
	}
}

func TestNFA2DFAMove(t *testing.T) {
	nfa := buildTestNFA()

	n2d := &nfa2dfa{nfa: nfa}

	var s intset
	s.Add(0)
	s.Add(1)
	s.Add(2)
	s.Add(4)
	s.Add(7)

	aStates := n2d.move(&s, 10)
	expectSet(t, aStates, 3, 8)

	bStates := n2d.move(&s, 20)
	expectSet(t, bStates, 5)
}

func TestNFA2DFA(t *testing.T) {
	nfa := buildTestNFA()

	dfa := NFAToDFA(nfa)
	dfaDot := dfaToDot(dfa, nil)

	const expected = `digraph G { rankdir="LR"; 0 [shape="circle"]; 0 -> 1 [label="10"]; 0 -> 2 [label="20"]; 1 [shape="circle"]; 1 -> 1 [label="10"]; 1 -> 3 [label="20"]; 2 [shape="circle"]; 2 -> 1 [label="10"]; 2 -> 2 [label="20"]; 3 [shape="circle"]; 3 -> 1 [label="10"]; 3 -> 4 [label="20"]; 4 [shape="doublecircle"]; 4 -> 1 [label="10"]; 4 -> 2 [label="20"]; }`

	if dfaDot != expected {
		t.Errorf("DFA doesn't match baseline.")
		t.Log("Expected:")
		t.Log(expected)
		t.Log("Actual:")
		t.Log(dfaDot)
	}
}
