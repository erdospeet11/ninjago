// Package inputs provides a set of utilities for handling input devices.
package inputs

import (
	"fmt"

	ebiten "github.com/hajimehoshi/ebiten/v2"
	inpututil "github.com/hajimehoshi/ebiten/v2/inpututil"
)

// VirtualGamepadButton represents a virtual gamepad button.
type VirtualGamepadButton int

// VirtualGamepadButton values.
const (
	VirtualGamepadButtonLeft VirtualGamepadButton = iota
	VirtualGamepadButtonRight
	VirtualGamepadButtonDown
	VirtualGamepadButtonUp
	VirtualGamepadButtonButtonA
	VirtualGamepadButtonButtonB
)

// VirtualGamepadButtons is a list of all virtual gamepad buttons.
var VirtualGamepadButtons = []VirtualGamepadButton{
	VirtualGamepadButtonLeft,
	VirtualGamepadButtonRight,
	VirtualGamepadButtonDown,
	VirtualGamepadButtonUp,
	VirtualGamepadButtonButtonA,
	VirtualGamepadButtonButtonB,
}

// StandardGamepadButton returns the standard gamepad button corresponding to the virtual gamepad button.
func (v VirtualGamepadButton) StandardGamepadButton() ebiten.StandardGamepadButton {
	switch v {
	case VirtualGamepadButtonLeft:
		return ebiten.StandardGamepadButtonLeftLeft
	case VirtualGamepadButtonRight:
		return ebiten.StandardGamepadButtonLeftRight
	case VirtualGamepadButtonDown:
		return ebiten.StandardGamepadButtonLeftBottom
	case VirtualGamepadButtonUp:
		return ebiten.StandardGamepadButtonLeftTop
	case VirtualGamepadButtonButtonA:
		return ebiten.StandardGamepadButtonRightBottom
	case VirtualGamepadButtonButtonB:
		return ebiten.StandardGamepadButtonRightRight
	default:
		panic("not reached")
	}
}

// GamepadConfig represents the configuration of a gamepad's buttons and axes.
const axisThreshold = 0.75

// axis represents an axis of a gamepad.
type axis struct {
	id       int  // The ID of the axis.
	positive bool // Indicates if the axis is positive.
}

// GamepadConfig represents the configuration of a gamepad's buttons and axes.
type GamepadConfig struct {
	gamepadID            ebiten.GamepadID // The ID of the gamepad to configure.
	gamepadIDInitialized bool             // Indicates if the gamepad ID has been initialized.

	current         VirtualGamepadButton                          // The current virtual gamepad button being configured.
	buttons         map[VirtualGamepadButton]ebiten.GamepadButton // The mapping of virtual gamepad buttons to physical buttons.
	axes            map[VirtualGamepadButton]axis                 // The mapping of virtual gamepad buttons to axes.
	assignedButtons map[ebiten.GamepadButton]struct{}             // The set of assigned buttons.
	assignedAxes    map[axis]struct{}                             // The set of assigned axes.

	defaultAxesValues map[int]float64 // The default values of the axes.
}

// SetGamepadID sets the ID of the gamepad to be configured.
//
// Parameters:
//   - id: The ID of the gamepad to set.
func (c *GamepadConfig) SetGamepadID(id ebiten.GamepadID) {
	c.gamepadID = id
	c.gamepadIDInitialized = true
}

// ResetGamepadID resets the ID of the gamepad, clearing any initialization status.
func (c *GamepadConfig) ResetGamepadID() {
	c.gamepadID = 0
	c.gamepadIDInitialized = false
}

// IsGamepadIDInitialized checks if the gamepad ID has been initialized.
//
// Returns:
//   - bool: True if the gamepad ID is initialized, otherwise false.
func (c *GamepadConfig) IsGamepadIDInitialized() bool {
	return c.gamepadIDInitialized
}

// NeedsConfiguration checks if the gamepad needs custom configuration.
//
// Returns:
//   - bool: True if the gamepad needs custom configuration, otherwise false.
func (c *GamepadConfig) NeedsConfiguration() bool {
	return !ebiten.IsStandardGamepadLayoutAvailable(c.gamepadID)
}

// initializeIfNeeded initializes the gamepad configuration if it has not been initialized yet.
func (c *GamepadConfig) initializeIfNeeded() {
	if !c.gamepadIDInitialized {
		panic("not reached")
	}

	if ebiten.IsStandardGamepadLayoutAvailable(c.gamepadID) {
		return
	}

	if c.buttons == nil {
		c.buttons = map[VirtualGamepadButton]ebiten.GamepadButton{}
	}
	if c.axes == nil {
		c.axes = map[VirtualGamepadButton]axis{}
	}
	if c.assignedButtons == nil {
		c.assignedButtons = map[ebiten.GamepadButton]struct{}{}
	}
	if c.assignedAxes == nil {
		c.assignedAxes = map[axis]struct{}{}
	}

	// Set default values.
	// It is assumed that all axes are not pressed here.
	//
	// These default values are used to detect if an axis is actually pressed.
	// For example, on PS4 controllers, L2/R2's axes value can be -1.0.
	if c.defaultAxesValues == nil {
		c.defaultAxesValues = map[int]float64{}
		na := int(ebiten.GamepadAxisCount(c.gamepadID))
		for a := int(0); a < na; a++ {
			c.defaultAxesValues[a] = ebiten.GamepadAxisValue(c.gamepadID, a)
		}
	}
}

// Reset clears the current configuration of the gamepad, including buttons and axes assignments.
func (c *GamepadConfig) Reset() {
	c.buttons = nil
	c.axes = nil
	c.assignedButtons = nil
	c.assignedAxes = nil
}

