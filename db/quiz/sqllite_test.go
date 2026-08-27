package quiz_test

import (
	"context"
	"database/sql"
	"fmt"
	"slices"
	"testing"
	"time"

	"github.com/d3akhtar/tfc/db"
	"github.com/d3akhtar/tfc/db/flashcard_set"
	"github.com/d3akhtar/tfc/db/quiz"
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

	t.Cleanup(func() {
		db.Close()
	})

	return db
}

func initTestSchema(database *sql.DB) error {
	_, err := database.Exec(db.Schema())
	return err
}

func TestQuizRepository_Count(t *testing.T) {
	ctx := context.Background()
	database := testDB(t)
	repo := quiz.NewQuizRepository(database)
	flashcardSetRepo := flashcard_set.NewFlashcardSetRepository(database)

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

		err := flashcardSetRepo.Create(ctx, flashcardSet)
		if err != nil {
			t.Fatalf("Unexpected error %v", err)
		}

		flashcardSets = append(flashcardSets, flashcardSet)
	}

	quizzes := []*domain.Quiz{}

	for _, flashcardSet := range flashcardSets {
		quiz := domain.NewQuiz(flashcardSet.Id, flashcardSet.Flashcards)

		quizzes = append(quizzes, quiz)

		err := repo.Create(ctx, quiz)
		if err != nil {
			t.Fatalf("Unexpected error while inserting quiz %v: %v", quiz, err)
		}
	}

	count, err := repo.Count(ctx)

	if err != nil {
		t.Fatalf("Unexpected error %v", err)
	}

	if count != int64(len(quizzes)) {
		t.Fatalf("count expected=%v, got=%v", len(quizzes), count)
	}
}

func TestQuizRepository_Create(t *testing.T) {
	ctx := context.Background()
	database := testDB(t)
	repo := quiz.NewQuizRepository(database)
	flashcardSetRepo := flashcard_set.NewFlashcardSetRepository(database)

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

	err := flashcardSetRepo.Create(ctx, flashcardSet)
	if err != nil {
		t.Fatalf("Unexpected test setup error %v", err)
	}

	quiz := &domain.Quiz{
		FlashcardSetId:         flashcardSet.Id,
		Flashcards:             flashcards,
		Unknown:                []int{1, 2},
		CurrentlySelectedIndex: 3,
	}

	err = repo.Create(ctx, quiz)
	if err != nil {
		t.Fatalf("Unexpected error %v", err)
	}

	if quiz.Id == 0 {
		t.Fatalf("Expected quiz Id to be set")
	}
}

func TestQuizRepository_Delete(t *testing.T) {
	ctx := context.Background()
	database := testDB(t)
	repo := quiz.NewQuizRepository(database)
	flashcardSetRepo := flashcard_set.NewFlashcardSetRepository(database)

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

	err := flashcardSetRepo.Create(ctx, flashcardSet)
	if err != nil {
		t.Fatalf("Unexpected test setup error %v", err)
	}

	quiz := &domain.Quiz{
		FlashcardSetId: flashcardSet.Id,
		Flashcards:     flashcards,
	}

	err = repo.Create(ctx, quiz)
	if err != nil {
		t.Fatalf("Unexpected error %v", err)
	}

	err = repo.Delete(ctx, quiz.Id)
	if err != nil {
		t.Fatalf("Unexpected error %v", err)
	}

	err = repo.Delete(ctx, quiz.Id)
	if err != db.ErrNotFound {
		t.Fatalf("expected=%v, got=%v", db.ErrNotFound, err)
	}
}

func TestQuizRepository_GetById(t *testing.T) {
	ctx := context.Background()
	database := testDB(t)
	repo := quiz.NewQuizRepository(database)
	flashcardSetRepo := flashcard_set.NewFlashcardSetRepository(database)

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

	err := flashcardSetRepo.Create(ctx, flashcardSet)
	if err != nil {
		t.Fatalf("Unexpected test setup error %v", err)
	}

	quiz := &domain.Quiz{
		FlashcardSetId:         flashcardSet.Id,
		Flashcards:             flashcards,
		Unknown:                []int{1, 2},
		CurrentlySelectedIndex: 3,
	}

	err = repo.Create(ctx, quiz)
	if err != nil {
		t.Fatalf("Unexpected error %v", err)
	}

	got, err := repo.GetById(ctx, quiz.Id)
	if err != nil {
		t.Fatalf("Unexpected error %v", err)
	}

	err = testCmpQuiz(quiz, got)
	if err != nil {
		t.Fatalf("%v", err)
	}
}

