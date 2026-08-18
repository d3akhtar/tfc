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

func InitLibraryUi(appState *app.State) {
	library := tview.NewGrid().
		SetRows(-1, -10, -10).
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
			selectedFolder.LastAccessed = time.Now()
			appState.SelectedFolder = selectedFolder
			appState.Navigation.GoToView(app.VIEW_NAMES.Folder)
		})

	flashcardSetList := tview.NewTable().SetSelectable(true, false).
		SetSelectedFunc(func(row, _ int) {
			pos := row
			selectedFlashcardSet := filteredFlashcardSetList[pos]
			selectedFlashcardSet.LastAccessed = time.Now()
			appState.SelectedFlashcardSet = selectedFlashcardSet
			appState.Navigation.GoToView(app.VIEW_NAMES.FlashcardEdit)
		})

	sortDropdown.
		SetSelectedFunc(func(_ string, option int) {
			sortFlashcardSetCollection(option, filteredFlashcardSetList)
			sortFolderCollection(option, filteredFolderList)

			flashcardSetList.Clear()
			folderList.Clear()

			for i, flashcardSet := range filteredFlashcardSetList {
				flashcardSetList.SetCell(
					i,
					0,
					tview.NewTableCell(flashcardSet.String()).SetExpansion(1),
				)
			}

			for i, folder := range filteredFolderList {
				folderList.SetCell(
					i,
					0,
					tview.NewTableCell(folder.String()).SetExpansion(1),
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
		SetTitle("Flashcard Sets").
		SetTitleAlign(tview.AlignLeft)

	folderList.
		SetBorder(true).
		SetFocusFunc(func() {
			folderList.SetBorderColor(tcell.ColorGreen)
			folderList.SetTitleColor(tcell.ColorGreen)
		}).
		SetBlurFunc(func() {
			folderList.SetBorderColor(tcell.ColorWhite)
			folderList.SetTitleColor(tcell.ColorWhite)
		}).
		SetBorderColor(tcell.ColorGreen).
		SetTitleColor(tcell.ColorGreen).
		SetTitle("Folders").
		SetTitleAlign(tview.AlignLeft)

	searchFolderInputField := tview.NewInputField().
		SetLabel("").
		SetPlaceholder("Search library...").
		SetChangedFunc(func(text string) {
			flashcardSetList.Clear()
			folderList.Clear()

			filteredFlashcardSetList = []db.FlashcardSet{}
			filteredFolderList = []db.Folder{}

			for _, flashcardSet := range appState.FlashcardSets {
				if text == "" || strings.Contains(strings.ToLower(flashcardSet.Name), strings.ToLower(text)) {
					filteredFlashcardSetList = append(filteredFlashcardSetList, flashcardSet)
				}
			}

			for i, flashcardSet := range filteredFlashcardSetList {
				flashcardSetList.SetCell(
					i,
					0,
					tview.NewTableCell(flashcardSet.String()).SetExpansion(1),
				)
			}

			for _, folder := range appState.Folders {
				if text == "" || strings.Contains(strings.ToLower(folder.Name), strings.ToLower(text)) {
					filteredFolderList = append(filteredFolderList, folder)
				}
			}

			for i, folder := range filteredFolderList {
				folderList.SetCell(
					i,
					0,
					tview.NewTableCell(folder.String()).SetExpansion(1),
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
			appState.App.SetFocus(flashcardSetList)
			return nil
		}

		if event.Key() == tcell.KeyRune && event.Rune() == '/' {
			appState.App.SetFocus(searchFolderInputField)
			return nil
		}

		return event
	})

	flashcardSetList.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
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
		AddItem(flashcardSetList, 2, 0, 1, 3, 0, 0, false)

	appState.Navigation.AddView(app.VIEW_NAMES.Library, library, false, func() {
		searchFolderInputField.SetText("")

		filteredFlashcardSetList = slices.Clone(appState.FlashcardSets)
		filteredFolderList = slices.Clone(appState.Folders)

		option, _ := sortDropdown.GetCurrentOption()

		sortFlashcardSetCollection(option, filteredFlashcardSetList)
		sortFolderCollection(option, filteredFolderList)

		flashcardSetList.Clear()
		folderList.Clear()

		for i, flashcardSet := range filteredFlashcardSetList {
			flashcardSetList.SetCell(
				i,
				0,
				tview.NewTableCell(flashcardSet.String()).SetExpansion(1),
			)
		}

		for i, folder := range filteredFolderList {
			folderList.SetCell(
				i,
				0,
				tview.NewTableCell(folder.String()).SetExpansion(1),
			)
		}
	})
}
