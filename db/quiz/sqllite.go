package quiz

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/d3akhtar/tfc/db"
	"github.com/d3akhtar/tfc/db/utils"
	"github.com/d3akhtar/tfc/domain"
	"github.com/mattn/go-sqlite3"
)

type QuizRepository struct {
	db *sql.DB
}

// Ensure interface is implemented at compile-time
var _ QuizRepo = (*QuizRepository)(nil)

func NewQuizRepository(db *sql.DB) *QuizRepository {
	return &QuizRepository{db}
}

func (r *QuizRepository) Count(ctx context.Context) (int64, error) {
	return utils.CountQuery(r.db, ctx, "Quizzes")
}

func (r *QuizRepository) Create(ctx context.Context, entity *domain.Quiz) error {
	return db.WithTransaction(ctx, r.db, func(tx *sql.Tx) error {
		stmt, err := tx.PrepareContext(ctx, `INSERT INTO Quizzes (FlashcardSetId, CurrentlySelectedIndex) VALUES (?, ?)`)

		if err != nil {
			return fmt.Errorf("Error while preparing statement for inserting quiz %w", err)
		}

		result, err := stmt.ExecContext(ctx, entity.FlashcardSetId, entity.CurrentlySelectedIndex)

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

		quizId, _ := result.LastInsertId()

		entity.Id = int(quizId)

		stmt, err = tx.PrepareContext(ctx, `INSERT INTO QuizFlashcards (QuizId, FlashcardId, Position, IsUnknown) VALUES (?, ?, ?, ?)`)

		if err != nil {
			return fmt.Errorf("Error while preparing statement for inserting quiz %w", err)
		}

		for i, flashcard := range entity.Flashcards {
			isUnknown := entity.Unknown.Contains(i)
			result, err = stmt.ExecContext(ctx, entity.Id, flashcard.Id, i, isUnknown)

			if err != nil {
				return fmt.Errorf("Error while adding unknown quiz flashcard %w", err)
			}
		}

		return nil
	})
}

func (r *QuizRepository) Delete(ctx context.Context, id int) error {
	return utils.DeleteQuery(r.db, ctx, "Quizzes", id)
}

func (r *QuizRepository) GetById(ctx context.Context, id int) (*domain.Quiz, error) {
	query := `
		SELECT
			q.Id,
			q.CurrentlySelectedIndex,
			fc.Id,
			fc.Question,
			fc.Answer,
			fc.FlashcardSetId,
			fc.Position,
			qf.Position,
			qf.IsUnknown
		FROM Quizzes q
		JOIN FlashcardSets fs ON fs.Id = q.FlashcardSetId
		JOIN Flashcards fc ON fc.FlashcardSetId = fs.Id
		JOIN QuizFlashcards qf ON qf.QuizId = q.Id AND qf.FlashcardId = fc.Id
		WHERE q.Id = $1
		ORDER BY qf.Position ASC
	`

	rows, err := r.db.QueryContext(ctx, query, id)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var quizId, quizCurrentlySelected int
	var isUnknown bool
	var pos *int
	flashcards := []domain.Flashcard{}
	unknown := []int{}
	for rows.Next() {
		fc := domain.Flashcard{}
		err := rows.Scan(
			&quizId,
			&quizCurrentlySelected,
			&fc.Id,
			&fc.Question,
			&fc.Answer,
			&fc.FlashcardSetId,
			&fc.Position,
			&pos,
			&isUnknown,
		)

		if err != nil {
			return nil, err
		}

		flashcards = append(flashcards, fc)

		if isUnknown {
			unknown = append(unknown, *pos)
		}
	}

	quiz := domain.NewQuiz(flashcards[0].FlashcardSetId, flashcards)
	quiz.Id = quizId
	quiz.CurrentlySelectedIndex = quizCurrentlySelected
	quiz.SetUnknownCardPositions(unknown)

	return quiz, nil
}

