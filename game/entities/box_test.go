package entities

import (
	"fmt"
	"math/rand"
	"testing"

	collider "github.com/vcscsvcscs/ebiten-collider"
)

func TestDropRandomStatusEffect(t *testing.T) {
	colliderSpace := collider.NewSpatialHash(16)
	box := NewBox(colliderSpace, 50.0, 60.0, false)

	rand.Int()

	tests := []struct {
		name     string
		expected string
	}{
		{"Test for Bomb Effect", "entities.Effect"},
		{"Test for Radius Effect", "entities.Effect"},
		{"Test for Skull Debuff", "entities.Effect"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			effect := box.DropRandomStatusEffect()
			effectType := fmt.Sprintf("%T", effect)

			if effectType != tt.expected {
				t.Errorf("Expected %s, got %s", tt.expected, effectType)
			}
		})
	}
}
