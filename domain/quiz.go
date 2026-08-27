package domain

type Quiz struct {
	Id                     int
	FlashcardSetId         int
	Flashcards             []Flashcard
	CurrentlySelectedIndex int

	Unknown []int

	nKnown, nUnknown int
}

func NewQuiz(flashcardSetId int, flashcards []Flashcard) *Quiz {
	return &Quiz{
		FlashcardSetId: flashcardSetId,
		Flashcards:     flashcards,
		Unknown:        make([]int, 0, len(flashcards)),
	}
}

func NewQuizFromOldQuiz(quiz *Quiz) *Quiz {
	flashcards := quiz.GetUnknownFlashcards()

	return &Quiz{
		Id:             quiz.Id,
		FlashcardSetId: quiz.FlashcardSetId,
		Flashcards:     flashcards,
		Unknown:        make([]int, 0, len(flashcards)),
	}
}

func (q *Quiz) SetUnknownCardPositions(unknown []int) {
	q.Unknown = unknown
	q.nUnknown = len(unknown)
	q.nKnown = q.CurrentlySelectedIndex - q.nUnknown
}

func (q *Quiz) AreCardsLeft() bool {
	return q.CurrentlySelectedIndex < len(q.Flashcards)
}

func (q *Quiz) Finished() bool {
	return q.CurrentlySelectedIndex >= len(q.Flashcards)
}

func (q *Quiz) CanUndo() bool {
	return q.CurrentlySelectedIndex > 0
}

func (q *Quiz) GoToNextCard(currentCardKnown bool) {
	if currentCardKnown {
		q.nKnown++
	} else {
		q.markCurrentlySelectedCardAsUnknown()
	}

	q.CurrentlySelectedIndex++
}

func (q *Quiz) Undo() Flashcard {
	if len(q.Unknown) > 0 && q.Unknown[len(q.Unknown)-1] == q.CurrentlySelectedIndex-1 {
		q.nUnknown--
		q.Unknown = q.Unknown[:len(q.Unknown)-1]
	} else {
		q.nKnown--
	}

	q.CurrentlySelectedIndex--
	return q.CurrentlySelectedCard()
}

func (q *Quiz) CurrentlySelectedCard() Flashcard {
	return q.Flashcards[q.CurrentlySelectedIndex]
}

func (q *Quiz) GetKnownAndUnknownCount() (int, int) {
	return q.nKnown, q.nUnknown
}

func (q *Quiz) GetUnknownFlashcards() []Flashcard {
	flashcards := make([]Flashcard, 0, len(q.Flashcards))
	for _, i := range q.Unknown {
		flashcards = append(flashcards, q.Flashcards[i])
	}

	return flashcards
}

func (q *Quiz) markCurrentlySelectedCardAsUnknown() {
	q.nUnknown++
	q.Unknown = append(q.Unknown, q.CurrentlySelectedIndex)
}
