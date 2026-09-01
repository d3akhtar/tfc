package app

import (
	"context"

	"github.com/d3akhtar/tfc/domain"
	"github.com/rivo/tview"
)

type State struct {
	App *tview.Application

	SelectedFlashcardSet *domain.FlashcardSet
	SelectedFolder       *domain.Folder

	Context context.Context

	Navigation *Navigation
}

func NewAppState(app *tview.Application) *State {
	nav := NewNavigation(app)
	return &State{
		App:        app,
		Navigation: nav,
		Context:    context.Background(),
	}
}

func (s *State) SetFocus(primitive tview.Primitive) {
	s.App.SetFocus(primitive)
}
