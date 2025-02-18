package main

import (
	"bytes"
	"image"

	_ "image/jpeg"
	_ "image/png"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/examples/resources/images"
	"github.com/setanarut/anim"
)

var animPlayer *anim.AnimationPlayer
var DIO *ebiten.DrawImageOptions = &ebiten.DrawImageOptions{}
var hud = `
CONTROLS

| Key   | State        |
| ----- | ------------ |
| D     | Run right    |
| A     | Run (flipX)  |
| Space | Jump         |
`

type Game struct {
}

func (g *Game) Update() error {
	animPlayer.Update()
	DIO.GeoM.Reset()

	if ebiten.IsKeyPressed(ebiten.KeySpace) {
		animPlayer.SetAnim("jump")
	} else if ebiten.IsKeyPressed(ebiten.KeyD) {
		animPlayer.SetAnim("run")
	} else if ebiten.IsKeyPressed(ebiten.KeyA) {
		animPlayer.SetAnim("run")
		// FlipX
		DIO.GeoM.Scale(-1, 1)
		// Align to zero
		DIO.GeoM.Translate(float64(animPlayer.CurrentFrame.Bounds().Dx()), 0)
	} else {
		animPlayer.SetAnim("idle")
	}
	DIO.GeoM.Scale(8, 8)

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.DrawImage(animPlayer.CurrentFrame, DIO)
	ebitenutil.DebugPrintAt(screen, hud, 220, 4)
	ebitenutil.DebugPrintAt(screen, animPlayer.String(), 220, 130)
}

func (g *Game) Layout(w, h int) (int, int) {
	return 400, 300
}

func main() {
	img, _, err := image.Decode(bytes.NewReader(images.Runner_png))
	if err != nil {
		log.Fatal(err)
	}

	spriteSheet := &anim.Atlas{
		Name:  "Default",
		Image: ebiten.NewImageFromImage(img),
	}

	animPlayer = anim.NewAnimationPlayer(spriteSheet)

	animPlayer.NewAnim("idle", 0, 0, 32, 32, 5, false, false, 5)
	animPlayer.NewAnim("run", 0, 32, 32, 32, 8, false, false, 12)
	animPlayer.NewAnim("jump", 0, 32*2, 32, 32, 4, false, false, 15)
	animPlayer.SetAnim("idle")

	ebiten.SetWindowSize(400, 300)
	if err := ebiten.RunGame(&Game{}); err != nil {
		log.Fatal(err)
	}
}
