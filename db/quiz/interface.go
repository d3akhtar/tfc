package quiz

import (
	"context"

	"github.com/d3akhtar/tfc/db"
	"github.com/d3akhtar/tfc/domain"
)

type QuizRepo interface {
	db.Repository[domain.Quiz]

	GetFlashcardSetForQuiz(ctx context.Context, entity *domain.Quiz) (*domain.FlashcardSet, error)
	GetUnknownFlashcardsForQuiz(ctx context.Context, entity *domain.Quiz) ([]domain.Flashcard, error)
	ReplaceQuiz(ctx context.Context, entity *domain.Quiz) error
}
