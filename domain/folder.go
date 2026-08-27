package domain

import (
	"fmt"
	"time"
)

type Folder struct {
	Id            int
	Name          string
	FlashcardSets []FlashcardSet
	LastAccessed  time.Time

	FlashcardSetCount int
}

func (f Folder) String() string {
	return fmt.Sprintf("○ %s | %d Sets | Last Accessed: %s", f.Name, len(f.FlashcardSets), f.LastAccessed.Format(time.RFC822))
}
