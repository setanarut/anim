[![GoDoc](https://godoc.org/github.com/setanarut/anim?status.svg)](https://pkg.go.dev/github.com/setanarut/anim)

# Anim

![demo](https://github.com/user-attachments/assets/15fb4ef9-00a6-412d-af5f-67ff9fa1fb4c)

Anim is an easy to use animation player package for Ebitengine v2

```Go
import "github.com/setanarut/anim"
```

## Tutorial

Let's declare a global variable for the animation player

```Go
var animPlayer *anim.AnimationPlayer
```

Make new animation player from sprite atlas

![runner](https://github.com/user-attachments/assets/54871498-ae7b-4107-adf4-e292aaff47e7)

```Go
spriteSheet := anim.Atlas{
	Name:  "Default",
	Image: ebiten.NewImageFromImage(img),
}
animPlayer = anim.NewAnimationPlayer(spriteSheet)
```

Let's specify the coordinates of the animations for the player states.
The figure shows the coordinates for "run" state.

![diag](https://github.com/user-attachments/assets/316be3e7-102f-4d3f-b126-637cda387253)


```Go
animPlayer.NewAnim("idle", 0, 0, 32, 32, 5, false, false, 5)
animPlayer.NewAnim("run", 0, 32, 32, 32, 8, false, false, 12)
animPlayer.NewAnim("jump", 0, 32*2, 32, 32, 4, false, false, 15)
```

Let's set the initial animation state.

```Go
animPlayer.SetAnim("idle)
```

Update animation player

```Go
func (g *Game) Update() error {
	animPlayer.Update()
```

Let's update the states according to the character's movement.

```Go
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
```

Finally let's draw Animation player

```Go
func (g *Game) Draw(screen *ebiten.Image) {
	screen.DrawImage(animPlayer.CurrentFrame, DIO)
```

## Examples

### Simple demo

Run [demo](./examples/demo/) on your local machine

```zsh
go run github.com/setanarut/anim/examples/demo@latest
```
### Multiple sprite sheet example

Example of alternative sprite sheets with exactly the same coordinates.  
Run [atlases](./examples/atlases/) on your local machine;

```zsh
go run github.com/setanarut/anim/examples/atlases@latest
```
