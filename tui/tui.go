package tui

import (
	"slices"
	"strings"

	"github.com/d3akhtar/tfc/app"
	"github.com/d3akhtar/tfc/db"
	"github.com/rivo/tview"
)

var (
	filteredFlashcardSetList []db.FlashcardSet
	filteredFolderList       []db.Folder
)

var (
	recentOption       = 0
	alphabeticalOption = 1
)

var (
	mainPageName      = "main"
	newFolderPageName = "folder"
)

var formNewFolderName = ""

func Init(appState *app.State) {
	InitHomeUi(appState)
	InitLibraryUi(appState)
	InitFlashcardEditUi(appState)
	InitFolderUi(appState)
}

func NewPaddedFrameAllSides(amount int) *tview.Frame {
	return tview.NewFrame(nil).SetBorders(amount, amount, 0, 0, amount, amount)
}

func NewPaddedFrameXY(x, y int) *tview.Frame {
	return tview.NewFrame(nil).SetBorders(y, y, 0, 0, x, x)
}

func NewPaddedFrame(top, bottom, left, right int) *tview.Frame {
	return tview.NewFrame(nil).SetBorders(top, bottom, 0, 0, left, right)
}

func NewModal(p tview.Primitive, width, height int) tview.Primitive {
	return tview.NewFlex().
		AddItem(nil, 0, 1, false).
		AddItem(tview.NewFlex().SetDirection(tview.FlexRow).
			AddItem(nil, 0, 1, false).
			AddItem(p, height, 1, true).
			AddItem(nil, 0, 1, false), width, 1, true).
		AddItem(nil, 0, 1, false)
}

func sortFolderCollection(option int, folders []db.Folder) {
	switch option {
	case recentOption:
		slices.SortFunc(folders, func(a, b db.Folder) int {
			return -a.LastAccessed.Compare(b.LastAccessed)
		})
	case alphabeticalOption:
		slices.SortFunc(folders, func(a, b db.Folder) int {
			return strings.Compare(a.Name, b.Name)
		})
	}
}

func sortFlashcardSetCollection(option int, flashcardSets []db.FlashcardSet) {
	switch option {
	case recentOption:
		slices.SortFunc(flashcardSets, func(a, b db.FlashcardSet) int {
			return -a.LastAccessed.Compare(b.LastAccessed)
		})
	case alphabeticalOption:
		slices.SortFunc(flashcardSets, func(a, b db.FlashcardSet) int {
			return strings.Compare(a.Name, b.Name)
		})
	}
}
