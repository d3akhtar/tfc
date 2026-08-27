package app

import (
	"time"

	"github.com/d3akhtar/tfc/domain"
	"github.com/rivo/tview"
)

type State struct {
	App *tview.Application

	RecentlyStudied      []domain.FlashcardSet
	FlashcardSets        []domain.FlashcardSet
	Folders              []domain.Folder
	SelectedFlashcardSet *domain.FlashcardSet
	SelectedFolder       *domain.Folder

	Navigation *Navigation
}

func NewAppState(app *tview.Application) *State {
	nav := NewNavigation(app)
	return &State{
		App: app,

		RecentlyStudied: []domain.FlashcardSet{
			{
				Name: "Stew C",
				Flashcards: []domain.Flashcard{
					{
						Question: "What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? What is 1+1? ",
						Answer:   "2",
					},
					{
						Question: "What is 1+2?",
						Answer:   "3",
					},
					{
						Question: "What is 1+4?",
						Answer:   "5",
					},
					{
						Question: "What is 1+5?",
						Answer:   "6",
					},
					{
						Question: "What is 1+6?",
						Answer:   "7",
					},
					{
						Question: "What is 1+7?",
						Answer:   "8",
					},
				},
				LastAccessed: time.Now(),
			},
			{
				Name: "Stew B",
				Flashcards: []domain.Flashcard{
					{},
				},
				LastAccessed: time.Now().Add(-15 * time.Minute),
			},
			{
				Name: "Stew A",
				Flashcards: []domain.Flashcard{
					{}, {},
				},
				LastAccessed: time.Now().Add(-30 * time.Minute),
			},
		},
		Folders: []domain.Folder{
			{
				Name:          "The Stew",
				FlashcardSets: []domain.FlashcardSet{},
			},
			{
				Name: "The Selected Stew",
				FlashcardSets: []domain.FlashcardSet{
					{
						Name:         "Stew C",
						Flashcards:   []domain.Flashcard{},
						LastAccessed: time.Now(),
					},
					{
						Name: "Stew B",
						Flashcards: []domain.Flashcard{
							{},
						},
						LastAccessed: time.Now().Add(-15 * time.Minute),
					},
					{
						Name: "Stew A",
						Flashcards: []domain.Flashcard{
							{}, {},
						},
						LastAccessed: time.Now().Add(-30 * time.Minute),
					},
				},
			},
		},
		FlashcardSets: []domain.FlashcardSet{
			{
				Name:         "Stew C",
				Flashcards:   []domain.Flashcard{},
				LastAccessed: time.Now(),
			},
			{
				Name: "Stew B",
				Flashcards: []domain.Flashcard{
					{},
				},
				LastAccessed: time.Now().Add(-15 * time.Minute),
			},
			{
				Name: "Stew A",
				Flashcards: []domain.Flashcard{
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
