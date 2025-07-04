// Package entities provides the definition and implementation of game entities
// and effects that can be applied to them. It includes various functions to create
// different effects.
package entities

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	collider "github.com/vcscsvcscs/ebiten-collider"
	"szofttech.inf.elte.hu/szofttech-c-2024/group-08/ninjago/assets"
)

// Effect represents a game effect that includes an entity and a status effect.
type Effect struct {
	entity
	StatusEffect StatusEffect
}

// NewBombEffect creates a new bomb effect at the specified position.
// It initializes the effect with a collider and an idle animation.
func NewBombEffect(colliderSpace *collider.SpatialHash, start_pos_x float64, start_pos_y float64) Effect {
	idleSprite, _, err := ebitenutil.NewImageFromFileSystem(assets.EmbeddedAssets, "assets/powerup/BombIncrease.png")
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
		frameCount:     1,
		currentFrame:   0,
		animationSpeed: 24,
	}

	e := Effect{
		entity: entity{
			collider:         colliderSpace.NewCircleShape(start_pos_x, start_pos_y, float64(idleSpriteHeight)/2),
			speed:            0,
			currentAnimation: "idle",
			animations: map[string]*Animation{
				"idle": idleAnimation,
			},
		},
		StatusEffect: NewBombInc(),
	}
	e.collider.SetParent(e)

	return e
}

// NewRadiusEffect creates a new radius effect at the specified position.
// It initializes the effect with a collider and an idle animation.
func NewRadiusEffect(collisionSpace *collider.SpatialHash, start_pos_x float64, start_pos_y float64) Effect {
	idleSprite, _, err := ebitenutil.NewImageFromFileSystem(assets.EmbeddedAssets, "assets/powerup/radiusIncrease.png")
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
		frameCount:     1,
		currentFrame:   0,
		animationSpeed: 24,
	}

	e := Effect{
		entity: entity{
			collider:         collisionSpace.NewCircleShape(start_pos_x, start_pos_y, float64(idleSpriteHeight)/2),
			speed:            0,
			currentAnimation: "idle",
			animations: map[string]*Animation{
				"idle": idleAnimation,
			},
		},
		StatusEffect: NewRadiusInc(),
	}
	e.collider.SetParent(e)

	return e
}

// NewSkullDebuff creates a new skull debuff effect at the specified position.
// It initializes the effect with a collider and an idle animation.
func NewSkullDebuff(collisionSpace *collider.SpatialHash, start_pos_x float64, start_pos_y float64) Effect {
	idleSprite, _, err := ebitenutil.NewImageFromFileSystem(assets.EmbeddedAssets, "assets/powerup/SkullDecrease.png")
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
		frameCount:     1,
		currentFrame:   0,
		animationSpeed: 24,
	}

	e := Effect{
		entity: entity{
			collider:         collisionSpace.NewCircleShape(start_pos_x, start_pos_y, float64(idleSpriteHeight)/2),
			speed:            0,
			currentAnimation: "idle",
			animations: map[string]*Animation{
				"idle": idleAnimation,
			},
		},
		StatusEffect: NewSkullDeb(),
	}
	e.collider.SetParent(e)

	return e
}

// NewRollerEffect creates a new roller effect at the specified position.
// It initializes the effect with a collider and an idle animation.
func NewRollerEffect(collisionSpace *collider.SpatialHash, start_pos_x float64, start_pos_y float64) Effect {
	idleSprite, _, err := ebitenutil.NewImageFromFileSystem(assets.EmbeddedAssets, "assets/powerup/Roller.png")
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
		frameCount:     1,
		currentFrame:   0,
		animationSpeed: 24,
	}

	e := Effect{
		entity: entity{
			collider:         collisionSpace.NewCircleShape(start_pos_x, start_pos_y, float64(idleSpriteHeight)/2),
			speed:            0,
			currentAnimation: "idle",
			animations: map[string]*Animation{
				"idle": idleAnimation,
			},
		},
		StatusEffect: NewRollerInc(),
	}
	e.collider.SetParent(e)

	return e
}

// NewObstacleEffect creates a new obstacle effect at the specified position.
// It initializes the effect with a collider and an idle animation.
func NewObstacleEffect(collisionSpace *collider.SpatialHash, start_pos_x float64, start_pos_y float64) Effect {
	idleSprite, _, err := ebitenutil.NewImageFromFileSystem(assets.EmbeddedAssets, "assets/powerup/Obstacle.png")
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
		frameCount:     1,
		currentFrame:   0,
		animationSpeed: 24,
	}

	e := Effect{
		entity: entity{
			collider:         collisionSpace.NewCircleShape(start_pos_x, start_pos_y, float64(idleSpriteHeight)/2),
			speed:            0,
			currentAnimation: "idle",
			animations: map[string]*Animation{
				"idle": idleAnimation,
			},
		},
		StatusEffect: NewObstacleInc(),
	}
	e.collider.SetParent(e)

	return e
}

// NewDetonatorEffect creates a new detonator effect at the specified position.
// It initializes the effect with a collider and an idle animation.
func NewDetonatorEffect(collisionSpace *collider.SpatialHash, start_pos_x float64, start_pos_y float64) Effect {
	idleSprite, _, err := ebitenutil.NewImageFromFileSystem(assets.EmbeddedAssets, "assets/powerup/Detonator.png")
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
		frameCount:     1,
		currentFrame:   0,
		animationSpeed: 24,
	}

	e := Effect{
		entity: entity{
			collider:         collisionSpace.NewCircleShape(start_pos_x, start_pos_y, float64(idleSpriteHeight)/2),
			speed:            0,
			currentAnimation: "idle",
			animations: map[string]*Animation{
				"idle": idleAnimation,
			},
		},
		StatusEffect: NewDetonatorInc(),
	}
	e.collider.SetParent(e)

	return e
}

// NewGhostEffect creates a new ghost effect at the specified position.
// It initializes the effect with a collider and an idle animation.
func NewGhostEffect(collisionSpace *collider.SpatialHash, start_pos_x float64, start_pos_y float64) Effect {
	idleSprite, _, err := ebitenutil.NewImageFromFileSystem(assets.EmbeddedAssets, "assets/powerup/Ghost.png")
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
		frameCount:     1,
		currentFrame:   0,
		animationSpeed: 24,
	}

	e := Effect{
		entity: entity{
			collider:         collisionSpace.NewCircleShape(start_pos_x, start_pos_y, float64(idleSpriteHeight)/2),
			speed:            0,
			currentAnimation: "idle",
			animations: map[string]*Animation{
				"idle": idleAnimation,
			},
		},
		StatusEffect: NewGhostInc(),
	}
	e.collider.SetParent(e)

	return e
}

// NewInvincibilityEffect creates a new invincibility effect at the specified position.
// It initializes the effect with a collider and an idle animation.
func NewInvincibilityEffect(collisionSpace *collider.SpatialHash, start_pos_x float64, start_pos_y float64) Effect {
	idleSprite, _, err := ebitenutil.NewImageFromFileSystem(assets.EmbeddedAssets, "assets/powerup/Invincibility.png")
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
		frameCount:     1,
		currentFrame:   0,
		animationSpeed: 24,
	}

	e := Effect{
		entity: entity{
			collider:         collisionSpace.NewCircleShape(start_pos_x, start_pos_y, float64(idleSpriteHeight)/2),
			speed:            0,
			currentAnimation: "idle",
			animations: map[string]*Animation{
				"idle": idleAnimation,
			},
		},
		StatusEffect: NewInvincibilityInc(),
	}
	e.collider.SetParent(e)

	return e
}