func TestQuizRepository_GetFlashcardSetForQuiz(t *testing.T) {
	ctx := context.Background()
	database := testDB(t)
	repo := quiz.NewQuizRepository(database)
	flashcardSetRepo := flashcard_set.NewFlashcardSetRepository(database)

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

	err := flashcardSetRepo.Create(ctx, flashcardSet)
	if err != nil {
		t.Fatalf("Unexpected test setup error %v", err)
	}

	quiz := &domain.Quiz{
		FlashcardSetId:         flashcardSet.Id,
		Flashcards:             flashcards,
		Unknown:                []int{1, 2},
		CurrentlySelectedIndex: 3,
	}

	err = repo.Create(ctx, quiz)
	if err != nil {
		t.Fatalf("Unexpected error %v", err)
	}

	got, err := repo.GetFlashcardSetForQuiz(ctx, quiz)
	err = testCmpFlashcardSet(flashcardSet, got)
	if err != nil {
		t.Fatalf("%v", err)
	}

	for i, gotFlashcard := range got.Flashcards {
		if flashcardSet.Flashcards[i] != gotFlashcard {
			t.Errorf("(%d) expected=%v, got=%v", i, flashcardSet.Flashcards[i], gotFlashcard)
		}
	}
}

func TestQuizRepository_GetUnknownFlashcardsForQuiz(t *testing.T) {
	ctx := context.Background()
	database := testDB(t)
	repo := quiz.NewQuizRepository(database)
	flashcardSetRepo := flashcard_set.NewFlashcardSetRepository(database)

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

	err := flashcardSetRepo.Create(ctx, flashcardSet)
	if err != nil {
		t.Fatalf("Unexpected test setup error %v", err)
	}

	quiz := &domain.Quiz{
		FlashcardSetId:         flashcardSet.Id,
		Flashcards:             flashcards,
		Unknown:                []int{1, 2},
		CurrentlySelectedIndex: 3,
	}

	err = repo.Create(ctx, quiz)
	if err != nil {
		t.Fatalf("Unexpected error %v", err)
	}

	got, err := repo.GetUnknownFlashcardsForQuiz(ctx, quiz)
	if err != nil {
		t.Fatalf("Unexpected error %v", err)
	}

	expected := quiz.GetUnknownFlashcards()

	for i, flashcard := range expected {
		if got[i] != flashcard {
			t.Errorf("(%d) expected=%#v, got=%#v", i, flashcard, got[i])
		}
	}
}

func TestQuizRepository_List(t *testing.T) {
	ctx := context.Background()
	database := testDB(t)
	repo := quiz.NewQuizRepository(database)
	flashcardSetRepo := flashcard_set.NewFlashcardSetRepository(database)

	quizzes := []*domain.Quiz{}

	for i := range 10 {
		flashcardSet := &domain.FlashcardSet{
			Name:          fmt.Sprintf("set %d", i),
			Description:   fmt.Sprintf("desc %d", i),
			LastAccessed:  time.Now(),
			TrackProgress: false,
			Front:         domain.FlashcardFront(domain.Question),
			Shuffle:       false,
			ShuffleSeed:   0,
		}

		err := flashcardSetRepo.Create(ctx, flashcardSet)
		if err != nil {
			t.Fatalf("Unexpected error %v", err)
		}

		quiz := &domain.Quiz{
			FlashcardSetId:         flashcardSet.Id,
			CurrentlySelectedIndex: 0,
		}

		err = repo.Create(ctx, quiz)
		if err != nil {
			t.Fatalf("Unexpected error %v", err)
		}

		quizzes = append(quizzes, quiz)
	}

	got, err := repo.List(ctx, 0, 10)
	if err != nil {
		t.Fatalf("Unexpected error %v", err)
	}

	for i, expected := range quizzes {
		err = testCmpQuiz(expected, got[i])
		if err != nil {
			t.Errorf("(%d) %v", i, err)
		}
	}

	got, err = repo.List(ctx, 3, 5)
	if err != nil {
		t.Fatalf("Unexpected error %v", err)
	}

	for i, gotQuiz := range got {
		err = testCmpQuiz(quizzes[i+3], gotQuiz)
		if err != nil {
			t.Errorf("(%d) %v", i, err)
		}
	}
}

