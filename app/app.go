package main

import (
	"fmt"

	"github.com/d3akhtar/tfc/app/ui"
	"github.com/gdamore/tcell/v3"
)

type App struct {
	screen tcell.Screen
	state  *AppState
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
		screen: s,
		state:  state,
	}

	return app, nil
}

func (app *App) Run() error {
	ui.DrawBox(app.screen, 1, 1, 42, 7, ui.Styles.Box, "Click and drag to draw a box")
	ui.DrawBox(app.screen, 5, 9, 32, 14, ui.Styles.Box, "Press C to reset")

	quit := func() {
		maybePanic := recover()
		app.screen.Fini()
		if maybePanic != nil {
			panic(maybePanic)
		}
	}
	defer quit()

	ox, oy := -1, -1
	for {
		app.screen.Show()

		ev := <-app.screen.EventQ()

		switch ev := ev.(type) {
		case *tcell.EventResize:
			app.screen.Sync()
		case *tcell.EventKey:
			if ev.Key() == tcell.KeyEscape || ev.Key() == tcell.KeyCtrlC {
				break
			} else if ev.Key() == tcell.KeyCtrlL {
				app.screen.Sync()
			} else if ev.Str() == "C" || ev.Str() == "c" {
				app.screen.Clear()
			}
		case *tcell.EventMouse:
			x, y := ev.Position()

			if ox != -1 || oy != -1 {
				app.screen.Clear()
				label := fmt.Sprintf("%d,%d to %d,%d", ox, oy, x, y)
				ui.DrawBox(app.screen, ox, oy, x, y, ui.Styles.Box, label)
			}

			switch ev.Buttons() {
			case tcell.Button1, tcell.Button2:
				if ox < 0 {
					ox, oy = x, y
				}

			case tcell.ButtonNone:
				if ox >= 0 {
					label := fmt.Sprintf("%d,%d to %d,%d", ox, oy, x, y)
					ui.DrawBox(app.screen, ox, oy, x, y, ui.Styles.Box, label)
					ox, oy = -1, -1
				}
			}
		}
	}
}
