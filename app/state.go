package app

import (
	"time"

	"github.com/d3akhtar/tfc/db"
	"github.com/rivo/tview"
)

type State struct {
	App *tview.Application

	RecentlyStudied    []db.Collection
	Collection         []db.Collection
	Folders            []db.Folder
	SelectedCollection db.Collection
	SelectedFolder     db.Folder

	Navigation *Navigation

	onSelectedCollectionChange []func(db.Collection)
	onSelectedFolderChange     []func(db.Folder)
}

func NewAppState(app *tview.Application, views *tview.Pages) *State {
	nav := NewNavigation(views)
	return &State{
		App: app,

		RecentlyStudied: []db.Collection{},
		Folders: []db.Folder{
			{
				Name:        "The Stew",
				Collections: []db.Collection{},
			},
			{
				Name: "The Selected Stew",
				Collections: []db.Collection{
					{
						Name:         "Stew C",
						Flashcards:   []db.Flashcard{},
						LastAccessed: time.Now(),
					},
					{
						Name: "Stew B",
						Flashcards: []db.Flashcard{
							{},
						},
						LastAccessed: time.Now().Add(-15 * time.Minute),
					},
					{
						Name: "Stew A",
						Flashcards: []db.Flashcard{
							{}, {},
						},
						LastAccessed: time.Now().Add(-30 * time.Minute),
					},
				},
			},
		},

		Navigation: nav,
	}
}

func (s *State) SetSelectedFolder(folder db.Folder) {
	s.SelectedFolder = folder
	for _, callback := range s.onSelectedFolderChange {
		callback(s.SelectedFolder)
	}
}

func (s *State) OnSelectedFolderChange(callback func(db.Folder)) {
	s.onSelectedFolderChange = append(s.onSelectedFolderChange, callback)
}

func (s *State) SetSelectedCollection(collection db.Collection) {
	s.SelectedCollection = collection
	for _, callback := range s.onSelectedCollectionChange {
		callback(s.SelectedCollection)
	}
}

func (s *State) OnSelectedCollectionChange(callback func(db.Collection)) {
	s.onSelectedCollectionChange = append(s.onSelectedCollectionChange, callback)
}
