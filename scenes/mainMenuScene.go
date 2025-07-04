// Package scenes provides the implementation of various game scenes and their management,
// including the main menu scene for the game.
package scenes

import (
	"bytes"
	"image"
	"image/color"
	"log"
	"os"
	"unicode/utf8"

	"github.com/ebitenui/ebitenui"
	uiimage "github.com/ebitenui/ebitenui/image"
	"github.com/ebitenui/ebitenui/widget"
	ebiten "github.com/hajimehoshi/ebiten/v2"
	text "github.com/hajimehoshi/ebiten/v2/text/v2"
	"szofttech.inf.elte.hu/szofttech-c-2024/group-08/ninjago/assets"
	"szofttech.inf.elte.hu/szofttech-c-2024/group-08/ninjago/game/multiplayer"
	"szofttech.inf.elte.hu/szofttech-c-2024/group-08/ninjago/inputs"
)

// imageBackground is the background image used in the main menu scene.
var imageBackground *ebiten.Image

// init initializes the background image from embedded assets.
func init() {
	imageBackgroundInBytes, err := assets.EmbeddedAssets.ReadFile("assets/graphics/background.png")
	if err != nil {
		panic(err)
	}

	img, _, err := image.Decode(bytes.NewReader(imageBackgroundInBytes))
	if err != nil {
		panic(err)
	}
	imageBackground = ebiten.NewImageFromImage(img)
}

// mainMenuScene represents the main menu scene of the game, where players can host or join games.
type mainMenuScene struct {
	count             int               // The count of the main menu scene.
	screenHeight      int               // The height of the screen.
	screenWidth       int               // The width of the screen.
	ui                *ebitenui.UI      // The user interface for the main menu scene.
	nameInput         *widget.TextInput // The text input for the player's name.
	addressToJoin     string            // The address to join a game.
	map1ButtonPressed bool              // The flag indicating whether the map1 button is pressed.
	map2ButtonPressed bool              // The flag indicating whether the map2 button is pressed.
	map3ButtonPressed bool              // The flag indicating whether the map3 button is pressed.
	joinGame          bool              // The flag indicating whether the player is joining a game.
}

