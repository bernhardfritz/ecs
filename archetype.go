package ecs

import "math/bits"

type archetype struct {
	id       uint
	entities *sparseSet[Entity]
	columns  []anyColumn
}

func (a *archetype) getEntities() []Entity {
	return a.entities.dense
}

func (a *archetype) getColumn(component AnyComponent) anyColumn {
	lowerMask := component.toUint() - 1

	return a.columns[bits.OnesCount(a.id&lowerMask)]
}

func (a *archetype) hasComponent(component AnyComponent) bool {
	return (a.id & component.toUint()) == component.toUint()
}
