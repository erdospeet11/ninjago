package main

import (
	ebiten "github.com/hajimehoshi/ebiten/v2"
	"szofttech.inf.elte.hu/szofttech-c-2024/group-08/ninjago/inputs"
	"szofttech.inf.elte.hu/szofttech-c-2024/group-08/ninjago/scenes"
	"szofttech.inf.elte.hu/szofttech-c-2024/group-08/ninjago/userinfo"
)

// Game represents the main game struct, including the scene manager, user information, and input manager.
type Game struct {
	sceneManager scenes.SceneManager // The scene manager used to switch between scenes.
	userInfo     userinfo.UserInfo   // The user information for the current player.
	input        inputs.Input        // The input manager used to handle player controls.
}

// Layout sets the screen dimensions for the game.
func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return ScreenWidth, ScreenHeight
}

// Update handles the game logic and updates the current scene.
func (g *Game) Update() error {
	if g.sceneManager == nil {
		g.sceneManager = scenes.NewSceneManager(ScreenWidth, ScreenHeight)
		g.sceneManager.GoTo(scenes.NewMainMenuScene(ScreenWidth, ScreenHeight))
	}

	g.input.Update()
	if err := g.sceneManager.Update(&g.input, &g.userInfo); err != nil {
		return err
	}
	return nil
}

// Draw renders the current scene onto the screen.
func (g *Game) Draw(screen *ebiten.Image) {
	g.sceneManager.Draw(screen)
}
