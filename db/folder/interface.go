package folder

import (
	"context"

	"github.com/d3akhtar/tfc/db"
	"github.com/d3akhtar/tfc/domain"
)

type FolderRepo interface {
	db.Repository[domain.Folder]

	GetFlashcardSetsForFolder(ctx context.Context, entity *domain.Folder) ([]*domain.FlashcardSet, error)
	AddFlashcardSetsToFolder(ctx context.Context, entity *domain.Folder, flashcardSets []domain.FlashcardSet) error
	RemoveFlashcardSetsFromFolder(ctx context.Context, entity *domain.Folder, flashcardSets []domain.FlashcardSet) error
	FilterFolders(ctx context.Context, query string, limit, offset int) ([]*domain.Folder, error)
	FilterFlashcardSetsInFolder(ctx context.Context, entity *domain.Folder, query string, limit, offset int) ([]*domain.FlashcardSet, error)
}
