package compiler

import (
	"bytes"
	"encoding/binary"
	"fmt"
)
import "sort"

type intset struct {
	key  string
	ints map[int]struct{}
}

func (s *intset) Reset() {
	s.key = ""
	s.ints = nil
}

func (s *intset) Add(i int) {
	s.key = ""
	if s.ints == nil {
		s.ints = make(map[int]struct{})
	}
	s.ints[i] = struct{}{}
}

func (s *intset) Has(i int) bool {
	if s.ints == nil {
		return false
	}
	_, ok := s.ints[i]
	return ok
}

func (s *intset) Len() int {
	return len(s.ints)
}

func (s *intset) ForEach(f func(i int)) {
	for i, _ := range s.ints {
		f(i)
	}
}

func (s *intset) Key() string {
	if s.key == "" {
		ints := make([]int, 0, len(s.ints))
		for i, _ := range s.ints {
			ints = append(ints, i)
		}

		sort.Sort(sort.IntSlice(ints))

		buf := make([]byte, len(ints)*binary.MaxVarintLen32)
		size := 0
		for _, i := range ints {
			size += binary.PutVarint(buf[size:], int64(i))
		}
		s.key = string(buf[:size])
	}
	return s.key
}

func (s *intset) Items() []int {
	ints := make([]int, 0, len(s.ints))
	for i, _ := range s.ints {
		ints = append(ints, i)
	}
	sort.Ints(ints)
	return ints
}

func (s *intset) String() string {
	var buf bytes.Buffer

	fmt.Fprintf(&buf, "{")
	sep := ""
	for _, i := range s.Items() {
		fmt.Fprintf(&buf, "%v%v", sep, i)
		sep = " "
	}
	fmt.Fprintf(&buf, "}")

	return buf.String()
}
