package anim_test

import (
	"testing"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/setanarut/anim"
)

func TestNewAnimationPlayer(t *testing.T) {
	t.Run("symetry", func(t *testing.T) {
		a := anim.NewAnimationPlayer(ebiten.NewImage(100, 100))
		if a.CurrentFrameIndex != 0 {
			t.Error("CurrentFrameIndex is not zero")
		}
	})
}
