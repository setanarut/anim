package anim

import (
	"image"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
)

// AnimationPlayer plays and manages animations.
type AnimationPlayer struct {
	// Current frame of the current animation.
	//
	// The frame is dynamically updated with `AnimationPlayer.Update()`.
	CurrentFrame *ebiten.Image
	// Current sprite sheet atlas
	CurrentAtlas string
	// Current state name (animation name)
	CurrentState string
	// Slice for Sprite sheet variations
	Atlases []*Atlas
	// Animations and alternative sprite sheet atlases
	Animations map[string]map[string]*Animation
	// If true, animations will be paused.
	Paused bool
	// Animation tick
	Tick float64
	// Current animation frame index
	CurrentIndex int
}

// Atlas is a sprite sheet state for animation player.
//
// It is used to easily switch between different sprite sheet variations
// that share the same coordinates.
type Atlas struct {
	Name  string
	Image *ebiten.Image
}

// NewAnimationPlayer returns new AnimationPlayer with spriteSheet
func NewAnimationPlayer(atlases ...*Atlas) *AnimationPlayer {
	ap := &AnimationPlayer{
		Atlases:    atlases,
		Animations: make(map[string]map[string]*Animation),
	}
	for _, atlas := range atlases {
		ap.Animations[atlas.Name] = make(map[string]*Animation)
	}
	ap.CurrentAtlas = atlases[0].Name
	return ap
}

// NewState appends a new Animation to the AnimationPlayer
// and returns the Animation.
//
// Parameters:
//
//   - x, y - Top-left coordinates of the first frame's rectangle.
//   - w, h - Width and height of the first frame's rectangle.
//   - frameCount - Animation frame count
//   - vertical - If true, frames are appended vertically, otherwise horizontally.
//   - pingPong - If true, arranges the animation indexes to play back and forth. [0 1 2 3 2 1]
//   - FPS - Playback FPS
func (ap *AnimationPlayer) NewState(
	stateName string,
	x, y,
	w, h,
	frameCount int,
	pingPong bool,
	vertical bool,
	FPS float64,
) {
	for _, atlas := range ap.Atlases {
		subImages := SubImages(atlas.Image, x, y, w, h, frameCount, vertical)
		if pingPong {
			subImages = MakePingPong(subImages)
		}
		animation := NewAnimation(stateName, subImages, FPS)
		ap.Animations[atlas.Name][stateName] = animation
	}
	ap.CurrentState = stateName
	ap.CurrentFrame = ap.Animations[ap.CurrentAtlas][stateName].Frames[ap.CurrentIndex]
}

// SetAllFPS overwrites the FPS of all animations.
func (ap *AnimationPlayer) SetAllFPS(FPS float64) {
	for _, atlas := range ap.Atlases {
		for _, anim := range ap.Animations[atlas.Name] {
			anim.FPS = FPS
		}

	}
}

// AddAnimation adds the given animation to this player.
// Adds the name of the animation as a map key.
func (ap *AnimationPlayer) AddAnimation(a *Animation) {
	ap.Animations[ap.CurrentAtlas][a.Name] = a
}

// CurrentStateFPS returns FPS of the current animation state
func (ap *AnimationPlayer) CurrentStateFPS() float64 {
	return ap.Animations[ap.CurrentAtlas][ap.CurrentState].FPS
}

// SetStateFPS sets FPS of the animation state.
//
//	// Shortcut func for
//	AnimationPlayer.Animations[stateName].FPS = 15
//	Animation.FPS = 15
func (ap *AnimationPlayer) SetStateFPS(stateName string, FPS float64) {
	ap.Animations[ap.CurrentAtlas][stateName].FPS = FPS
}

// SetStateAndReset sets the animation state and resets to the first frame.
func (ap *AnimationPlayer) SetStateAndReset(state string) {
	ap.CurrentState = state
	ap.Tick = 0
	ap.CurrentIndex = 0
}

// PauseAtFrame pauses the current animation at the frame. If index is out of range it does nothing.
func (ap *AnimationPlayer) PauseAtFrame(index int) {
	if index < len(ap.Animations[ap.CurrentAtlas][ap.CurrentState].Frames) && index >= 0 {
		ap.Paused = true
		ap.CurrentIndex = index
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
		ap.Tick += ap.Animations[ap.CurrentAtlas][ap.CurrentState].FPS / 60.0
		ap.CurrentIndex = int(math.Floor(ap.Tick))
		if ap.CurrentIndex >= len(ap.Animations[ap.CurrentAtlas][ap.CurrentState].Frames) {
			ap.Tick = 0
			ap.CurrentIndex = 0
		}
	}
	// update current frame
	ap.CurrentFrame = ap.Animations[ap.CurrentAtlas][ap.CurrentState].Frames[ap.CurrentIndex]
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