// NewMainMenuScene creates a new main menu scene with the specified screen dimensions.
//
// Parameters:
//   - ScreenWidth: The width of the screen.
//   - ScreenHeight: The height of the screen.
//
// Returns:
//   - Scene: The initialized main menu scene.
func NewMainMenuScene(ScreenWidth int, ScreenHeight int) Scene {
	rootContainer := widget.NewContainer(
		widget.ContainerOpts.Layout(
			widget.NewGridLayout(
				widget.GridLayoutOpts.Columns(1),
				widget.GridLayoutOpts.Stretch([]bool{true}, []bool{true, true}),
				widget.GridLayoutOpts.Spacing(0, 5),
				widget.GridLayoutOpts.Padding(widget.Insets{
					Top:    20,
					Bottom: 20,
				}),
			),
		),
	)

	mainMenu := &mainMenuScene{
		screenWidth:  ScreenWidth,
		screenHeight: ScreenHeight,
	}
	mainMenu.ui = &ebitenui.UI{
		Container: rootContainer,
	}

	//JoinWindow
	joinwindowContainer := widget.NewContainer(
		widget.ContainerOpts.BackgroundImage(uiimage.NewNineSliceColor(color.NRGBA{19, 128, 30, 255})),
		widget.ContainerOpts.Layout(widget.NewRowLayout(
			widget.RowLayoutOpts.Direction(widget.DirectionVertical),
			widget.RowLayoutOpts.Padding(widget.NewInsetsSimple(5)),
			widget.RowLayoutOpts.Spacing(10),
		)),
	)
	joinwindowContainer.AddChild(widget.NewTextInput(
		widget.TextInputOpts.WidgetOpts(
			widget.WidgetOpts.LayoutData(widget.RowLayoutData{
				Position:  widget.RowLayoutPositionCenter,
				Stretch:   true,
				MaxHeight: 25,
				MaxWidth:  200,
			}),
		),
		widget.TextInputOpts.Image(&widget.TextInputImage{
			Idle:     uiimage.NewNineSliceColor(color.NRGBA{R: 75, G: 105, B: 47, A: 255}),
			Disabled: uiimage.NewNineSliceColor(color.NRGBA{R: 100, G: 100, B: 100, A: 255}),
		}),

		widget.TextInputOpts.Face(assets.EbitenUIFont(10)),

		widget.TextInputOpts.Color(&widget.TextInputColor{
			Idle:          color.NRGBA{255, 255, 255, 255},
			Disabled:      color.NRGBA{R: 200, G: 200, B: 200, A: 255},
			Caret:         color.NRGBA{254, 255, 255, 255},
			DisabledCaret: color.NRGBA{R: 200, G: 200, B: 200, A: 255},
		}),

		widget.TextInputOpts.Padding(widget.NewInsetsSimple(5)),

		widget.TextInputOpts.SubmitHandler(func(args *widget.TextInputChangedEventArgs) {
			if args.InputText != "" {
				mainMenu.addressToJoin = args.InputText

				mainMenu.joinGame = true
			}
		}),

		widget.TextInputOpts.CaretOpts(
			widget.CaretOpts.Size(assets.EbitenUIFont(10), 2),
		),

		widget.TextInputOpts.Placeholder("IP Address"),
	))

	joinWindowContent := widget.NewContainer(
		widget.ContainerOpts.BackgroundImage(uiimage.NewNineSliceColor(color.NRGBA{75, 105, 47, 255})),
		widget.ContainerOpts.Layout(widget.NewAnchorLayout()),
	)
	joinWindowContent.AddChild(widget.NewText(
		widget.TextOpts.Text("Enter IP Address to join a server", assets.EbitenUIFont(12), color.NRGBA{254, 255, 255, 255}),
		widget.TextOpts.WidgetOpts(widget.WidgetOpts.LayoutData(widget.AnchorLayoutData{
			HorizontalPosition: widget.AnchorLayoutPositionCenter,
			VerticalPosition:   widget.AnchorLayoutPositionCenter,
		})),
	))

	joinwindow := widget.NewWindow(
		widget.WindowOpts.Contents(joinwindowContainer),
		widget.WindowOpts.TitleBar(joinWindowContent, 25),
		widget.WindowOpts.Modal(),
		widget.WindowOpts.CloseMode(widget.CLICK_OUT),
		widget.WindowOpts.Draggable(),
		widget.WindowOpts.Resizeable(),
		widget.WindowOpts.MinSize(200, 100),
		widget.WindowOpts.MaxSize(300, 300),
	)

	//hostWindow
	hostwindowContainer := widget.NewContainer(
		widget.ContainerOpts.BackgroundImage(uiimage.NewNineSliceColor(color.NRGBA{19, 128, 30, 255})),
		widget.ContainerOpts.Layout(widget.NewRowLayout(
			widget.RowLayoutOpts.Direction(widget.DirectionVertical),
			widget.RowLayoutOpts.Padding(widget.NewInsetsSimple(5)),
			widget.RowLayoutOpts.Spacing(5),
		)),
	)

	innerContainer1 := widget.NewContainer(
		widget.ContainerOpts.Layout(widget.NewRowLayout(
			widget.RowLayoutOpts.Direction(widget.DirectionHorizontal),
			widget.RowLayoutOpts.Padding(widget.NewInsetsSimple(5)),
			widget.RowLayoutOpts.Spacing(75),
		)),
		widget.ContainerOpts.WidgetOpts(
			widget.WidgetOpts.LayoutData(widget.RowLayoutData{
				Position: widget.RowLayoutPositionCenter,
				Stretch:  true,
			}),
		),
	)

	button1 := widget.NewButton(
		widget.ButtonOpts.WidgetOpts(
			widget.WidgetOpts.LayoutData(widget.RowLayoutData{
				Position:  widget.RowLayoutPositionCenter,
				Stretch:   true,
				MaxHeight: 50,
				MaxWidth:  75,
			}),
		),
		widget.ButtonOpts.Image(assets.LoadButtonImage()),
		widget.ButtonOpts.Text("Select Map1", assets.EbitenUIFont(10), &widget.ButtonTextColor{
			Idle: color.NRGBA{0xdf, 0xf4, 0xff, 0xff},
		}),
		widget.ButtonOpts.TextPadding(widget.Insets{
			Left:   30,
			Right:  30,
			Top:    5,
			Bottom: 5,
		}),
		widget.ButtonOpts.ClickedHandler(func(args *widget.ButtonClickedEventArgs) {
			mainMenu.map1ButtonPressed = true
		}),
	)

	button2 := widget.NewButton(
		widget.ButtonOpts.WidgetOpts(
			widget.WidgetOpts.LayoutData(widget.RowLayoutData{
				Position:  widget.RowLayoutPositionCenter,
				Stretch:   true,
				MaxHeight: 50,
				MaxWidth:  75,
			}),
		),
		widget.ButtonOpts.Image(assets.LoadButtonImage()),
		widget.ButtonOpts.Text("Select Map2", assets.EbitenUIFont(10), &widget.ButtonTextColor{
			Idle: color.NRGBA{0xdf, 0xf4, 0xff, 0xff},
		}),
		widget.ButtonOpts.TextPadding(widget.Insets{
			Left:   30,
			Right:  30,
			Top:    5,
			Bottom: 5,
		}),
		widget.ButtonOpts.ClickedHandler(func(args *widget.ButtonClickedEventArgs) {
			mainMenu.map2ButtonPressed = true
		}),
	)

	button3 := widget.NewButton(
		widget.ButtonOpts.WidgetOpts(
			widget.WidgetOpts.LayoutData(widget.RowLayoutData{
				Position:  widget.RowLayoutPositionCenter,
				Stretch:   true,
				MaxHeight: 50,
				MaxWidth:  75,
			}),
		),
		widget.ButtonOpts.Image(assets.LoadButtonImage()),
		widget.ButtonOpts.Text("Select Map3", assets.EbitenUIFont(10), &widget.ButtonTextColor{
			Idle: color.NRGBA{0xdf, 0xf4, 0xff, 0xff},
		}),
		widget.ButtonOpts.TextPadding(widget.Insets{
			Left:   30,
			Right:  30,
			Top:    5,
			Bottom: 5,
		}),
		widget.ButtonOpts.ClickedHandler(func(args *widget.ButtonClickedEventArgs) {
			mainMenu.map3ButtonPressed = true
		}),
	)

	innerContainer1.AddChild(button1)
	innerContainer1.AddChild(button2)
	innerContainer1.AddChild(button3)

	innerContainer2 := widget.NewContainer(
		widget.ContainerOpts.Layout(widget.NewRowLayout(
			widget.RowLayoutOpts.Direction(widget.DirectionHorizontal),
			widget.RowLayoutOpts.Padding(widget.NewInsetsSimple(5)),
			widget.RowLayoutOpts.Spacing(40),
		)),
		widget.ContainerOpts.WidgetOpts(
			widget.WidgetOpts.LayoutData(widget.RowLayoutData{
				Position: widget.RowLayoutPositionCenter,
				Stretch:  true,
			}),
		),
	)

	innerContainer2.AddChild(widget.NewGraphic(
		widget.GraphicOpts.Image(assets.LoadAndScaleImage("assets/assets/graphics/map.png", 100, 100)),
		widget.GraphicOpts.WidgetOpts(
			widget.WidgetOpts.LayoutData(widget.RowLayoutData{
				Position: widget.RowLayoutPositionCenter,
				Stretch:  false,
			}),
		),
	))

	innerContainer2.AddChild(widget.NewGraphic(
		widget.GraphicOpts.Image(assets.LoadAndScaleImage("assets/assets/graphics/map2.png", 100, 100)),
		widget.GraphicOpts.WidgetOpts(
			widget.WidgetOpts.LayoutData(widget.RowLayoutData{
				Position: widget.RowLayoutPositionCenter,
				Stretch:  false,
			}),
		),
	))

	innerContainer2.AddChild(widget.NewGraphic(
		widget.GraphicOpts.Image(assets.LoadAndScaleImage("assets/assets/graphics/map3.png", 100, 100)),
		widget.GraphicOpts.WidgetOpts(
			widget.WidgetOpts.LayoutData(widget.RowLayoutData{
				Position: widget.RowLayoutPositionCenter,
				Stretch:  false,
			}),
		),
	))

	innerContainer3 := widget.NewContainer(
		widget.ContainerOpts.Layout(widget.NewRowLayout(
			widget.RowLayoutOpts.Direction(widget.DirectionVertical),
			widget.RowLayoutOpts.Padding(widget.NewInsetsSimple(5)),
			widget.RowLayoutOpts.Spacing(40),
		)),
		widget.ContainerOpts.WidgetOpts(
			widget.WidgetOpts.LayoutData(widget.RowLayoutData{
				Position: widget.RowLayoutPositionCenter,
				Stretch:  true,
			}),
		),
	)

	button4 := widget.NewButton(
		widget.ButtonOpts.WidgetOpts(
			widget.WidgetOpts.LayoutData(widget.RowLayoutData{
				Position:  widget.RowLayoutPositionCenter,
				Stretch:   false,
				MaxWidth:  100,
				MaxHeight: 50,
			}),
		),
		widget.ButtonOpts.Image(assets.LoadButtonImage()),
		widget.ButtonOpts.Text("Close", assets.EbitenUIFont(10), &widget.ButtonTextColor{
			Idle: color.NRGBA{0xdf, 0xf4, 0xff, 0xff},
		}),
		widget.ButtonOpts.TextPadding(widget.Insets{
			Left:   30,
			Right:  30,
			Top:    5,
			Bottom: 5,
		}),
	)

	innerContainer3.AddChild(button4)

	hostwindowContainer.AddChild(innerContainer1)
	hostwindowContainer.AddChild(innerContainer2)
	hostwindowContainer.AddChild(innerContainer3)

	hostWindowContent := widget.NewContainer(
		widget.ContainerOpts.BackgroundImage(uiimage.NewNineSliceColor(color.NRGBA{75, 105, 47, 255})),
		widget.ContainerOpts.Layout(widget.NewAnchorLayout()),
	)
	hostWindowContent.AddChild(widget.NewText(
		widget.TextOpts.Text("Select map for new game", assets.EbitenUIFont(12), color.NRGBA{254, 255, 255, 255}),
		widget.TextOpts.WidgetOpts(widget.WidgetOpts.LayoutData(widget.AnchorLayoutData{
			HorizontalPosition: widget.AnchorLayoutPositionCenter,
			VerticalPosition:   widget.AnchorLayoutPositionCenter,
		})),
	))

	hostwindow := widget.NewWindow(
		widget.WindowOpts.Contents(hostwindowContainer),
		widget.WindowOpts.TitleBar(hostWindowContent, 25),
		widget.WindowOpts.Modal(),
		widget.WindowOpts.CloseMode(widget.CLICK_OUT),
		widget.WindowOpts.Draggable(),
		widget.WindowOpts.Resizeable(),
		widget.WindowOpts.MinSize(250, 150),
		widget.WindowOpts.MaxSize(450, 750),
	)

	mainMenu.nameInput = widget.NewTextInput(
		widget.TextInputOpts.WidgetOpts(
			widget.WidgetOpts.LayoutData(widget.GridLayoutData{
				MaxWidth:           175,
				MaxHeight:          25,
				VerticalPosition:   widget.GridLayoutPositionEnd,
				HorizontalPosition: widget.GridLayoutPositionCenter,
			}),
		),
		widget.TextInputOpts.Image(&widget.TextInputImage{
			Idle:     uiimage.NewNineSliceColor(color.NRGBA{R: 75, G: 105, B: 47, A: 255}),
			Disabled: uiimage.NewNineSliceColor(color.NRGBA{R: 100, G: 100, B: 100, A: 255}),
		}),

		widget.TextInputOpts.Face(assets.EbitenUIFont(10)),

		widget.TextInputOpts.Color(&widget.TextInputColor{
			Idle:          color.NRGBA{254, 255, 255, 255},
			Disabled:      color.NRGBA{R: 200, G: 200, B: 200, A: 255},
			Caret:         color.NRGBA{254, 255, 255, 255},
			DisabledCaret: color.NRGBA{R: 200, G: 200, B: 200, A: 255},
		}),

		widget.TextInputOpts.Padding(widget.Insets{
			Left:   30,
			Right:  30,
			Top:    5,
			Bottom: 5,
		}),

		widget.TextInputOpts.CaretOpts(
			widget.CaretOpts.Size(assets.EbitenUIFont(10), 1),
		),

		widget.TextInputOpts.Validation(func(newInputText string) (bool, *string) {
			if utf8.RuneCountInString(newInputText) > 20 {
				return false, nil
			}
			return true, nil
		}),

		widget.TextInputOpts.Placeholder("Name"),
	)

	rootContainer.AddChild(mainMenu.nameInput)

	rootContainer.AddChild(widget.NewButton(
		widget.ButtonOpts.WidgetOpts(
			widget.WidgetOpts.LayoutData(widget.GridLayoutData{
				MaxWidth:           100,
				MaxHeight:          50,
				VerticalPosition:   widget.GridLayoutPositionEnd,
				HorizontalPosition: widget.GridLayoutPositionCenter,
			}),
		),

		widget.ButtonOpts.Image(assets.LoadButtonImage()),

		widget.ButtonOpts.Text("Host", assets.EbitenUIFont(20), &widget.ButtonTextColor{
			Idle: color.NRGBA{0xdf, 0xf4, 0xff, 0xff},
		}),

		widget.ButtonOpts.TextPadding(widget.Insets{
			Left:   30,
			Right:  30,
			Top:    5,
			Bottom: 5,
		}),

		widget.ButtonOpts.ClickedHandler(func(args *widget.ButtonClickedEventArgs) {
			x, y := hostwindow.Contents.PreferredSize()

			r := image.Rect(0, 0, x, y)

			r = r.Add(image.Point{100, 50})

			hostwindow.SetLocation(r)

			mainMenu.ui.AddWindow(hostwindow)
		}),
	))

	rootContainer.AddChild(widget.NewButton(
		widget.ButtonOpts.WidgetOpts(
			widget.WidgetOpts.LayoutData(widget.GridLayoutData{
				MaxWidth:           100,
				MaxHeight:          50,
				VerticalPosition:   widget.GridLayoutPositionEnd,
				HorizontalPosition: widget.GridLayoutPositionCenter,
			}),
		),

		widget.ButtonOpts.Image(assets.LoadButtonImage()),

		widget.ButtonOpts.Text("Join", assets.EbitenUIFont(20), &widget.ButtonTextColor{
			Idle: color.NRGBA{0xdf, 0xf4, 0xff, 0xff},
		}),

		widget.ButtonOpts.TextPadding(widget.Insets{
			Left:   30,
			Right:  30,
			Top:    5,
			Bottom: 5,
		}),

		widget.ButtonOpts.ClickedHandler(func(args *widget.ButtonClickedEventArgs) {
			x, y := joinwindow.Contents.PreferredSize()

			r := image.Rect(0, 0, x, y)

			r = r.Add(image.Point{100, 50})

			joinwindow.SetLocation(r)

			mainMenu.ui.AddWindow(joinwindow)
		}),
	))

	rootContainer.AddChild(widget.NewButton(
		widget.ButtonOpts.WidgetOpts(
			widget.WidgetOpts.LayoutData(widget.GridLayoutData{
				MaxWidth:           100,
				MaxHeight:          50,
				VerticalPosition:   widget.GridLayoutPositionStart,
				HorizontalPosition: widget.GridLayoutPositionCenter,
			}),
		),

		widget.ButtonOpts.Image(assets.LoadButtonImage()),

		widget.ButtonOpts.Text("Exit", assets.EbitenUIFont(20), &widget.ButtonTextColor{
			Idle: color.NRGBA{0xdf, 0xf4, 0xff, 0xff},
		}),

		widget.ButtonOpts.TextPadding(widget.Insets{
			Left:   30,
			Right:  30,
			Top:    5,
			Bottom: 5,
		}),

		widget.ButtonOpts.ClickedHandler(func(args *widget.ButtonClickedEventArgs) {
			os.Exit(0)
		}),
	))

	return mainMenu
}

