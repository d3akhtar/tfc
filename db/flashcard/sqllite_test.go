package flashcard_test

import (
	"context"
	"database/sql"
	"fmt"
	"testing"
	"time"

	"github.com/d3akhtar/tfc/db"
	"github.com/d3akhtar/tfc/db/flashcard"
	"github.com/d3akhtar/tfc/domain"
)

func testDB(t *testing.T) *sql.DB {
	t.Helper()

	db, err := sql.Open("sqlite3", ":memory:?_foreign_keys=ON")
	if err != nil {
		t.Fatalf("Failed to open test database: %v", err)
	}

	if err := initTestSchema(db); err != nil {
		db.Close()
		t.Fatalf("Failed to initialize schema: %v", err)
	}

	_, err = db.Exec(
		`INSERT INTO FlashcardSets
		(Name, Description, LastAccessed, TrackProgress, Shuffle, ShuffleSeed)
		VALUES ('thing', 'desc', CURRENT_TIMESTAMP, 0, 0, 32)`,
	)

	if err != nil {
		db.Close()
		t.Fatalf("Unexpected test db setup error %v", err)
	}

	t.Cleanup(func() {
		db.Close()
	})

	return db
}

func initTestSchema(database *sql.DB) error {
	_, err := database.Exec(db.Schema())
	return err
}

func TestFlashcardRepository_Create(t *testing.T) {
	ctx := context.Background()
	database := testDB(t)
	repo := flashcard.NewFlashcardRepository(database)

	flashcard := &domain.Flashcard{
		Question:       "q1",
		Answer:         "a1",
		FlashcardSetId: 1,
		Position:       0,
	}

	err := repo.Create(ctx, flashcard)
	if err != nil {
		t.Fatalf("Unexpected error %v", err)
	}

	if flashcard.Id == 0 {
		t.Fatal("Expected flashcard id to be set")
	}

	duplicate := &domain.Flashcard{
		Question:       "q2",
		Answer:         "a2",
		FlashcardSetId: 1,
		Position:       0,
	}

	err = repo.Create(ctx, duplicate)

	if err != db.ErrDuplicate {
		t.Fatalf("expected=%v, got=%v", db.ErrDuplicate, err)
	}
}

func TestFlashcardRepository_GetById(t *testing.T) {
	ctx := context.Background()
	database := testDB(t)
	repo := flashcard.NewFlashcardRepository(database)

	flashcard := &domain.Flashcard{
		Question:       "q1",
		Answer:         "a1",
		FlashcardSetId: 1,
		Position:       0,
	}

	err := repo.Create(ctx, flashcard)
	if err != nil {
		t.Fatalf("Unexpected error %v", err)
	}

	retrieved, err := repo.GetById(ctx, flashcard.Id)
	if err != nil {
		t.Fatalf("Unexpected error %v", err)
	}

	if *retrieved != *flashcard {
		t.Fatalf("expected=%v, actual=%v", *flashcard, *retrieved)
	}
}

func TestFlashcardRepository_Update(t *testing.T) {
	ctx := context.Background()
	database := testDB(t)
	repo := flashcard.NewFlashcardRepository(database)

	flashcard := &domain.Flashcard{
		Question:       "q1",
		Answer:         "a1",
		FlashcardSetId: 1,
		Position:       0,
	}

	err := repo.Create(ctx, flashcard)
	if err != nil {
		t.Fatalf("Unexpected error %v", err)
	}

	flashcard.Question = "q2"
	flashcard.Answer = "a2"
	flashcard.Position = 1

	err = repo.Update(ctx, flashcard)
	if err != nil {
		t.Fatalf("Unexpected error %v", err)
	}

	retrieved, err := repo.GetById(ctx, flashcard.Id)
	if err != nil {
		t.Fatalf("Unexpected error %v", err)
	}

	if *retrieved != *flashcard {
		t.Fatalf("expected=%v, actual=%v", *flashcard, *retrieved)
	}

	flashcard.Id = 2
	err = repo.Update(ctx, flashcard)

	if err != db.ErrNotFound {
		t.Fatalf("expected=%v, got=%v", db.ErrNotFound, err)
	}
}

func TestFlashcardRepository_Delete(t *testing.T) {
	ctx := context.Background()
	database := testDB(t)
	repo := flashcard.NewFlashcardRepository(database)

	flashcard := &domain.Flashcard{
		Question:       "q1",
		Answer:         "a1",
		FlashcardSetId: 1,
		Position:       0,
	}

	err := repo.Create(ctx, flashcard)
	if err != nil {
		t.Fatalf("Unexpected error %v", err)
	}

	err = repo.Delete(ctx, flashcard.Id)
	if err != nil {
		t.Fatalf("Unexpected error %v", err)
	}

	_, err = repo.GetById(ctx, flashcard.Id)

	if err != db.ErrNotFound {
		t.Fatalf("expected=%v, got=%v", db.ErrNotFound, err)
	}
}