// Scan scans the current input state and assigns the given virtual gamepad button b
// to the current (physical) pressed buttons of the gamepad.
//
// Parameters:
//   - b: The virtual gamepad button to configure.
//
// Returns:
//   - bool: True if a button or axis is successfully assigned, otherwise false.
func (c *GamepadConfig) Scan(b VirtualGamepadButton) bool {
	if !c.gamepadIDInitialized {
		panic("not reached")
	}

	c.initializeIfNeeded()

	delete(c.buttons, b)
	delete(c.axes, b)

	ebn := ebiten.GamepadButton(ebiten.GamepadButtonCount(c.gamepadID))
	for eb := ebiten.GamepadButton(0); eb < ebn; eb++ {
		if _, ok := c.assignedButtons[eb]; ok {
			continue
		}
		if inpututil.IsGamepadButtonJustPressed(c.gamepadID, eb) {
			c.buttons[b] = eb
			c.assignedButtons[eb] = struct{}{}
			return true
		}
	}

	na := int(ebiten.GamepadAxisCount(c.gamepadID))
	for a := int(0); a < na; a++ {
		v := ebiten.GamepadAxisValue(c.gamepadID, a)
		const delta = 0.25

		// Check |v| < 1.0 because there is a bug that a button returns
		// an axis value wrongly and the value may be over 1 on some platforms.
		if axisThreshold <= v && v <= 1.0 &&
			(v < c.defaultAxesValues[a]-delta || c.defaultAxesValues[a]+delta < v) {
			if _, ok := c.assignedAxes[axis{a, true}]; !ok {
				c.axes[b] = axis{a, true}
				c.assignedAxes[axis{a, true}] = struct{}{}
				return true
			}
		}
		if -1.0 <= v && v <= -axisThreshold &&
			(v < c.defaultAxesValues[a]-delta || c.defaultAxesValues[a]+delta < v) {
			if _, ok := c.assignedAxes[axis{a, false}]; !ok {
				c.axes[b] = axis{a, false}
				c.assignedAxes[axis{a, false}] = struct{}{}
				return true
			}
		}
	}

	return false
}

// IsButtonPressed reports whether the given virtual button b is pressed.
//
// Parameters:
//   - b: The virtual gamepad button to check.
//
// Returns:
//   - bool: True if the button is pressed, otherwise false.
func (c *GamepadConfig) IsButtonPressed(b VirtualGamepadButton) bool {
	if !c.gamepadIDInitialized {
		panic("not reached")
	}

	if ebiten.IsStandardGamepadLayoutAvailable(c.gamepadID) {
		if ebiten.IsStandardGamepadButtonPressed(c.gamepadID, b.StandardGamepadButton()) {
			return true
		}

		const threshold = 0.7
		switch b {
		case VirtualGamepadButtonLeft:
			return ebiten.StandardGamepadAxisValue(c.gamepadID, ebiten.StandardGamepadAxisLeftStickHorizontal) < -threshold
		case VirtualGamepadButtonRight:
			return ebiten.StandardGamepadAxisValue(c.gamepadID, ebiten.StandardGamepadAxisLeftStickHorizontal) > threshold
		case VirtualGamepadButtonDown:
			return ebiten.StandardGamepadAxisValue(c.gamepadID, ebiten.StandardGamepadAxisLeftStickVertical) > threshold
		case VirtualGamepadButtonUp:
			return ebiten.StandardGamepadAxisValue(c.gamepadID, ebiten.StandardGamepadAxisLeftStickVertical) < -threshold
		default:
			return false
		}
	}

	c.initializeIfNeeded()

	bb, ok := c.buttons[b]
	if ok {
		return ebiten.IsGamepadButtonPressed(c.gamepadID, bb)
	}

	a, ok := c.axes[b]
	if ok {
		v := ebiten.GamepadAxisValue(c.gamepadID, a.id)
		if a.positive {
			return axisThreshold <= v && v <= 1.0
		}
		return -1.0 <= v && v <= -axisThreshold
	}
	return false
}

// IsButtonJustPressed reports whether the given virtual button b started to be pressed now.
//
// Parameters:
//   - b: The virtual gamepad button to check.
//
// Returns:
//   - bool: True if the button was just pressed, otherwise false.
func (c *GamepadConfig) IsButtonJustPressed(b VirtualGamepadButton) bool {
	if !c.gamepadIDInitialized {
		panic("not reached")
	}

	if ebiten.IsStandardGamepadLayoutAvailable(c.gamepadID) {
		return inpututil.IsStandardGamepadButtonJustPressed(c.gamepadID, b.StandardGamepadButton())
	}

	c.initializeIfNeeded()

	bb, ok := c.buttons[b]
	if ok {
		return inpututil.IsGamepadButtonJustPressed(c.gamepadID, bb)
	}
	return false
}

// ButtonName returns the physical button's name for the given virtual button.
//
// Parameters:
//   - b: The virtual gamepad button to get the name for.
//
// Returns:
//   - string: The name of the physical button or axis assigned to the virtual button.
func (c *GamepadConfig) ButtonName(b VirtualGamepadButton) string {
	if !c.gamepadIDInitialized {
		panic("not reached")
	}

	c.initializeIfNeeded()

	bb, ok := c.buttons[b]
	if ok {
		return fmt.Sprintf("Button %d", bb)
	}

	a, ok := c.axes[b]
	if ok {
		if a.positive {
			return fmt.Sprintf("Axis %d+", a.id)
		}
		return fmt.Sprintf("Axis %d-", a.id)
	}

	return ""
}
