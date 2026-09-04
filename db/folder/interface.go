package folder

import (
	"context"

	"github.com/d3akhtar/tfc/db"
	"github.com/d3akhtar/tfc/domain"
)

type FolderRepo interface {
	db.Repository[domain.Folder]

	ListRecentlyAccessedFolders(ctx context.Context) ([]*domain.Folder, error)
	GetFlashcardSetsForFolder(ctx context.Context, entity *domain.Folder) ([]*domain.FlashcardSet, error)
	FilterFolders(ctx context.Context, query string, limit, offset int, sort db.FilterSortCriteria) ([]*domain.Folder, error)
	FilterFlashcardSetsInFolder(ctx context.Context, entity *domain.Folder, query string, limit, offset int, sort db.FilterSortCriteria) ([]*domain.FlashcardSet, error)
}
