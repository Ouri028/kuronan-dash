// Copy from github.com/hajimehoshi/ebiten/v2/example/blocks

package scenes

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/kemokemo/kuronan-dash/internal/view"
)

const transitionMaxCount = 20

// SceneManager manages all scenes.
type SceneManager struct {
	current         Scene
	next            Scene
	transitionFrom  *ebiten.Image
	transitionTo    *ebiten.Image
	op              *ebiten.DrawImageOptions
	transitionCount int
}

// versionInfo is the version info of this game.
var versionInfo string

// NewSceneManager returns a new SceneManager.
func NewSceneManager(ver string) *SceneManager {
	sm := &SceneManager{}
	versionInfo = ver
	sm.transitionFrom = ebiten.NewImage(view.ScreenWidth, view.ScreenHeight)
	sm.transitionTo = ebiten.NewImage(view.ScreenWidth, view.ScreenHeight)
	sm.op = &ebiten.DrawImageOptions{}
	return sm
}

// Update updates the status of this scene.
func (s *SceneManager) Update() error {
	if s.transitionCount == 0 {
		return s.current.Update(&GameState{
			SceneManager: s,
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

// Draw draws background and characters. This function play music too.
func (s *SceneManager) Draw(r *ebiten.Image) {
	if s.transitionCount == 0 {
		s.current.Draw(r)
		return
	}

	s.transitionFrom.Clear()
	s.current.Draw(s.transitionFrom)
	s.transitionTo.Clear()
	s.next.Draw(s.transitionTo)
	r.DrawImage(s.transitionFrom, nil)

	alpha := 1 - float64(s.transitionCount)/float64(transitionMaxCount)
	s.op.ColorM.Scale(1, 1, 1, alpha)
	r.DrawImage(s.transitionTo, s.op)
}

// GoTo sets resources to the new scene and change the current scene
// to the new scene. This stops the music of the current and starts
// the music of the next.
func (s *SceneManager) GoTo(scene Scene) error {
	err := scene.Initialize()
	if err != nil {
		return err
	}

	if s.current == nil {
		s.current = scene
		s.current.StartMusic()
	} else {
		err = s.current.StopMusic()
		if err != nil {
			return err
		}
		s.next = scene
		s.next.StartMusic()
		s.transitionCount = transitionMaxCount
	}

	return nil
}
