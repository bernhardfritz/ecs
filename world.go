package ecs

import (
	"iter"
	"slices"
)

type World struct {
	capacity             int
	nextEntity           int
	nextComponent        uint
	nextSystem           int
	archetypes           map[uint]*archetype
	systemComponentsMask []uint
	systemArchetypes     [][]*archetype
	entityArchetype      []*archetype
}

func NewWorld(capacity int) *World {
	return &World{
		capacity:             capacity,
		nextEntity:           0,
		nextComponent:        1,
		nextSystem:           0,
		archetypes:           make(map[uint]*archetype),
		systemComponentsMask: make([]uint, 0),
		systemArchetypes:     make([][]*archetype, 0),
		entityArchetype:      make([]*archetype, capacity),
	}
}

func (w *World) RegisterComponent(component ComponentPointer) {
	component.set(w.nextComponent)
	w.nextComponent <<= 1
}

func (w *World) CreateEntity(componentValues ...ComponentValue) Entity {
	entity := Entity(w.nextEntity)
	w.nextEntity++
	var archetypeId uint = 0
	for _, componentValue := range componentValues {
		archetypeId |= componentValue.GetComponent().toUint()
	}
	arch := w.archetypes[archetypeId]
	if arch == nil {
		slices.SortFunc(componentValues, compareComponentValues)
		columns := make([]anyColumn, 0, len(componentValues))
		for _, component := range componentValues {
			columns = append(columns, component.GetComponent().createColumn())
		}
		arch = &archetype{
			id:       archetypeId,
			entities: newSparseSet[Entity](w.capacity),
			columns:  columns,
		}
		w.archetypes[archetypeId] = arch
		for system, componentsMask := range w.systemComponentsMask {
			if (archetypeId & componentsMask) != componentsMask {
				continue
			}
			w.systemArchetypes[system] = append(w.systemArchetypes[system], arch)
		}
	}
	// components are not necessarily sorted if archetype already exists but that's fine as long as code below doesn't require components to be sorted
	arch.entities.add(entity)
	for _, componentValue := range componentValues {
		column := arch.getColumn(componentValue.GetComponent())
		column.add(componentValue)
	}
	w.entityArchetype[entity] = arch

	return entity
}

func (w *World) CreateSystem(fn func(*World, iter.Seq[Entity]), components ...AnyComponent) func() {
	system := w.nextSystem
	w.nextSystem++
	var componentsMask uint = 0
	for _, component := range components {
		componentsMask |= component.toUint()
	}
	w.systemComponentsMask = append(w.systemComponentsMask, componentsMask)
	w.systemArchetypes = append(w.systemArchetypes, make([]*archetype, 0))
	for _, arch := range w.archetypes {
		if (arch.id & componentsMask) != componentsMask {
			continue
		}
		w.systemArchetypes[system] = append(w.systemArchetypes[system], arch)
	}

	return func() {
		fn(w, func(yield func(Entity) bool) {
			for _, arch := range w.systemArchetypes[system] {
				for _, entity := range arch.getEntities() {
					if !yield(entity) {
						return
					}
				}
			}
		})
	}
}

// TODO removing entities
