package tui

import (
	"slices"
	"strings"
	"time"

	"github.com/d3akhtar/tfc/app"
	"github.com/d3akhtar/tfc/db/flashcard_set"
	"github.com/d3akhtar/tfc/db/folder"
	"github.com/d3akhtar/tfc/domain"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func InitLibraryUi(appState *app.State, flashcardSetRepository flashcard_set.FlashcardSetRepo, folderRepository folder.FolderRepo) {
	flashcardSets := []*domain.FlashcardSet{}
	folders := []*domain.Folder{}

	library := tview.NewGrid().
		SetRows(3, -10, -10).
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
			appState.Navigation.GoToView(app.VIEW_NAMES.FlashcardSetPreview)
		})

	sortDropdown.
		SetSelectedFunc(func(_ string, option int) {
			sortFlashcardSetPointerCollection(option, filteredFlashcardSetList)
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

	SetBorderFocusAndBlurCallbacks(flashcardSetList.Box)

	flashcardSetList.
		SetTitle("Flashcard Sets").
		SetTitleAlign(tview.AlignLeft)

	SetBorderFocusAndBlurCallbacks(folderList.Box)

	folderList.
		SetTitle("Folders").
		SetTitleAlign(tview.AlignLeft).
		SetBorderColor(Focused).
		SetTitleColor(Focused)

	searchFolderInputField := tview.NewInputField().
		SetChangedFunc(func(text string) {
			flashcardSetList.Clear()
			folderList.Clear()

			filteredFlashcardSetList = []*domain.FlashcardSet{}
			filteredFolderList = []*domain.Folder{}

			for _, flashcardSet := range flashcardSets {
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

			for _, folder := range folders {
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
		})

	SetBorderFocusAndBlurCallbacks(searchFolderInputField.Box)
	searchFolderInputField.
		SetFieldBackgroundColor(Background).
		SetTitle("Search").
		SetTitleAlign(tview.AlignLeft)

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
		AddItem(
			tview.NewFlex().
				SetDirection(tview.FlexRow).
				AddItem(nil, 0, 1, false).
				AddItem(sortDropdown, 0, 1, true).
				AddItem(nil, 0, 1, false),
			0, 2, 1, 1, 0, 0, false).
		AddItem(folderList, 1, 0, 1, 3, 0, 0, true).
		AddItem(flashcardSetList, 2, 0, 1, 3, 0, 0, false)

	appState.Navigation.AddView(app.VIEW_NAMES.Library, library, false, func() {
		sets, err := flashcardSetRepository.List(appState.Context, 0, 50)
		if err != nil {
			return
		}

		flashcardSets = sets

		folders, err = folderRepository.List(appState.Context, 0, 50)
		if err != nil {
			return
		}

		searchFolderInputField.SetText("")

		filteredFlashcardSetList = slices.Clone(flashcardSets)
		filteredFolderList = slices.Clone(folders)

		option, _ := sortDropdown.GetCurrentOption()

		sortFlashcardSetPointerCollection(option, filteredFlashcardSetList)
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
