package db

import (
	"fmt"
	"time"
)

type Folder struct {
	Name          string         `json:"name"`
	FlashcardSets []FlashcardSet `json:"collections"`
	LastAccessed  time.Time
}

func (f Folder) String() string {
	return fmt.Sprintf("%s | %d Sets", f.Name, len(f.FlashcardSets))
}
