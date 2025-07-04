package entities

import (
	"testing"

	"github.com/hajimehoshi/ebiten/v2"
)

func TestAnimation_Update(t *testing.T) {
	type fields struct {
		sprite         *ebiten.Image
		frameOX        int
		frameOY        int
		frameWidth     int
		frameHeight    int
		frameCount     int
		currentFrame   int
		animationSpeed int
	}
	tests := []struct {
		name   string
		fields fields
		want   int
	}{
		{
			name: "Increment frame",
			fields: fields{
				sprite:         ebiten.NewImage(100, 100),
				frameOX:        0,
				frameOY:        0,
				frameWidth:     100,
				frameHeight:    100,
				frameCount:     5,
				currentFrame:   0,
				animationSpeed: 1,
			},
			want: 1,
		},
		{
			name: "Wraparound frame",
			fields: fields{
				sprite:         ebiten.NewImage(100, 100),
				frameOX:        0,
				frameOY:        0,
				frameWidth:     100,
				frameHeight:    100,
				frameCount:     5,
				currentFrame:   4,
				animationSpeed: 1,
			},
			want: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &Animation{
				sprite:         tt.fields.sprite,
				frameOX:        tt.fields.frameOX,
				frameOY:        tt.fields.frameOY,
				frameWidth:     tt.fields.frameWidth,
				frameHeight:    tt.fields.frameHeight,
				frameCount:     tt.fields.frameCount,
				currentFrame:   tt.fields.currentFrame,
				animationSpeed: tt.fields.animationSpeed,
			}
			a.Update()
			got := (a.currentFrame / a.animationSpeed) % a.frameCount
			if got != tt.want {
				t.Errorf("Animation.Update() got = %v, want %v", got, tt.want)
			}
		})
	}
}
