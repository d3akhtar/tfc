package main

import (
	tfcapp "github.com/d3akhtar/tfc/app"
	"github.com/d3akhtar/tfc/tui"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func main() {
	tview.Borders.TopLeft = '╭'
	tview.Borders.TopRight = '╮'
	tview.Borders.BottomLeft = '╰'
	tview.Borders.BottomRight = '╯'

	app := tview.NewApplication()
	views := tview.NewPages()

	state := tfcapp.NewAppState(app, views)

	home := tui.InitHomeUi(state)
	library := tui.InitLibraryUi(state)
	flashcardEdit := tui.InitFlashcardEditUi(state)
	folder := tui.InitFolderUi(state)

	views.AddPage(tfcapp.VIEW_NAMES.Home, home, true, true)
	views.AddPage(tfcapp.VIEW_NAMES.Library, library, true, false)
	views.AddPage(tfcapp.VIEW_NAMES.FlashcardEdit, flashcardEdit, true, false)
	views.AddPage(tfcapp.VIEW_NAMES.Folder, folder, true, false)

	views.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyCtrlBackslash:
			state.Navigation.RevertView()
		}

		return event
	})

	if err := app.SetRoot(views, true).EnableMouse(true).Run(); err != nil {
		panic(err)
	}
}
