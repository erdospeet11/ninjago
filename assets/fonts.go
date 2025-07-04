package assets

import (
	"bytes"
	"image/color"
	"log"

	"github.com/golang/freetype/truetype"
	ebiten "github.com/hajimehoshi/ebiten/v2"
	text "github.com/hajimehoshi/ebiten/v2/text/v2"
	"golang.org/x/image/font"
)

const (
	arcadeFontBaseSize = 8
)

var (
	arcadeFaceSource *text.GoTextFaceSource
)

func init() {
	font, err := EmbeddedAssets.ReadFile("assets/fonts/NotoSans-Bold.ttf")
	if err != nil {
		log.Panicln(err)
	}

	s, err := text.NewGoTextFaceSource(bytes.NewReader(font))
	if err != nil {
		log.Panicln(err)
	}
	arcadeFaceSource = s
}

var (
	shadowColor = color.RGBA{0, 0, 0, 0x80}
)

func DrawTextWithShadow(rt *ebiten.Image, str string, x, y, scale int, clr color.Color, primaryAlign, secondaryAlign text.Align) {
	op := &text.DrawOptions{}
	op.GeoM.Translate(float64(x)+1, float64(y)+1)
	op.ColorScale.ScaleWithColor(shadowColor)
	op.LineSpacing = arcadeFontBaseSize * float64(scale)
	op.PrimaryAlign = primaryAlign
	op.SecondaryAlign = secondaryAlign
	text.Draw(rt, str, &text.GoTextFace{
		Source: arcadeFaceSource,
		Size:   arcadeFontBaseSize * float64(scale),
	}, op)

	op.GeoM.Reset()
	op.GeoM.Translate(float64(x), float64(y))
	op.ColorScale.Reset()
	op.ColorScale.ScaleWithColor(clr)
	text.Draw(rt, str, &text.GoTextFace{
		Source: arcadeFaceSource,
		Size:   arcadeFontBaseSize * float64(scale),
	}, op)
}

func EbitenUIFont(size float64) font.Face {
	fontFile, err := EmbeddedAssets.ReadFile("assets/fonts/NotoSans-Bold.ttf")
	if err != nil {
		log.Panicln(err)
	}

	ttFont, err := truetype.Parse(fontFile)
	if err != nil {
		log.Panicln(err)
	}
	return truetype.NewFace(ttFont, &truetype.Options{
		Size:    size,
		DPI:     72,
		Hinting: font.HintingFull,
	})
}
