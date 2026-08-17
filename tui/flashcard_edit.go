package tui

import "github.com/rivo/tview"

func InitFlashcardEditUi(app *tview.Application, pages *tview.Pages) tview.Primitive {
	return tview.NewBox().
		SetBorder(true).
		SetTitle(VIEW_NAMES.FlashcardEdit).
		SetTitleAlign(tview.AlignLeft)
}
