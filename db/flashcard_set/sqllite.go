package flashcard_set

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/d3akhtar/tfc/db"
	"github.com/d3akhtar/tfc/db/utils"
	"github.com/d3akhtar/tfc/domain"
	_ "github.com/mattn/go-sqlite3"
)

type FlashcardSetRepository struct {
	db *sql.DB
}

// Ensure interface is implemented at compile-time
var _ FlashcardSetRepo = (*FlashcardSetRepository)(nil)

func NewFlashcardSetRepository(db *sql.DB) *FlashcardSetRepository {
	return &FlashcardSetRepository{db}
}

func (r *FlashcardSetRepository) Count(ctx context.Context) (int64, error) {
	return utils.CountQuery(r.db, ctx, "FlashcardSets")
}

func (r *FlashcardSetRepository) Create(ctx context.Context, entity *domain.FlashcardSet) error {
	return db.WithTransaction(ctx, r.db, func(tx *sql.Tx) error {
		stmt, err := tx.PrepareContext(ctx, `
			INSERT INTO FlashcardSets (Name, Description, LastAccessed, TrackProgress, Front, Shuffle, ShuffleSeed)
			VALUES (?, ?, ?, ?, ?, ?, ?)
		`)

		if err != nil {
			return fmt.Errorf("Error while preparing statement for flashcard set creation %w", err)
		}

		result, err := stmt.ExecContext(ctx,
			entity.Name,
			entity.Description,
			entity.LastAccessed,
			entity.TrackProgress,
			entity.Front,
			entity.Shuffle,
			entity.ShuffleSeed,
		)

		if err != nil {
			return fmt.Errorf("Error while adding flashcard set '%s'", entity.Name)
		}

		flashcardSetId, _ := result.LastInsertId()

		entity.Id = int(flashcardSetId)

		stmt.Close()

		if len(entity.Flashcards) == 0 {
			return nil
		}

		stmt, err = tx.PrepareContext(ctx, `
			INSERT INTO Flashcards (Question, Answer, FlashcardSetId, Position)
			VALUES (?, ?, ?, ?)
		`)

		if err != nil {
			return fmt.Errorf("Error while preparing statement for flashcard set flashcard creation %w", err)
		}

		for i, flashcard := range entity.Flashcards {
			flashcard.FlashcardSetId = int(flashcardSetId)

			result, err = stmt.ExecContext(ctx,
				flashcard.Question,
				flashcard.Answer,
				flashcard.FlashcardSetId,
				flashcard.Position,
			)

			if err != nil {
				return fmt.Errorf("Error while adding flashcard for flashcard set with Id %d: %w", entity.Id, err)
			}

			flashcardId, _ := result.LastInsertId()

			entity.Flashcards[i].Id = int(flashcardId)
		}

		return nil
	})
}

func (r *FlashcardSetRepository) Delete(ctx context.Context, id int) error {
	return utils.DeleteQuery(r.db, ctx, "FlashcardSets", id)
}

