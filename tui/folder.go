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
		SetRows(3, -20, -2).
		SetColumns(-1, -3, 3, -1)

	folderNameLabel := tview.
		NewTextView().
		SetScrollable(false).
		SetSize(2, 30).
		SetTextAlign(tview.AlignCenter)

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
			appState.SelectedFlashcardSet = &selectedFlashcardSet
			appState.Navigation.GoToView(app.VIEW_NAMES.FlashcardSetPreview)
		})

	addFlashcardSetButton := NewButton("Add Flashcard Set")

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
		SetBorderColor(Focused).
		SetTitleColor(Focused)

	SetBorderFocusAndBlurCallbacks(flashcardSetList.Box)

	searchFolderInputField := tview.NewInputField().
		SetChangedFunc(func(text string) {
			flashcardSetList.Clear()
			filteredFlashcardSetList = []db.FlashcardSet{}

			for _, collection := range appState.SelectedFolder.FlashcardSets {
				if text == "" || strings.Contains(strings.ToLower(collection.Name), strings.ToLower(text)) {
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
		})

	SetBorderFocusAndBlurCallbacks(searchFolderInputField.Box)

	searchFolderInputField.
		SetFieldBackgroundColor(Background).
		SetTitle("Search folder").
		SetTitleAlign(tview.AlignLeft)

	flashcardSetList.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyBacktab:
			appState.App.SetFocus(sortDropdown)
			return nil
		case tcell.KeyTab:
			appState.App.SetFocus(addFlashcardSetButton)
			return nil
		case tcell.KeyRune:
			if event.Rune() == '/' {
				appState.App.SetFocus(searchFolderInputField)
				return nil
			}
		}

		return event
	})

	addFlashcardSetButton.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyBacktab:
			appState.App.SetFocus(flashcardSetList)
			return nil
		case tcell.KeyRune:
			if event.Rune() == '/' {
				appState.App.SetFocus(searchFolderInputField)
				return nil
			}
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
		AddItem(
			tview.NewFlex().
				SetDirection(tview.FlexRow).
				AddItem(nil, 0, 1, false).
				AddItem(folderNameLabel, 0, 1, true).
				AddItem(nil, 0, 1, false),
			0, 0, 1, 1, 0, 0, false).
		AddItem(searchFolderInputField, 0, 1, 1, 1, 0, 0, false).
		AddItem(tview.NewBox(), 0, 2, 1, 1, 0, 0, false).
		AddItem(
			tview.NewFlex().
				SetDirection(tview.FlexRow).
				AddItem(nil, 0, 1, false).
				AddItem(sortDropdown, 0, 1, true).
				AddItem(nil, 0, 1, false),
			0, 3, 1, 1, 0, 0, false).
		AddItem(flashcardSetList, 1, 0, 1, 4, 0, 0, true).
		AddItem(addFlashcardSetButton, 2, 0, 1, 4, 0, 0, false)

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
