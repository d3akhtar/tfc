package tui

import (
	"database/sql"

	"github.com/d3akhtar/tfc/app"
	"github.com/d3akhtar/tfc/db/flashcard_set"
	"github.com/d3akhtar/tfc/db/folder"
	"github.com/d3akhtar/tfc/db/quiz"
	"github.com/d3akhtar/tfc/domain"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

var (
	filteredFolderFlashcardSetList []*domain.FlashcardSet
	filteredFlashcardSetList       []domain.FlashcardSet
	filteredFolderList             []*domain.Folder
)

var (
	mainPageName      = "main"
	newFolderPageName = "folder"
)

var formNewFolderName = ""

func SetDefaults() {
	tview.Styles.PrimitiveBackgroundColor = Background
	tview.Styles.BorderColor = BoxBorder
	tview.Styles.TitleColor = BoxBorder
	tview.Styles.PrimaryTextColor = Text
}

func Init(appState *app.State, db *sql.DB) {
	flashcardSetRepository := flashcard_set.NewFlashcardSetRepository(db)
	folderRepository := folder.NewFolderRepository(db)
	quizRepository := quiz.NewQuizRepository(db)

	appState.AddCallbackForSelectedFlashcardSetChange(func(fs *domain.FlashcardSet) {
		flashcardSetRepository.UpdateLastAccessedTime(appState.Context, fs)
	})

	appState.AddCallbackForSelectedFolderChange(func(f *domain.Folder) {
		folderRepository.UpdateLastAccessedTime(appState.Context, f)
	})

	InitHomeUi(appState, flashcardSetRepository, folderRepository)
	InitLibraryUi(appState, flashcardSetRepository, folderRepository)
	InitFlashcardEditUi(appState, flashcardSetRepository)
	InitFlashcardSetPreview(appState, flashcardSetRepository, quizRepository)
	InitFolderUi(appState, folderRepository, flashcardSetRepository)
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

func NewConfirmActionModal(appState *app.State, confirm, cancel func()) tview.Primitive {
	grid := tview.NewGrid().
		SetRows(-1, -1).
		SetColumns(-1, 1, -1)

	header := tview.NewTextView().
		SetTextAlign(tview.AlignCenter).
		SetText("Are you sure?").
		SetLabelWidth(0)

	confirmButton := NewButton("Yes").SetSelectedFunc(confirm)
	cancelButton := NewButton("No").SetSelectedFunc(cancel)

	confirmButton.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyRight, tcell.KeyTab:
			appState.SetFocus(cancelButton)
			return nil
		}

		return event
	})

	cancelButton.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyLeft, tcell.KeyBacktab:
			appState.SetFocus(confirmButton)
			return nil
		}

		return event
	})

	grid.
		AddItem(header, 0, 0, 1, 3, 0, 0, false).
		AddItem(confirmButton, 1, 0, 1, 1, 0, 0, true).
		AddItem(cancelButton, 1, 2, 1, 1, 0, 0, false)

	SetBorderFocusAndBlurCallbacks(grid.Box)

	modal := NewModal(grid, 30, 7)

	return modal
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
