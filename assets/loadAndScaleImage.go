package assets

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

func LoadAndScaleImage(path string, targetWidth, targetHeight int) *ebiten.Image {
	img, _, err := ebitenutil.NewImageFromFile(path)
	if err != nil {
		log.Fatalf("failed to load image: %v", err)
	}

	originalWidth, originalHeight := img.Bounds().Dx(), img.Bounds().Dy()

	scaleX := float64(targetWidth) / float64(originalWidth)
	scaleY := float64(targetHeight) / float64(originalHeight)

	resizedImg := ebiten.NewImage(targetWidth, targetHeight)

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(scaleX, scaleY)
	resizedImg.DrawImage(img, op)

	return resizedImg
}
