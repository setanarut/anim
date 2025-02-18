package main

import (
	"bytes"
	"image"
	"math"

	"image/color"
	"image/draw"
	_ "image/jpeg"
	_ "image/png"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/examples/resources/images"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/setanarut/anim"
)

var animPlayer *anim.AnimationPlayer
var DIO *ebiten.DrawImageOptions = &ebiten.DrawImageOptions{}
var hud = `
CONTROLS
| Key   | State        |
| ----- | ------------ |
| F     | Change atlas |
| D     | Run right    |
| A     | Run (flipX)  |
| Space | Jump         |
`

type Game struct {
}

func (g *Game) Update() error {
	animPlayer.Update()
	DIO.GeoM.Reset()

	if inpututil.IsKeyJustPressed(ebiten.KeyF) {
		switch animPlayer.Atlas() {
		case "ShiftHue":
			animPlayer.SetAtlas("Default")
		case "Default":
			animPlayer.SetAtlas("ShiftHue")
		}
	}

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
	hueImg := image.NewRGBA(img.Bounds())
	// clone runner atlas for colorization
	draw.Draw(hueImg, img.Bounds(), img, img.Bounds().Min, draw.Src)
	hueImg = ShiftHue(hueImg, 120).(*image.RGBA)
	defaultAtlas := &anim.Atlas{
		Name:  "Default",
		Image: ebiten.NewImageFromImage(img),
	}
	HueAtlas := &anim.Atlas{
		Name:  "ShiftHue",
		Image: ebiten.NewImageFromImage(hueImg),
	}
	animPlayer = anim.NewAnimationPlayer(defaultAtlas, HueAtlas)

	animPlayer.NewAnim("idle", 0, 0, 32, 32, 5, false, false, 5)
	animPlayer.NewAnim("run", 0, 32, 32, 32, 8, false, false, 12)
	animPlayer.NewAnim("jump", 0, 32*2, 32, 32, 4, false, false, 15)
	animPlayer.SetAnim("idle")

	ebiten.SetWindowSize(400, 300)
	if err := ebiten.RunGame(&Game{}); err != nil {
		log.Fatal(err)
	}
}

func ShiftHue(img image.Image, shiftAngle float64) image.Image {
	bounds := img.Bounds()
	shifted := image.NewRGBA(bounds)

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r, g, b, a := img.At(x, y).RGBA()

			rf := float64(r >> 8)
			gf := float64(g >> 8)
			bf := float64(b >> 8)

			max := math.Max(math.Max(rf, gf), bf)
			min := math.Min(math.Min(rf, gf), bf)
			delta := max - min

			var hue float64
			if delta == 0 {
				hue = 0
			} else if max == rf {
				hue = 60 * math.Mod((gf-bf)/delta, 6)
			} else if max == gf {
				hue = 60 * ((bf-rf)/delta + 2)
			} else {
				hue = 60 * ((rf-gf)/delta + 4)
			}

			hue = math.Mod(hue+shiftAngle, 360)
			if hue < 0 {
				hue += 360
			}

			c := max - min
			x2 := c * (1 - math.Abs(math.Mod(hue/60, 2)-1))

			var r2, g2, b2 float64
			switch {
			case hue >= 0 && hue < 60:
				r2, g2, b2 = c, x2, 0
			case hue >= 60 && hue < 120:
				r2, g2, b2 = x2, c, 0
			case hue >= 120 && hue < 180:
				r2, g2, b2 = 0, c, x2
			case hue >= 180 && hue < 240:
				r2, g2, b2 = 0, x2, c
			case hue >= 240 && hue < 300:
				r2, g2, b2 = x2, 0, c
			case hue >= 300 && hue < 360:
				r2, g2, b2 = c, 0, x2
			}

			m := max - c
			r2 = (r2 + m) * 255
			g2 = (g2 + m) * 255
			b2 = (b2 + m) * 255

			shifted.Set(x, y, color.RGBA{
				R: uint8(r2),
				G: uint8(g2),
				B: uint8(b2),
				A: uint8(a >> 8),
			})
		}
	}

	return shifted
}
