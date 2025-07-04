package scenes

import (
	"testing"

	"github.com/hajimehoshi/ebiten/v2"
)

func TestNewSceneManager(t *testing.T) {
	sm := NewSceneManager(800, 600).(*sceneManager)
	if sm == nil {
		t.Fatal("SceneManager should not be nil")
	}
	if sm.transitionFrom == nil || sm.transitionTo == nil {
		t.Error("Transition surfaces should be initialized")
	}
}

type mockScene struct {
	updated bool
	drawn   bool
}

func (m *mockScene) Update(gs *GameState) error {
	m.updated = true
	return nil
}

func (m *mockScene) Draw(img *ebiten.Image) {
	m.drawn = true
}

func TestSceneManager_GoTo(t *testing.T) {
	sm := NewSceneManager(800, 600).(*sceneManager)
	initialScene := &mockScene{}
	nextScene := &mockScene{}

	sm.GoTo(initialScene)
	if sm.current != initialScene {
		t.Error("Initial scene was not set correctly")
	}

	sm.GoTo(nextScene)
	if sm.next != nextScene {
		t.Error("Next scene was not set correctly")
	}

	for i := 0; i < transitionMaxCount; i++ {
		sm.Update(nil, nil)

		if i < transitionMaxCount-1 && sm.current != initialScene {
			t.Errorf("Current scene should still be initialScene until transition completes")
		}
	}

	if sm.current != nextScene {
		t.Error("Scene did not transition correctly after transition count finished")
	}
}

func TestSceneManager_UpdateDraw(t *testing.T) {
	sm := NewSceneManager(800, 600).(*sceneManager)
	scene1 := &mockScene{}
	scene2 := &mockScene{}

	sm.GoTo(scene1)
	sm.Update(nil, nil)
	if !scene1.updated {
		t.Error("Current scene should have received update")
	}

	sm.GoTo(scene2)
	sm.Draw(ebiten.NewImage(800, 600))

	if !scene1.drawn || !scene2.drawn {
		t.Error("Both scenes should have been drawn during transition")
	}
}