func TestFlashcardRepository_List(t *testing.T) {
	ctx := context.Background()
	database := testDB(t)
	repo := flashcard.NewFlashcardRepository(database)

	expected := []*domain.Flashcard{}

	for i := range 10 {
		flashcard := &domain.Flashcard{
			Question:       fmt.Sprintf("q%d", i),
			Answer:         fmt.Sprintf("a%d", i),
			FlashcardSetId: 1,
			Position:       i,
		}

		expected = append(expected, flashcard)

		err := repo.Create(ctx, flashcard)
		if err != nil {
			t.Fatalf("Unexpected error while inserting flashcard %v: %v", *flashcard, err)
		}
	}

	actual, err := repo.List(ctx, 0, 10)

	if err != nil {
		t.Fatalf("Unexpected error %v", err)
	}

	for i, a := range actual {
		if *a != *expected[i] {
			t.Errorf("expected=%v, got=%v", *expected[i], *a)
		}
	}

	actual, err = repo.List(ctx, 20, 10)
	if len(actual) > 0 {
		t.Fatalf("len(actual) expected=%v, got=%v", 0, len(actual))
	}

	actual, err = repo.List(ctx, 4, 5)
	if err != nil {
		t.Fatalf("Unexpected error %v", err)
	}

	if len(actual) != 5 {
		t.Fatalf("len(actual) expected=%v, got=%v", 5, len(actual))
	}

	for i, a := range actual {
		if *a != *expected[i+4] {
			t.Errorf("expected=%v, got=%v", *expected[i], *a)
		}
	}
}

func TestFlashcardRepository_Count(t *testing.T) {
	ctx := context.Background()
	database := testDB(t)
	repo := flashcard.NewFlashcardRepository(database)

	flashcards := []*domain.Flashcard{}

	for i := range 10 {
		flashcard := &domain.Flashcard{
			Question:       fmt.Sprintf("q%d", i),
			Answer:         fmt.Sprintf("a%d", i),
			FlashcardSetId: 1,
			Position:       i,
		}

		flashcards = append(flashcards, flashcard)

		err := repo.Create(ctx, flashcard)
		if err != nil {
			t.Fatalf("Unexpected error while inserting flashcard %v: %v", *flashcard, err)
		}
	}

	count, err := repo.Count(ctx)

	if err != nil {
		t.Fatalf("Unexpected error %v", err)
	}

	if count != int64(len(flashcards)) {
		t.Fatalf("count expected=%v, got=%v", len(flashcards), count)
	}
}

func TestFlashcardRepository_GetFlashcardSetForFlashcard(t *testing.T) {
	ctx := context.Background()
	database := testDB(t)
	repo := flashcard.NewFlashcardRepository(database)

	flashcard := &domain.Flashcard{
		Question:       "q1",
		Answer:         "a1",
		FlashcardSetId: 1,
		Position:       0,
	}

	err := repo.Create(ctx, flashcard)
	if err != nil {
		t.Fatalf("Unexpected error %v", err)
	}

	flashcardSet, err := repo.GetFlashcardSetForFlashcard(ctx, flashcard)
	if err != nil {
		t.Fatalf("Unexpected error %v", err)
	}

	lastAccessed, _ := time.Parse(time.DateTime, "2026-08-22 14:30:00")

	expected := domain.FlashcardSet{
		Id:            1,
		Name:          "thing",
		Description:   "desc",
		LastAccessed:  lastAccessed,
		TrackProgress: false,
		Shuffle:       false,
		ShuffleSeed:   32,
	}

	if expected.Id != flashcardSet.Id ||
		expected.Name != flashcardSet.Name ||
		expected.Description != flashcardSet.Description ||
		expected.TrackProgress != flashcardSet.TrackProgress ||
		expected.Shuffle != flashcardSet.Shuffle ||
		expected.ShuffleSeed != flashcardSet.ShuffleSeed {
		t.Fatalf("expected=%#v, got=%#v", expected, flashcardSet)
	}
}

func TestFlashcardRepository_UpdateFlashcardPosition(t *testing.T) {
	ctx := context.Background()
	database := testDB(t)
	repo := flashcard.NewFlashcardRepository(database)

	flashcard := &domain.Flashcard{
		Question:       "q1",
		Answer:         "a1",
		FlashcardSetId: 1,
		Position:       0,
	}

	err := repo.Create(ctx, flashcard)
	if err != nil {
		t.Fatalf("Unexpected error %v", err)
	}

	err = repo.UpdateFlashcardPosition(ctx, flashcard, 1)
	if err != nil {
		t.Fatalf("Unexpected error %v", err)
	}

	retrieved, err := repo.GetById(ctx, flashcard.Id)
	if err != nil {
		t.Fatalf("Unexpected error %v", err)
	}

	if retrieved.Position != 1 {
		t.Fatalf("expected=%v, actual=%v", 1, retrieved.Position)
	}

	flashcard.Id = 2

	err = repo.UpdateFlashcardPosition(ctx, flashcard, 2)

	if err != db.ErrNotFound {
		t.Fatalf("expected=%v, got=%v", db.ErrNotFound, err)
	}
}
