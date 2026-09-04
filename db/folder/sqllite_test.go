package folder_test

import (
	"context"
	"database/sql"
	"fmt"
	"slices"
	"strings"
	"testing"
	"time"

	"github.com/d3akhtar/tfc/db"
	"github.com/d3akhtar/tfc/db/flashcard_set"
	"github.com/d3akhtar/tfc/db/folder"
	"github.com/d3akhtar/tfc/db/utils/test"
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
		err = test_utils.TestCmpFlashcardSet(&flashcardSets[i], got)
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

	slices.SortFunc(flashcardSets, func(a, b domain.FlashcardSet) int {
		return a.LastAccessed.Compare(b.LastAccessed)
	})

	gotFolderFlashcardSets, err := repo.FilterFlashcardSetsInFolder(ctx, folder, "name", 10, 0, db.Recent)
	if err != nil {
		t.Fatalf("Unexpected error %v", err)
	}

	for i, got := range gotFolderFlashcardSets {
		err = test_utils.TestCmpFlashcardSet(&flashcardSets[i], got)
		if err != nil {
			t.Errorf("(%d) %v", i, err)
		}
	}

	slices.SortFunc(flashcardSets, func(a, b domain.FlashcardSet) int {
		return strings.Compare(a.Name, b.Name)
	})

	gotFolderFlashcardSets, err = repo.FilterFlashcardSetsInFolder(ctx, folder, "name", 10, 0, db.Alphabetical)
	if err != nil {
		t.Fatalf("Unexpected error %v", err)
	}

	for i, got := range gotFolderFlashcardSets {
		err = test_utils.TestCmpFlashcardSet(&flashcardSets[i], got)
		if err != nil {
			t.Errorf("(%d) %v", i, err)
		}
	}

	slices.SortFunc(flashcardSets, func(a, b domain.FlashcardSet) int {
		return strings.Compare(a.Name, b.Name)
	})

	gotFolderFlashcardSets, err = repo.FilterFlashcardSetsInFolder(ctx, folder, "name0", 1, 0, db.Alphabetical)
	if err != nil {
		t.Fatalf("Unexpected error %v", err)
	}

	if len(gotFolderFlashcardSets) != 1 {
		t.Fatalf("len(filtered) expected=%v, got=%v", 1, len(gotFolderFlashcardSets))
	}

	err = test_utils.TestCmpFlashcardSet(&flashcardSets[0], gotFolderFlashcardSets[0])
	if err != nil {
		t.Errorf("%v", err)
	}

	gotFolderFlashcardSets, err = repo.FilterFlashcardSetsInFolder(ctx, folder, "name", 5, 4, db.Alphabetical)
	if err != nil {
		t.Fatalf("Unexpected error %v", err)
	}

	for i, got := range gotFolderFlashcardSets {
		expected := flashcardSets[i+4]
		err = test_utils.TestCmpFlashcardSet(&expected, got)
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

	err = test_utils.TestCmpFolder(folder, got)
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
		err = test_utils.TestCmpFlashcardSet(&flashcardSets[i], got)
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
		err := test_utils.TestCmpFolder(folders[i], a)
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
		err := test_utils.TestCmpFolder(expected, a)
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

	folder.Name = "new folder name"
	err = repo.Update(ctx, folder)
	if err != nil {
		t.Fatalf("Unexpected error %v", err)
	}

	got, err := repo.GetById(ctx, folder.Id)
	if err != nil {
		t.Fatalf("Unexpected error %v", err)
	}

	err = test_utils.TestCmpFolder(folder, got)
	if err != nil {
		t.Fatalf("%v", err)
	}

	gotFlashcardSets, err := repo.GetFlashcardSetsForFolder(ctx, got)
	for i, gotFlashcardSet := range gotFlashcardSets {
		err = test_utils.TestCmpFlashcardSet(&flashcardSets[i], gotFlashcardSet)
		if err != nil {
			t.Errorf("(%d) %v", i, err)
		}
	}

	for i := 11; i <= 20; i++ {
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

	got.FlashcardSets = flashcardSets

	err = repo.Update(ctx, got)
	if err != nil {
		t.Fatalf("Unexpected error %v", err)
	}

	gotFlashcardSets, err = repo.GetFlashcardSetsForFolder(ctx, got)
	for i, gotFlashcardSet := range gotFlashcardSets {
		err = test_utils.TestCmpFlashcardSet(&flashcardSets[i], gotFlashcardSet)
		if err != nil {
			t.Errorf("(%d) %v", i, err)
		}
	}

	folder.Id = 2
	err = repo.Update(ctx, folder)
	if err != db.ErrNotFound {
		t.Fatalf("expected=%v, got=%v", db.ErrNotFound, err)
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

	slices.SortFunc(folders, func(a, b *domain.Folder) int {
		return a.LastAccessed.Compare(b.LastAccessed)
	})

	got, err := repo.FilterFolders(ctx, "folder", 0, 10, db.Recent)
	if err != nil {
		t.Fatalf("Unexpected error %v", err)
	}

	for i, a := range got {
		err := test_utils.TestCmpFolder(folders[i], a)
		if err != nil {
			t.Errorf("(%d) %v", i, err)
		}

		if a.FlashcardSetCount != len(folders[i].FlashcardSets) {
			t.Errorf("flashcard count expected=%d, got=%d", len(folders[i].FlashcardSets), a.FlashcardSetCount)
		}
	}

	slices.SortFunc(folders, func(a, b *domain.Folder) int {
		return strings.Compare(a.Name, b.Name)
	})

	got, err = repo.FilterFolders(ctx, "folder", 0, 10, db.Alphabetical)
	if err != nil {
		t.Fatalf("Unexpected error %v", err)
	}

	for i, a := range got {
		err := test_utils.TestCmpFolder(folders[i], a)
		if err != nil {
			t.Errorf("(%d) %v", i, err)
		}

		if a.FlashcardSetCount != len(folders[i].FlashcardSets) {
			t.Errorf("flashcard count expected=%d, got=%d", len(folders[i].FlashcardSets), a.FlashcardSetCount)
		}
	}

	got, err = repo.FilterFolders(ctx, "folder", 5, 3, db.Alphabetical)
	if err != nil {
		t.Fatalf("Unexpected error %v", err)
	}

	for i, a := range got {
		expected := folders[i+3]
		err := test_utils.TestCmpFolder(expected, a)
		if err != nil {
			t.Errorf("(%d) %v", i, err)
		}

		if a.FlashcardSetCount != len(expected.FlashcardSets) {
			t.Errorf("flashcard count expected=%d, got=%d", len(expected.FlashcardSets), a.FlashcardSetCount)
		}
	}

	got, err = repo.FilterFolders(ctx, "folder0", 10, 0, db.Alphabetical)
	if err != nil {
		t.Fatalf("Unexpected error %v", err)
	}

	if len(got) != 1 {
		t.Fatalf("len(filtered) expected=%v, got=%v", 1, len(got))
	}

	err = test_utils.TestCmpFolder(folders[0], got[0])
	if err != nil {
		t.Errorf("%v", err)
	}
}

func TestFolderRepository_ListRecentlyAccessedFolders(t *testing.T) {
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

	slices.SortFunc(folders, func(a, b *domain.Folder) int {
		return a.LastAccessed.Compare(b.LastAccessed)
	})

	got, err := repo.ListRecentlyAccessedFolders(ctx)
	if err != nil {
		t.Fatalf("Unexpected error %v", err)
	}

	for i, a := range got {
		err := test_utils.TestCmpFolder(folders[i], a)
		if err != nil {
			t.Errorf("(%d) %v", i, err)
		}

		if a.FlashcardSetCount != len(folders[i].FlashcardSets) {
			t.Errorf("flashcard count expected=%d, got=%d", len(folders[i].FlashcardSets), a.FlashcardSetCount)
		}
	}
}
