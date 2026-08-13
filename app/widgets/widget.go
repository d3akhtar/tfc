package widgets

import "github.com/gdamore/tcell/v3"

type Widget interface {
	Draw(screen tcell.Screen)
	HandleEvent(event tcell.Event)
}
