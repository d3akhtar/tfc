package app

import "github.com/rivo/tview"

type Navigation struct {
	views   *tview.Pages
	history []string

	refresh map[string]func()
}

func NewNavigation() *Navigation {
	return &Navigation{
		history: []string{},
		views:   tview.NewPages(),

		refresh: make(map[string]func()),
	}
}

func (n *Navigation) Views() *tview.Pages {
	return n.views
}

func (n *Navigation) AddView(pageName string, primitive tview.Primitive, visible bool, refresh func()) {
	n.views.AddPage(pageName, primitive, true, visible)
	n.refresh[pageName] = refresh
}

func (n *Navigation) GoToView(pageName string) {
	n.showPage(pageName)
	n.hidePage(n.MostRecentlyVisitedViewName())
	n.history = append(n.history, pageName)
}

func (n *Navigation) RevertView() {
	if len(n.history) == 0 {
		return
	}

	n.views.HidePage(n.MostRecentlyVisitedViewName())
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

func (n *Navigation) showPage(pageName string) {
	n.views.ShowPage(pageName)
	n.refresh[pageName]()
}

func (n *Navigation) hidePage(pageName string) {
	n.views.HidePage(pageName)
}
