package folder_test

import (
	"context"
	"database/sql"
	"fmt"
	"slices"
	"testing"
	"time"

	"github.com/d3akhtar/tfc/db"
	"github.com/d3akhtar/tfc/db/flashcard_set"
	"github.com/d3akhtar/tfc/db/folder"
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

func TestFolderRepository_AddFlashcardSetsToFolder(t *testing.T) {
	ctx := context.Background()
	database := testDB(t)
	repo := folder.NewFolderRepository(database)
	flashcardSetRepo := flashcard_set.NewFlashcardSetRepository(database)

	flashcardSets := []domain.FlashcardSet{}
	for i := range 10 {
		flashcardSet := domain.FlashcardSet{
			Id:            i,
			Name:          fmt.Sprintf("name%d", i),
			Description:   fmt.Sprintf("desc%d", i),
			LastAccessed:  time.Now(),
			TrackProgress: i%2 == 0,
			Front:         domain.FlashcardFront((i + 1) % 2),
			Shuffle:       i%2 == 0,
			ShuffleSeed:   i + 20,
		}

		err := flashcardSetRepo.Create(ctx, &flashcardSet)
		if err != nil {
			t.Fatalf("Unexpected error %v", err)
		}

		flashcardSets = append(flashcardSets, flashcardSet)
	}

	folder := &domain.Folder{
		Name:         "folder",
		LastAccessed: time.Now(),
	}

	err := repo.Create(ctx, folder)
	if err != nil {
		t.Fatalf("Unexpected error %v", err)
	}

	err = repo.AddFlashcardSetsToFolder(ctx, folder, flashcardSets)
	if err != nil {
		t.Fatalf("Unexpected error %v", err)
	}

	for i, fs := range folder.FlashcardSets {
		err = testCmpFlashcardSet(&flashcardSets[i], &fs)
		if err != nil {
			t.Errorf("(%d) %v", i, err)
		}
	}

	gotFlashcardSets, err := repo.GetFlashcardSetsForFolder(ctx, folder)
	for i, got := range gotFlashcardSets {
		err = testCmpFlashcardSet(&flashcardSets[i], got)
		if err != nil {
			t.Errorf("(%d) %v", i, err)
		}
	}
}

func TestFolderRepository_Count(t *testing.T) {
	ctx := context.Background()
	database := testDB(t)
	repo := folder.NewFolderRepository(database)

	folders := []*domain.Folder{}

	for i := range 10 {
		folder := &domain.Folder{
			Name:         fmt.Sprintf("folder %d", i),
			LastAccessed: time.Now(),
		}

		folders = append(folders, folder)

		err := repo.Create(ctx, folder)
		if err != nil {
			t.Fatalf("Unexpected error while inserting folder %v: %v", *folder, err)
		}
	}

	count, err := repo.Count(ctx)

	if err != nil {
		t.Fatalf("Unexpected error %v", err)
	}

	if count != int64(len(folders)) {
		t.Fatalf("count expected=%v, got=%v", len(folders), count)
	}
}

func TestFolderRepository_Create(t *testing.T) {
	ctx := context.Background()
	database := testDB(t)
	repo := folder.NewFolderRepository(database)
	flashcardSetRepo := flashcard_set.NewFlashcardSetRepository(database)

	flashcardSets := []domain.FlashcardSet{}
	for i := range 10 {
		flashcardSet := domain.FlashcardSet{
			Id:            i,
			Name:          fmt.Sprintf("name%d", i),
			Description:   fmt.Sprintf("desc%d", i),
			LastAccessed:  time.Now(),
			TrackProgress: i%2 == 0,
			Front:         domain.FlashcardFront((i + 1) % 2),
			Shuffle:       i%2 == 0,
			ShuffleSeed:   i + 20,
		}

		err := flashcardSetRepo.Create(ctx, &flashcardSet)
		if err != nil {
			t.Fatalf("Unexpected error %v", err)
		}

		flashcardSets = append(flashcardSets, flashcardSet)
	}

	folder := &domain.Folder{
		Name:          "folder",
		LastAccessed:  time.Now(),
		FlashcardSets: flashcardSets,
	}

	err := repo.Create(ctx, folder)
	if err != nil {
		t.Fatalf("Unexpected error %v", err)
	}

	if folder.Id == 0 {
		t.Fatalf("Expected folder Id to be set")
	}

	gotFolderFlashcardSets, err := repo.GetFlashcardSetsForFolder(ctx, folder)
	if err != nil {
		t.Fatalf("Unexpected error %v", err)
	}

	for i, got := range gotFolderFlashcardSets {
		err = testCmpFlashcardSet(&flashcardSets[i], got)
		if err != nil {
			t.Errorf("(%d) %v", i, err)
		}
	}
}

func TestFolderRepository_Delete(t *testing.T) {
	ctx := context.Background()
	database := testDB(t)
	repo := folder.NewFolderRepository(database)

	folder := &domain.Folder{
		Name:         "folder",
		LastAccessed: time.Now(),
	}

	err := repo.Create(ctx, folder)
	if err != nil {
		t.Fatalf("Unexpected error %v", err)
	}

	err = repo.Delete(ctx, folder.Id)
	if err != nil {
		t.Fatalf("Unexpected error %v", err)
	}

	_, err = repo.GetById(ctx, folder.Id)
	if err != db.ErrNotFound {
		t.Fatalf("expected=%v, got=%v", db.ErrNotFound, err)
	}
}

func TestFolderRepository_FilterFlashcardSetsInFolder(t *testing.T) {
	ctx := context.Background()
	database := testDB(t)
	repo := folder.NewFolderRepository(database)
	flashcardSetRepo := flashcard_set.NewFlashcardSetRepository(database)

	flashcardSets := []domain.FlashcardSet{}
	for i := range 10 {
		flashcardSet := domain.FlashcardSet{
			Id:            i,
			Name:          fmt.Sprintf("name%d", i),
			Description:   fmt.Sprintf("desc%d", i),
			LastAccessed:  time.Now(),
			TrackProgress: i%2 == 0,
			Front:         domain.FlashcardFront((i + 1) % 2),
			Shuffle:       i%2 == 0,
			ShuffleSeed:   i + 20,
		}

		err := flashcardSetRepo.Create(ctx, &flashcardSet)
		if err != nil {
			t.Fatalf("Unexpected error %v", err)
		}

		flashcardSets = append(flashcardSets, flashcardSet)
	}

	folder := &domain.Folder{
		Name:          "folder",
		LastAccessed:  time.Now(),
		FlashcardSets: flashcardSets,
	}

	err := repo.Create(ctx, folder)
	if err != nil {
		t.Fatalf("Unexpected error %v", err)
	}

	gotFolderFlashcardSets, err := repo.FilterFlashcardSetsInFolder(ctx, folder, "name", 10, 0)
	if err != nil {
		t.Fatalf("Unexpected error %v", err)
	}

	for i, got := range gotFolderFlashcardSets {
		err = testCmpFlashcardSet(&flashcardSets[i], got)
		if err != nil {
			t.Errorf("(%d) %v", i, err)
		}
	}

	gotFolderFlashcardSets, err = repo.FilterFlashcardSetsInFolder(ctx, folder, "name0", 1, 0)
	if err != nil {
		t.Fatalf("Unexpected error %v", err)
	}

	if len(gotFolderFlashcardSets) != 1 {
		t.Fatalf("len(filtered) expected=%v, got=%v", 1, len(gotFolderFlashcardSets))
	}

	err = testCmpFlashcardSet(&flashcardSets[0], gotFolderFlashcardSets[0])
	if err != nil {
		t.Errorf("%v", err)
	}

	gotFolderFlashcardSets, err = repo.FilterFlashcardSetsInFolder(ctx, folder, "name", 5, 4)
	if err != nil {
		t.Fatalf("Unexpected error %v", err)
	}

	for i, got := range gotFolderFlashcardSets {
		expected := flashcardSets[i+4]
		err = testCmpFlashcardSet(&expected, got)
		if err != nil {
			t.Errorf("(%d) %v", i, err)
		}
	}
}

func TestFolderRepository_GetById(t *testing.T) {
	ctx := context.Background()
	database := testDB(t)
	repo := folder.NewFolderRepository(database)

	folder := &domain.Folder{
		Name:         "folder1",
		LastAccessed: time.Now(),
	}

	err := repo.Create(ctx, folder)
	if err != nil {
		t.Fatalf("Unexpected error %v", err)
	}

	got, err := repo.GetById(ctx, folder.Id)
	if err != nil {
		t.Fatalf("Unexpected error %v", err)
	}

	err = testCmpFolder(folder, got)
	if err != nil {
		t.Fatalf("%v", err)
	}
}

func TestFolderRepository_GetFlashcardSetsForFolder(t *testing.T) {
	ctx := context.Background()
	database := testDB(t)
	repo := folder.NewFolderRepository(database)
	flashcardSetRepo := flashcard_set.NewFlashcardSetRepository(database)

	flashcardSets := []domain.FlashcardSet{}
	for i := range 10 {
		flashcardSet := domain.FlashcardSet{
			Id:            i,
			Name:          fmt.Sprintf("name%d", i),
			Description:   fmt.Sprintf("desc%d", i),
			LastAccessed:  time.Now(),
			TrackProgress: i%2 == 0,
			Front:         domain.FlashcardFront((i + 1) % 2),
			Shuffle:       i%2 == 0,
			ShuffleSeed:   i + 20,
		}

		err := flashcardSetRepo.Create(ctx, &flashcardSet)
		if err != nil {
			t.Fatalf("Unexpected error %v", err)
		}

		flashcardSets = append(flashcardSets, flashcardSet)
	}

	folder := &domain.Folder{
		Name:          "folder",
		LastAccessed:  time.Now(),
		FlashcardSets: flashcardSets,
	}

	err := repo.Create(ctx, folder)
	if err != nil {
		t.Fatalf("Unexpected error %v", err)
	}

	gotFolderFlashcardSets, err := repo.GetFlashcardSetsForFolder(ctx, folder)
	if err != nil {
		t.Fatalf("Unexpected error %v", err)
	}

	for i, got := range gotFolderFlashcardSets {
		err = testCmpFlashcardSet(&flashcardSets[i], got)
		if err != nil {
			t.Errorf("(%d) %v", i, err)
		}
	}
}

func TestFolderRepository_List(t *testing.T) {
	ctx := context.Background()
	database := testDB(t)
	repo := folder.NewFolderRepository(database)
	flashcardSetRepo := flashcard_set.NewFlashcardSetRepository(database)

	folders := []*domain.Folder{}
	flashcardSets := []domain.FlashcardSet{}
	for i := range 10 {
		flashcardSet := domain.FlashcardSet{
			Id:            i,
			Name:          fmt.Sprintf("name%d", i),
			Description:   fmt.Sprintf("desc%d", i),
			LastAccessed:  time.Now(),
			TrackProgress: i%2 == 0,
			Front:         domain.FlashcardFront((i + 1) % 2),
			Shuffle:       i%2 == 0,
			ShuffleSeed:   i + 20,
		}

		err := flashcardSetRepo.Create(ctx, &flashcardSet)
		if err != nil {
			t.Fatalf("Unexpected error %v", err)
		}

		flashcardSets = append(flashcardSets, flashcardSet)

		folder := &domain.Folder{
			Name:          fmt.Sprintf("folder%d", i),
			LastAccessed:  time.Now(),
			FlashcardSets: slices.Clone(flashcardSets),
		}

		err = repo.Create(ctx, folder)
		if err != nil {
			t.Fatalf("Unexpected error %v", err)
		}

		folders = append(folders, folder)
	}

	got, err := repo.List(ctx, 0, 10)
	if err != nil {
		t.Fatalf("Unexpected error %v", err)
	}

	for i, a := range got {
		err := testCmpFolder(folders[i], a)
		if err != nil {
			t.Errorf("(%d) %v", i, err)
		}

		if a.FlashcardSetCount != len(folders[i].FlashcardSets) {
			t.Errorf("flashcard count expected=%d, got=%d", len(folders[i].FlashcardSets), a.FlashcardSetCount)
		}
	}

	got, err = repo.List(ctx, 3, 5)
	if err != nil {
		t.Fatalf("Unexpected error %v", err)
	}

	for i, a := range got {
		expected := folders[i+3]
		err := testCmpFolder(expected, a)
		if err != nil {
			t.Errorf("(%d) %v", i, err)
		}

		if a.FlashcardSetCount != len(expected.FlashcardSets) {
			t.Errorf("flashcard count expected=%d, got=%d", len(expected.FlashcardSets), a.FlashcardSetCount)
		}
	}
}

func TestFolderRepository_Update(t *testing.T) {
	ctx := context.Background()
	database := testDB(t)
	repo := folder.NewFolderRepository(database)

	folder := &domain.Folder{
		Name:         "folder",
		LastAccessed: time.Now(),
	}

	err := repo.Create(ctx, folder)
	if err != nil {
		t.Fatalf("Unexpected error %v", err)
	}

	folder.Name = "new folder name"
	err = repo.Update(ctx, folder)
	if err != nil {
		t.Fatalf("Unexpected error %v", err)
	}

	got, err := repo.GetById(ctx, folder.Id)
	if err != nil {
		t.Fatalf("Unexpected error %v", err)
	}

	err = testCmpFolder(folder, got)
	if err != nil {
		t.Fatalf("%v", err)
	}

	folder.Id = 2
	err = repo.Update(ctx, folder)
	if err != db.ErrNotFound {
		t.Fatalf("expected=%v, got=%v", db.ErrNotFound, err)
	}
}

func TestFolderRepository_RemoveFlashcardSetsFromFolder(t *testing.T) {
	ctx := context.Background()
	database := testDB(t)
	repo := folder.NewFolderRepository(database)
	flashcardSetRepo := flashcard_set.NewFlashcardSetRepository(database)

	flashcardSets := []domain.FlashcardSet{}
	for i := range 10 {
		flashcardSet := domain.FlashcardSet{
			Id:            i,
			Name:          fmt.Sprintf("name%d", i),
			Description:   fmt.Sprintf("desc%d", i),
			LastAccessed:  time.Now(),
			TrackProgress: i%2 == 0,
			Front:         domain.FlashcardFront((i + 1) % 2),
			Shuffle:       i%2 == 0,
			ShuffleSeed:   i + 20,
		}

		err := flashcardSetRepo.Create(ctx, &flashcardSet)
		if err != nil {
			t.Fatalf("Unexpected error %v", err)
		}

		flashcardSets = append(flashcardSets, flashcardSet)
	}

	folder := &domain.Folder{
		Name:          "folder",
		LastAccessed:  time.Now(),
		FlashcardSets: flashcardSets,
	}

	err := repo.Create(ctx, folder)
	if err != nil {
		t.Fatalf("Unexpected error %v", err)
	}

	err = repo.RemoveFlashcardSetsFromFolder(ctx, folder, flashcardSets[:4])

	gotFlashcardSets, err := repo.GetFlashcardSetsForFolder(ctx, folder)

	expectedRemainingFlashcardSets := flashcardSets[4:]

	for i, got := range gotFlashcardSets {
		err = testCmpFlashcardSet(&expectedRemainingFlashcardSets[i], got)
		if err != nil {
			t.Errorf("(%d) %v", i, err)
		}
	}
}

func TestFolderRepository_FilterFolders(t *testing.T) {
	ctx := context.Background()
	database := testDB(t)
	repo := folder.NewFolderRepository(database)
	flashcardSetRepo := flashcard_set.NewFlashcardSetRepository(database)

	folders := []*domain.Folder{}
	flashcardSets := []domain.FlashcardSet{}
	for i := range 10 {
		flashcardSet := domain.FlashcardSet{
			Id:            i,
			Name:          fmt.Sprintf("name%d", i),
			Description:   fmt.Sprintf("desc%d", i),
			LastAccessed:  time.Now(),
			TrackProgress: i%2 == 0,
			Front:         domain.FlashcardFront((i + 1) % 2),
			Shuffle:       i%2 == 0,
			ShuffleSeed:   i + 20,
		}

		err := flashcardSetRepo.Create(ctx, &flashcardSet)
		if err != nil {
			t.Fatalf("Unexpected error %v", err)
		}

		flashcardSets = append(flashcardSets, flashcardSet)

		folder := &domain.Folder{
			Name:          fmt.Sprintf("folder%d", i),
			LastAccessed:  time.Now(),
			FlashcardSets: slices.Clone(flashcardSets),
		}

		err = repo.Create(ctx, folder)
		if err != nil {
			t.Fatalf("Unexpected error %v", err)
		}

		folders = append(folders, folder)
	}

	got, err := repo.FilterFolders(ctx, "folder", 0, 10)
	if err != nil {
		t.Fatalf("Unexpected error %v", err)
	}

	for i, a := range got {
		err := testCmpFolder(folders[i], a)
		if err != nil {
			t.Errorf("(%d) %v", i, err)
		}

		if a.FlashcardSetCount != len(folders[i].FlashcardSets) {
			t.Errorf("flashcard count expected=%d, got=%d", len(folders[i].FlashcardSets), a.FlashcardSetCount)
		}
	}

	got, err = repo.FilterFolders(ctx, "folder", 5, 3)
	if err != nil {
		t.Fatalf("Unexpected error %v", err)
	}

	for i, a := range got {
		expected := folders[i+3]
		err := testCmpFolder(expected, a)
		if err != nil {
			t.Errorf("(%d) %v", i, err)
		}

		if a.FlashcardSetCount != len(expected.FlashcardSets) {
			t.Errorf("flashcard count expected=%d, got=%d", len(expected.FlashcardSets), a.FlashcardSetCount)
		}
	}

	got, err = repo.FilterFolders(ctx, "folder0", 10, 0)
	if err != nil {
		t.Fatalf("Unexpected error %v", err)
	}

	if len(got) != 1 {
		t.Fatalf("len(filtered) expected=%v, got=%v", 1, len(got))
	}

	err = testCmpFolder(folders[0], got[0])
	if err != nil {
		t.Errorf("%v", err)
	}
}

func testCmpFolder(expected, actual *domain.Folder) error {
	if expected.Id != actual.Id {
		return fmt.Errorf("folder Id expected=%v, got=%v", expected.Id, actual.Id)
	}

	if expected.Name != actual.Name {
		return fmt.Errorf("folder Name expected=%v, got=%v", expected.Name, actual.Name)
	}

	if expected.LastAccessed.UTC().Truncate(time.Second).Compare(actual.LastAccessed.Truncate(time.Second)) != 0 {
		return fmt.Errorf(
			"folder LastAccessed expected=%v, got=%v",
			expected.LastAccessed.UTC().Truncate(time.Second),
			actual.LastAccessed.UTC().Truncate(time.Second))
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
