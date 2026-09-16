package domain

import "fmt"

type Flashcard struct {
	Id       int    `json:"-"`
	Question string `json:"question"`
	Answer   string `json:"answer"`

	FlashcardSetId int `json:"flashcardSetId"`
	Position       int `json:"position"`
}

func (f Flashcard) String() string {
	return fmt.Sprintf("(Flashcard) \n\tQuestion: %s\n\tAnswer: %s", f.Question, f.Answer)
}
