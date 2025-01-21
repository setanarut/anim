package main

import (
	"bytes"
	"fmt"
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

STATS

Current atlas: %v
Current state: %v
Current state FPS: %v
Current index %v
`

type Game struct {
}

func (g *Game) Update() error {
	animPlayer.Update()
	DIO.GeoM.Reset()

	if ebiten.IsKeyPressed(ebiten.KeySpace) {
		animPlayer.CurrentState = "jump"
	} else if ebiten.IsKeyPressed(ebiten.KeyD) {
		animPlayer.CurrentState = "run"
	} else if ebiten.IsKeyPressed(ebiten.KeyA) {
		animPlayer.CurrentState = "run"
		// FlipX
		DIO.GeoM.Scale(-1, 1)
		// Align to zero
		DIO.GeoM.Translate(float64(animPlayer.CurrentFrame.Bounds().Dx()), 0)
	} else {
		animPlayer.CurrentState = "idle"
	}
	DIO.GeoM.Scale(8, 8)

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.DrawImage(animPlayer.CurrentFrame, DIO)
	ebitenutil.DebugPrintAt(
		screen,
		fmt.Sprintf(
			hud,
			animPlayer.CurrentAtlas,
			animPlayer.CurrentState,
			animPlayer.CurrentStateFPS(),
			animPlayer.CurrentIndex,
		),
		220,
		25)
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

	animPlayer.NewState("idle", 0, 0, 32, 32, 5, false, false, 5)
	animPlayer.NewState("run", 0, 32, 32, 32, 8, false, false, 12)
	animPlayer.NewState("jump", 0, 32*2, 32, 32, 4, false, false, 15)
	animPlayer.CurrentState = "idle"

	ebiten.SetWindowSize(400, 300)
	if err := ebiten.RunGame(&Game{}); err != nil {
		log.Fatal(err)
	}
}
