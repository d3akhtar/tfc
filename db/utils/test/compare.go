package test_utils

import (
	"fmt"
	"slices"
	"time"

	"github.com/d3akhtar/tfc/domain"
)

func TestCmpQuiz(expected, actual *domain.Quiz) error {
	if expected.Id != actual.Id {
		return fmt.Errorf("quiz Id expected=%v, got=%v", expected.Id, actual.Id)
	}

	if expected.FlashcardSetId != actual.FlashcardSetId {
		return fmt.Errorf("quiz FlashcardSetId expected=%v, got=%v", expected.FlashcardSetId, actual.FlashcardSetId)
	}

	if !slices.Equal(expected.Flashcards, actual.Flashcards) {
		return fmt.Errorf("quiz Flashcards expected=%v, got=%v", expected.Flashcards, actual.Flashcards)
	}

	if expected.CurrentlySelectedIndex != actual.CurrentlySelectedIndex {
		return fmt.Errorf("quiz CurrentlySelectedIndex expected=%v, got=%v", expected.CurrentlySelectedIndex, actual.CurrentlySelectedIndex)
	}

	if expected.Unknown.Length() != actual.Unknown.Length() {
		return fmt.Errorf("quiz.Unknown.Length() expected=%v, got=%v", expected.Unknown.Length(), actual.Unknown.Length())
	}

	nonMatching := make([]int, 0, expected.Unknown.Length())
	for u := range expected.Unknown.Items() {
		if !actual.Unknown.Contains(u) {
			nonMatching = append(nonMatching, u)
		}
	}

	if len(nonMatching) != 0 {
		return fmt.Errorf("quiz.Unknown unmatched values %v", nonMatching)
	}

	return nil
}

func TestCmpFlashcardSet(expected, actual *domain.FlashcardSet) error {
	if expected.Id != actual.Id {
		return fmt.Errorf("flashcard set Id expected=%v, got=%v", expected.Id, actual.Id)
	}

	if expected.Name != actual.Name {
		return fmt.Errorf("flashcard set Name expected=%v, got=%v", expected.Name, actual.Name)
	}

	if expected.Description != actual.Description {
		return fmt.Errorf("flashcard set Description expected=%v, got=%v", expected.Description, actual.Description)
	}

	if expected.LastAccessed.UTC() != actual.LastAccessed.UTC() {
		return fmt.Errorf("flashcard set LastAccessedexpected=%v, got=%v", expected.LastAccessed.UTC(), actual.LastAccessed.UTC())
	}

	if expected.TrackProgress != actual.TrackProgress {
		return fmt.Errorf("flashcard set TrackProgress expected=%v, got=%v", expected.TrackProgress, actual.TrackProgress)
	}

	if expected.Front != actual.Front {
		return fmt.Errorf("flashcard set Front expected=%v, got=%v", expected.Front, actual.Front)
	}

	if expected.Shuffle != actual.Shuffle {
		return fmt.Errorf("flashcard set Shuffle expected=%v, got=%v", expected.Shuffle, actual.Shuffle)
	}

	if expected.ShuffleSeed != actual.ShuffleSeed {
		return fmt.Errorf("flashcard set ShuffleSeed expected=%v, got=%v", expected.ShuffleSeed, actual.ShuffleSeed)
	}

	return nil
}

func TestCmpFolder(expected, actual *domain.Folder) error {
	if expected.Id != actual.Id {
		return fmt.Errorf("folder Id expected=%v, got=%v", expected.Id, actual.Id)
	}

	if expected.Name != actual.Name {
		return fmt.Errorf("folder Name expected=%v, got=%v", expected.Name, actual.Name)
	}

	if expected.LastAccessed.UTC().Truncate(time.Second).Compare(actual.LastAccessed.Truncate(time.Second)) != 0 {
		return fmt.Errorf(
			"folder LastAccessed expected=%v, got=%v",
			expected.LastAccessed.UTC().Truncate(time.Second),
			actual.LastAccessed.UTC().Truncate(time.Second))
	}

	return nil
}
