package anim

import (
	"image"

	"github.com/hajimehoshi/ebiten/v2"
)

// SubImages returns sub-images from spriteSheet image
func SubImages(spriteSheet *ebiten.Image, x, y, w, h, subImageCount int, vertical bool) []*ebiten.Image {

	subImages := []*ebiten.Image{}
	frameRect := image.Rect(x, y, x+w, y+h)

	for range subImageCount {
		subImages = append(subImages, spriteSheet.SubImage(frameRect).(*ebiten.Image))
		if vertical {
			frameRect.Min.Y += h
			frameRect.Max.Y += h
		} else {
			frameRect.Min.X += w
			frameRect.Max.X += w
		}
	}
	return subImages

}

// MakePingPong arranges the animation indexes to play back and forth.
// [0 1 2 3] -> [0 1 2 3 2 1]
func MakePingPong(frames []*ebiten.Image) []*ebiten.Image {
	for i := len(frames) - 2; i > 0; i-- {
		frames = append(frames, frames[i])
	}
	return frames
}
