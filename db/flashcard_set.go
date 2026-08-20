package db

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
	Name          string `json:"name"`
	LastAccessed  time.Time
	TrackProgress bool
	Front         FlashcardFront

	Flashcards []Flashcard `json:"flashcards"`
	shuffled   []Flashcard

	shuffle     bool
	shuffleSeed int

	Quiz *Quiz
}

func (f FlashcardSet) String() string {
	return fmt.Sprintf(
		"○ %s | %d flashcards | Last Accessed: %s",
		f.Name,
		len(f.Flashcards),
		f.LastAccessed.Format(time.RFC822),
	)
}

func (f *FlashcardSet) Shuffled() bool {
	return f.shuffle
}

func (f *FlashcardSet) SetShuffle(val bool) {
	f.shuffle = val
	if f.shuffle {
		f.shuffleSeed = rand.Int()
	} else {
		f.shuffled = nil
	}
}

func (f *FlashcardSet) GetFlashcards() []Flashcard {
	if f.shuffle {
		if f.shuffled == nil {
			pcg := rand.NewPCG(uint64(f.shuffleSeed), 0)
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
		f.Quiz = NewQuiz(flashcards)
	} else {
		f.Quiz = NewQuizFromOldQuiz(f.Quiz)
	}
}
