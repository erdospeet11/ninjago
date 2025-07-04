package entities

import (
	"image/color"
	"testing"

	collider "github.com/vcscsvcscs/ebiten-collider"
	"szofttech.inf.elte.hu/szofttech-c-2024/group-08/ninjago/userinfo"
)

func TestNewPlayer(t *testing.T) {
	colliderSpace := collider.NewSpatialHash(16)
	userData := &userinfo.UserInfo{
		Username: "Player",
		UserID:   "localhost",
	}

	player := NewPlayer(colliderSpace, 0, 0, userData, color.Black)

	if player == nil {
		t.Fatal("Expected player to be initialized, got nil")
	}

	if player.speed != 0.6 {
		t.Errorf("Expected initial speed to be 0.6, got %f", player.speed)
	}

	if len(player.animations) == 0 {
		t.Error("Expected animations to be loaded")
	}
}

func TestPlayerUpdate_Movement(t *testing.T) {
	colliderSpace := collider.NewSpatialHash(16)
	userData := &userinfo.UserInfo{
		Username: "Player",
		UserID:   "localhost",
	}

	player := NewPlayer(colliderSpace, 0, 0, userData, color.Black)
	player.Control.Up = true

	_, _, err := player.Update()
	if err != nil {
		t.Errorf("Update failed: %v", err)
	}

	if player.currentAnimation != "walkUp" {
		t.Errorf("Expected current animation to be 'walkUp', got '%s'", player.currentAnimation)
	}
}

func TestPlayerPlaceBomb(t *testing.T) {
	colliderSpace := collider.NewSpatialHash(16)
	userData := &userinfo.UserInfo{
		Username: "Player",
		UserID:   "localhost",
	}

	player := NewPlayer(colliderSpace, 0, 0, userData, color.Black)
	player.Control.Ability1 = true
	player.canPlaceBomb = true
	player.NumberOfBombs = 1

	bomb, _, _ := player.Update()

	if bomb == nil {
		t.Error("Expected bomb to be placed")
	}

	if player.NumberOfBombs != 0 {
		t.Errorf("Expected number of bombs to decrement, remaining: %d", player.NumberOfBombs)
	}
}
