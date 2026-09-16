package importing_test

import (
	"fmt"
	"os"
	"path"
	"regexp"
	"slices"
	"strings"
	"testing"

	"github.com/d3akhtar/tfc/domain"
	"github.com/d3akhtar/tfc/importing"
)

func TestImportFlashcardSet(t *testing.T) {
	wd, _ := os.Getwd()
	testFilePath := path.Join(wd, "test.tfcfs")
	file, _ := os.Create(testFilePath)
	defer file.Close()
	defer os.Remove(testFilePath)

	contents := `
		{
			"id": 1,
			"name": "name",
			"description": "desc",
			"flashcards": [
				{
					"question": "q1",
					"answer": "a1",
					"flashcardSetId": 1,
					"position": 0
				},
				{
					"question": "q2",
					"answer": "a2",
					"flashcardSetId": 1,
					"position": 1
				}
			]
		}
	`

	contents = strings.TrimSpace(contents)
	re := regexp.MustCompile(`\s+`)
	contents = re.ReplaceAllString(contents, "")
	os.WriteFile(testFilePath, []byte(contents), 0644)

	got, err := importing.ImportFlashcardSet(testFilePath)
	if err != nil {
		t.Fatalf("Error while importing flashcard set %v", err)
	}

	expected := domain.FlashcardSet{
		Id:          1,
		Name:        "name",
		Description: "desc",
	}

	expected.AddFlashcard("q1", "a1")
	expected.AddFlashcard("q2", "a2")

	err = testCmpFlashcardSet(&expected, got)
	if err != nil {
		t.Fatalf("%v", err)
	}
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

	equal := func(a, b domain.Flashcard) bool {
		return a.Question == b.Question && a.Answer == b.Answer && a.Position == b.Position
	}

	if !slices.EqualFunc(expected.Flashcards, actual.Flashcards, equal) {
		return fmt.Errorf("flashcard set flashcards expected=%v, got=%v", expected.Flashcards, actual.Flashcards)
	}

	return nil
}
