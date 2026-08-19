package app

import (
	"time"

	"github.com/d3akhtar/tfc/db"
	"github.com/rivo/tview"
)

type State struct {
	App *tview.Application

	RecentlyStudied      []db.FlashcardSet
	FlashcardSets        []db.FlashcardSet
	Folders              []db.Folder
	SelectedFlashcardSet db.FlashcardSet
	SelectedFolder       db.Folder

	Navigation *Navigation
}

func NewAppState(app *tview.Application) *State {
	nav := NewNavigation(app)
	return &State{
		App: app,

		RecentlyStudied: []db.FlashcardSet{
			{
				Name: "Stew C",
				Flashcards: []db.Flashcard{
					{
						Question: "What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? ",
						Answer:   "2",
					},
					{
						Question: "What is 1+2?",
						Answer:   "3",
					},
					{
						Question: "What is 1+3?",
						Answer:   "4",
					},
					{
						Question: "What is 1+3?",
						Answer:   "4",
					},
				},
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
		Folders: []db.Folder{
			{
				Name:          "The Stew",
				FlashcardSets: []db.FlashcardSet{},
			},
			{
				Name: "The Selected Stew",
				FlashcardSets: []db.FlashcardSet{
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
		FlashcardSets: []db.FlashcardSet{
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

		Navigation: nav,
	}
}

func (s *State) SetFocus(primitive tview.Primitive) {
	s.App.SetFocus(primitive)
}
