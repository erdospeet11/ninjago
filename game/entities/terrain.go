// Package entities provides the definition and implementation of game entities
// and their interactions within the game world, including various types of terrain.
package entities

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	collider "github.com/vcscsvcscs/ebiten-collider"
	"szofttech.inf.elte.hu/szofttech-c-2024/group-08/ninjago/assets"
)

// terrain represents a basic terrain entity with solidity and destructibility properties.
type terrain struct {
	entity           // The terrain entity inherits from the base entity.
	solid       bool // Indicates if the terrain is solid and blocks movement.
	destroyable bool // Indicates if the terrain can be destroyed by bombs.
}

// Terrain defines the behavior of a terrain entity in the game.
// It extends the Entity interface with additional methods to check solidity and destructibility.
type Terrain interface {
	Entity               // The terrain entity extends the base entity.
	IsSolid() bool       // IsSolid returns true if the terrain is solid and blocks movement.
	IsDestroyable() bool // IsDestroyable returns true if the terrain can be destroyed by bombs.
}

// IsSolid returns true if the terrain is solid.
func (t *terrain) IsSolid() bool {
	return t.solid
}

// IsDestroyable returns true if the terrain is destroyable.
func (t *terrain) IsDestroyable() bool {
	return t.destroyable
}

// NewWall creates a new wall terrain entity at the specified position with the provided image path.
// It initializes the wall as a solid and non-destroyable terrain.
func NewWall(collisionSpace *collider.SpatialHash, x, y float64, imagePath string) Terrain {
	sprite, _, err := ebitenutil.NewImageFromFileSystem(assets.EmbeddedAssets, imagePath)
	if err != nil {
		log.Panic(err)
	}
	t := &terrain{
		entity: entity{
			collider:         collisionSpace.NewRectangleShape(x, y, 16, 16),
			speed:            1,
			currentAnimation: "idle",
			animations: map[string]*Animation{
				"idle": {
					sprite:         sprite,
					frameOX:        0,
					frameOY:        0,
					frameWidth:     32,
					frameHeight:    32,
					frameCount:     1,
					currentFrame:   0,
					animationSpeed: 10,
				},
			},
		},
		solid:       true,
		destroyable: false,
	}
	t.collider.SetParent((Terrain)(t))

	return t
}

// NewGrass creates a new grass terrain entity at the specified position with the provided image path.
// It initializes the grass as a non-solid and non-destroyable terrain.
func NewGrass(collisionSpace *collider.SpatialHash, x, y float64, imagePath string) Terrain {
	sprite, _, err := ebitenutil.NewImageFromFileSystem(assets.EmbeddedAssets, imagePath)
	if err != nil {
		log.Panic(err)
	}
	t := &terrain{
		entity: entity{
			collider:         collisionSpace.NewRectangleShape(x, y, 16, 16),
			speed:            0,
			currentAnimation: "idle",
			animations: map[string]*Animation{
				"idle": {
					sprite:         sprite,
					frameOX:        0,
					frameOY:        0,
					frameWidth:     32,
					frameHeight:    32,
					frameCount:     1,
					currentFrame:   0,
					animationSpeed: 10,
				},
			},
		},
		solid:       false,
		destroyable: false,
	}
	t.collider.SetParent((Terrain)(t))

	return t
}

// NewDestroyAbleSolid creates a new destroyable solid terrain entity at the specified position with the provided image path.
// It initializes the terrain as a solid and destroyable entity.
func NewDestroyAbleSolid(collisionSpace *collider.SpatialHash, x, y float64, imagePath string) Terrain {
	sprite, _, err := ebitenutil.NewImageFromFileSystem(assets.EmbeddedAssets, imagePath)
	if err != nil {
		log.Panic(err)
	}

	t := &terrain{
		entity: entity{
			collider:         collisionSpace.NewRectangleShape(x, y, 16, 16),
			speed:            0,
			currentAnimation: "idle",
			animations: map[string]*Animation{
				"idle": {
					sprite:         sprite,
					frameOX:        0,
					frameOY:        0,
					frameWidth:     32,
					frameHeight:    32,
					frameCount:     1,
					currentFrame:   0,
					animationSpeed: 10,
				},
			},
		},
		solid:       true,
		destroyable: true,
	}
	t.collider.SetParent((Terrain)(t))

	return t
}
