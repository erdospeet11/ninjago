// Package scenes provides the implementation of various game scenes and their management,
// including the lobby scene for multiplayer games.
package scenes

import (
	"bytes"
	"image"
	"image/color"
	"log"
	"math"

	"github.com/ebitenui/ebitenui"
	uiimage "github.com/ebitenui/ebitenui/image"
	"github.com/ebitenui/ebitenui/widget"
	"github.com/hajimehoshi/ebiten/v2"
	"szofttech.inf.elte.hu/szofttech-c-2024/group-08/ninjago/assets"
	"szofttech.inf.elte.hu/szofttech-c-2024/group-08/ninjago/game/multiplayer"
)

// lobbyBackground is the background image used in the lobby scene.
var lobbyBackground *ebiten.Image

// init initializes the lobby background image from embedded assets.
func init() {
	lobbyBackgroundInBytes, err := assets.EmbeddedAssets.ReadFile("assets/graphics/lobbybackground-export.png")
	if err != nil {
		log.Println("Error reading file")
		panic(err)
	}

	img, _, err := image.Decode(bytes.NewReader(lobbyBackgroundInBytes))
	if err != nil {
		log.Println("Error decoding file")
		panic(err)
	}
	lobbyBackground = ebiten.NewImageFromImage(img)
}

// newPlayerWidget creates a new widget for displaying player information in the lobby.
//
// Parameters:
//   - username: The username of the player.
//   - pColor: The color associated with the player.
//
// Returns:
//   - *widget.Container: The container widget displaying the player information.
func newPlayerWidget(username string, pColor color.Color) *widget.Container {
	playerWidget := widget.NewContainer(
		widget.ContainerOpts.BackgroundImage(uiimage.NewNineSliceColor(pColor)),
		widget.ContainerOpts.Layout(widget.NewRowLayout(
			widget.RowLayoutOpts.Direction(widget.DirectionHorizontal),
			widget.RowLayoutOpts.Padding(widget.NewInsetsSimple(5)),
			widget.RowLayoutOpts.Spacing(10),
		)),
	)

	playerWidget.AddChild(
		widget.NewGraphic(
			widget.GraphicOpts.Image(assets.LoadAndScaleImage("assets/assets/player/character-down.png", 50, 50)),
			widget.GraphicOpts.WidgetOpts(
				widget.WidgetOpts.LayoutData(widget.RowLayoutData{
					Position: widget.RowLayoutPositionStart,
					Stretch:  true,
				}),
			),
		),
	)

	playerWidget.AddChild(widget.NewText(
		widget.TextOpts.WidgetOpts(
			widget.WidgetOpts.LayoutData(widget.RowLayoutData{
				Position:  widget.RowLayoutPositionCenter,
				Stretch:   true,
				MaxHeight: 25,
				MaxWidth:  200,
			}),
		),

		widget.TextOpts.TextFace(assets.EbitenUIFont(10)),

		widget.TextOpts.TextColor(color.NRGBA{254, 255, 255, 255}),

		widget.TextOpts.Insets(widget.NewInsetsSimple(5)),

		widget.TextOpts.TextLabel(username),
	))
	return playerWidget
}

// lobbyScene represents the scene for the multiplayer lobby, where players can join and prepare for the game.
type lobbyScene struct {
	players           *[]multiplayer.ProtoPlayer // The list of players in the lobby.
	count             int                        // The count for the background animation.
	screenHeight      int                        // The height of the screen.
	screenWidth       int                        // The width of the screen.
	ui                *ebitenui.UI               // The user interface for the lobby scene.
	playerList        *widget.ScrollContainer    // The scroll container for the player list.
	content           *widget.Container          // The container for the player list content.
	playButtonPressed bool                       // Indicates if the play button has been pressed.
	backButtonPressed bool                       // Indicates if the back button has been pressed.
	Server            *multiplayer.GameServer    // The game server for hosting the game.
	Client            *multiplayer.GameClient    // The game client for joining the game.
}

