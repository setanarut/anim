# Anim

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

```Go
spriteSheet := ebiten.NewImageFromImage(img)
animPlayer = anim.NewAnimationPlayer(spriteSheet)
```

Let's specify the coordinates of the animations for the player states.

```Go
animPlayer.NewAnimationState("idle", 0, 0, 32, 32, 5, false, false).FPS = 5
animPlayer.NewAnimationState("run", 0, 32, 32, 32, 8, false, false)
animPlayer.NewAnimationState("jump", 0, 32*2, 32, 32, 4, false, false)
```

Let's set the initial state.

```Go
animPlayer.SetState("idle")
```

## Example

Run demo on your local machine

```zsh
go run github.com/setanarut/anim/examples/demo@latest
```