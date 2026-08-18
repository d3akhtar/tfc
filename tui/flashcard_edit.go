package tui

import (
	"github.com/d3akhtar/tfc/app"
	"github.com/rivo/tview"
)

func InitFlashcardEditUi(appState *app.State) {
	flashcardEdit := tview.NewBox().
		SetBorder(true).
		SetTitle(app.VIEW_NAMES.FlashcardEdit).
		SetTitleAlign(tview.AlignLeft)

	appState.Navigation.AddView(app.VIEW_NAMES.FlashcardEdit, flashcardEdit, false, nil)
}
