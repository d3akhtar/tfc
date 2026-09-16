package exporting_test

import (
	"os"
	"path"
	"regexp"
	"strings"
	"testing"

	"github.com/d3akhtar/tfc/domain"
	"github.com/d3akhtar/tfc/exporting"
)

func TestExportFlashcardSet(t *testing.T) {
	wd, _ := os.Getwd()
	testFilePath := path.Join(wd, "test.tfcfs")
	file, _ := os.Create(testFilePath)
	defer file.Close()
	defer os.Remove(testFilePath)

	testFlashcardSet := domain.FlashcardSet{
		Id:          1,
		Name:        "name",
		Description: "desc",
	}

	testFlashcardSet.AddFlashcard("q1", "a1")
	testFlashcardSet.AddFlashcard("q2", "a2")

	err := exporting.ExportFlashcardSet(&testFlashcardSet, testFilePath)
	if err != nil {
		t.Fatalf("Error while exporting flashcard set: %v", err)
	}

	fileJsonContentsBytes, err := os.ReadFile(testFilePath)
	fileJsonContents := string(fileJsonContentsBytes)

	expected := `
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

	expected = strings.TrimSpace(expected)
	re := regexp.MustCompile(`\s+`)
	expected = re.ReplaceAllString(expected, "")

	if fileJsonContents != expected {
		t.Fatalf("expected:\n%s\nactual\n%s\n", expected, fileJsonContents)
	}
}
