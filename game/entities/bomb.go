// Package entities provides the definition and implementation of game entities
// and their interactions within the game world.
package entities

import (
	collider "github.com/vcscsvcscs/ebiten-collider"
)

// Bomb represents a bomb entity in the game, which can be placed by a player
// and has a defined explosion range and timer.
type Bomb struct {
	entity
	Owner          *Player // The player who placed the bomb.
	ExplosionRange int     // The range of the bomb's explosion.
	time           int     // The timer for the bomb's explosion countdown.
	manualDetonate bool    // Indicates if the bomb is set to manual detonation.
}

// NewBomb creates a new bomb entity at the specified position with the given owner and explosion range.
// It initializes the bomb with a circular collider and an idle animation.
func NewBomb(collisionSpace *collider.SpatialHash, Owner *Player, ExplosionRange int, start_pos_x float64, start_pos_y float64) *Bomb {
	b := &Bomb{
		entity: entity{
			collider:         collisionSpace.NewCircleShape(start_pos_x, start_pos_y, 7),
			speed:            0,
			currentAnimation: "idle",
			animations:       LoadAnimations(8, 24, map[string]string{"idle": "assets/bomb/bomb_explosion.png"}),
		},
		ExplosionRange: ExplosionRange,
		time:           180,
		Owner:          Owner,
		manualDetonate: false,
	}
	b.collider.SetParent(b)

	return b
}

// Update updates the bomb's state, decrementing the timer if the bomb is not set to manual detonation.
// It returns true if the bomb's timer has reached zero, indicating that the bomb should explode.
func (b *Bomb) Update() bool {
	if !b.manualDetonate {
		b.time--
		if b.time > 0 {
			b.entity.Update()
		}
		return b.time <= 0
	}
	return false
}

// Detonate sets the bomb's timer to zero, causing it to explode immediately.
func (b *Bomb) Detonate() {
	b.time = 0
	b.manualDetonate = false
}
