package folder

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/d3akhtar/tfc/db"
	"github.com/d3akhtar/tfc/domain"
	"github.com/mattn/go-sqlite3"
)

type FolderRepository struct {
	db *sql.DB
}

// Ensure interface is implemented at compile-time
var _ FolderRepo = (*FolderRepository)(nil)

func NewFolderRepository(db *sql.DB) *FolderRepository {
	return &FolderRepository{db}
}

func (r *FolderRepository) AddFlashcardSetsToFolder(ctx context.Context, entity *domain.Folder, flashcardSets []domain.FlashcardSet) error {
	return db.WithTransaction(ctx, r.db, func(tx *sql.Tx) error {
		stmt, err := tx.PrepareContext(ctx, `
			INSERT INTO FolderFlashcardSet (FolderId, FlashcardSetId)
			VALUES (?, ?)
		`)

		if err != nil {
			return fmt.Errorf("Failed to prepare task statement: %w", err)
		}

		defer stmt.Close()

		for _, flashcardSet := range flashcardSets {
			_, err := stmt.ExecContext(ctx,
				entity.Id,
				flashcardSet.Id,
			)

			if err != nil {
				return fmt.Errorf("Failed to add flashcard set with Id %d to folder with Id %d", flashcardSet.Id, entity.Id)
			}
		}

		entity.FlashcardSets = flashcardSets

		return nil
	})
}

func (r *FolderRepository) Count(ctx context.Context) (int64, error) {
	return db.CountQuery(r.db, ctx, "Folders")
}

func (r *FolderRepository) Create(ctx context.Context, entity *domain.Folder) error {
	query := `
		INSERT INTO Folders (Name)
		VALUES ($1)
	`

	result, err := r.db.ExecContext(ctx, query, entity.Name)

	if err != nil {
		var sqliteErr sqlite3.Error
		if errors.As(err, &sqliteErr) {
			switch sqliteErr.ExtendedCode {
			case sqlite3.ErrConstraintPrimaryKey:
				return db.ErrDuplicate
			}
		}

		return err
	}

	folderId, _ := result.LastInsertId()
	entity.Id = int(folderId)

	return r.AddFlashcardSetsToFolder(ctx, entity, entity.FlashcardSets)
}

func (r *FolderRepository) Delete(ctx context.Context, id int) error {
	return db.DeleteQuery(r.db, ctx, "Folders", id)
}

func (r *FolderRepository) FilterFlashcardSetsInFolder(ctx context.Context, entity *domain.Folder, filter string, limit, offset int) ([]*domain.FlashcardSet, error) {
	query := fmt.Sprintf(`
		SELECT
			fs.Id,
			fs.Name,
			fs.Description,
			fs.LastAccessed,
			fs.TrackProgress,
			fs.Front,
			fs.Shuffle,
			fs.ShuffleSeed
		FROM Folders f
		INNER JOIN FolderFlashcardSet ffs ON ffs.FolderId = f.Id
		INNER JOIN FlashcardSets fs ON fs.Id = ffs.FlashcardSetId
		WHERE f.Id = $1
		AND fs.Name LIKE '%%%s%%'
		ORDER BY fs.Id ASC
		LIMIT $2 OFFSET $3`,
		filter,
	)

	rows, err := r.db.QueryContext(ctx, query, entity.Id, limit, offset)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	flashcardSets := []*domain.FlashcardSet{}
	for rows.Next() {
		fs := &domain.FlashcardSet{}
		err := rows.Scan(
			&fs.Id,
			&fs.Name,
			&fs.Description,
			&fs.LastAccessed,
			&fs.TrackProgress,
			&fs.Front,
			&fs.Shuffle,
			&fs.ShuffleSeed,
		)

		if err != nil {
			return nil, err
		}

		flashcardSets = append(flashcardSets, fs)
	}

	if errors.Is(err, sql.ErrNoRows) {
		return nil, db.ErrNotFound
	}

	return flashcardSets, err
}

func (r *FolderRepository) GetById(ctx context.Context, id int) (*domain.Folder, error) {
	query := `
		SELECT Id, Name, LastAccessed FROM Folders
		WHERE Id = $1
	`

	folder := &domain.Folder{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&folder.Id,
		&folder.Name,
		&folder.LastAccessed,
	)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, db.ErrNotFound
	}

	if err != nil {
		return nil, err
	}

	return folder, nil
}

