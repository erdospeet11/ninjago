package entities

import (
	"testing"

	collider "github.com/vcscsvcscs/ebiten-collider"
)

func TestNewBomb(t *testing.T) {
	colliderSpace := collider.NewSpatialHash(16)
	owner := &Player{}

	bomb := NewBomb(colliderSpace, owner, 5, 50.0, 60.0)

	if bomb == nil {
		t.Fatal("Expected bomb to be initialized, got nil")
	}

	if bomb.Owner != owner {
		t.Error("Bomb owner is not set correctly")
	}

	if bomb.ExplosionRange != 5 {
		t.Errorf("Expected ExplosionRange to be 5, got %d", bomb.ExplosionRange)
	}

	if bomb.time != 180 {
		t.Errorf("Expected initial time to be 180, got %d", bomb.time)
	}
}

func TestBombUpdate(t *testing.T) {
	colliderSpace := collider.NewSpatialHash(16)
	owner := &Player{}

	bomb := NewBomb(colliderSpace, owner, 5, 50.0, 60.0)

	for i := 0; i < 179; i++ {
		if bomb.Update() {
			t.Fatal("Bomb should not have exploded yet")
		}
	}

	if !bomb.Update() {
		t.Error("Bomb should explode now")
	}
}
