package main

import (
	tfcapp "github.com/d3akhtar/tfc/app"
	"github.com/d3akhtar/tfc/db"
	"github.com/d3akhtar/tfc/tui"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func main() {
	tui.SetDefaults()

	app := tview.NewApplication()

	database, err := db.InitializeSchema()
	if err != nil {
		panic(err)
	}

	state := tfcapp.NewAppState(app)

	tui.Init(state, database)

	views := state.Navigation.Views()

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
