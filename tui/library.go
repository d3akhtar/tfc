package tui

import (
	"slices"
	"strings"
	"time"

	"github.com/d3akhtar/tfc/app"
	"github.com/d3akhtar/tfc/db"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func InitLibraryUi(appState *app.State) tview.Primitive {
	library := tview.NewGrid().
		SetRows(-1, -10, -10, -1).
		SetColumns(-5, 3, -1)

	sortDropdown := tview.NewDropDown().
		SetLabel("Sort: ").
		SetLabelWidth(7).
		AddOption("Recent", nil).
		AddOption("Alphabetical", nil).
		SetFieldWidth(0).
		SetFieldBackgroundColor(tcell.ColorGray).
		SetCurrentOption(0)

	folderList := tview.NewTable().SetSelectable(true, false).
		SetSelectedFunc(func(row, _ int) {
			pos := row
			selectedFolder := filteredFolderList[pos]
			appState.SetSelectedFolder(selectedFolder)
		})

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
		})

	folderList.
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

	folderList.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyBacktab:
			appState.App.SetFocus(sortDropdown)
			return nil
		case tcell.KeyTab:
			appState.App.SetFocus(collectionList)
			return nil
		}

		if event.Key() == tcell.KeyRune && event.Rune() == '/' {
			appState.App.SetFocus(searchFolderInputField)
			return nil
		}

		return event
	})

	collectionList.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyBacktab:
			appState.App.SetFocus(folderList)
			return nil
		}

		if event.Key() == tcell.KeyRune && event.Rune() == '/' {
			appState.App.SetFocus(searchFolderInputField)
			return nil
		}

		return event
	})

	sortDropdown.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyTab:
			appState.App.SetFocus(folderList)
			return nil
		case tcell.KeyBacktab:
			appState.App.SetFocus(searchFolderInputField)
			return nil
		}

		return event
	})

	searchFolderInputField.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyTab:
			appState.App.SetFocus(sortDropdown)
			return nil
		case tcell.KeyEnter:
			appState.App.SetFocus(folderList)
			return nil
		}

		return event
	})

	library.
		AddItem(searchFolderInputField, 0, 0, 1, 1, 0, 0, false).
		AddItem(tview.NewBox(), 0, 1, 1, 1, 0, 0, false).
		AddItem(sortDropdown, 0, 2, 1, 1, 0, 0, false).
		AddItem(folderList, 1, 0, 1, 3, 0, 0, true).
		AddItem(collectionList, 2, 0, 1, 3, 0, 0, false)

	return library
}
