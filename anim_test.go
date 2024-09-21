package anim_test

import (
	"bytes"
	"image"
	_ "image/png"
	"testing"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/examples/resources/images"
	"github.com/setanarut/anim"
)

var im *ebiten.Image

func TestAnim(t *testing.T) {
	t.Run("Load Image", LoadImage)
	t.Run("Ping Pong", PingPong)
}

func LoadImage(t *testing.T) {
	img, _, err := image.Decode(bytes.NewReader(images.Runner_png))
	if err != nil {
		t.Skip("image load error")
	}
	im = ebiten.NewImageFromImage(img)
}

func PingPong(t *testing.T) {

	img := ebiten.NewImageFromImage(im)
	spriteSheet := ebiten.NewImageFromImage(img)
	animPlayer := anim.NewAnimationPlayer(spriteSheet)
	animPlayer.NewAnimationState("idle", 0, 0, 32, 32, 5, true, false)

	frames := animPlayer.Animations["idle"].Frames
	secondFrame := frames[1]
	lastFrame := frames[len(frames)-1]

	if secondFrame != lastFrame {
		t.Fatalf("Ping Pong index error")
	}
}