// anyGamepadVirtualButtonJustPressed checks if any gamepad virtual button is just pressed.
//
// Parameters:
//   - i: The current input state.
//
// Returns:
//   - bool: True if any gamepad virtual button is just pressed, otherwise false.
func anyGamepadVirtualButtonJustPressed(i *inputs.Input) bool {
	if !i.GamepadConfig.IsGamepadIDInitialized() {
		return false
	}

	for _, b := range inputs.VirtualGamepadButtons {
		if i.GamepadConfig.IsButtonJustPressed(b) {
			return true
		}
	}
	return false
}

// Update updates the main menu scene based on the current game state, handling input for hosting or joining games.
//
// Parameters:
//   - state: The current game state.
//
// Returns:
//   - error: An error if the update fails.
func (s *mainMenuScene) Update(state *GameState) error {
	s.ui.Update()
	s.count++

	if s.nameInput.GetText() != "" {
		state.UserInfo.Username = s.nameInput.GetText()
	}

	if s.map1ButtonPressed && state.UserInfo.Username != "" {
		log.Println("Map1 button pressed")
		state.SceneManager.GoTo(NewHostLobby(400, 300, "assets/levels/level1.txt"))
	}

	if s.map2ButtonPressed && state.UserInfo.Username != "" {
		log.Println("Map2 button pressed")
		state.SceneManager.GoTo(NewHostLobby(400, 300, "assets/levels/level2.txt"))
	}

	if s.map3ButtonPressed && state.UserInfo.Username != "" {
		log.Println("Map3 button pressed")
		state.SceneManager.GoTo(NewHostLobby(400, 300, "assets/levels/level3.txt"))
	}

	if s.joinGame && state.UserInfo.Username != "" {
		log.Println("Join Game at address:", s.addressToJoin)
		client := multiplayer.NewGameClient(s.addressToJoin, state.UserInfo, func(code int, text string) error {
			state.SceneManager.GoTo(NewMainMenuScene(400, 300))

			return nil
		})
		state.SceneManager.GoTo(NewClientLobby(400, 300, client, state))
	}

	return nil
}

