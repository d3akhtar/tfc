package domain

import (
	"fmt"
	"math/rand/v2"
	"slices"
	"time"
)

type FlashcardFront int

const (
	Question FlashcardFront = iota
	Answer
)

type FlashcardSet struct {
	Id            int
	Name          string
	Description   string
	LastAccessed  time.Time
	TrackProgress bool
	Front         FlashcardFront
	Shuffle       bool
	ShuffleSeed   int

	Flashcards []Flashcard
	shuffled   []Flashcard

	FlashcardCount int

	Quiz *Quiz
}

func NewFlashcardSet(name string) *FlashcardSet {
	return &FlashcardSet{
		Name: name,
		Flashcards: []Flashcard{
			{
				Question: "Question 1",
				Answer:   "Answer 1",
			},
		},
	}
}

func (f FlashcardSet) String() string {
	return fmt.Sprintf(
		"○ %s | %d flashcards | Last Accessed: %s",
		f.Name,
		f.count(),
		f.LastAccessed.Local().Format(time.RFC822),
	)
}

func (f *FlashcardSet) SwitchCardPosition(oldPos, newPos int) bool {
	if newPos < 0 || newPos > len(f.Flashcards)-1 {
		return false
	}

	f.Flashcards[newPos].Position = oldPos
	f.Flashcards[oldPos].Position = newPos

	temp := f.Flashcards[newPos]
	f.Flashcards[newPos] = f.Flashcards[oldPos]
	f.Flashcards[oldPos] = temp

	return true
}

func (f *FlashcardSet) Shuffled() bool {
	return f.Shuffle
}

func (f *FlashcardSet) SetShuffle(val bool) {
	if f.Shuffle == val {
		return
	}

	f.Shuffle = val
	if f.Shuffle {
		f.ShuffleSeed = rand.Int()
	} else {
		f.ShuffleSeed = 0
		f.shuffled = nil
	}
}

func (f *FlashcardSet) GetFlashcards() []Flashcard {
	if f.Shuffle {
		if f.shuffled == nil {
			pcg := rand.NewPCG(uint64(f.ShuffleSeed), 0)
			r := rand.New(pcg)
			f.shuffled = slices.Clone(f.Flashcards)
			r.Shuffle(len(f.shuffled), func(i, j int) {
				f.shuffled[i], f.shuffled[j] = f.shuffled[j], f.shuffled[i]
			})
		}

		return f.shuffled
	} else {
		return f.Flashcards
	}
}

func (f *FlashcardSet) StartQuiz() {
	flashcards := f.GetFlashcards()
	if f.Quiz == nil {
		f.Quiz = NewQuiz(f.Id, flashcards)
	} else {
		f.Quiz = NewQuizFromOldQuiz(f.Quiz)
	}
}

func (f *FlashcardSet) AddFlashcard(question, answer string) {
	flashcard := Flashcard{
		Question:       question,
		Answer:         answer,
		FlashcardSetId: f.Id,
		Position:       len(f.Flashcards),
	}

	f.Flashcards = append(f.Flashcards, flashcard)
}

func (f *FlashcardSet) RemoveById(id int) {
	f.Flashcards = slices.DeleteFunc(f.Flashcards, func(fc Flashcard) bool {
		return fc.Id == id
	})
}

func (f *FlashcardSet) ResetQuizProgress() {
	lastQuizId := f.Quiz.Id
	f.Quiz = NewQuiz(f.Id, f.GetFlashcards())
	f.Quiz.Id = lastQuizId
}

func (f *FlashcardSet) count() int {
	if len(f.Flashcards) == 0 {
		return f.FlashcardCount
	} else {
		return len(f.Flashcards)
	}
}
