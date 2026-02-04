package ecs

import "cmp"

type Component[T any] uint

func (c Component[T]) Get(w *World, entity Entity) *T {
	arch := w.entityArchetype[entity]
	if !arch.hasComponent(c) {
		return nil
	}
	column := arch.getColumn(c)

	return column.get(arch.entities.indexOf(entity)).(*T)
}

type ComponentPointer interface {
	set(value uint)
}

func (c *Component[T]) set(value uint) {
	*c = Component[T](value)
}

type AnyComponent interface {
	toUint() uint
	createColumn() anyColumn
}

func (c Component[T]) toUint() uint {
	return uint(c)
}

func (c Component[T]) createColumn() anyColumn {
	return newColumn[T]()
}

type ComponentValue interface {
	GetComponent() AnyComponent
}

func compareComponentValues(a, b ComponentValue) int {
	return cmp.Compare(a.GetComponent().toUint(), b.GetComponent().toUint())
}
