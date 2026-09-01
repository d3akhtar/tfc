package app

import (
	"context"

	"github.com/d3akhtar/tfc/domain"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type State struct {
	App *tview.Application

	SelectedFlashcardSet *domain.FlashcardSet
	SelectedFolder       *domain.Folder

	Context context.Context

	Navigation *Navigation
}

func NewApp(app *tview.Application) *State {
	nav := NewNavigation(app)

	views := nav.Views()

	views.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyCtrlBackslash:
			nav.RevertView()
		}

		return event
	})

	state := &State{
		App:        app,
		Navigation: nav,
		Context:    context.Background(),
	}

	return state
}

func (s *State) SetFocus(primitive tview.Primitive) {
	s.App.SetFocus(primitive)
}
