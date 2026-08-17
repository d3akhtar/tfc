package main

import (
	"github.com/d3akhtar/tfc/tui"
	"github.com/rivo/tview"
)

func main() {
	app := tview.NewApplication()

	pages := tview.NewPages()

	home := tui.InitHomeUi(app, pages)
	library := tui.InitLibraryUi(app, pages)

	pages.AddPage(tui.PAGE_NAMES.Home, home, true, true)
	pages.AddPage(tui.PAGE_NAMES.Library, library, true, false)

	if err := app.SetRoot(pages, true).EnableMouse(true).Run(); err != nil {
		panic(err)
	}
}
