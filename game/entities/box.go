// Package entities provides the definition and implementation of game entities
// and their interactions within the game world.
package entities

import (
	"math/rand"

	collider "github.com/vcscsvcscs/ebiten-collider"
)

// Box represents a box entity in the game, which can be either blank or contain
// a random status effect that can be dropped.
type Box struct {
	entity       // The box entity inherits the entity properties.
	IsBlank bool // IsBlank represents whether the box is blank or contains a status effect.
}

// NewBox creates a new box entity at the specified position with the given blank status.
// It initializes the box with a rectangular collider and an idle animation.
func NewBox(colliderSpace *collider.SpatialHash, start_pos_x float64, start_pos_y float64, isBlank bool) *Box {
	b := &Box{
		entity: entity{
			collider:         colliderSpace.NewRectangleShape(start_pos_x, start_pos_y, 16, 16),
			currentAnimation: "idle",
			animations:       LoadAnimations(1, 1, map[string]string{"idle": "assets/map/box.png"}),
		},
		IsBlank: isBlank,
	}
	b.collider.SetParent(b)

	return b
}

// DropRandomStatusEffect drops a random status effect from the box with a 40% chance.
// It returns an Effect representing the dropped status effect.
func (b *Box) DropRandomStatusEffect() Effect {
	randomEffect := rand.Intn(8) + 1
	randomNumber := rand.Intn(100) + 1
	if randomNumber <= 40 {
		switch randomEffect {
		case 1:
			return NewObstacleEffect(b.collider.GetHash(), b.collider.GetPosition().X, b.collider.GetPosition().Y)
		case 2:
			return NewSkullDebuff(b.collider.GetHash(), b.collider.GetPosition().X, b.collider.GetPosition().Y)
		case 3:
			return NewRollerEffect(b.collider.GetHash(), b.collider.GetPosition().X, b.collider.GetPosition().Y)
		case 4:
			return NewRadiusEffect(b.collider.GetHash(), b.collider.GetPosition().X, b.collider.GetPosition().Y)
		case 5:
			return NewGhostEffect(b.collider.GetHash(), b.collider.GetPosition().X, b.collider.GetPosition().Y)
		case 6:
			return NewDetonatorEffect(b.collider.GetHash(), b.collider.GetPosition().X, b.collider.GetPosition().Y)
		case 7:
			return NewBombEffect(b.collider.GetHash(), b.collider.GetPosition().X, b.collider.GetPosition().Y)
		case 8:
			return NewInvincibilityEffect(b.collider.GetHash(), b.collider.GetPosition().X, b.collider.GetPosition().Y)
		}
	}
	return Effect{}
}
