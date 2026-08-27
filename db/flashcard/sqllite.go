package flashcard

import (
	"context"
	"database/sql"
	"errors"

	"github.com/d3akhtar/tfc/db"
	"github.com/d3akhtar/tfc/domain"
	"github.com/mattn/go-sqlite3"
)

type FlashcardRepository struct {
	db *sql.DB
}

// Ensure interface is implemented at compile-time
var _ FlashcardRepo = (*FlashcardRepository)(nil)

func NewFlashcardRepository(db *sql.DB) *FlashcardRepository {
	return &FlashcardRepository{db}
}

func (r *FlashcardRepository) Create(ctx context.Context, entity *domain.Flashcard) error {
	query := `
		INSERT INTO Flashcards (Question, Answer, FlashcardSetId, Position)
		VALUES ($1, $2, $3, $4)
	`

	result, err := r.db.ExecContext(ctx, query,
		entity.Question,
		entity.Answer,
		entity.FlashcardSetId,
		entity.Position,
	)

	if err != nil {
		var sqliteErr sqlite3.Error
		if errors.As(err, &sqliteErr) {
			switch sqliteErr.ExtendedCode {
			case sqlite3.ErrConstraintPrimaryKey, sqlite3.ErrConstraintUnique:
				return db.ErrDuplicate
			}
		}

		return err
	}

	flashcardId, _ := result.LastInsertId()
	entity.Id = int(flashcardId)

	return nil
}

func (r *FlashcardRepository) GetById(ctx context.Context, id int) (*domain.Flashcard, error) {
	query := `
		SELECT Id, Question, Answer, FlashcardSetId, Position FROM Flashcards WHERE Id = $1
	`

	flashcard := &domain.Flashcard{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&flashcard.Id,
		&flashcard.Question,
		&flashcard.Answer,
		&flashcard.FlashcardSetId,
		&flashcard.Position,
	)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, db.ErrNotFound
	}

	if err != nil {
		return nil, err
	}

	return flashcard, nil
}

func (r *FlashcardRepository) Update(ctx context.Context, entity *domain.Flashcard) error {
	query := `
		UPDATE Flashcards
		SET Question = $1, Answer = $2, Position = $3
		WHERE Id = $4
	`

	result, err := r.db.ExecContext(ctx, query, entity.Question, entity.Answer, entity.Position, entity.Id)
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

func (r *FlashcardRepository) Delete(ctx context.Context, id int) error {
	return db.DeleteQuery(r.db, ctx, "Flashcards", id)
}

func (r *FlashcardRepository) List(ctx context.Context, offset, limit int) ([]*domain.Flashcard, error) {
	query := `
		SELECT Id, Question, Answer, FlashcardSetId, Position
		FROM Flashcards
		ORDER BY Id ASC
		LIMIT $1 OFFSET $2
	`

	rows, err := r.db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var flashcards []*domain.Flashcard
	for rows.Next() {
		fc := &domain.Flashcard{}
		err := rows.Scan(
			&fc.Id,
			&fc.Question,
			&fc.Answer,
			&fc.FlashcardSetId,
			&fc.Position,
		)

		if err != nil {
			return nil, err
		}

		flashcards = append(flashcards, fc)
	}

	return flashcards, nil
}

func (r *FlashcardRepository) Count(ctx context.Context) (int64, error) {
	return db.CountQuery(r.db, ctx, "Flashcards")
}

func (r *FlashcardRepository) GetFlashcardSetForFlashcard(ctx context.Context, entity *domain.Flashcard) (*domain.FlashcardSet, error) {
	query := `
		SELECT fs.Id, fs.Name, fs.Description, fs.LastAccessed, fs.TrackProgress, fs.Front, fs.Shuffle, fs.ShuffleSeed
		FROM Flashcards fc
		INNER JOIN FlashcardSets fs ON fs.Id = fc.FlashcardSetId
		WHERE fc.Id = $1
	`

	flashcardSet := &domain.FlashcardSet{}
	err := r.db.QueryRowContext(ctx, query, entity.Id).Scan(
		&flashcardSet.Id,
		&flashcardSet.Name,
		&flashcardSet.Description,
		&flashcardSet.LastAccessed,
		&flashcardSet.TrackProgress,
		&flashcardSet.Front,
		&flashcardSet.Shuffle,
		&flashcardSet.ShuffleSeed,
	)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, db.ErrNotFound
	}

	if err != nil {
		return nil, err
	}

	return flashcardSet, nil
}

func (r *FlashcardRepository) UpdateFlashcardPosition(ctx context.Context, entity *domain.Flashcard, newPos int) error {
	query := `
		UPDATE Flashcards
		SET Position = $1
		WHERE Id = $2
	`

	result, err := r.db.ExecContext(ctx, query, newPos, entity.Id)
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
