package entities

import (
	"testing"

	collider "github.com/vcscsvcscs/ebiten-collider"
)

func TestNewBombEffect(t *testing.T) {
	colliderSpace := collider.NewSpatialHash(16)
	x, y := 10.0, 20.0
	effect := NewBombEffect(colliderSpace, x, y)

	if effect.StatusEffect == nil {
		t.Error("Expected StatusEffect to be initialized")
	}
	if _, ok := effect.StatusEffect.(*bombInc); !ok {
		t.Errorf("Expected StatusEffect to be of type *bombInc, got %T", effect.StatusEffect)
	}
	if effect.collider == nil {
		t.Error("Expected collider to be initialized")
	}
	if effect.entity.animations["idle"] == nil {
		t.Error("Expected idle animation to be initialized")
	}
}

func TestNewRadiusEffect(t *testing.T) {
	colliderSpace := collider.NewSpatialHash(16)
	x, y := 30.0, 40.0
	effect := NewRadiusEffect(colliderSpace, x, y)

	if effect.StatusEffect == nil {
		t.Error("Expected StatusEffect to be initialized")
	}
	if _, ok := effect.StatusEffect.(*radiusInc); !ok {
		t.Errorf("Expected StatusEffect to be of type *radiusInc, got %T", effect.StatusEffect)
	}
	if effect.collider == nil {
		t.Error("Expected collider to be initialized")
	}
	if effect.entity.animations["idle"] == nil {
		t.Error("Expected idle animation to be initialized")
	}
}

func TestNewSkullDebuff(t *testing.T) {
	colliderSpace := collider.NewSpatialHash(16)
	x, y := 50.0, 60.0
	effect := NewSkullDebuff(colliderSpace, x, y)

	if effect.StatusEffect == nil {
		t.Error("Expected StatusEffect to be initialized")
	}
	if _, ok := effect.StatusEffect.(*skullDeb); !ok {
		t.Errorf("Expected StatusEffect to be of type *skullDeb, got %T", effect.StatusEffect)
	}
	if effect.collider == nil {
		t.Error("Expected collider to be initialized")
	}
	if effect.entity.animations["idle"] == nil {
		t.Error("Expected idle animation to be initialized")
	}
}
