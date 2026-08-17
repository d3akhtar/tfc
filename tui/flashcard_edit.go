package tui

import (
	"github.com/d3akhtar/tfc/app"
	"github.com/rivo/tview"
)

func InitFlashcardEditUi(appState *app.State) tview.Primitive {
	return tview.NewBox().
		SetBorder(true).
		SetTitle(app.VIEW_NAMES.FlashcardEdit).
		SetTitleAlign(tview.AlignLeft)
}
