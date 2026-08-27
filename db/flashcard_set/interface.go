package flashcard_set

import (
	"context"

	"github.com/d3akhtar/tfc/db"
	"github.com/d3akhtar/tfc/domain"
)

type FlashcardSetRepo interface {
	db.Repository[domain.FlashcardSet]

	GetAllFlashcardsForSet(ctx context.Context, entity *domain.FlashcardSet) ([]domain.Flashcard, error)
	GetQuizForFlashcardSet(ctx context.Context, entity *domain.FlashcardSet) (*domain.Quiz, error)
	FilterFlashcardSets(ctx context.Context, query string, limit, offset int) ([]domain.FlashcardSet, error)
}
