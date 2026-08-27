package domain

import "fmt"

type Flashcard struct {
	Id       int
	Question string
	Answer   string

	FlashcardSetId int
	Position       int
}

func (f Flashcard) String() string {
	return fmt.Sprintf("(Flashcard) \n\tQuestion: %s\n\tAnswer: %s", f.Question, f.Answer)
}