func TestQuizRepository_Update(t *testing.T) {
	ctx := context.Background()
	database := testDB(t)
	repo := quiz.NewQuizRepository(database)
	flashcardSetRepo := flashcard_set.NewFlashcardSetRepository(database)

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

	err := flashcardSetRepo.Create(ctx, flashcardSet)
	if err != nil {
		t.Fatalf("Unexpected test setup error %v", err)
	}

	quiz := &domain.Quiz{
		FlashcardSetId:         flashcardSet.Id,
		Flashcards:             flashcards,
		Unknown:                []int{1, 2},
		CurrentlySelectedIndex: 3,
	}

	err = repo.Create(ctx, quiz)
	if err != nil {
		t.Fatalf("Unexpected error %v", err)
	}

	quiz.CurrentlySelectedIndex = 6
	quiz.Unknown = append(quiz.Unknown, 4)

	err = repo.Update(ctx, quiz)
	if err != nil {
		t.Fatalf("Unexpected error %v", err)
	}

	got, err := repo.GetUnknownFlashcardsForQuiz(ctx, quiz)
	if err != nil {
		t.Fatalf("Unexpected error %v", err)
	}

	expected := quiz.GetUnknownFlashcards()

	for i, flashcard := range expected {
		if got[i] != flashcard {
			t.Errorf("(%d) expected=%#v, got=%#v", i, flashcard, got[i])
		}
	}
}

func testCmpQuiz(expected, actual *domain.Quiz) error {
	if expected.Id != actual.Id {
		return fmt.Errorf("quiz Id expected=%v, got=%v", expected.Id, actual.Id)
	}

	if expected.FlashcardSetId != actual.FlashcardSetId {
		return fmt.Errorf("quiz FlashcardSetId expected=%v, got=%v", expected.FlashcardSetId, actual.FlashcardSetId)
	}

	if !slices.Equal(expected.Flashcards, actual.Flashcards) {
		return fmt.Errorf("quiz Flashcards expected=%v, got=%v", expected.Flashcards, actual.Flashcards)
	}

	if expected.CurrentlySelectedIndex != actual.CurrentlySelectedIndex {
		return fmt.Errorf("quiz CurrentlySelectedIndex expected=%v, got=%v", expected.CurrentlySelectedIndex, actual.CurrentlySelectedIndex)
	}

	if !slices.Equal(expected.Unknown, actual.Unknown) {
		return fmt.Errorf("quiz.Unknown expected=%v, got=%v", expected.Unknown, actual.Unknown)
	}

	return nil
}

func testCmpFlashcardSet(expected, actual *domain.FlashcardSet) error {
	if expected.Id != actual.Id {
		return fmt.Errorf("flashcard set Id expected=%v, got=%v", expected.Id, actual.Id)
	}

	if expected.Name != actual.Name {
		return fmt.Errorf("flashcard set Name expected=%v, got=%v", expected.Name, actual.Name)
	}

	if expected.Description != actual.Description {
		return fmt.Errorf("flashcard set Description expected=%v, got=%v", expected.Description, actual.Description)
	}

	if expected.LastAccessed.UTC() != actual.LastAccessed.UTC() {
		return fmt.Errorf("flashcard set LastAccessedexpected=%v, got=%v", expected.LastAccessed.UTC(), actual.LastAccessed.UTC())
	}

	if expected.TrackProgress != actual.TrackProgress {
		return fmt.Errorf("flashcard set TrackProgress expected=%v, got=%v", expected.TrackProgress, actual.TrackProgress)
	}

	if expected.Front != actual.Front {
		return fmt.Errorf("flashcard set Front expected=%v, got=%v", expected.Front, actual.Front)
	}

	if expected.Shuffle != actual.Shuffle {
		return fmt.Errorf("flashcard set Shuffle expected=%v, got=%v", expected.Shuffle, actual.Shuffle)
	}

	if expected.ShuffleSeed != actual.ShuffleSeed {
		return fmt.Errorf("flashcard set ShuffleSeed expected=%v, got=%v", expected.ShuffleSeed, actual.ShuffleSeed)
	}

	return nil
}
