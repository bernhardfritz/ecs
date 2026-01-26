package ecs

type isColumn interface {
	add(value any)
	get(index int) any
}

type column[T any] []T

func newColumn[T any]() *column[T] {
	return new(column[T])
}

func (c *column[T]) add(value any) {
	*c = append(*c, value.(T))
}

func (c *column[T]) get(index int) any {
	return &(*c)[index]
}
