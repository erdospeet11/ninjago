// Package scenes provides the implementation of various game scenes,
// including the client lobby scene for joining and managing multiplayer games.
package scenes

import (
	"image/color"

	"github.com/ebitenui/ebitenui/widget"
	"szofttech.inf.elte.hu/szofttech-c-2024/group-08/ninjago/assets"
	"szofttech.inf.elte.hu/szofttech-c-2024/group-08/ninjago/game/multiplayer"
)

// NewClientLobby creates a new client lobby scene for joining a multiplayer game.
// It initializes the game client, sets up the user interface, and manages player interactions.
//
// Parameters:
//   - ScreenWidth: The width of the game screen.
//   - ScreenHeight: The height of the game screen.
//   - client: The game client used to connect to the multiplayer game.
//   - State: The current game state.
//
// Returns:
//   - Scene: The initialized client lobby scene.
func NewClientLobby(ScreenWidth int, ScreenHeight int, client *multiplayer.GameClient, State *GameState) Scene {
	s := newLobbyScene(800, 600)
	s.Client = client
	s.players = &s.Client.GameInfo.Players

	navigationButtons := widget.NewContainer(
		widget.ContainerOpts.Layout(widget.NewRowLayout(
			widget.RowLayoutOpts.Direction(widget.DirectionHorizontal),
			widget.RowLayoutOpts.Spacing(195),
			widget.RowLayoutOpts.Padding(widget.NewInsetsSimple(5)),
		)),
	)

	navigationButtons.AddChild(widget.NewButton(
		widget.ButtonOpts.WidgetOpts(
			widget.WidgetOpts.LayoutData(widget.RowLayoutData{
				MaxWidth:  100,
				MaxHeight: 50,
				Position:  widget.RowLayoutPositionStart,
			}),
		),

		widget.ButtonOpts.Image(assets.LoadButtonImage()),

		widget.ButtonOpts.Text("Back", assets.EbitenUIFont(20), &widget.ButtonTextColor{
			Idle: color.NRGBA{0xdf, 0xf4, 0xff, 0xff},
		}),

		widget.ButtonOpts.TextPadding(widget.Insets{
			Left:   30,
			Right:  30,
			Top:    5,
			Bottom: 5,
		}),

		widget.ButtonOpts.ClickedHandler(func(args *widget.ButtonClickedEventArgs) {
			s.backButtonPressed = true
			s.Client.Close()
		}),
	))

	s.ui.Container.AddChild(navigationButtons)

	return s
}
