package entities

import (
	"testing"

	collider "github.com/vcscsvcscs/ebiten-collider"
)

func TestNewExplosion(t *testing.T) {
	collisionSpace := collider.NewSpatialHash(16)
	explosion := NewExplosion(collisionSpace, 50.0, 60.0)

	if explosion == nil {
		t.Fatal("Explosion should not be nil")
	}
	if explosion.time != 90 {
		t.Errorf("Expected explosion time to be 90, got %d", explosion.time)
	}
	if len(explosion.entity.animations) == 0 {
		t.Error("No animations loaded for explosion")
	}
	if explosion.entity.currentAnimation != "idle" {
		t.Errorf("Expected current animation to be 'idle', got %s", explosion.entity.currentAnimation)
	}
}

func TestExplosionUpdate(t *testing.T) {
	collisionSpace := collider.NewSpatialHash(16)
	explosion := NewExplosion(collisionSpace, 50.0, 60.0)

	for i := 0; i < 89; i++ {
		if explosion.Update() {
			t.Error("Explosion should not have finished yet")
		}
	}

	if !explosion.Update() {
		t.Error("Explosion should finish on the 90th update")
	}
}
