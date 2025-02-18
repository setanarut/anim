package anim

import (
	"fmt"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
)

const str string = `
Playback state;
Current Atlas: %v
Current Anim: %v
Paused: %v
Tick: %v
Current Anim Index %v
Current Anim FPS: %v
`

// Animation for AnimationPlayer
type Animation struct {
	Name   string          // Name of the aimation state
	Frames []*ebiten.Image // Animation frames
	FPS    float64         // Animation playback speed (Frames Per Second).
}

// PlaybackData is AnimationPlayer's playback data.
// With this structure, the playback state can be saved to disk and reloaded.
type PlaybackData struct {
	// Current sprite sheet atlas
	CurrentAtlas string
	// CurrentAnim name (animation name)
	CurrentAnim string
	// If true, animations will be paused.
	Paused bool
	// Animation tick
	Tick float64
	// Current animation frame index
	CurrentIndex int
}

// AnimationPlayer plays and manages animations.
type AnimationPlayer struct {
	Data *PlaybackData
	// Current frame of the current animation.
	//
	// The frame is dynamically updated with `AnimationPlayer.Update()`.
	CurrentFrame *ebiten.Image
	Atlases      []*Atlas
	// Animations and alternative sprite sheet atlases
	Animations map[string]map[string]*Animation
}

func (ap *AnimationPlayer) String() string {
	return fmt.Sprintf(
		str,
		ap.Data.CurrentAtlas,
		ap.Data.CurrentAnim,
		ap.Data.Paused,
		ap.Data.Tick,
		ap.Data.CurrentIndex,
		ap.CurrentAnimFPS(),
	)
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
		Data:       &PlaybackData{},
		Atlases:    atlases,
		Animations: make(map[string]map[string]*Animation),
	}
	for _, atlas := range atlases {
		ap.Animations[atlas.Name] = make(map[string]*Animation)
	}
	ap.Data.CurrentAtlas = atlases[0].Name
	return ap
}

// NewAnim appends a new Animation to the AnimationPlayer
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
func (ap *AnimationPlayer) NewAnim(
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
		// animation := NewAnimation(stateName, subImages, FPS)
		animation := &Animation{
			Name:   stateName,
			Frames: subImages,
			FPS:    FPS,
		}
		ap.Animations[atlas.Name][stateName] = animation
	}
	ap.Data.CurrentAnim = stateName
	ap.CurrentFrame = ap.Animations[ap.Data.CurrentAtlas][stateName].Frames[ap.Data.CurrentIndex]
}

// SetAllFPS overwrites the FPS of all animations.
func (ap *AnimationPlayer) SetAllFPS(FPS float64) {
	for _, atlas := range ap.Atlases {
		for _, anim := range ap.Animations[atlas.Name] {
			anim.FPS = FPS
		}

	}
}

// CurrentAnimFPS returns FPS of the current animation
func (ap *AnimationPlayer) CurrentAnimFPS() float64 {
	return ap.Animations[ap.Data.CurrentAtlas][ap.Data.CurrentAnim].FPS
}

// SetAnimFPS sets FPS of the animation state.
func (ap *AnimationPlayer) SetAnimFPS(animName string, FPS float64) {
	ap.Animations[ap.Data.CurrentAtlas][animName].FPS = FPS
}

// SetAnim sets the animation and resets to the first frame.
//
// If you assign ap.Data.CurrentAnim = "animName" directly, the animation will not be reset.
func (ap *AnimationPlayer) SetAnim(animName string) {
	if ap.Data.CurrentAnim != animName {
		ap.Data.CurrentAnim = animName
		ap.Data.Tick = 0
		ap.Data.CurrentIndex = 0
	}
}

func (ap *AnimationPlayer) SetAtlas(atlasName string) {
	ap.Data.CurrentAtlas = atlasName
}

// Atlas returns current Atlas
func (ap *AnimationPlayer) Atlas() string {
	return ap.Data.CurrentAtlas
}

// Anim returns current animation state name
func (ap *AnimationPlayer) Anim() string {
	return ap.Data.CurrentAnim
}

// PauseAtFrame pauses the current animation at the frame. If index is out of range it does nothing.
func (ap *AnimationPlayer) PauseAtFrame(index int) {
	if index < len(ap.Animations[ap.Data.CurrentAtlas][ap.Data.CurrentAnim].Frames) && index >= 0 {
		ap.Data.Paused = true
		ap.Data.CurrentIndex = index
	}
}

// Update updates AnimationPlayer. Place this func inside Ebitengine `Game.Update()`.
//
//	// example
//	func (g *Game) Update() error {
//	animPlayer.Update()
//	...
func (ap *AnimationPlayer) Update() {
	if !ap.Data.Paused {
		ap.Data.Tick += ap.Animations[ap.Data.CurrentAtlas][ap.Data.CurrentAnim].FPS / 60.0
		ap.Data.CurrentIndex = int(math.Floor(ap.Data.Tick))
		if ap.Data.CurrentIndex >= len(ap.Animations[ap.Data.CurrentAtlas][ap.Data.CurrentAnim].Frames) {
			ap.Data.Tick = 0
			ap.Data.CurrentIndex = 0
		}
	}
	// update current frame
	ap.CurrentFrame = ap.Animations[ap.Data.CurrentAtlas][ap.Data.CurrentAnim].Frames[ap.Data.CurrentIndex]
}
