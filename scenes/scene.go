// Package scenes provides the definition of the Scene interface
// for managing and rendering different game scenes.
package scenes

import ebiten "github.com/hajimehoshi/ebiten/v2"

// Scene defines the behavior of a game scene, including methods for updating and drawing the scene.
type Scene interface {
	Update(state *GameState) error // Update updates the scene's state based on the game state.
	Draw(screen *ebiten.Image)     // Draw renders the scene on the screen.
}
