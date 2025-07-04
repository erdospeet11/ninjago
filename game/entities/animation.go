// Package entities provides the definition and implementation of game entities
// and their animations within the game world.
package entities

import (
	"image"
	"log"

	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"szofttech.inf.elte.hu/szofttech-c-2024/group-08/ninjago/assets"

	"github.com/hajimehoshi/ebiten/v2"
)

// Animation represents a sprite animation with multiple frames.
// It includes the sprite image, frame dimensions, and animation speed.
type Animation struct {
	sprite         *ebiten.Image // The sprite image containing all frames of the animation.
	frameOX        int           // The x offset of the first frame in the sprite.
	frameOY        int           // The y offset of the first frame in the sprite.
	frameWidth     int           // The width of each frame.
	frameHeight    int           // The height of each frame.
	frameCount     int           // The total number of frames in the animation.
	currentFrame   int           // The current frame index in the animation.
	animationSpeed int           // The speed of the animation (frames per update).
}

// Update advances the animation to the next frame based on the animation speed.
func (a *Animation) Update() {
	a.currentFrame++
}

// Draw renders the current frame of the animation at the specified position on the screen.
// It takes into account the screen width and height for positioning.
func (a *Animation) Draw(screen *ebiten.Image, x, y float64, screenWidth float64, screenHeight float64) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(x, y)
	i := (a.currentFrame / a.animationSpeed) % a.frameCount
	sx, sy := a.frameOX+i*a.frameWidth, a.frameOY
	screen.DrawImage(a.sprite.SubImage(image.Rect(sx, sy, sx+a.frameWidth, sy+a.frameHeight)).(*ebiten.Image), op)
}

// LoadAnimations loads sprite animations from the provided paths and initializes them with the specified frame count and animation speed.
// It returns a map where the keys are animation names and the values are pointers to the loaded animations.
func LoadAnimations(frameCount int, animationSpeed int, spritePaths map[string]string) map[string]*Animation {
	animations := make(map[string]*Animation)

	for animationName, spritePath := range spritePaths {
		sprite, _, err := ebitenutil.NewImageFromFileSystem(assets.EmbeddedAssets, spritePath)
		if err != nil {
			log.Fatalf("Failed to load sprite: %v", err)
		}

		_, spriteHeight := sprite.Bounds().Dx(), sprite.Bounds().Dy()

		animation := &Animation{
			sprite:         sprite,
			frameOX:        0,
			frameOY:        0,
			frameWidth:     spriteHeight,
			frameHeight:    spriteHeight,
			frameCount:     frameCount,
			currentFrame:   0,
			animationSpeed: animationSpeed,
		}

		animations[animationName] = animation
	}

	return animations
}