func (r *FlashcardSetRepository) FilterFlashcardSets(ctx context.Context, filter string, limit, offset int) ([]domain.FlashcardSet, error) {
	query := fmt.Sprintf(`
		SELECT Id, Name, Description, LastAccessed, TrackProgress, Front, Shuffle, ShuffleSeed
		FROM FlashcardSets
		WHERE Name LIKE '%%%s%%'
		ORDER BY Id ASC
		LIMIT $1 OFFSET $2
		`,
		filter,
	)

	rows, err := r.db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	flashcardSets := []domain.FlashcardSet{}
	for rows.Next() {
		fs := domain.FlashcardSet{}
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

func (r *FlashcardSetRepository) GetAllFlashcardsForSet(ctx context.Context, entity *domain.FlashcardSet) ([]domain.Flashcard, error) {
	query := `
		SELECT
			fc.Id,
			fc.Question,
			fc.Answer,
			fc.FlashcardSetId,
			fc.Position
		FROM FlashcardSets fs
		JOIN Flashcards fc ON fc.FlashcardSetId = fs.Id
		WHERE fs.Id = $1
	`

	rows, err := r.db.QueryContext(ctx, query, entity.Id)
	if err != nil {
		return nil, err
	}

	defer rows.Close()
	flashcards := []domain.Flashcard{}
	for rows.Next() {
		fc := domain.Flashcard{}
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

func (r *FlashcardSetRepository) GetById(ctx context.Context, id int) (*domain.FlashcardSet, error) {
	query := `
		SELECT Id, Name, Description, LastAccessed, TrackProgress, Front, Shuffle, ShuffleSeed
		FROM FlashcardSets
		WHERE Id = $1
	`

	flashcardSet := &domain.FlashcardSet{}
	var shuffleSeed sql.NullInt32
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&flashcardSet.Id,
		&flashcardSet.Name,
		&flashcardSet.Description,
		&flashcardSet.LastAccessed,
		&flashcardSet.TrackProgress,
		&flashcardSet.Front,
		&flashcardSet.Shuffle,
		&shuffleSeed,
	)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, db.ErrNotFound
	}

	if err != nil {
		return nil, err
	}

	if shuffleSeed.Valid {
		flashcardSet.ShuffleSeed = int(shuffleSeed.Int32)
	}

	return flashcardSet, nil
}

func (r *FlashcardSetRepository) GetQuizForFlashcardSet(ctx context.Context, entity *domain.FlashcardSet) (*domain.Quiz, error) {
	query := `
		SELECT
			q.Id,
			q.CurrentlySelectedIndex,
			quf.Position
		FROM FlashcardSets fs
		JOIN Quizzes q ON q.FlashcardSetId = fs.Id
		JOIN QuizzesUnknownFlashcard quf ON quf.QuizId = q.Id
		WHERE fs.Id = $1
	`

	quiz := domain.NewQuiz(entity.Id, entity.Flashcards)

	unknown := []int{}
	rows, err := r.db.QueryContext(ctx, query, entity.Id)

	for rows.Next() {
		var u int
		err := rows.Scan(
			&quiz.Id,
			&quiz.CurrentlySelectedIndex,
			&u,
		)

		if err != nil {
			return nil, err
		}

		unknown = append(unknown, u)
	}

	quiz.Unknown = unknown

	return quiz, err
}

func (r *FlashcardSetRepository) List(ctx context.Context, offset int, limit int) ([]*domain.FlashcardSet, error) {
	query := `
		SELECT fs.Id, fs.Name, fs.Description, fs.LastAccessed, fs.TrackProgress, fs.Front, fs.Shuffle, fs.ShuffleSeed, COUNT(fc.Id)
		FROM FlashcardSets fs
		LEFT JOIN Flashcards fc ON fc.FlashcardSetId = fs.Id
		GROUP BY fs.Id
		ORDER BY fs.Id ASC
		LIMIT $1 OFFSET $2
	`

	rows, err := r.db.QueryContext(ctx, query, limit, offset)
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
			&fs.FlashcardCount,
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

func (r *FlashcardSetRepository) Update(ctx context.Context, entity *domain.FlashcardSet) error {
	return db.WithTransaction(ctx, r.db, func(tx *sql.Tx) error {
		stmt, err := tx.PrepareContext(ctx, `
			UPDATE FlashcardSets
			SET
				Name = $1,
				Description = $2,
				LastAccessed = $3,
				TrackProgress = $4,
				Front = $5,
				Shuffle = $6,
				ShuffleSeed = $7
			WHERE Id = $8
		`)

		if err != nil {
			return fmt.Errorf("Error while preparing statement for flashcard set update %w", err)
		}

		result, err := stmt.ExecContext(ctx,
			entity.Name,
			entity.Description,
			entity.LastAccessed,
			entity.TrackProgress,
			entity.Front,
			entity.Shuffle,
			entity.ShuffleSeed,
			entity.Id,
		)

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

		stmt.Close()

		stmt, err = tx.PrepareContext(ctx, `
			INSERT INTO Flashcards (Question, Answer, Position, FlashcardSetId) VALUES(?, ?, ?, ?)
			ON CONFLICT (FlashcardSetId, Position)
			DO UPDATE SET Question = ?, Answer = ?, Position = ?
		`)

		if err != nil {
			return fmt.Errorf("Error while preparing statement for flashcard set flashcard upsert %w", err)
		}

		for i, flashcard := range entity.Flashcards {
			result, err = stmt.ExecContext(ctx,
				flashcard.Question,
				flashcard.Answer,
				flashcard.Position,
				entity.Id,
				flashcard.Question,
				flashcard.Answer,
				flashcard.Position,
			)

			if err != nil {
				return fmt.Errorf("Error while updating flashcard in flashcard set %v", err)
			}

			flashcardId, _ := result.LastInsertId()

			if entity.Flashcards[i].Id == 0 {
				entity.Flashcards[i].Id = int(flashcardId)
			}
		}

		return nil
	})
}