// newLobbyScene creates a new lobby scene with the specified screen dimensions.
//
// Parameters:
//   - ScreenWidth: The width of the screen.
//   - ScreenHeight: The height of the screen.
//
// Returns:
//   - *lobbyScene: The initialized lobby scene.
func newLobbyScene(ScreenWidth int, ScreenHeight int) *lobbyScene {
	rootContainer := widget.NewContainer(
		widget.ContainerOpts.Layout(widget.NewGridLayout(
			widget.GridLayoutOpts.Columns(2),
			widget.GridLayoutOpts.Spacing(2, 10),
			widget.GridLayoutOpts.Stretch([]bool{true, false}, []bool{true}),
		)),
	)

	content := widget.NewContainer(widget.ContainerOpts.Layout(widget.NewRowLayout(
		widget.RowLayoutOpts.Direction(widget.DirectionVertical),
		widget.RowLayoutOpts.Spacing(5),
	)))

	playerlist := widget.NewScrollContainer(
		widget.ScrollContainerOpts.Content(content),
		widget.ScrollContainerOpts.StretchContentWidth(),
		widget.ScrollContainerOpts.Image(&widget.ScrollContainerImage{
			Idle: uiimage.NewNineSliceColor(color.NRGBA{0x13, 0x1a, 0x22, 0}),
			Mask: uiimage.NewNineSliceColor(color.NRGBA{0x13, 0x1a, 0x22, 0xff}),
		}),
	)

	rootContainer.AddChild(playerlist)

	lobby := &lobbyScene{
		screenWidth:  ScreenWidth,
		screenHeight: ScreenHeight,
		playerList:   playerlist,
		content:      content,
	}
	lobby.ui = &ebitenui.UI{
		Container: rootContainer,
	}

	pageSizeFunc := func() int {
		return int(math.Round(float64(playerlist.ViewRect().Dy()) / float64(content.GetWidget().Rect.Dy()) * 1000))
	}

	vSlider := widget.NewSlider(
		widget.SliderOpts.Direction(widget.DirectionVertical),
		widget.SliderOpts.MinMax(0, 1000),
		widget.SliderOpts.PageSizeFunc(pageSizeFunc),

		widget.SliderOpts.ChangedHandler(func(args *widget.SliderChangedEventArgs) {
			playerlist.ScrollTop = float64(args.Slider.Current) / 1000
		}),
		widget.SliderOpts.Images(
			&widget.SliderTrackImage{
				Idle:  uiimage.NewNineSliceColor(color.NRGBA{55, 148, 110, 255}),
				Hover: uiimage.NewNineSliceColor(color.NRGBA{55, 148, 110, 255}),
			},
			&widget.ButtonImage{
				Idle:    uiimage.NewNineSliceColor(color.NRGBA{5, 66, 47, 255}),
				Hover:   uiimage.NewNineSliceColor(color.NRGBA{5, 66, 47, 255}),
				Pressed: uiimage.NewNineSliceColor(color.NRGBA{5, 66, 47, 255}),
			},
		),
	)

	playerlist.GetWidget().ScrolledEvent.AddHandler(func(args interface{}) {
		a := args.(*widget.WidgetScrolledEventArgs)
		p := pageSizeFunc() / 3
		if p < 1 {
			p = 1
		}
		vSlider.Current -= int(math.Round(a.Y * float64(p)))
	})

	rootContainer.AddChild(vSlider)

	return lobby
}

// Update updates the lobby scene, handling player list updates, server/client interactions,
// and transitioning to the game scene when the game starts.
//
// Parameters:
//   - state: The current game state.
//
// Returns:
//   - error: An error if the update fails.
func (s *lobbyScene) Update(state *GameState) error {
	s.ui.Update()

	for i, player := range *s.players {
		if len(s.content.Children()) < i+1 {
			s.content.AddChild(newPlayerWidget(player.Username, player.Color))
		}
	}
	if len(s.content.Children()) > len(*s.players) {
		s.content.RemoveChildren()
		for _, player := range *s.players {
			s.content.AddChild(newPlayerWidget(player.Username, player.Color))
		}
	}

	if s.Server != nil {
		s.Server.GameInfo.Players[0].Username = state.UserInfo.Username
	}

	if s.playButtonPressed && s.Server != nil {
		state.SceneManager.GoTo(NewMultiPlayerGameSceneHost(s.Server, s.Server.GameInfo.Level, state.UserInfo))

		return nil
	}

	if s.Client != nil && s.Client.GameInfo.GameState == multiplayer.GameStateRunning {
		log.Println("Game started")
		state.SceneManager.GoTo(NewMultiPlayerGameSceneJoin(s.Client))

		return nil
	}

	if s.backButtonPressed {
		state.SceneManager.GoTo(NewMainMenuScene(s.screenWidth, s.screenHeight))

		return nil
	}

	s.count++

	return nil
}

// Draw renders the lobby scene onto the provided screen image.
//
// Parameters:
//   - r: The image to which the scene is drawn.
func (s *lobbyScene) Draw(r *ebiten.Image) {
	s.drawTitleBackground(r, s.count)

	s.ui.Draw(r)
}

// drawTitleBackground draws the scrolling background for the lobby scene.
//
// Parameters:
//   - r: The image to which the background is drawn.
//   - c: The counter used to animate the background.
func (s *lobbyScene) drawTitleBackground(r *ebiten.Image, c int) {
	w, h := lobbyBackground.Bounds().Dx(), lobbyBackground.Bounds().Dy()
	op := &ebiten.DrawImageOptions{}
	for i := 0; i < (s.screenWidth/w+1)*(s.screenHeight/h+2); i++ {
		op.GeoM.Reset()
		dx := -(c / 4) % w
		dy := (c / 4) % h
		dstX := (i%(s.screenWidth/w+1))*w + dx
		dstY := (i/(s.screenWidth/w+1)-1)*h + dy
		op.GeoM.Translate(float64(dstX), float64(dstY))
		r.DrawImage(lobbyBackground, op)
	}

	// Loop the background vertically
	if c/4%h != 0 {
		op.GeoM.Reset()
		dx := -(c / 4) % w
		dy := (c/4)%h - h
		dstX := (s.screenWidth/w)*w + dx
		dstY := ((s.screenHeight/h+2)-1)*h + dy
		op.GeoM.Translate(float64(dstX), float64(dstY))
		r.DrawImage(lobbyBackground, op)
	}
}
