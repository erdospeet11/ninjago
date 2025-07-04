package inputs

import (
	ebiten "github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

// Input manages the input state including gamepads and keyboards.
type Input struct {
	gamepadIDs                 []ebiten.GamepadID           // The gamepad IDs for the input.
	VirtualGamepadButtonStates map[VirtualGamepadButton]int // The state for the virtual gamepad buttons.
	GamepadConfig              GamepadConfig                // The gamepad configuration for the input.
}

// PlayerControl represents the player controls for the game.
type PlayerControl interface {
	IsAbilityOneJustPressed() bool // IsAbilityOneJustPressed returns true if the ability one button is just pressed.
	IsAbilityTwoJustPressed() bool // IsAbilityTwoJustPressed returns true if the ability two button is just pressed.
	StateForLeft() int             // StateForLeft returns the state for the left direction.
	StateForRight() int            // StateForRight returns the state for the right direction.
	StateForDown() int             // StateForDown returns the state for the down direction.
	StateForUp() int               // StateForUp returns the state for the up direction.
}

// GamepadIDButtonPressed returns a gamepad ID where at least one button is pressed.
// If no button is pressed, GamepadIDButtonPressed returns -1.
func (i *Input) GamepadIDButtonPressed() ebiten.GamepadID {
	i.gamepadIDs = ebiten.AppendGamepadIDs(i.gamepadIDs[:0])
	for _, id := range i.gamepadIDs {
		for b := ebiten.GamepadButton(0); b <= ebiten.GamepadButtonMax; b++ {
			if ebiten.IsGamepadButtonPressed(id, b) {
				return id
			}
		}
	}

	return -1
}

// stateForVirtualGamepadButton returns the current state of the specified virtual gamepad button.
//
// Parameters:
//   - b: The virtual gamepad button to check.
//
// Returns:
//   - int: The current state of the virtual gamepad button.
func (i *Input) stateForVirtualGamepadButton(b VirtualGamepadButton) int {
	if i.VirtualGamepadButtonStates == nil {
		return 0
	}
	return i.VirtualGamepadButtonStates[b]
}

// Update updates the state of all virtual gamepad buttons.
func (i *Input) Update() {
	if !i.GamepadConfig.IsGamepadIDInitialized() {
		return
	}

	if i.VirtualGamepadButtonStates == nil {
		i.VirtualGamepadButtonStates = map[VirtualGamepadButton]int{}
	}
	for _, b := range VirtualGamepadButtons {
		if !i.GamepadConfig.IsButtonPressed(b) {
			i.VirtualGamepadButtonStates[b] = 0
			continue
		}
		i.VirtualGamepadButtonStates[b]++
	}
}

// IsAbilityOneJustPressed checks if the first ability button was just pressed.
//
// Returns:
//   - bool: True if the first ability button was just pressed, otherwise false.
func (i *Input) IsAbilityOneJustPressed() bool {
	if inpututil.IsKeyJustPressed(ebiten.KeySpace) || inpututil.IsKeyJustPressed(ebiten.KeyX) {
		return true
	}
	return i.stateForVirtualGamepadButton(VirtualGamepadButtonButtonB) == 1
}

// IsAbilityTwoJustPressed checks if the second ability button was just pressed.
//
// Returns:
//   - bool: True if the second ability button was just pressed, otherwise false.
func (i *Input) IsAbilityTwoJustPressed() bool {
	if inpututil.IsKeyJustPressed(ebiten.KeyZ) || inpututil.IsKeyJustPressed(ebiten.KeyShiftLeft) {
		return true
	}
	return i.stateForVirtualGamepadButton(VirtualGamepadButtonButtonA) == 1
}

// StateForLeft returns the duration for which the left direction button has been pressed.
//
// Returns:
//   - int: The duration for which the left direction button has been pressed.
func (i *Input) StateForLeft() int {
	if v := inpututil.KeyPressDuration(ebiten.KeyArrowLeft); 0 < v {
		return v
	} else if v := inpututil.KeyPressDuration(ebiten.KeyA); 0 < v {
		return v
	}

	return i.stateForVirtualGamepadButton(VirtualGamepadButtonLeft)
}

// StateForRight returns the duration for which the right direction button has been pressed.
//
// Returns:
//   - int: The duration for which the right direction button has been pressed.
func (i *Input) StateForRight() int {
	if v := inpututil.KeyPressDuration(ebiten.KeyArrowRight); 0 < v {
		return v
	} else if v := inpututil.KeyPressDuration(ebiten.KeyD); 0 < v {
		return v
	}

	return i.stateForVirtualGamepadButton(VirtualGamepadButtonRight)
}

// StateForDown returns the duration for which the down direction button has been pressed.
//
// Returns:
//   - int: The duration for which the down direction button has been pressed.
func (i *Input) StateForDown() int {
	if v := inpututil.KeyPressDuration(ebiten.KeyArrowDown); 0 < v {
		return v
	} else if v := inpututil.KeyPressDuration(ebiten.KeyS); 0 < v {
		return v
	}

	return i.stateForVirtualGamepadButton(VirtualGamepadButtonDown)
}

// StateForUp returns the duration for which the up direction button has been pressed.
//
// Returns:
//   - int: The duration for which the up direction button has been pressed.
func (i *Input) StateForUp() int {
	if v := inpututil.KeyPressDuration(ebiten.KeyArrowUp); 0 < v {
		return v
	} else if v := inpututil.KeyPressDuration(ebiten.KeyW); 0 < v {
		return v
	}

	return i.stateForVirtualGamepadButton(VirtualGamepadButtonUp)
}
