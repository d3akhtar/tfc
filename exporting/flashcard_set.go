package exporting

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/d3akhtar/tfc/domain"
)

func ExportFlashcardSet(flashcardSet *domain.FlashcardSet, path string) error {
	exportedJson, err := json.Marshal(flashcardSet)
	if err != nil {
		return fmt.Errorf("Error occurred while marshalling flashcard set")
	}

	err = os.WriteFile(path, exportedJson, 0644)
	if err != nil {
		return fmt.Errorf("Error occurred while marshalling flashcard set")
	}

	return nil
}
