package ecs

import (
	"cmp"
	"iter"
	"math/bits"
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
		nextComponent:        0,
		nextSystem:           0,
		archetypes:           make(map[uint]*archetype),
		systemComponentsMask: make([]uint, 0),
		systemArchetypes:     make([][]*archetype, 0),
		entityArchetype:      make([]*archetype, capacity),
	}
}

func (w *World) RegisterComponent(componentPointer IsComponentPointer) {
	componentPointer.set(w.nextComponent)
	w.nextComponent++
}

func CreateComponent[T any](w *World) Component[T] {
	component := Component[T](w.nextComponent)
	w.nextComponent++

	return component
}

type HasComponent interface {
	Component(w *World) IsComponent
}

func (w *World) compareHasComponents(a, b HasComponent) int {
	return cmp.Compare(a.Component(w).id(), b.Component(w).id())
}

func (w *World) CreateEntity(components ...HasComponent) Entity {
	entity := Entity(w.nextEntity)
	w.nextEntity++
	var archetypeId uint = 0
	for _, component := range components {
		archetypeId |= component.Component(w).flag()
	}
	arch := w.archetypes[archetypeId]
	if arch == nil {
		slices.SortFunc(components, w.compareHasComponents)
		columns := make([]isColumn, 0, len(components))
		for _, component := range components {
			columns = append(columns, component.Component(w).createColumn())
		}
		arch = &archetype{
			id:       archetypeId,
			entities: newSparseSet(w.capacity),
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
	for _, component := range components {
		arch.columns[countLowerSetBits(component.Component(w).id(), component.Component(w).flag())].add(component)
	}
	w.entityArchetype[entity] = arch

	return entity
}

func (w *World) CreateSystem(fn func(*World, iter.Seq[Entity]), components ...IsComponent) func() {
	system := w.nextSystem
	w.nextSystem++
	var componentsMask uint = 0
	for _, component := range components {
		componentsMask |= component.flag()
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

func countLowerSetBits(number uint, flag uint) int {
	lowerMask := flag - 1

	return bits.OnesCount(number & lowerMask)
}

// TODO removing entities
