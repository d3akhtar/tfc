package views

import "github.com/gdamore/tcell/v3"

type View interface {
	Draw(screen tcell.Screen)
	HandleEvent(event tcell.Event) View
}
