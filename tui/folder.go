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

func InitFolderUi(appState *app.State) {
	folder := tview.NewGrid().
		SetRows(-1, -23).
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

	flashcardSetList := tview.NewTable().SetSelectable(true, false).
		SetSelectedFunc(func(row, _ int) {
			pos := row
			selectedFlashcardSet := appState.SelectedFolder.FlashcardSets[pos]
			selectedFlashcardSet.LastAccessed = time.Now()
			appState.SelectedFlashcardSet = selectedFlashcardSet
		})

	sortDropdown.
		SetSelectedFunc(func(_ string, option int) {
			sortFlashcardSetCollection(option, filteredFlashcardSetList)
			flashcardSetList.Clear()
			for i, flashcardSet := range filteredFlashcardSetList {
				flashcardSetList.SetCell(
					i,
					0,
					tview.NewTableCell(flashcardSet.String()).SetExpansion(1),
				)
			}
		})

	flashcardSetList.
		SetBorder(true).
		SetFocusFunc(func() {
			flashcardSetList.SetBorderColor(tcell.ColorGreen)
			flashcardSetList.SetTitleColor(tcell.ColorGreen)
		}).
		SetBlurFunc(func() {
			flashcardSetList.SetBorderColor(tcell.ColorWhite)
			flashcardSetList.SetTitleColor(tcell.ColorWhite)
		}).
		SetBorderColor(tcell.ColorGreen).
		SetTitleColor(tcell.ColorGreen)

	searchFolderInputField := tview.NewInputField().
		SetLabel("").
		SetPlaceholder("Search folder...").
		SetChangedFunc(func(text string) {
			flashcardSetList.Clear()
			filteredFlashcardSetList = []db.FlashcardSet{}

			for _, collection := range appState.SelectedFolder.FlashcardSets {
				if text == "" || strings.Contains(collection.Name, text) {
					filteredFlashcardSetList = append(filteredFlashcardSetList, collection)
				}
			}

			for i, collection := range filteredFlashcardSetList {
				flashcardSetList.SetCell(
					i,
					0,
					tview.NewTableCell(collection.String()).SetExpansion(1),
				)
			}
		}).
		SetFieldWidth(0)

	flashcardSetList.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyBacktab:
			appState.App.SetFocus(sortDropdown)
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
			appState.App.SetFocus(flashcardSetList)
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
			appState.App.SetFocus(flashcardSetList)
			return nil
		}

		return event
	})

	folder.
		AddItem(folderNameLabel, 0, 0, 1, 1, 0, 0, false).
		AddItem(searchFolderInputField, 0, 1, 1, 1, 0, 0, false).
		AddItem(tview.NewBox(), 0, 2, 1, 1, 0, 0, false).
		AddItem(sortDropdown, 0, 3, 1, 1, 0, 0, false).
		AddItem(flashcardSetList, 1, 0, 1, 4, 0, 0, true)

	appState.Navigation.AddView(app.VIEW_NAMES.Folder, folder, false, func() {
		searchFolderInputField.SetText("")

		filteredFlashcardSetList = slices.Clone(appState.SelectedFolder.FlashcardSets)

		option, _ := sortDropdown.GetCurrentOption()
		sortFlashcardSetCollection(option, filteredFlashcardSetList)

		folderNameLabel.SetText(fmt.Sprintf("[ %s ]", appState.SelectedFolder.Name))

		flashcardSetList.Clear()

		for i, flashcardSet := range filteredFlashcardSetList {
			flashcardSetList.SetCell(
				i,
				0,
				tview.NewTableCell(flashcardSet.String()).SetExpansion(1),
			)
		}
	})
}
