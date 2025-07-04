package entities

import (
	"testing"
)

func TestDirectionTowardsTarget(t *testing.T) {
	tests := []struct {
		name     string
		monsterX float64
		monsterY float64
		targetX  float64
		targetY  float64
		want     int
	}{
		{"Jobbra", 0, 0, 1, 0, 3},
		{"Balra", 0, 0, -1, 0, 2},
		{"Lefelé", 0, 0, 0, 1, 1},
		{"Felfelé", 0, 0, 0, -1, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := directionTowardsTarget(tt.monsterX, tt.monsterY, tt.targetX, tt.targetY); got != tt.want {
				t.Errorf("DirectionTowardsTarget() = %v, want %v", got, tt.want)
			}
		})
	}
}
