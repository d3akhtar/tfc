package importing_test

import (
	"os"
	"path"
	"regexp"
	"strings"
	"testing"

	"github.com/d3akhtar/tfc/domain"
	"github.com/d3akhtar/tfc/importing"
)

func TestImportFolder(t *testing.T) {
	wd, _ := os.Getwd()
	testFilePath := path.Join(wd, "test.tfcf")
	file, _ := os.Create(testFilePath)
	defer file.Close()
	defer os.Remove(testFilePath)

	contents :=
		`
		{
			"id": 1,
			"name": "folder",
			"flashcardSets": [
				{
					"id": 1,
					"name": "name1",
					"description": "desc1",
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
				},
				{
					"id": 2,
					"name": "name2",
					"description": "desc2",
					"flashcards": [
						{
							"question": "q3",
							"answer": "a3",
							"flashcardSetId": 2,
							"position": 0
						},
						{
							"question": "q4",
							"answer": "a4",
							"flashcardSetId": 2,
							"position": 1
						}
					]
				}
			]
		}
	`

	contents = strings.TrimSpace(contents)
	re := regexp.MustCompile(`\s+`)
	contents = re.ReplaceAllString(contents, "")
	os.WriteFile(testFilePath, []byte(contents), 0644)

	got, err := importing.ImportFolder(testFilePath)
	if err != nil {
		t.Fatalf("Error while importing folder %v", err)
	}

	fc1 := domain.FlashcardSet{
		Id:          1,
		Name:        "name1",
		Description: "desc1",
	}

	fc1.AddFlashcard("q1", "a1")
	fc1.AddFlashcard("q2", "a2")

	fc2 := domain.FlashcardSet{
		Id:          2,
		Name:        "name2",
		Description: "desc2",
	}

	fc2.AddFlashcard("q3", "a3")
	fc2.AddFlashcard("q4", "a4")

	expected := domain.Folder{
		Id:   1,
		Name: "folder",
		FlashcardSets: []domain.FlashcardSet{
			fc1, fc2,
		},
	}

	if got.Id != expected.Id {
		t.Fatalf("folder Id expected=%v, got=%v", expected.Id, got.Id)
	}

	if got.Name != expected.Name {
		t.Fatalf("folder Name expected=%v, got=%v", expected.Name, got.Name)
	}

	for i, gotFlashcardSet := range got.FlashcardSets {
		err = testCmpFlashcardSet(&expected.FlashcardSets[i], &gotFlashcardSet)
		if err != nil {
			t.Errorf("(%d) %v", i, err)
		}
	}
}
