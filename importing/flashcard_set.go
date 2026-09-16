package importing

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/d3akhtar/tfc/domain"
)

func ImportFlashcardSet(path string) (*domain.FlashcardSet, error) {
	if !strings.HasSuffix(path, ".tfcfs") {
		return nil, fmt.Errorf("Path should end with .tfcfc for imported flashcard sets. Path: %s", path)
	}

	bytes, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("Error while reading file while importing flashcard set %v", err)
	}

	res := domain.FlashcardSet{}
	err = json.Unmarshal(bytes, &res)
	if err != nil {
		return nil, fmt.Errorf("Error while unmarshalling imported flashcard set %v", err)
	}

	return &res, nil
}
