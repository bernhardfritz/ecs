package ecs

type IsComponentPointer interface {
	set(value uint)
}

type IsComponent interface {
	id() uint
	flag() uint
	createColumn() isColumn
}

type Component[T any] uint

func (c Component[T]) Get(w *World, entity Entity) *T {
	arch := w.entityArchetype[entity]
	if !arch.has(c) {
		return nil
	}
	columnIndex := countLowerSetBits(arch.id, c.flag())

	return arch.columns[columnIndex].get(arch.entities.indexOf(entity)).(*T)
}

func (c *Component[T]) set(value uint) {
	*c = Component[T](value)
}

func (c Component[T]) id() uint {
	return uint(c)
}

func (c Component[T]) flag() uint {
	return 1 << uint(c)
}

func (c Component[T]) createColumn() isColumn {
	return newColumn[T]()
}
