package main

import (
	"bytes"
	"fmt"
	"image"

	_ "image/png"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/examples/resources/images"
	"github.com/setanarut/anim"
)

var animPlayer *anim.AnimationPlayer
var DIO *ebiten.DrawImageOptions = &ebiten.DrawImageOptions{}
var keys []ebiten.Key
var hud = `

CONTROLS

| Key   | State       |
| ----- | ---------   |
| D     | Run right   |
| A     | Run (flipX) |
| Space | Jump        |

STATS

Current state: %v
Current state FPS: %v
`

type Game struct {
}

func (g *Game) Update() error {
	animPlayer.Update()
	DIO.GeoM.Reset()

	if ebiten.IsKeyPressed(ebiten.KeySpace) {
		animPlayer.SetState("jump")
	} else if ebiten.IsKeyPressed(ebiten.KeyD) {
		animPlayer.SetState("run")
	} else if ebiten.IsKeyPressed(ebiten.KeyA) {
		animPlayer.SetState("run")
		// FlipX
		DIO.GeoM.Scale(-1, 1)
		// Align to zero
		DIO.GeoM.Translate(float64(animPlayer.CurrentFrame.Bounds().Dx()), 0)
	} else {
		animPlayer.SetState("idle")
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
			animPlayer.CurrentState,
			animPlayer.CurrentStateFPS()),
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
	spriteSheet := ebiten.NewImageFromImage(img)
	animPlayer = anim.NewAnimationPlayer(spriteSheet)
	animPlayer.AddStateAnimation("idle", 0, 0, 32, 32, 5, false, false).FPS = 5
	animPlayer.AddStateAnimation("run", 0, 32, 32, 32, 8, false, false)
	animPlayer.AddStateAnimation("jump", 0, 32*2, 32, 32, 4, false, false)
	animPlayer.SetState("idle")

	ebiten.SetWindowSize(400, 300)
	if err := ebiten.RunGame(&Game{}); err != nil {
		log.Fatal(err)
	}
}