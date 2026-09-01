package flashcard_set_test

import (
	"context"
	"database/sql"
	"fmt"
	"slices"
	"testing"
	"time"

	"github.com/d3akhtar/tfc/db"
	"github.com/d3akhtar/tfc/db/flashcard"
	"github.com/d3akhtar/tfc/db/flashcard_set"
	"github.com/d3akhtar/tfc/db/utils/test"
	"github.com/d3akhtar/tfc/domain"
)

func testDB(t *testing.T, init bool) *sql.DB {
	t.Helper()

	db, err := sql.Open("sqlite3", ":memory:?_foreign_keys=ON")
	if err != nil {
		t.Fatalf("Failed to open test database: %v", err)
	}

	if err := initTestSchema(db); err != nil {
		db.Close()
		t.Fatalf("Failed to initialize schema: %v", err)
	}

	if init {
		_, err = db.Exec(
			`
				INSERT INTO FlashcardSets (Name, Description, LastAccessed, TrackProgress, Shuffle)
				VALUES ('thing', 'desc', CURRENT_TIMESTAMP, 0, 0);

				INSERT INTO Flashcards (Question, Answer, FlashcardSetId, Position) VALUES ('q1', 'a1', 1, 0);
				INSERT INTO Flashcards (Question, Answer, FlashcardSetId, Position) VALUES ('q2', 'a2', 1, 1);
				INSERT INTO Flashcards (Question, Answer, FlashcardSetId, Position) VALUES ('q3', 'a3', 1, 2);
				INSERT INTO Flashcards (Question, Answer, FlashcardSetId, Position) VALUES ('q4', 'a4', 1, 3);
				INSERT INTO Flashcards (Question, Answer, FlashcardSetId, Position) VALUES ('q5', 'a5', 1, 4);

				INSERT INTO Quizzes (FlashcardSetId) VALUES (1);

				INSERT INTO QuizzesUnknownFlashcard (QuizId, FlashcardId, Position) VALUES (1, 2, 1);
				INSERT INTO QuizzesUnknownFlashcard (QuizId, FlashcardId, Position) VALUES (1, 3, 2);
			`,
		)

		if err != nil {
			db.Close()
			t.Fatalf("Unexpected test setup error %v", err)
		}
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

func TestFlashcardSetRepository_Count(t *testing.T) {
	ctx := context.Background()
	database := testDB(t, false)
	repo := flashcard_set.NewFlashcardSetRepository(database)

	for i := range 10 {
		flashcardSet := &domain.FlashcardSet{
			Id:            i,
			Name:          fmt.Sprintf("name%d", i),
			Description:   fmt.Sprintf("desc%d", i),
			LastAccessed:  time.Now(),
			TrackProgress: i%2 == 0,
			Front:         domain.FlashcardFront((i + 1) % 2),
			Shuffle:       i%2 == 0,
			ShuffleSeed:   i + 20,
		}

		err := repo.Create(ctx, flashcardSet)
		if err != nil {
			t.Fatalf("Unexpected error %v", err)
		}
	}

	count, err := repo.Count(ctx)

	if err != nil {
		t.Fatalf("Unexpected error %v", err)
	}

	if count != 10 {
		t.Fatalf("count expected=%v, got=%v", 10, count)
	}
}

func TestFlashcardSetRepository_Create(t *testing.T) {
	ctx := context.Background()
	database := testDB(t, false)
	repo := flashcard_set.NewFlashcardSetRepository(database)

	flashcards := []domain.Flashcard{}

	for i := range 10 {
		flashcard := domain.Flashcard{
			Question:       fmt.Sprintf("q%d", i),
			Answer:         fmt.Sprintf("a%d", i),
			FlashcardSetId: 1,
			Position:       i,
		}

		flashcards = append(flashcards, flashcard)
	}

	flashcardSet := &domain.FlashcardSet{
		Name:          "set 1",
		Description:   "desc 1",
		LastAccessed:  time.Now(),
		TrackProgress: false,
		Front:         domain.FlashcardFront(domain.Question),
		Flashcards:    flashcards,
		Shuffle:       false,
		ShuffleSeed:   0,
	}

	err := repo.Create(ctx, flashcardSet)
	if err != nil {
		t.Fatalf("Unexpected error %v", err)
	}

	if flashcardSet.Id == 0 {
		t.Fatal("Expected flashcard set id to be set")
	}

	for i, flashcard := range flashcardSet.Flashcards {
		if flashcard.Id == 0 {
			t.Errorf("(%d) Expected flashcard id to be set", i)
		}
	}
}

func TestFlashcardSetRepository_Delete(t *testing.T) {
	ctx := context.Background()
	database := testDB(t, false)
	repo := flashcard_set.NewFlashcardSetRepository(database)
	flashcardRepo := flashcard.NewFlashcardRepository(database)

	flashcards := []domain.Flashcard{}

	for i := range 10 {
		flashcard := domain.Flashcard{
			Question:       fmt.Sprintf("q%d", i),
			Answer:         fmt.Sprintf("a%d", i),
			FlashcardSetId: 1,
			Position:       i,
		}

		flashcards = append(flashcards, flashcard)
	}

	accessedTime := time.Now()

	flashcardSet := &domain.FlashcardSet{
		Name:          "set 1",
		Description:   "desc 1",
		LastAccessed:  accessedTime,
		TrackProgress: false,
		Front:         domain.FlashcardFront(domain.Question),
		Flashcards:    flashcards,
		Shuffle:       false,
		ShuffleSeed:   0,
	}

	err := repo.Create(ctx, flashcardSet)
	if err != nil {
		t.Fatalf("Unexpected error %v", err)
	}

	err = repo.Delete(ctx, flashcardSet.Id)
	if err != nil {
		t.Fatalf("Unexpected error %v", err)
	}

	_, err = repo.GetById(ctx, flashcardSet.Id)

	if err != db.ErrNotFound {
		t.Fatalf("expected=%v, got=%v", db.ErrNotFound, err)
	}

	for _, flashcard := range flashcardSet.Flashcards {
		_, err = flashcardRepo.GetById(ctx, flashcard.Id)

		if err != db.ErrNotFound {
			t.Fatalf("expected=%v, got=%v", db.ErrNotFound, err)
		}
	}
}

func TestFlashcardSetRepository_FilterFlashcardSets(t *testing.T) {
	ctx := context.Background()
	database := testDB(t, false)
	repo := flashcard_set.NewFlashcardSetRepository(database)

	flashcardSets := []*domain.FlashcardSet{}
	for i := range 10 {
		flashcardSet := &domain.FlashcardSet{
			Id:            i,
			Name:          fmt.Sprintf("name%d", i),
			Description:   fmt.Sprintf("desc%d", i),
			LastAccessed:  time.Now(),
			TrackProgress: i%2 == 0,
			Front:         domain.FlashcardFront((i + 1) % 2),
			Shuffle:       i%2 == 0,
			ShuffleSeed:   i + 20,
		}

		err := repo.Create(ctx, flashcardSet)
		if err != nil {
			t.Fatalf("Unexpected error %v", err)
		}

		flashcardSets = append(flashcardSets, flashcardSet)
	}

	filtered, err := repo.FilterFlashcardSets(ctx, "name", 10, 0)
	if err != nil {
		t.Fatalf("Unexpected error %v", err)
	}

	for i, got := range filtered {
		err = test_utils.TestCmpFlashcardSet(flashcardSets[i], &got)
		if err != nil {
			t.Errorf("(%d) %v", i, err)
		}
	}

	filtered, err = repo.FilterFlashcardSets(ctx, "name0", 10, 0)
	if len(filtered) != 1 {
		t.Fatalf("len(filtered) expected=%v, got=%v", 1, len(filtered))
	}

	err = test_utils.TestCmpFlashcardSet(flashcardSets[0], &filtered[0])
	if err != nil {
		t.Errorf("%v", err)
	}

	filtered, err = repo.FilterFlashcardSets(ctx, "name", 5, 2)
	if err != nil {
		t.Fatalf("Unexpected error %v", err)
	}

	for i, got := range filtered {
		err = test_utils.TestCmpFlashcardSet(flashcardSets[i+2], &got)
		if err != nil {
			t.Errorf("(%d) %v", i, err)
		}
	}
}

func TestFlashcardSetRepository_GetAllFlashcardsForSet(t *testing.T) {
	ctx := context.Background()
	database := testDB(t, false)
	repo := flashcard_set.NewFlashcardSetRepository(database)

	flashcards := []domain.Flashcard{}

	for i := range 10 {
		flashcard := domain.Flashcard{
			Question:       fmt.Sprintf("q%d", i),
			Answer:         fmt.Sprintf("a%d", i),
			FlashcardSetId: 1,
			Position:       i,
		}

		flashcards = append(flashcards, flashcard)
	}

	accessedTime := time.Now()

	flashcardSet := &domain.FlashcardSet{
		Name:          "set 1",
		Description:   "desc 1",
		LastAccessed:  accessedTime,
		TrackProgress: false,
		Front:         domain.FlashcardFront(domain.Question),
		Flashcards:    flashcards,
		Shuffle:       false,
		ShuffleSeed:   0,
	}

	err := repo.Create(ctx, flashcardSet)
	if err != nil {
		t.Fatalf("Unexpected error %v", err)
	}

	got, err := repo.GetAllFlashcardsForSet(ctx, flashcardSet)
	if err != nil {
		t.Fatalf("Unexpected error %v", err)
	}

	slices.SortFunc(got, func(a, b domain.Flashcard) int {
		return a.Id - b.Id
	})

	for i, a := range got {
		if a != flashcards[i] {
			t.Errorf("(%d) expected=%v, got=%v", i, flashcards[i], a)
		}
	}
}

func TestFlashcardSetRepository_GetById(t *testing.T) {
	ctx := context.Background()
	database := testDB(t, false)
	repo := flashcard_set.NewFlashcardSetRepository(database)

	expected := &domain.FlashcardSet{
		Name:          "set 1",
		Description:   "desc 1",
		LastAccessed:  time.Now(),
		TrackProgress: false,
		Front:         domain.FlashcardFront(domain.Question),
		Shuffle:       false,
		ShuffleSeed:   0,
	}

	err := repo.Create(ctx, expected)
	if err != nil {
		t.Fatalf("Unexpected error %v", err)
	}

	got, err := repo.GetById(ctx, expected.Id)
	if err != nil {
		t.Fatalf("Unexpected error %v", err)
	}

	err = test_utils.TestCmpFlashcardSet(expected, got)
	if err != nil {
		t.Errorf("%v", err)
	}
}

func TestFlashcardSetRepository_GetQuizForFlashcardSet(t *testing.T) {
	ctx := context.Background()
	database := testDB(t, true)
	repo := flashcard_set.NewFlashcardSetRepository(database)

	flashcardSet := &domain.FlashcardSet{Id: 1, Flashcards: []domain.Flashcard{}}
	flashcardSet.AddFlashcard("q1", "a1")
	flashcardSet.AddFlashcard("q2", "a2")
	flashcardSet.AddFlashcard("q3", "a3")
	flashcardSet.AddFlashcard("q4", "a4")

	quiz, err := repo.GetQuizForFlashcardSet(ctx, flashcardSet)
	if err != nil {
		t.Fatalf("Unexpected error %v", err)
	}

	if quiz.Id != 1 {
		t.Fatalf("quiz Id expected=%v, got=%v", 1, quiz.Id)
	}

	if !slices.Equal(quiz.Unknown, []int{1, 2}) {
		t.Fatalf("quiz Unknown expected=%v, got=%v", []int{1, 2}, quiz.Unknown)
	}
}

func TestFlashcardSetRepository_List(t *testing.T) {
	ctx := context.Background()
	database := testDB(t, false)
	repo := flashcard_set.NewFlashcardSetRepository(database)

	flashcardSets := []*domain.FlashcardSet{}
	for i := range 10 {
		flashcardSet := &domain.FlashcardSet{
			Id:            i,
			Name:          fmt.Sprintf("name%d", i),
			Description:   fmt.Sprintf("desc%d", i),
			LastAccessed:  time.Now(),
			TrackProgress: i%2 == 0,
			Front:         domain.FlashcardFront((i + 1) % 2),
			Shuffle:       i%2 == 0,
			ShuffleSeed:   i + 20,
		}

		err := repo.Create(ctx, flashcardSet)
		if err != nil {
			t.Fatalf("Unexpected error %v", err)
		}

		flashcardSets = append(flashcardSets, flashcardSet)
	}

	actual, err := repo.List(ctx, 0, 10)
	if err != nil {
		t.Fatalf("Unexpected error %v", err)
	}

	slices.SortFunc(actual, func(a, b *domain.FlashcardSet) int {
		return a.Id - b.Id
	})

	for i, a := range actual {
		err = test_utils.TestCmpFlashcardSet(flashcardSets[i], a)
		if err != nil {
			t.Errorf("(%d) %v", i, err)
		}
	}

	actual, err = repo.List(ctx, 2, 4)
	if err != nil {
		t.Fatalf("Unexpected error %v", err)
	}

	slices.SortFunc(actual, func(a, b *domain.FlashcardSet) int {
		return a.Id - b.Id
	})

	for i := range 4 {
		err = test_utils.TestCmpFlashcardSet(flashcardSets[i+2], actual[i])
		if err != nil {
			t.Errorf("(%d) %v", i, err)
		}
	}
}

func TestFlashcardSetRepository_Update(t *testing.T) {
	ctx := context.Background()
	database := testDB(t, true)
	repo := flashcard_set.NewFlashcardSetRepository(database)

	flashcardSet, err := repo.GetById(ctx, 1)
	if err != nil {
		t.Fatalf("Unexpected error %v", err)
	}

	flashcards, err := repo.GetAllFlashcardsForSet(ctx, flashcardSet)
	if err != nil {
		t.Fatalf("Unexpected error %v", err)
	}

	flashcardSet.Flashcards = flashcards

	flashcardSet.Flashcards[0].Question = "new q1"
	flashcardSet.Flashcards[0].Answer = "new a1"

	flashcardSet.AddFlashcard("q6", "a6")
	flashcardSet.AddFlashcard("q7", "a7")

	err = repo.Update(ctx, flashcardSet)
	if err != nil {
		t.Fatalf("Unexpected error %v", err)
	}

	updated, err := repo.GetById(ctx, flashcardSet.Id)
	if err != nil {
		t.Fatalf("Unexpected error %v", err)
	}

	for i, updatedFlashcard := range updated.Flashcards {
		if updatedFlashcard != flashcardSet.Flashcards[i] {
			t.Errorf("(%d) expected=%v, got=%v", i, flashcardSet.Flashcards[i], updatedFlashcard)
		}
	}

	updated.Id = 5
	err = repo.Update(ctx, updated)

	if err != db.ErrNotFound {
		t.Fatalf("expected=%v, got=%v", db.ErrNotFound, err)
	}
}
