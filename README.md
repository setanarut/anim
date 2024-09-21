[![GoDoc](https://godoc.org/github.com/setanarut/anim?status.svg)](https://pkg.go.dev/github.com/setanarut/anim)

# Anim

![demo](https://github.com/user-attachments/assets/fa7cd502-9bf9-4a88-a83c-df0e23c80e41)

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
spriteSheet := ebiten.NewImageFromImage(img)
animPlayer = anim.NewAnimationPlayer(spriteSheet)
```

Let's specify the coordinates of the animations for the player states.
The figure shows the coordinates for "run" state. `NewAnimationState("run", 0, 32, 32, 32, 8, false, false)`

![diag](https://github.com/user-attachments/assets/316be3e7-102f-4d3f-b126-637cda387253)


```Go
animPlayer.NewAnimationState("idle", 0, 0, 32, 32, 5, false, false).FPS = 5
animPlayer.NewAnimationState("run", 0, 32, 32, 32, 8, false, false)
animPlayer.NewAnimationState("jump", 0, 32*2, 32, 32, 4, false, false)
```

Let's set the initial state.

```Go
animPlayer.SetState("idle")
```

Update animation player

```Go
func (g *Game) Update() error {
	animPlayer.Update()
```

Let's update the states according to the character's movement.

```Go
if ebiten.IsKeyPressed(ebiten.KeySpace) {
    animPlayer.SetState("jump")
} else if ebiten.IsKeyPressed(ebiten.KeyD) {
    animPlayer.SetState("run")
} else if ebiten.IsKeyPressed(ebiten.KeyA) {
    animPlayer.SetState("run")
    // Flip X axis with *ebiten.DrawImageOptions.GeoM
    DIO.GeoM.Scale(-1, 1)
} else {
    animPlayer.SetState("idle")
}
```

Finally let's draw Animation player

```Go
func (g *Game) Draw(screen *ebiten.Image) {
	screen.DrawImage(animPlayer.CurrentFrame, DIO)
```

## Example

Run [demo](./examples/demo/) on your local machine

```zsh
go run github.com/setanarut/anim/examples/demo@latest
```
