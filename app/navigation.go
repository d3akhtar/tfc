package app

import "github.com/rivo/tview"

type Navigation struct {
	views   *tview.Pages
	history []string
}

func NewNavigation(pages *tview.Pages) *Navigation {
	return &Navigation{
		views:   pages,
		history: []string{},
	}
}

func (n *Navigation) GoToView(pageName string) {
	n.views.ShowPage(pageName)
	n.views.HidePage(n.MostRecentlyVisitedViewName())
	n.history = append(n.history, pageName)
}

func (n *Navigation) RevertView() {
	if len(n.history) == 0 {
		return
	}

	n.views.HidePage(n.MostRecentlyVisitedViewName())
	n.history = n.history[:len(n.history)-1]
	n.views.ShowPage(n.MostRecentlyVisitedViewName())
}

func (n *Navigation) MostRecentlyVisitedViewName() string {
	if len(n.history) > 0 {
		return n.history[len(n.history)-1]
	} else {
		return VIEW_NAMES.Home
	}
}
