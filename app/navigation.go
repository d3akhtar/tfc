package app

import "github.com/rivo/tview"

type Navigation struct {
	app *tview.Application

	pages   *tview.Pages
	history []string

	views map[string]*View
}

func NewNavigation(app *tview.Application) *Navigation {
	return &Navigation{
		app: app,

		history: []string{},
		pages:   tview.NewPages(),

		views: map[string]*View{},
	}
}

func (n *Navigation) Views() *tview.Pages {
	return n.pages
}

func (n *Navigation) AddView(pageName string, page *tview.Pages, visible bool, refresh RefreshCallback, exit ExitCallback) {
	view := NewView(pageName, page, refresh, exit, n.app)

	n.pages.AddPage(pageName, view.page, true, visible)
	n.views[pageName] = view

	view.lastFocusedPrimitive = page
}

func (n *Navigation) GoToView(pageName string) {
	n.hidePage(n.MostRecentlyVisitedViewName())
	if !n.showPage(pageName) {
		n.showPage(n.MostRecentlyVisitedViewName())
		return
	}

	n.history = append(n.history, pageName)
}

func (n *Navigation) RevertView() {
	if len(n.history) == 0 {
		return
	}

	err := n.currentView().exit()
	if err != nil {
		n.currentView().Error(err)
		return
	}

	n.pages.HidePage(n.MostRecentlyVisitedViewName())
	n.history = n.history[:len(n.history)-1]
	n.showPage(n.MostRecentlyVisitedViewName())
}

func (n *Navigation) MostRecentlyVisitedViewName() string {
	if len(n.history) > 0 {
		return n.history[len(n.history)-1]
	} else {
		return VIEW_NAMES.Home
	}
}

func (n *Navigation) showPage(pageName string) bool {
	err := n.views[pageName].Show()
	if err != nil {
		n.currentView().Error(err)
		return false
	}

	n.pages.ShowPage(pageName)

	return true
}

func (n *Navigation) hidePage(pageName string) bool {
	err := n.views[pageName].Exit()
	if err != nil {
		n.currentView().Error(err)
		return false
	}

	n.pages.HidePage(pageName)

	return true
}

func (n *Navigation) currentView() *View {
	return n.views[n.MostRecentlyVisitedViewName()]
}
