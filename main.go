package main

import (
	"github.com/d3akhtar/tfc/app"
	"github.com/d3akhtar/tfc/tui"
	"github.com/rivo/tview"
)

func main() {
	state := app.NewAppState()

	tview.Borders.TopLeft = '╭'
	tview.Borders.TopRight = '╮'
	tview.Borders.BottomLeft = '╰'
	tview.Borders.BottomRight = '╯'

	app := tview.NewApplication()

	views := tview.NewPages()

	home := tui.InitHomeUi(app, views, state)
	library := tui.InitLibraryUi(app, views)
	flashcardEdit := tui.InitFlashcardEditUi(app, views)
	folder := tui.InitFolderUi(app, views, state)

	views.AddPage(tui.VIEW_NAMES.Home, home, true, true)
	views.AddPage(tui.VIEW_NAMES.Library, library, true, false)
	views.AddPage(tui.VIEW_NAMES.FlashcardEdit, flashcardEdit, true, false)
	views.AddPage(tui.VIEW_NAMES.Folder, folder, true, false)

	if err := app.SetRoot(views, true).EnableMouse(true).Run(); err != nil {
		panic(err)
	}
}
