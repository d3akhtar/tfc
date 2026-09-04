package tui

import (
	"time"

	"github.com/d3akhtar/tfc/app"
	"github.com/d3akhtar/tfc/db"
	"github.com/d3akhtar/tfc/db/flashcard_set"
	"github.com/d3akhtar/tfc/db/folder"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func InitLibraryUi(appState *app.State, flashcardSetRepository flashcard_set.FlashcardSetRepo, folderRepository folder.FolderRepo) {
	library := tview.NewPages()

	libraryGrid := tview.NewGrid().
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
			appState.SetSelectedFolder(selectedFolder)
			appState.Navigation.GoToView(app.VIEW_NAMES.Folder)
		})

	flashcardSetList := tview.NewTable().SetSelectable(true, false).
		SetSelectedFunc(func(row, _ int) {
			pos := row
			selectedFlashcardSet := &filteredFlashcardSetList[pos]
			appState.SetSelectedFlashcardSet(selectedFlashcardSet)

			flashcards, err := flashcardSetRepository.GetAllFlashcardsForSet(appState.Context, appState.SelectedFlashcardSet())
			if err != nil {
				return
			}

			appState.SelectedFlashcardSet().Flashcards = flashcards

			appState.Navigation.GoToView(app.VIEW_NAMES.FlashcardSetPreview)
		})

	sortDropdown.
		SetSelectedFunc(func(_ string, option int) {
			sort := db.FilterSortCriteria(option)

			sets, err := flashcardSetRepository.FilterFlashcardSets(appState.Context, "", 500, 0, sort)
			if err != nil {
				return
			}

			filteredFlashcardSetList = sets

			folders, err := folderRepository.FilterFolders(appState.Context, "", 500, 0, sort)
			if err != nil {
				return
			}

			filteredFolderList = folders

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

	searchFolderInputField := tview.NewInputField()
	searchFolderInputField.
		SetDoneFunc(func(key tcell.Key) {
			if key != tcell.KeyEnter {
				return
			}

			flashcardSetList.Clear()
			folderList.Clear()

			text := searchFolderInputField.GetText()
			sortOption, _ := sortDropdown.GetCurrentOption()
			sort := db.FilterSortCriteria(sortOption)

			sets, err := flashcardSetRepository.FilterFlashcardSets(appState.Context, text, 500, 0, sort)
			if err != nil {
				return
			}

			filteredFlashcardSetList = sets

			folders, err := folderRepository.FilterFolders(appState.Context, text, 500, 0, sort)
			if err != nil {
				return
			}

			filteredFolderList = folders

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
		}

		return event
	})

	libraryGrid.
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

	refresh := func() error {
		sortOption, _ := sortDropdown.GetCurrentOption()
		option := db.FilterSortCriteria(sortOption)

		sets, err := flashcardSetRepository.FilterFlashcardSets(appState.Context, "", 500, 0, option)
		if err != nil {
			return err
		}

		filteredFlashcardSetList = sets

		folders, err := folderRepository.FilterFolders(appState.Context, "", 500, 0, option)
		if err != nil {
			return err
		}

		filteredFolderList = folders

		searchFolderInputField.SetText("")

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

		return nil
	}

	library.AddPage("main", libraryGrid, true, true)

	appState.Navigation.AddView(app.VIEW_NAMES.Library, library, false, refresh, nil)
}
