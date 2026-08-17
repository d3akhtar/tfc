package db

import (
	"bytes"
	"fmt"
)

type Collection struct {
	Name       string      `json:"name"`
	Flashcards []Flashcard `json:"flashcards"`
}

func (c Collection) String() string {
	var out bytes.Buffer
	for _, f := range c.Flashcards {
		fmt.Fprintf(&out, "%s", f)
		out.WriteString("\n")
	}

	return out.String()
}
