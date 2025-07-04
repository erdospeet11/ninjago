package entities

import (
	"image/color"
	"testing"

	collider "github.com/vcscsvcscs/ebiten-collider"
	"szofttech.inf.elte.hu/szofttech-c-2024/group-08/ninjago/userinfo"
)

func setupPlayer() *Player {

	userData := &userinfo.UserInfo{
		Username: "Player",
		UserID:   "localhost",
	}

	colliderSpace := collider.NewSpatialHash(16)
	player := NewPlayer(colliderSpace, 0, 0, userData, color.Black)

	return player
}

func TestSkullDebEffect(t *testing.T) {
	player := setupPlayer()
	effect := NewSkullDeb()
	player.State = effect

	for i := 0; i < 1800; i++ {
		if i == 1799 {
			if player.speed != 0.3 {
				t.Errorf("Expected speed 0.3, got %f", player.speed)
			}
		}
		effect.Update(player)
	}

	if player.speed != 0.6 {
		t.Errorf("Expected speed to reset to 0.6 after effect, got %f", player.speed)
	}
}

func TestSkateEffect(t *testing.T) {
	player := setupPlayer()
	effect := NewSkate()
	player.State = effect

	for i := 0; i < 60; i++ {
		effect.Update(player)
	}

	if player.speed != 0.6 {
		t.Errorf("Expected speed to reset to 0.6 after effect, got %f", player.speed)
	}
}

func TestRadiusIncEffect(t *testing.T) {
	player := setupPlayer()
	effect := NewRadiusInc()
	player.State = effect

	for i := 0; i < 1799; i++ {
		effect.Update(player)
	}

	if player.BombRange != 5 {
		t.Errorf("Expected BombRange 5 during effect, got %d", player.BombRange)
	}

	effect.Update(player)
	if player.BombRange != 2 {
		t.Errorf("Expected BombRange to reset to 2 after effect, got %d", player.BombRange)
	}
}
