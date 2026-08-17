package tui

import (
	"github.com/d3akhtar/tfc/db"
	"github.com/rivo/tview"
)

var (
	filteredCollectionList []db.Collection
	filteredFolderList     []db.Folder
)

var (
	recentOption       = 0
	alphabeticalOption = 1
)

var (
	mainPageName      = "main"
	newFolderPageName = "folder"
)

var formNewFolderName = ""

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
