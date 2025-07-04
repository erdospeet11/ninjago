// Package scenes provides the implementation of various game scenes and their management,
// including the configuration of gamepad controls.
package scenes

import (
	"fmt"
	"image/color"
	"strings"

	ebiten "github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	text "github.com/hajimehoshi/ebiten/v2/text/v2"
	"szofttech.inf.elte.hu/szofttech-c-2024/group-08/ninjago/assets"
	"szofttech.inf.elte.hu/szofttech-c-2024/group-08/ninjago/inputs"
)

// GamepadScene represents a scene for configuring gamepad controls.
type GamepadScene struct {
	previusScene      Scene            // The previous scene to return to after configuring gamepad controls.
	gamepadID         ebiten.GamepadID // The gamepad ID for the current gamepad.
	currentIndex      int              // The current index of the button being configured.
	countAfterSetting int              // The count after setting the gamepad controls.
	buttonStates      []string         // The current state of the gamepad buttons.
}

// Update updates the gamepad configuration scene based on the current game state.
// It handles input for setting gamepad buttons and transitions back to the previous scene.
//
// Parameters:
//   - state: The current game state.
//
// Returns:
//   - error: An error if the update fails.
func (s *GamepadScene) Update(state *GameState) error {
	if s.currentIndex == 0 {
		state.Input.GamepadConfig.Reset()
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		state.Input.GamepadConfig.Reset()
		state.Input.GamepadConfig.ResetGamepadID()
		state.SceneManager.GoTo(s.previusScene)
		return nil
	}

	if s.buttonStates == nil {
		s.buttonStates = make([]string, len(inputs.VirtualGamepadButtons))
	}
	for i, b := range inputs.VirtualGamepadButtons {
		if i < s.currentIndex {
			s.buttonStates[i] = strings.ToUpper(state.Input.GamepadConfig.ButtonName(b))
			continue
		}
		if s.currentIndex == i {
			s.buttonStates[i] = "_"
			continue
		}
		s.buttonStates[i] = ""
	}

	if 0 < s.countAfterSetting {
		s.countAfterSetting--
		if s.countAfterSetting <= 0 {
			state.SceneManager.GoTo(s.previusScene)
		}
		return nil
	}

	b := inputs.VirtualGamepadButtons[s.currentIndex]
	if state.Input.GamepadConfig.Scan(b) {
		s.currentIndex++
		if s.currentIndex == len(inputs.VirtualGamepadButtons) {
			s.countAfterSetting = ebiten.TPS()
		}
	}
	return nil
}

// Draw renders the gamepad configuration scene onto the provided screen image.
//
// Parameters:
//   - screen: The image to which the scene is drawn.
func (s *GamepadScene) Draw(screen *ebiten.Image) {
	screen.Fill(color.Black)

	if s.buttonStates == nil {
		return
	}

	f := `GAMEPAD CONFIGURATION
(PRESS ESC TO CANCEL)


MOVE LEFT:    %s

MOVE RIGHT:   %s

MOVE UP:    %s

MOVE DOWN:   %s


PLACE BOMB:   %s

USE ABILITY:  %s



%s`
	msg := ""
	if s.currentIndex == len(inputs.VirtualGamepadButtons) {
		msg = "OK!"
	}
	str := fmt.Sprintf(f, s.buttonStates[0], s.buttonStates[1], s.buttonStates[2], s.buttonStates[3], s.buttonStates[4], s.buttonStates[5], msg)
	assets.DrawTextWithShadow(screen, str, 16, 16, 1, color.White, text.AlignStart, text.AlignStart)
}
