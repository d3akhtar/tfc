package db

type Quiz struct {
	Flashcards        []Flashcard
	CurrentlySelected int

	unknown []int

	nKnown, nUnknown int
}

func NewQuiz(flashcards []Flashcard) *Quiz {
	return &Quiz{
		Flashcards: flashcards,
		unknown:    make([]int, 0, len(flashcards)),
	}
}

func NewQuizFromOldQuiz(quiz *Quiz) *Quiz {
	cards := make([]Flashcard, 0, len(quiz.Flashcards))
	for _, i := range quiz.unknown {
		cards = append(cards, quiz.Flashcards[i])
	}

	return &Quiz{
		Flashcards: cards,
		unknown:    make([]int, 0, len(cards)),
	}
}

func (q *Quiz) AreCardsLeft() bool {
	return q.CurrentlySelected < len(q.Flashcards)
}

func (q *Quiz) IsLastCard() bool {
	return q.CurrentlySelected == len(q.Flashcards)-1
}

func (q *Quiz) CanUndo() bool {
	return q.CurrentlySelected > 0
}

func (q *Quiz) NextCard(known bool) Flashcard {
	if known {
		q.nKnown++
	} else {
		q.nUnknown++
		q.unknown = append(q.unknown, q.CurrentlySelected)
	}

	q.CurrentlySelected++
	return q.CurrentlySelectedCard()
}

func (q *Quiz) Undo() Flashcard {
	if len(q.unknown) > 0 && q.unknown[len(q.unknown)-1] == q.CurrentlySelected-1 {
		q.nUnknown--
		q.unknown = q.unknown[:len(q.unknown)-1]
	} else {
		q.nKnown--
	}

	q.CurrentlySelected--
	return q.CurrentlySelectedCard()
}

func (q *Quiz) CurrentlySelectedCard() Flashcard {
	return q.Flashcards[q.CurrentlySelected]
}

func (q *Quiz) GetKnownAndUnknownCount() (int, int) {
	return q.nKnown, q.nUnknown
}
