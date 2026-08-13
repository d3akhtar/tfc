package models

import (
	"bytes"
	"fmt"
)

type Collection struct {
	Flashcards []Flashcard
}

func (c Collection) String() string {
	var out bytes.Buffer
	for _, f := range c.Flashcards {
		fmt.Fprintf(&out, "%s", f)
		out.WriteString("\n")
	}

	return out.String()
}
