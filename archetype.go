package ecs

type archetype struct {
	id       uint
	entities *sparseSet
	columns  []isColumn
}

func (a *archetype) getEntities() []Entity {
	return a.entities.dense
}

func (a *archetype) has(component IsComponent) bool {
	return (a.id & component.flag()) == component.flag()
}
