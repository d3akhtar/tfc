package main

import (
	"fmt"

	"github.com/d3akhtar/tfc/app/ui"
	"github.com/gdamore/tcell/v3"

	"github.com/d3akhtar/tfc/app/views"
)

type App struct {
	screen      tcell.Screen
	state       *AppState
	currentView views.View
}

func NewApp() (*App, error) {
	s, err := tcell.NewScreen()
	if err != nil {
		return nil, fmt.Errorf("%+v", err)
	}

	if err := s.Init(); err != nil {
		return nil, fmt.Errorf("%+v", err)
	}

	s.SetStyle(ui.Styles.Background)
	s.EnableMouse()
	s.EnablePaste()
	s.Clear()

	state := &AppState{
		Running: false,
		Settings: SettingsState{
			DarkMode:   false,
			ShowHidden: false,
		},
	}

	app := &App{
		screen:      s,
		state:       state,
		currentView: &views.Home{},
	}

	return app, nil
}

func (app *App) Run() error {
	quit := func() {
		maybePanic := recover()
		app.screen.Fini()
		if maybePanic != nil {
			panic(maybePanic)
		}
	}
	defer quit()

	for {
		app.screen.Clear()

		app.currentView.Draw(app.screen)

		app.screen.Show()

		event := <-app.screen.EventQ()

		switch event.(type) {
		case *tcell.EventResize:
			app.screen.Sync()
		}

		if newView := app.currentView.HandleEvent(event); newView != nil {
			switch newView.(type) {
			case *views.Null:
				return nil
			default:
				app.currentView = newView
			}
		}
	}
}
