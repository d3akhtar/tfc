package views

import (
	"github.com/d3akhtar/tfc/app/ui"
	"github.com/gdamore/tcell/v3"
)

type Home struct{}

func (v *Home) Draw(screen tcell.Screen) {
	ui.DrawText(screen, 6, 10, 31, 13, ui.Styles.Title, "Home")
	ui.DrawText(screen, 6, 16, 31, 19, ui.Styles.Title, "1 - My Library")
	ui.DrawText(screen, 6, 22, 31, 25, ui.Styles.Title, "Q - Quit")
}

func (v *Home) HandleEvent(event tcell.Event) View {
	switch ev := event.(type) {
	case *tcell.EventKey:
		switch ev.Str() {
		case "1":
			return &Library{}
		case "q", "Q":
			return &Null{}
		}
	}

	return nil
}
