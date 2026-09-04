package app

import (
	"context"
	"time"

	"github.com/d3akhtar/tfc/domain"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type State struct {
	App        *tview.Application
	Context    context.Context
	Navigation *Navigation

	selectedFlashcardSet *domain.FlashcardSet
	selectedFolder       *domain.Folder

	onSelectedFlashcardSetChange []func(*domain.FlashcardSet)
	onSelectedFolderChange       []func(*domain.Folder)
}

func NewApp(app *tview.Application) *State {
	nav := NewNavigation(app)

	views := nav.Views()

	app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyCtrlC:
			err := nav.views[nav.MostRecentlyVisitedViewName()].exit()
			if err != nil {
				panic(err)
			}
		}

		return event
	})

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

func (s *State) SelectedFlashcardSet() *domain.FlashcardSet {
	return s.selectedFlashcardSet
}

func (s *State) SelectedFolder() *domain.Folder {
	return s.selectedFolder
}

func (s *State) SetSelectedFlashcardSet(flashcardSet *domain.FlashcardSet) {
	s.selectedFlashcardSet = flashcardSet

	if s.selectedFlashcardSet != nil {
		s.selectedFlashcardSet.LastAccessed = time.Now()
		for _, callback := range s.onSelectedFlashcardSetChange {
			callback(s.selectedFlashcardSet)
		}
	}
}

func (s *State) SetSelectedFolder(folder *domain.Folder) {
	s.selectedFolder = folder

	if s.selectedFolder != nil {
		s.selectedFolder.LastAccessed = time.Now()
		for _, callback := range s.onSelectedFolderChange {
			callback(s.selectedFolder)
		}
	}
}

func (s *State) AddCallbackForSelectedFlashcardSetChange(callback func(*domain.FlashcardSet)) {
	s.onSelectedFlashcardSetChange = append(s.onSelectedFlashcardSetChange, callback)
}

func (s *State) AddCallbackForSelectedFolderChange(callback func(*domain.Folder)) {
	s.onSelectedFolderChange = append(s.onSelectedFolderChange, callback)
}

func (s *State) SetFocus(primitive tview.Primitive) {
	s.App.SetFocus(primitive)
}
