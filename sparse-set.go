package ecs

import "math/bits"

type sparseSet struct {
	sparse []int
	dense  []Entity
}

func newSparseSet(capacity int) *sparseSet {
	sparse := make([]int, max(capacity, 1)) // TODO allow sparseSet to be initialized with 0 capacity, requires custom code for add

	return &sparseSet{
		sparse: sparse,
		dense:  make([]Entity, 0, capacity),
	}
}

func (s *sparseSet) contains(e Entity) bool {
	return int(e) < len(s.sparse) && s.sparse[e] != 0
}

func (s *sparseSet) add(e Entity) {
	if s.contains(e) {
		return
	}
	i := len(s.dense)
	s.dense = append(s.dense, e)
	if uint(e) >= uint(len(s.sparse)) {
		ratio := (uint(len(s.sparse)) + uint(e) - 1) / uint(e)
		n := bits.Len(ratio - 1)
		newLen := uint(e) << n
		newSparse := make([]int, newLen) // grow sparse array by the smallest multiple of 2 in order to fit entity e
		copy(newSparse, s.sparse)
		s.sparse = newSparse
	}
	s.sparse[e] = i + 1 // intentionally store index + 1
}

func (s *sparseSet) remove(e Entity) {
	if !s.contains(e) {
		return
	}
	last := s.dense[len(s.dense)-1]
	s.dense[len(s.dense)-1], s.dense[s.indexOf(e)] = s.dense[s.indexOf(e)], s.dense[len(s.dense)-1]
	s.sparse[last], s.sparse[e] = s.sparse[e], s.sparse[last]
	s.dense = s.dense[:len(s.dense)-1]
	s.sparse[e] = 0
}

func (s *sparseSet) indexOf(e Entity) int {
	return s.sparse[e] - 1 // intentionally return index - 1
}
