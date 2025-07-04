// Package scenes provides the implementation of various game scenes,
// including the host lobby scene for setting up and managing multiplayer games.
package scenes

import (
	"image/color"
	"log"

	"github.com/ebitenui/ebitenui/widget"
	"szofttech.inf.elte.hu/szofttech-c-2024/group-08/ninjago/assets"
	"szofttech.inf.elte.hu/szofttech-c-2024/group-08/ninjago/game/multiplayer"
)

// NewHostLobby creates a new host lobby scene for setting up a multiplayer game.
// It initializes the game server, sets up the user interface, and manages player interactions.
//
// Parameters:
//   - ScreenWidth: The width of the game screen.
//   - ScreenHeight: The height of the game screen.
//   - mapPath: The path to the map file to be used in the game.
//
// Returns:
//   - Scene: The initialized host lobby scene.
func NewHostLobby(ScreenWidth int, ScreenHeight int, mapPath string) Scene {
	s := newLobbyScene(800, 600)
	log.Println("Init server")
	s.Server = multiplayer.NewGameServer()
	log.Println("Starting server")
	s.Server.GameInfo.Level = mapPath
	closeServer := s.Server.Run()
	log.Println("Server started")
	s.Server.GameInfo.Players[0] = multiplayer.ProtoPlayer{Username: "Host", X: 1, Y: 1, Color: color.RGBA{R: 0, G: 155, B: 150, A: 255}}
	s.players = &s.Server.GameInfo.Players

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
			closeServer()
		}),
	))

	navigationButtons.AddChild(widget.NewButton(
		widget.ButtonOpts.WidgetOpts(
			widget.WidgetOpts.LayoutData(widget.RowLayoutData{
				MaxWidth:  100,
				MaxHeight: 50,
				Position:  widget.RowLayoutPositionEnd,
			}),
		),

		widget.ButtonOpts.Image(assets.LoadButtonImage()),

		widget.ButtonOpts.Text("Play", assets.EbitenUIFont(20), &widget.ButtonTextColor{
			Idle: color.NRGBA{0xdf, 0xf4, 0xff, 0xff},
		}),

		widget.ButtonOpts.TextPadding(widget.Insets{
			Left:   30,
			Right:  30,
			Top:    5,
			Bottom: 5,
		}),

		widget.ButtonOpts.ClickedHandler(func(args *widget.ButtonClickedEventArgs) {
			s.playButtonPressed = true
		}),
	))

	s.ui.Container.AddChild(navigationButtons)

	return s
}
