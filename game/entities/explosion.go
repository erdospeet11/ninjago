// Package entities provides the definition and implementation of game entities
// and their interactions within the game world.
package entities

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	collider "github.com/vcscsvcscs/ebiten-collider"
	"szofttech.inf.elte.hu/szofttech-c-2024/group-08/ninjago/assets"
)

type Explosion struct {
	entity     // The explosion entity inherits from the base entity.
	time   int // The lifetime of the explosion in update ticks.
}

// NewExplosion creates a new explosion entity at the specified position.
// It initializes the explosion with a circular collider and an idle animation.
func NewExplosion(collisionSpace *collider.SpatialHash, start_pos_x float64, start_pos_y float64) *Explosion {
	idleSprite, _, err := ebitenutil.NewImageFromFileSystem(assets.EmbeddedAssets, "assets/graphics/explosion_animation.png")
	if err != nil {
		log.Fatalf("Failed to load idle sprite: %v", err)
	}

	idleSpriteHeight := idleSprite.Bounds().Dy()

	idleAnimation := &Animation{
		sprite:         idleSprite,
		frameOX:        0,
		frameOY:        0,
		frameWidth:     idleSpriteHeight,
		frameHeight:    idleSpriteHeight,
		frameCount:     8,
		currentFrame:   0,
		animationSpeed: 24,
	}

	e := &Explosion{
		entity: entity{
			collider:         collisionSpace.NewCircleShape(start_pos_x, start_pos_y, float64(idleSpriteHeight)/2),
			speed:            0,
			currentAnimation: "idle",
			animations: map[string]*Animation{
				"idle": idleAnimation,
			},
		},
		time: 90,
	}
	e.collider.SetParent(e)

	return e
}

// Update updates the explosion's state, decrementing its lifetime.
// It returns true if the explosion's lifetime has reached zero, indicating that the explosion should be removed.
func (e *Explosion) Update() bool {
	e.time--
	if e.time != 0 {
		e.entity.Update()
	}
	return e.time == 0
}
