// Package entities provides the definition and implementation of game entities
// that can be updated, drawn, and animated. It includes the Entity interface and
// the concrete implementation of the entity struct.
package entities

import (
	"fmt"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	collider "github.com/vcscsvcscs/ebiten-collider"
)

// Entity defines the behavior of a game entity which can be updated, drawn,
// and have animations. It provides methods to get the entity's position and collider.
type Entity interface {
	// Update updates the entity's state.
	Update() error
	// Draw draws the entity on the screen.
	Draw(screen *ebiten.Image)
	// SetCurrentAnimation sets the current animation of the entity.
	SetCurrentAnimation(animationKey string) error
	// TilePosition returns the entity's position in terms of tile coordinates.
	TilePosition() (int, int)
	// GetPosition returns the entity's position in the game world.
	GetPosition() (float64, float64)
	// GetCollider returns the entity's collider shape.
	GetCollider() collider.Shape
}

// entity is a concrete implementation of the Entity interface.
// It holds the entity's collider, speed, current animation, and available animations.
type entity struct {
	collider         collider.Shape        // The entity's collider shape.
	speed            float64               // The entity's speed.
	currentAnimation string                // The key of the current animation.
	animations       map[string]*Animation // The available animations for the entity.
}

// GetCollider returns the entity's collider shape.
func (e *entity) GetCollider() collider.Shape {
	return e.collider
}

// TilePosition returns the entity's position in terms of tile coordinates.
// The tile size is assumed to be 16x16.
func (e *entity) TilePosition() (int, int) {
	x := int(math.Floor(e.collider.GetPosition().X / float64(16)))
	y := int(math.Floor(e.collider.GetPosition().Y / float64(16)))

	return x, y
}

// GetPosition returns the entity's position in the game world as (x, y) coordinates.
func (e *entity) GetPosition() (float64, float64) {
	return e.collider.GetPosition().X, e.collider.GetPosition().Y
}

// Update updates the entity's state by advancing its current animation.
// If the current animation is not found, an error is returned.
func (e *entity) Update() error {

	if animation, ok := e.animations[e.currentAnimation]; ok {
		animation.Update()

		return nil
	}

	return fmt.Errorf("animation with key %s not found", e.currentAnimation)
}

// Draw draws the entity's current animation on the provided screen at the entity's position.
// The drawing takes into account the screen dimensions.
func (e *entity) Draw(screen *ebiten.Image) {
	if animation, ok := e.animations[e.currentAnimation]; ok {
		animation.Draw(screen, e.collider.GetPosition().X, e.collider.GetPosition().Y, float64(screen.Bounds().Dx()), float64(screen.Bounds().Dy()))
	}
}

// SetCurrentAnimation sets the current animation of the entity to the animation specified by animationKey.
// If the animation key is not found, an error is returned.
func (e *entity) SetCurrentAnimation(animationKey string) error {
	if _, ok := e.animations[animationKey]; ok {
		e.currentAnimation = animationKey

		return nil
	}

	return fmt.Errorf("animation with key %s not found", animationKey)
}
