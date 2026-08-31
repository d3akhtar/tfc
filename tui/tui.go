package tui

import (
	"database/sql"
	"slices"
	"strings"

	"github.com/d3akhtar/tfc/app"
	"github.com/d3akhtar/tfc/db/flashcard_set"
	"github.com/d3akhtar/tfc/db/folder"
	"github.com/d3akhtar/tfc/db/quiz"
	"github.com/d3akhtar/tfc/domain"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

var (
	filteredFolderFlashcardSetList []domain.FlashcardSet
	filteredFlashcardSetList       []*domain.FlashcardSet
	filteredFolderList             []*domain.Folder
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

func SetDefaults() {
	tview.Borders.TopLeft = '╭'
	tview.Borders.TopRight = '╮'
	tview.Borders.BottomLeft = '╰'
	tview.Borders.BottomRight = '╯'

	tview.Borders.TopLeftFocus = '╭'
	tview.Borders.TopRightFocus = '╮'
	tview.Borders.BottomLeftFocus = '╰'
	tview.Borders.BottomRightFocus = '╯'

	tview.Borders.HorizontalFocus = tview.Borders.Horizontal
	tview.Borders.VerticalFocus = tview.Borders.Vertical

	tview.Styles.PrimitiveBackgroundColor = Background
	tview.Styles.BorderColor = BoxBorder
	tview.Styles.TitleColor = BoxBorder
	tview.Styles.PrimaryTextColor = Text
}

func Init(appState *app.State, db *sql.DB) {
	flashcardSetRepository := flashcard_set.NewFlashcardSetRepository(db)
	folderRepositry := folder.NewFolderRepository(db)
	quizRepository := quiz.NewQuizRepository(db)

	InitHomeUi(appState, flashcardSetRepository, folderRepositry)
	InitLibraryUi(appState, flashcardSetRepository, folderRepositry)
	InitFlashcardEditUi(appState, flashcardSetRepository)
	InitFlashcardSetPreview(appState)
	InitFolderUi(appState, folderRepositry, flashcardSetRepository)
	InitQuizNormalUi(appState)
	InitQuizProgressTrackUi(appState, quizRepository)
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

func NewVerticallyAlignedTextView(text string) *tview.TextView {
	t := tview.NewTextView().
		SetText(text).
		SetTextAlign(tview.AlignCenter)

	t.SetDrawFunc(func(screen tcell.Screen, x, y, width, height int) (int, int, int, int) {
		y += height / 2
		return x, y, width, height
	})

	return t
}

func SetBorderFocusAndBlurCallbacks(p *tview.Box) {
	p.
		SetBorder(true).
		SetFocusFunc(func() {
			p.SetBorderColor(Focused)
			p.SetTitleColor(Focused)
		}).
		SetBlurFunc(func() {
			p.SetBorderColor(BoxBorder).SetTitleColor(BoxBorder)
		}).SetBorderColor(BoxBorder).SetTitleColor(BoxBorder)
}

func NewButton(title string) *tview.Button {
	return tview.NewButton(title).
		SetStyle(PrimaryButtonStyle)
}

func sortFolderCollection(option int, folders []*domain.Folder) {
	switch option {
	case recentOption:
		slices.SortFunc(folders, func(a, b *domain.Folder) int {
			return -a.LastAccessed.Compare(b.LastAccessed)
		})
	case alphabeticalOption:
		slices.SortFunc(folders, func(a, b *domain.Folder) int {
			return strings.Compare(a.Name, b.Name)
		})
	}
}

func sortFlashcardSetPointerCollection(option int, flashcardSets []*domain.FlashcardSet) {
	switch option {
	case recentOption:
		slices.SortFunc(flashcardSets, func(a, b *domain.FlashcardSet) int {
			return -a.LastAccessed.Compare(b.LastAccessed)
		})
	case alphabeticalOption:
		slices.SortFunc(flashcardSets, func(a, b *domain.FlashcardSet) int {
			return strings.Compare(a.Name, b.Name)
		})
	}
}

func sortFlashcardSetCollection(option int, flashcardSets []domain.FlashcardSet) {
	switch option {
	case recentOption:
		slices.SortFunc(flashcardSets, func(a, b domain.FlashcardSet) int {
			return -a.LastAccessed.Compare(b.LastAccessed)
		})
	case alphabeticalOption:
		slices.SortFunc(flashcardSets, func(a, b domain.FlashcardSet) int {
			return strings.Compare(a.Name, b.Name)
		})
	}
}
