package main

import (
	tfcapp "github.com/d3akhtar/tfc/app"
	"github.com/d3akhtar/tfc/db"
	"github.com/d3akhtar/tfc/tui"
	"github.com/rivo/tview"
)

func main() {
	tui.SetDefaults()

	app := tview.NewApplication()

	database, err := db.InitializeSchema()
	if err != nil {
		panic(err)
	}

	state := tfcapp.NewApp(app)

	tui.Init(state, database)

	state.Navigation.GoToView(tfcapp.VIEW_NAMES.Home)

	if err := app.SetRoot(state.Navigation.Views(), true).EnableMouse(true).Run(); err != nil {
		panic(err)
	}
}
