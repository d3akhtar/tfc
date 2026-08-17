package tui

import "github.com/rivo/tview"

type App struct {
	*tview.Application
}

func NewPaddedFrameAllSides(amount int) *tview.Frame {
	return tview.NewFrame(nil).SetBorders(amount, amount, 0, 0, amount, amount)
}

func NewPaddedFrameXY(x, y int) *tview.Frame {
	return tview.NewFrame(nil).SetBorders(y, y, 0, 0, x, x)
}

func NewPaddedFrame(top, bottom, left, right int) *tview.Frame {
	return tview.NewFrame(nil).SetBorders(top, bottom, 0, 0, left, right)
}

func NewModal(p tview.Primitive, width, height int) tview.Primitive {
	return tview.NewFlex().
		AddItem(nil, 0, 1, false).
		AddItem(tview.NewFlex().SetDirection(tview.FlexRow).
			AddItem(nil, 0, 1, false).
			AddItem(p, height, 1, true).
			AddItem(nil, 0, 1, false), width, 1, true).
		AddItem(nil, 0, 1, false)
}
