package ui

import (
	"github.com/gdamore/tcell/v3"
	"github.com/gdamore/tcell/v3/color"
)

var Styles = struct {
	Normal   tcell.Style
	Title    tcell.Style
	Subtitle tcell.Style

	Button        tcell.Style
	ButtonFocused tcell.Style

	Box tcell.Style

	Border tcell.Style
	Error  tcell.Style

	Background tcell.Style
}{
	Normal:   tcell.StyleDefault.Normal(),
	Title:    tcell.StyleDefault.Bold(true),
	Subtitle: tcell.StyleDefault.Bold(true).Foreground(color.Gray),

	Button:        tcell.StyleDefault.Foreground(color.White).Background(color.Green),
	ButtonFocused: tcell.StyleDefault.Foreground(color.White).Background(color.DarkGreen),

	Box: tcell.StyleDefault.Foreground(color.White).Background(color.Green),

	Border: tcell.StyleDefault.Foreground(color.Gray),
	Error:  tcell.StyleDefault.Foreground(color.Red),

	Background: tcell.StyleDefault.Background(color.Reset).Foreground(color.Reset),
}
