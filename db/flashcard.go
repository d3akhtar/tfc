package db

import "fmt"

type Flashcard struct {
	Question string `json:"question"`
	Answer   string `json:"answer"`
}

func (f Flashcard) String() string {
	return fmt.Sprintf("(Flashcard) \n\tQuestion: %s\n\tAnswer: %s", f.Question, f.Answer)
}
