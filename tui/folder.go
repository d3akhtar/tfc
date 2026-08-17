package tui

import (
	"fmt"
	"slices"
	"strings"
	"time"

	"github.com/d3akhtar/tfc/app"
	"github.com/d3akhtar/tfc/db"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

var (
	filteredCollectionList []db.Collection
)

var (
	recentOption       = 0
	alphabeticalOption = 1
)

func InitFolderUi(app *tview.Application, pages *tview.Pages, appState *app.State) tview.Primitive {
	folder := tview.NewGrid().
		SetRows(-1, -23, -1).
		SetColumns(-1, -4, 3, -1)

	folderNameLabel := tview.
		NewTextView().
		SetText(fmt.Sprintf("[ %s ]", appState.SelectedFolder.Name)).
		SetScrollable(false).
		SetSize(2, 30)

	sortDropdown := tview.NewDropDown().
		SetLabel("Sort: ").
		SetLabelWidth(7).
		AddOption("Recent", nil).
		AddOption("Alphabetical", nil).
		SetFieldWidth(0).
		SetFieldBackgroundColor(tcell.ColorGray).
		SetCurrentOption(0)

	collectionList := tview.NewTable().SetSelectable(true, false).
		SetSelectedFunc(func(row, _ int) {
			pos := row
			selectedCollection := appState.SelectedFolder.Collections[pos]
			selectedCollection.LastAccessed = time.Now()
			appState.SetSelectedCollection(selectedCollection)
		})

	sortDropdown.
		SetSelectedFunc(func(_ string, index int) {
			switch index {
			case recentOption:
				slices.SortFunc(filteredCollectionList, func(a, b db.Collection) int {
					return a.LastAccessed.Compare(b.LastAccessed)
				})
			case alphabeticalOption:
				slices.SortFunc(filteredCollectionList, func(a, b db.Collection) int {
					return strings.Compare(a.Name, b.Name)
				})
				collectionList.Clear()

				for i, collection := range filteredCollectionList {
					collectionList.SetCell(
						i,
						0,
						tview.NewTableCell(collection.String()).SetExpansion(1),
					)
				}
			}
		})

	collectionList.
		SetBorder(true).
		SetFocusFunc(func() {
			collectionList.SetBorderColor(tcell.ColorGreen)
			collectionList.SetTitleColor(tcell.ColorGreen)
		}).
		SetBlurFunc(func() {
			collectionList.SetBorderColor(tcell.ColorWhite)
			collectionList.SetTitleColor(tcell.ColorWhite)
		}).
		SetBorderColor(tcell.ColorGreen).
		SetTitleColor(tcell.ColorGreen)

	searchFolderInputField := tview.NewInputField().
		SetLabel("").
		SetPlaceholder("Search folder...").
		SetChangedFunc(func(text string) {
			collectionList.Clear()
			filteredCollectionList = []db.Collection{}

			for _, collection := range appState.SelectedFolder.Collections {
				if text == "" || strings.Contains(collection.Name, text) {
					filteredCollectionList = append(filteredCollectionList, collection)
				}
			}

			for i, collection := range filteredCollectionList {
				collectionList.SetCell(
					i,
					0,
					tview.NewTableCell(collection.String()).SetExpansion(1),
				)
			}
		}).
		SetFieldWidth(0)

	homeButton := tview.
		NewButton("Back To Home").
		SetSelectedFunc(func() {
			pages.ShowPage(VIEW_NAMES.Home)
			pages.HidePage(VIEW_NAMES.Folder)
		})

	homeButton.SetBorder(false)

	collectionList.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyTab:
			app.SetFocus(homeButton)
			return nil
		case tcell.KeyBacktab:
			app.SetFocus(sortDropdown)
			return nil
		}

		if event.Key() == tcell.KeyRune && event.Rune() == '/' {
			app.SetFocus(searchFolderInputField)
			return nil
		}

		return event
	})

	sortDropdown.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyTab:
			app.SetFocus(collectionList)
			return nil
		case tcell.KeyBacktab:
			app.SetFocus(searchFolderInputField)
			return nil
		}

		return event
	})

	searchFolderInputField.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyTab:
			app.SetFocus(sortDropdown)
			return nil
		case tcell.KeyEnter:
			app.SetFocus(collectionList)
			return nil
		}

		return event
	})

	homeButton.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyBacktab:
			app.SetFocus(collectionList)
			return nil
		}

		return event
	})

	folder.
		AddItem(folderNameLabel, 0, 0, 1, 1, 0, 0, false).
		AddItem(searchFolderInputField, 0, 1, 1, 1, 0, 0, false).
		AddItem(tview.NewBox(), 0, 2, 1, 1, 0, 0, false).
		AddItem(sortDropdown, 0, 3, 1, 1, 0, 0, false).
		AddItem(collectionList, 1, 0, 1, 4, 0, 0, true).
		AddItem(homeButton, 2, 0, 1, 4, 0, 0, false)

	appState.OnSelectedFolderChange(func(newFolder db.Folder) {
		filteredCollectionList = slices.Clone(newFolder.Collections)

		folderNameLabel.SetText(fmt.Sprintf("[ %s ]", newFolder.Name))

		collectionList.Clear()

		for i, collection := range newFolder.Collections {
			collectionList.SetCell(
				i,
				0,
				tview.NewTableCell(collection.String()).SetExpansion(1),
			)
		}
	})

	return folder
}
