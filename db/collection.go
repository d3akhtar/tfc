package db

import (
	"fmt"
	"time"
)

type Collection struct {
	Name         string      `json:"name"`
	Flashcards   []Flashcard `json:"flashcards"`
	LastAccessed time.Time
}

func (c Collection) String() string {
	return fmt.Sprintf(
		"○ %s | %d flashcards | Last Accessed: %s",
		c.Name,
		len(c.Flashcards),
		c.LastAccessed.Format(time.RFC822),
	)
}
