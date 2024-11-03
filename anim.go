package anim

import (
	"image"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
)

// AnimationPlayer plays and manages animations.
type AnimationPlayer struct {
	// Image atlas for all animation states
	SpriteSheet *ebiten.Image
	// Current frame of the current animation.
	//
	// The frame is dynamically updated with `AnimationPlayer.Update()`.
	CurrentFrame *ebiten.Image
	// Animations for states
	Animations map[string]*Animation
	// If true, animations will be paused.
	Paused bool
	// Current state name (animation name)
	CurrentState string

	tick              float64
	currentFrameIndex int
}

// NewAnimationPlayer returns new AnimationPlayer with spriteSheet
func NewAnimationPlayer(spriteSheet *ebiten.Image) *AnimationPlayer {
	return &AnimationPlayer{
		SpriteSheet:       spriteSheet,
		Paused:            false,
		Animations:        make(map[string]*Animation),
		currentFrameIndex: 0,
	}

}

// NewAnimationState appends a new Animation to the AnimationPlayer
// and returns the Animation.
//
// The default FPS is 15.
//
// Parameters:
//
//   - x, y - Top-left coordinates of the first frame's rectangle.
//   - w, h - Width and height of the first frame's rectangle.
//   - frameCount - Animation frame count
//   - vertical - If true, frames are appended vertically, otherwise horizontally.
//   - pingPong - If true, arranges the animation indexes to play back and forth. [0 1 2 3 2 1]
func (ap *AnimationPlayer) NewAnimationState(
	stateName string,
	x, y,
	w, h,
	frameCount int,
	pingPong bool,
	vertical bool) *Animation {

	subImages := SubImages(ap.SpriteSheet, x, y, w, h, frameCount, vertical)

	if pingPong {
		subImages = MakePingPong(subImages)
	}

	animation := NewAnimation(stateName, subImages, 15)
	ap.CurrentState = stateName
	ap.Animations[stateName] = animation
	ap.CurrentFrame = ap.Animations[ap.CurrentState].Frames[ap.currentFrameIndex]
	return animation
}

// CurrentFrameIndex returns current index.
// May be useful for debugging.
func (ap *AnimationPlayer) CurrentFrameIndex() int {
	return ap.currentFrameIndex
}

// SetAllFPS overwrites the FPS of all animations.
func (ap *AnimationPlayer) SetAllFPS(FPS float64) {
	for _, anim := range ap.Animations {
		anim.FPS = FPS
	}
}

// AddAnimation adds the given animation to this player.
// Adds the name of the animation as a map key.
func (ap *AnimationPlayer) AddAnimation(a *Animation) {
	ap.Animations[a.Name] = a
}

// State returns current active animation state
func (ap *AnimationPlayer) State() string {
	return ap.CurrentState
}

// CurrentStateFPS returns FPS of the current animation state
func (ap *AnimationPlayer) CurrentStateFPS() float64 {
	return ap.Animations[ap.State()].FPS
}

// SetStateFPS sets FPS of the animation state.
//
//	// Shortcut func for
//	AnimationPlayer.Animations[stateName].FPS = 15
//	Animation.FPS = 15
func (ap *AnimationPlayer) SetStateFPS(stateName string, FPS float64) {
	ap.Animations[stateName].FPS = FPS
}

// SetStateAndReset sets the animation state and resets to the first frame.
func (ap *AnimationPlayer) SetStateAndReset(state string) {
	if ap.CurrentState != state {
		ap.CurrentState = state
		ap.tick = 0
		ap.currentFrameIndex = 0
	}
}

// SetState sets the animation state.
func (ap *AnimationPlayer) SetState(state string) {
	if ap.CurrentState != state {
		ap.CurrentState = state
	}
}

// PauseAtFrame pauses the current animation at the frame. If index is out of range it does nothing.
func (ap *AnimationPlayer) PauseAtFrame(frameIndex int) {
	if frameIndex < len(ap.Animations[ap.State()].Frames) && frameIndex >= 0 {
		ap.Paused = true
		ap.currentFrameIndex = frameIndex
	}
}

// Update updates AnimationPlayer. Place this func inside Ebitengine `Game.Update()`.
//
//	// example
//	func (g *Game) Update() error {
//	animPlayer.Update()
//	...
func (ap *AnimationPlayer) Update() {
	if !ap.Paused {
		ap.tick += ap.Animations[ap.CurrentState].FPS / 60.0
		ap.currentFrameIndex = int(math.Floor(ap.tick))
		if ap.currentFrameIndex >= len(ap.Animations[ap.CurrentState].Frames) {
			ap.tick = 0
			ap.currentFrameIndex = 0
		}
	}
	// update current frame
	ap.CurrentFrame = ap.Animations[ap.CurrentState].Frames[ap.currentFrameIndex]
}

// Animation for AnimationPlayer
type Animation struct {
	Name   string          // Name of the aimation state
	Frames []*ebiten.Image // Animation frames
	FPS    float64         // Animation playback speed (Frames Per Second).
}

// NewAnimation returns new Animation
func NewAnimation(name string, frames []*ebiten.Image, FPS float64) *Animation {
	return &Animation{
		Name:   name,
		Frames: frames,
		FPS:    FPS,
	}
}

// SubImages returns sub-images from spriteSheet image
func SubImages(spriteSheet *ebiten.Image, x, y, w, h, subImageCount int, vertical bool) []*ebiten.Image {

	subImages := []*ebiten.Image{}
	frameRect := image.Rect(x, y, x+w, y+h)

	for i := 0; i < subImageCount; i++ {
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
