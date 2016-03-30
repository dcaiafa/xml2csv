package compiler

import (
	"bytes"
	"encoding/binary"
	"math"
	"reflect"
	"testing"
)

func decodeKey(k string) []int {
	var ints []int

	r := bytes.NewReader([]byte(k))

	for {
		i, err := binary.ReadVarint(r)
		if err != nil {
			break
		}
		ints = append(ints, int(i))
	}

	return ints
}

func TestIntSetKey(t *testing.T) {
	var s intset

	s.Add(9524)
	s.Add(1)
	s.Add(9524)
	s.Add(math.MaxInt32)
	s.Add(821974)
	s.Add(1)

	expected := []int{1, 9524, 821974, math.MaxInt32}
	actual := decodeKey(s.Key())
	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("expected decoded key %v, actual %v", expected, actual)
	}

	s.Add(857)
	s.Add(999)

	expected = []int{1, 857, 999, 9524, 821974, math.MaxInt32}
	actual = decodeKey(s.Key())
	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("expected decoded key %v, actual %v", expected, actual)
	}
}
