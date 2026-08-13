package views

import (
	"github.com/d3akhtar/tfc/app/ui"
	"github.com/gdamore/tcell/v3"
)

type Library struct {
	Titles []string
}

func (v *Library) Draw(screen tcell.Screen) {
	ui.DrawText(screen, 6, 10, 31, 13, ui.Styles.Title, "Library")

	for i, title := range v.Titles {
		ui.DrawText(screen, 6, 16+i*6, 31, 19+i*6, ui.Styles.Subtitle, title)
	}

	ui.DrawText(screen, 6, 100, 31, 125, ui.Styles.Title, "Q - Home")
}

func (v *Library) HandleEvent(event tcell.Event) View {
	switch ev := event.(type) {
	case *tcell.EventKey:
		switch ev.Str() {
		case "q", "Q":
			return &Home{}
		}
	}

	return nil
}
