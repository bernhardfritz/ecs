package ecs

import (
	"iter"
	"testing"

	"github.com/stretchr/testify/assert"
)

type Position struct {
	X float32
	Y float32
}

func (p Position) GetComponent() AnyComponent {
	return positionComponent
}

type Velocity struct {
	Vx float32
	Vy float32
}

func (v Velocity) GetComponent() AnyComponent {
	return velocityComponent
}

type Health int

func (h Health) GetComponent() AnyComponent {
	return healthComponent
}

var positionComponent Component[Position]
var velocityComponent Component[Velocity]
var healthComponent Component[Health]

func TestEcs(t *testing.T) {
	world := NewWorld(10)
	world.RegisterComponent(&positionComponent)
	world.RegisterComponent(&velocityComponent)
	world.RegisterComponent(&healthComponent)
	firstEntity := world.CreateEntity(Position{X: 12, Y: 34}, Velocity{Vx: 1, Vy: 2})
	secondEntity := world.CreateEntity(Position{X: 56, Y: 78}, Health(100))
	movementSystem := world.CreateSystem(move, positionComponent, velocityComponent)

	movementSystem()

	assert.Equal(t, Position{13, 36}, *positionComponent.Get(world, firstEntity))
	assert.Equal(t, Position{56, 78}, *positionComponent.Get(world, secondEntity))
	assert.Equal(t, (*Health)(nil), healthComponent.Get(world, firstEntity))
	assert.Equal(t, Health(100), *healthComponent.Get(world, secondEntity))
}

func move(w *World, entities iter.Seq[Entity]) {
	for entity := range entities {
		position, velocity := positionComponent.Get(w, entity), velocityComponent.Get(w, entity)
		position.X += velocity.Vx
		position.Y += velocity.Vy
	}
}
