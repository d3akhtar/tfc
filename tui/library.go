package tui

import "github.com/rivo/tview"

func InitLibraryUi(app *tview.Application, pages *tview.Pages) tview.Primitive {
	return tview.NewBox().
		SetBorder(true).
		SetTitle(VIEW_NAMES.Library).
		SetTitleAlign(tview.AlignLeft)
}
