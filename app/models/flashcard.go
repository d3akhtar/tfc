package models

import "fmt"

type Flashcard struct {
	Question string
	Answer   string
}

func (f Flashcard) String() string {
	return fmt.Sprint("(Flashcard) \n\tQuestion: %s\n\tAnswer: %s", f.Question, f.Answer)
}
