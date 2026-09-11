package domain

import (
	"fmt"
	"time"

	"github.com/d3akhtar/tfc/utils"
)

type Folder struct {
	Id            int
	Name          string
	FlashcardSets []FlashcardSet
	LastAccessed  time.Time

	FlashcardSetCount int
}

func (f Folder) String() string {
	return fmt.Sprintf("%s | %d Sets | Last Accessed: %s", f.Name, f.count(), f.LastAccessed.Local().Format(utils.TimeFormat))
}

func (f *Folder) count() int {
	if len(f.FlashcardSets) == 0 {
		return f.FlashcardSetCount
	} else {
		return len(f.FlashcardSets)
	}
}