func (r *QuizRepository) GetFlashcardSetForQuiz(ctx context.Context, entity *domain.Quiz) (*domain.FlashcardSet, error) {
	query := `
		SELECT
			fs.Id,
			fs.Name,
			fs.Description,
			fs.LastAccessed,
			fs.TrackProgress,
			fs.Front,
			fs.Shuffle,
			fs.ShuffleSeed,
			fc.Id,
			fc.Question,
			fc.Answer,
			fc.FlashcardSetId,
			fc.Position
		FROM Quizzes q
		JOIN FlashcardSets fs ON fs.Id = q.FlashcardSetId
		JOIN Flashcards fc ON fc.FlashcardSetId = fs.Id
		WHERE q.Id = $1
		ORDER BY fs.Id ASC, fc.Id ASC
	`

	rows, err := r.db.QueryContext(ctx, query, entity.Id)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	flashcardSet := &domain.FlashcardSet{}
	flashcards := []domain.Flashcard{}
	for rows.Next() {
		fc := domain.Flashcard{}
		err := rows.Scan(
			&flashcardSet.Id,
			&flashcardSet.Name,
			&flashcardSet.Description,
			&flashcardSet.LastAccessed,
			&flashcardSet.TrackProgress,
			&flashcardSet.Front,
			&flashcardSet.Shuffle,
			&flashcardSet.ShuffleSeed,
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

	flashcardSet.Flashcards = flashcards

	return flashcardSet, nil
}

func (r *QuizRepository) GetUnknownFlashcardsForQuiz(ctx context.Context, entity *domain.Quiz) ([]domain.Flashcard, error) {
	query := `
		SELECT
			fc.Id,
			fc.Question,
			fc.Answer,
			fc.FlashcardSetId,
			fc.Position,
			qf.Position,
			qf.IsUnknown
		FROM Quizzes q
		JOIN FlashcardSets fs ON fs.Id = q.FlashcardSetId
		JOIN Flashcards fc ON fc.FlashcardSetId = fs.Id
		JOIN QuizFlashcards qf ON qf.QuizId = q.Id AND qf.FlashcardId = fc.Id
		WHERE q.Id = $1 AND qf.IsUnknown = 1
		ORDER BY qf.Position ASC
	`

	rows, err := r.db.QueryContext(ctx, query, entity.Id)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var isUnknown bool
	var pos int
	flashcards := []domain.Flashcard{}
	unknown := []int{}
	for rows.Next() {
		fc := domain.Flashcard{}

		err := rows.Scan(
			&fc.Id,
			&fc.Question,
			&fc.Answer,
			&fc.FlashcardSetId,
			&fc.Position,
			&pos,
			&isUnknown,
		)

		if err != nil {
			return nil, err
		}

		if isUnknown {
			flashcards = append(flashcards, fc)
			unknown = append(unknown, pos)
		}
	}

	entity.SetUnknownCardPositions(unknown)

	return flashcards, nil
}

func (r *QuizRepository) List(ctx context.Context, offset int, limit int) ([]*domain.Quiz, error) {
	query := `
		SELECT q.Id, q.FlashcardSetId, q.CurrentlySelectedIndex
		FROM Quizzes q
		INNER JOIN FlashcardSets fs ON fs.Id = q.FlashcardSetId
		ORDER BY q.Id ASC
		LIMIT $1 OFFSET $2
	`

	rows, err := r.db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	quizzes := []*domain.Quiz{}
	for rows.Next() {
		quiz := &domain.Quiz{}
		err := rows.Scan(
			&quiz.Id,
			&quiz.FlashcardSetId,
			&quiz.CurrentlySelectedIndex,
		)

		if err != nil {
			return nil, err
		}

		quizzes = append(quizzes, quiz)
	}

	return quizzes, nil
}

func (r *QuizRepository) Update(ctx context.Context, entity *domain.Quiz) error {
	return db.WithTransaction(ctx, r.db, func(tx *sql.Tx) error {
		stmt, err := tx.PrepareContext(ctx, `
			UPDATE Quizzes
			SET CurrentlySelectedIndex = ?
			WHERE Id = ?
		`)

		if err != nil {
			return fmt.Errorf("Error while preparing statement for quiz update %w", err)
		}

		err = utils.ExecStmtUpdate(stmt, ctx, entity.CurrentlySelectedIndex, entity.Id)
		if err != nil {
			return err
		}

		stmt.Close()

		stmt, err = tx.PrepareContext(ctx, `
			INSERT INTO QuizFlashcards (QuizId, FlashcardId, Position, IsUnknown) VALUES ($1, $2, $3, $4)
			ON CONFLICT (QuizId, FlashcardId, Position)
			DO UPDATE SET Position = $3, IsUnknown = $4
		`)

		if err != nil {
			return fmt.Errorf("Error while preparing statement for quiz update %w", err)
		}

		defer stmt.Close()

		for i, flashcard := range entity.Flashcards {
			isUnknown := entity.Unknown.Contains(i)
			_, err := stmt.ExecContext(ctx,
				entity.Id,
				flashcard.Id,
				i,
				isUnknown,
			)

			if err != nil {
				return fmt.Errorf("Failed to add flashcard with Id %d for quiz with Id %d", entity.Flashcards[i].Id, entity.Id)
			}
		}

		return nil
	})
}

func (r *QuizRepository) CreateQuizForUnknownFlashcards(ctx context.Context, entity *domain.Quiz) error {
	return db.WithTransaction(ctx, r.db, func(tx *sql.Tx) error {
		stmt, err := tx.PrepareContext(ctx, `
			UPDATE Quizzes
			SET CurrentlySelectedIndex = ?
			WHERE Id = ?
		`)

		if err != nil {
			return fmt.Errorf("Error while preparing statement for quiz update %w", err)
		}

		err = utils.ExecStmtUpdate(stmt, ctx, entity.CurrentlySelectedIndex, entity.Id)
		if err != nil {
			return err
		}

		stmt.Close()

		stmt, err = tx.PrepareContext(ctx, `
			DELETE FROM QuizFlashcards WHERE QuizId = ?
		`)

		if err != nil {
			return fmt.Errorf("Error while preparing statement for quiz update %w", err)
		}

		err = utils.ExecStmtUpdate(stmt, ctx, entity.Id)
		if err != nil && err != db.ErrNotFound {
			return err
		}

		stmt.Close()

		stmt, err = tx.PrepareContext(ctx, `
			INSERT INTO QuizFlashcards (QuizId, FlashcardId, Position, IsUnknown) VALUES ($1, $2, $3, $4)
		`)

		if err != nil {
			return fmt.Errorf("Error while preparing statement for quiz update %w", err)
		}

		defer stmt.Close()

		for i, flashcard := range entity.Flashcards {
			isUnknown := entity.Unknown.Contains(i)
			_, err := stmt.ExecContext(ctx,
				entity.Id,
				flashcard.Id,
				i,
				isUnknown,
			)

			if err != nil {
				return fmt.Errorf("Failed to add flashcard with Id %d for quiz with Id %d", entity.Flashcards[i].Id, entity.Id)
			}
		}

		return nil
	})
}
