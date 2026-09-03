package domain

import "github.com/d3akhtar/tfc/utils"

type Quiz struct {
	Id                     int
	FlashcardSetId         int
	Flashcards             []Flashcard
	CurrentlySelectedIndex int

	Unknown utils.HashSet[int]
}

func NewQuiz(flashcardSetId int, flashcards []Flashcard) *Quiz {
	return &Quiz{
		FlashcardSetId: flashcardSetId,
		Flashcards:     flashcards,
		Unknown:        utils.NewHashSet[int](),
	}
}

func NewQuizFromOldQuiz(quiz *Quiz) *Quiz {
	flashcards := quiz.GetUnknownFlashcards()

	return &Quiz{
		Id:             quiz.Id,
		FlashcardSetId: quiz.FlashcardSetId,
		Flashcards:     flashcards,
		Unknown:        utils.NewHashSet[int](),
	}
}

func (q *Quiz) SetUnknownCardPositions(unknown []int) {
	q.Unknown = utils.NewHashSetWithItems(unknown)
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
	if !currentCardKnown {
		q.markCurrentlySelectedCardAsUnknown()
	}

	q.CurrentlySelectedIndex++
}

func (q *Quiz) Undo() Flashcard {
	q.Unknown.Remove(q.CurrentlySelectedIndex - 1)
	q.CurrentlySelectedIndex--
	return q.CurrentlySelectedCard()
}

func (q *Quiz) CurrentlySelectedCard() Flashcard {
	return q.Flashcards[q.CurrentlySelectedIndex]
}

func (q *Quiz) GetKnownAndUnknownCount() (int, int) {
	nUnknown := q.Unknown.Length()
	nKnown := q.CurrentlySelectedIndex - nUnknown
	return nKnown, nUnknown
}

func (q *Quiz) GetUnknownFlashcards() []Flashcard {
	flashcards := make([]Flashcard, 0, len(q.Flashcards))
	for i := range q.Unknown.Items() {
		flashcards = append(flashcards, q.Flashcards[i])
	}

	return flashcards
}

func (q *Quiz) markCurrentlySelectedCardAsUnknown() {
	q.Unknown.Add(q.CurrentlySelectedIndex)
}