func (r *FolderRepository) GetFlashcardSetsForFolder(ctx context.Context, entity *domain.Folder) ([]*domain.FlashcardSet, error) {
	query := `
		SELECT
			fs.Id,
			fs.Name,
			fs.Description,
			fs.LastAccessed,
			fs.TrackProgress,
			fs.Front,
			fs.Shuffle,
			fs.ShuffleSeed
		FROM Folders f
		INNER JOIN FolderFlashcardSet ffs ON ffs.FolderId = f.Id
		INNER JOIN FlashcardSets fs ON fs.Id = ffs.FlashcardSetId
		WHERE f.Id = $1
		ORDER BY fs.Id Asc
	`

	rows, err := r.db.QueryContext(ctx, query, entity.Id)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	flashcardSets := []*domain.FlashcardSet{}
	for rows.Next() {
		fc := &domain.FlashcardSet{}
		err := rows.Scan(
			&fc.Id,
			&fc.Name,
			&fc.Description,
			&fc.LastAccessed,
			&fc.TrackProgress,
			&fc.Front,
			&fc.Shuffle,
			&fc.ShuffleSeed,
		)

		if err != nil {
			return nil, err
		}

		flashcardSets = append(flashcardSets, fc)
	}

	if errors.Is(err, sql.ErrNoRows) {
		return nil, db.ErrNotFound
	}

	return flashcardSets, err
}

func (r *FolderRepository) List(ctx context.Context, offset int, limit int) ([]*domain.Folder, error) {
	query := `
		SELECT f.Id, f.Name, f.LastAccessed
		FROM Folders f
		GROUP BY f.Id
		ORDER BY f.Id ASC
		LIMIT $1 OFFSET $2
	`

	rows, err := r.db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	folders := []*domain.Folder{}
	for rows.Next() {
		folder := &domain.Folder{}
		err := rows.Scan(
			&folder.Id,
			&folder.Name,
			&folder.LastAccessed,
		)

		if err != nil {
			return nil, err
		}

		folders = append(folders, folder)
	}

	return folders, nil
}

func (r *FolderRepository) Update(ctx context.Context, entity *domain.Folder) error {
	query := `
		UPDATE Folders
		SET Name = $1, LastAccessed = $2
		WHERE Id = $3
	`
	result, err := r.db.ExecContext(ctx, query, entity.Name, entity.LastAccessed, entity.Id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return db.ErrNotFound
	}

	return nil
}

func (r *FolderRepository) RemoveFlashcardSetsFromFolder(ctx context.Context, entity *domain.Folder, flashcardSets []domain.FlashcardSet) error {
	return db.WithTransaction(ctx, r.db, func(tx *sql.Tx) error {
		stmt, err := tx.PrepareContext(ctx, `
			DELETE FROM FolderFlashcardSet
			WHERE FolderId = $1
			AND FlashcardSetId = $2
		`)

		if err != nil {
			return fmt.Errorf("Error while preparing statement: %w", err)
		}

		defer stmt.Close()

		for _, flashcardSet := range flashcardSets {
			_, err := stmt.ExecContext(ctx, entity.Id, flashcardSet.Id)

			if err != nil {
				return fmt.Errorf("Error while removing flashcard set with Id %d to folder with Id %d", flashcardSet.Id, entity.Id)
			}
		}

		return nil
	})
}

func (r *FolderRepository) FilterFolders(ctx context.Context, filter string, limit, offset int) ([]*domain.Folder, error) {
	query := fmt.Sprintf(`
		SELECT
			f.Id,
			f.Name,
			f.LastAccessed,
			COUNT(ffs.FlashcardSetId)
		FROM Folders f
		JOIN FolderFlashcardSet ffs ON ffs.FolderId = f.Id
		WHERE f.Name LIKE '%%%s%%'
		GROUP BY f.Id
		ORDER BY f.Id ASC
		LIMIT $1 OFFSET $2`,
		filter,
	)

	rows, err := r.db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	folders := []*domain.Folder{}
	for rows.Next() {
		f := &domain.Folder{}
		err := rows.Scan(
			&f.Id,
			&f.Name,
			&f.LastAccessed,
			&f.FlashcardSetCount,
		)

		if err != nil {
			return nil, err
		}

		folders = append(folders, f)
	}

	if errors.Is(err, sql.ErrNoRows) {
		return nil, db.ErrNotFound
	}

	return folders, err
}
