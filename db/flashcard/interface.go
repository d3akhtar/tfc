package flashcard

import (
	"context"

	"github.com/d3akhtar/tfc/db"
	"github.com/d3akhtar/tfc/domain"
)

type FlashcardRepo interface {
	db.Repository[domain.Flashcard]

	GetFlashcardSetForFlashcard(ctx context.Context, entity *domain.Flashcard) (*domain.FlashcardSet, error)
	UpdateFlashcardPosition(ctx context.Context, entity *domain.Flashcard, newPos int) error
}
