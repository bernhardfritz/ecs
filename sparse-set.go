package ecs

import "math/bits"

type sparseSet[T ~uint] struct {
	sparse []int
	dense  []T
}

func newSparseSet[T ~uint](capacity int) *sparseSet[T] {
	sparse := make([]int, max(capacity, 1)) // TODO allow sparseSet to be initialized with 0 capacity, requires custom code for add

	return &sparseSet[T]{
		sparse: sparse,
		dense:  make([]T, 0, capacity),
	}
}

func (s *sparseSet[T]) contains(value T) bool {
	return int(value) < len(s.sparse) && s.sparse[value] != 0
}

func (s *sparseSet[T]) add(value T) {
	if s.contains(value) {
		return
	}
	i := len(s.dense)
	s.dense = append(s.dense, value)
	if uint(value) >= uint(len(s.sparse)) {
		ratio := (uint(len(s.sparse)) + uint(value) - 1) / uint(value)
		n := bits.Len(ratio - 1)
		newLen := uint(value) << n
		newSparse := make([]int, newLen) // grow sparse array by the smallest multiple of 2 in order to fit value
		copy(newSparse, s.sparse)
		s.sparse = newSparse
	}
	s.sparse[value] = i + 1 // intentionally store index + 1
}

func (s *sparseSet[T]) remove(value T) {
	if !s.contains(value) {
		return
	}
	last := s.dense[len(s.dense)-1]
	s.dense[len(s.dense)-1], s.dense[s.indexOf(value)] = s.dense[s.indexOf(value)], s.dense[len(s.dense)-1]
	s.sparse[last], s.sparse[value] = s.sparse[value], s.sparse[last]
	s.dense = s.dense[:len(s.dense)-1]
	s.sparse[value] = 0
}

func (s *sparseSet[T]) indexOf(value T) int {
	return s.sparse[value] - 1 // intentionally return index - 1
}
