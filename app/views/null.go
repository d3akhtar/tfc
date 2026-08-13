package views

import "github.com/gdamore/tcell/v3"

type Null struct{}

func (v *Null) Draw(screen tcell.Screen)           {}
func (v *Null) HandleEvent(event tcell.Event) View { return nil }
