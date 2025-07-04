package entities

import (
	"testing"

	collider "github.com/vcscsvcscs/ebiten-collider"
)

func TestNewWall(t *testing.T) {
	collisionSpace := collider.NewSpatialHash(16)
	wall := NewWall(collisionSpace, 10.0, 20.0, "assets/map/solid_block.png")

	if !wall.IsSolid() {
		t.Error("Expected wall to be solid")
	}
	if wall.IsDestroyable() {
		t.Error("Expected wall not to be destroyable")
	}
}

func TestNewGrass(t *testing.T) {
	collisionSpace := collider.NewSpatialHash(16)
	grass := NewGrass(collisionSpace, 30.0, 40.0, "assets/map/grass_block.png")

	if grass.IsSolid() {
		t.Error("Expected grass not to be solid")
	}
	if grass.IsDestroyable() {
		t.Error("Expected grass not to be destroyable")
	}
}

func TestNewDestroyAbleSolid(t *testing.T) {
	collisionSpace := collider.NewSpatialHash(16)
	destructible := NewDestroyAbleSolid(collisionSpace, 50.0, 60.0, "assets/map/brick_block.png")

	if !destructible.IsSolid() {
		t.Error("Expected destructible object to be solid")
	}
	if !destructible.IsDestroyable() {
		t.Error("Expected destructible object to be destroyable")
	}
}