// Draw renders the main menu scene onto the provided screen image.
//
// Parameters:
//   - r: The image to which the scene is drawn.
func (s *mainMenuScene) Draw(r *ebiten.Image) {
	s.drawTitleBackground(r, s.count)
	drawLogo(r, s.screenWidth, 35, "NinjaGo Bomberman")

	s.ui.Draw(r)
}

// drawTitleBackground draws the scrolling background for the main menu scene.
//
// Parameters:
//   - r: The image to which the background is drawn.
//   - c: The counter used to animate the background.
func (s *mainMenuScene) drawTitleBackground(r *ebiten.Image, c int) {
	w, h := imageBackground.Bounds().Dx(), imageBackground.Bounds().Dy()
	op := &ebiten.DrawImageOptions{}
	dx := (c / 4) % w
	for x := -dx; x < s.screenWidth; x += w {
		for y := 0; y < s.screenHeight; y += h {
			op.GeoM.Reset()
			op.GeoM.Translate(float64(x), float64(y))
			r.DrawImage(imageBackground, op)
		}
	}
}

// drawLogo draws the game logo text on the screen.
//
// Parameters:
//   - r: The image to which the logo is drawn.
//   - screenWidth: The width of the screen.
//   - y: The y-coordinate of the logo.
//   - str: The logo text.
func drawLogo(r *ebiten.Image, screenWidth int, y int, str string) {
	const scale = 4
	x := screenWidth / 2
	assets.DrawTextWithShadow(r, str, x, y, scale, color.RGBA{0x00, 0x00, 0x80, 0xff}, text.AlignCenter, text.AlignStart)
}
