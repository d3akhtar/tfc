package app

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type RefreshCallback func() error
type ExitCallback func() error

var defaultExitCallback ExitCallback = func() error { return nil }

var errorPanelName = "error"

type View struct {
	Name string

	page                 *tview.Pages
	refresh              RefreshCallback
	exit                 ExitCallback
	lastFocusedPrimitive tview.Primitive

	app *tview.Application
}

func NewView(name string, pages *tview.Pages, refresh RefreshCallback, exit ExitCallback, app *tview.Application) *View {
	if exit == nil {
		exit = defaultExitCallback
	}

	view := &View{
		Name:                 name,
		page:                 pages,
		refresh:              refresh,
		exit:                 exit,
		lastFocusedPrimitive: nil,
		app:                  app,
	}

	view.addErrorTextPanelToView()

	return view
}

func (v *View) Show() error {
	if err := v.refresh(); err != nil {
		return err
	}

	v.app.SetFocus(v.lastFocusedPrimitive)
	return nil
}

func (v *View) Exit() error {
	newlastFocusedPrimitive := v.app.GetFocus()
	if err := v.exit(); err != nil {
		return err
	}

	v.lastFocusedPrimitive = newlastFocusedPrimitive
	return nil
}

func (v *View) Error(err error) {
	showErrorInView(v.page, err)
}

func (v *View) addErrorTextPanelToView() {
	errorTextPanel := tview.NewTextView().
		SetTextColor(tcell.ColorRed).
		SetTextAlign(tview.AlignLeft)

	errorTextPanel.
		SetBorder(true).
		SetBorderColor(tcell.ColorRed).
		SetTitle("Error").
		SetTitleColor(tcell.ColorRed).
		SetTitleAlign(tview.AlignLeft).
		SetBorderPadding(1, 1, 1, 1)

	errorTextPanel.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyEnter:
			v.page.HidePage(errorPanelName)
			v.app.SetFocus(v.lastFocusedPrimitive)
			return nil
		}

		return event
	})

	v.page.AddPage(errorPanelName, newModal(errorTextPanel, 60, 30), true, false)
}

func showErrorInView(view *tview.Pages, err error) {
	view.ShowPage(errorPanelName)
	view.GetPage(errorPanelName).(*tview.Flex).GetItem(1).(*tview.Flex).GetItem(1).(*tview.TextView).SetText(err.Error())
}

func newModal(p tview.Primitive, width, height int) tview.Primitive {
	return tview.NewFlex().
		AddItem(nil, 0, 1, false).
		AddItem(tview.NewFlex().SetDirection(tview.FlexRow).
			AddItem(nil, 0, 1, false).
			AddItem(p, height, 1, true).
			AddItem(nil, 0, 1, false), width, 1, true).
		AddItem(nil, 0, 1, false)
}
