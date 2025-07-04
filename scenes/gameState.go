// Package scenes provides the implementation of various game scenes and their management.
package scenes

import (
	"szofttech.inf.elte.hu/szofttech-c-2024/group-08/ninjago/inputs"
	"szofttech.inf.elte.hu/szofttech-c-2024/group-08/ninjago/userinfo"
)

// GameState represents the current state of the game, including the scene manager, input, and user information.
type GameState struct {
	SceneManager SceneManager       // The scene manager used to switch between scenes.
	Input        *inputs.Input      // The input manager used to handle player controls.
	UserInfo     *userinfo.UserInfo // The user information for the current player.
}
