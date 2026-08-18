package db

import (
	"fmt"
	"time"
)

type FlashcardSet struct {
	Name         string      `json:"name"`
	Flashcards   []Flashcard `json:"flashcards"`
	LastAccessed time.Time
}

func (f FlashcardSet) String() string {
	return fmt.Sprintf(
		"○ %s | %d flashcards | Last Accessed: %s",
		f.Name,
		len(f.Flashcards),
		f.LastAccessed.Format(time.RFC822),
	)
}
