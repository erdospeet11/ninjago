// Package scenes provides the implementation of scene management for transitioning between different game scenes.
package scenes

import (
	ebiten "github.com/hajimehoshi/ebiten/v2"
	"szofttech.inf.elte.hu/szofttech-c-2024/group-08/ninjago/inputs"
	"szofttech.inf.elte.hu/szofttech-c-2024/group-08/ninjago/userinfo"
)

const transitionMaxCount = 20

// SceneManager is the interface that defines the behavior for managing scenes in the game.
type SceneManager interface {
	Update(input *inputs.Input, userInfo *userinfo.UserInfo) error // Update updates the current scene based on the input and user information.
	Draw(r *ebiten.Image)                                          // Draw renders the current scene on the screen.
	GoTo(scene Scene)                                              // GoTo transitions to the specified scene.
}

// sceneManager is the concrete implementation of the SceneManager interface.
type sceneManager struct {
	current         Scene         // The current Scene
	next            Scene         // The next Scene
	transitionCount int           // The transition count
	transitionFrom  *ebiten.Image // The transition image from the current scene
	transitionTo    *ebiten.Image // The transition image to the next scene
}

// NewSceneManager creates a new scene manager with the specified screen width and height.
func NewSceneManager(ScreenWidth, ScreenHeight int) SceneManager {
	return &sceneManager{
		transitionFrom: ebiten.NewImage(ScreenWidth, ScreenHeight),
		transitionTo:   ebiten.NewImage(ScreenWidth, ScreenHeight),
	}
}

// Update updates the current scene or handles the transition between scenes if a transition is in progress.
// It takes the current input and user information as parameters.
func (s *sceneManager) Update(input *inputs.Input, userInfo *userinfo.UserInfo) error {
	if s.transitionCount == 0 {
		return s.current.Update(&GameState{
			SceneManager: s,
			UserInfo:     userInfo,
			Input:        input,
		})
	}

	s.transitionCount--
	if s.transitionCount > 0 {
		return nil
	}

	s.current = s.next
	s.next = nil
	return nil
}

// Draw renders the current scene or handles the transition effect between scenes if a transition is in progress.
func (s *sceneManager) Draw(r *ebiten.Image) {
	if s.transitionCount == 0 {
		s.current.Draw(r)
		return
	}

	s.transitionFrom.Clear()
	s.current.Draw(s.transitionFrom)

	s.transitionTo.Clear()
	s.next.Draw(s.transitionTo)

	r.DrawImage(s.transitionFrom, nil)

	alpha := 1 - float32(s.transitionCount)/float32(transitionMaxCount)
	op := &ebiten.DrawImageOptions{}
	op.ColorScale.ScaleAlpha(alpha)
	r.DrawImage(s.transitionTo, op)
}

// GoTo transitions to the specified scene, initiating a transition effect if a current scene is active.
func (s *sceneManager) GoTo(scene Scene) {
	if s.current == nil {
		s.current = scene
	} else {
		s.next = scene
		s.transitionCount = transitionMaxCount
	}
}
