package tui

import (
	"fmt"
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

func InitFolderUi(appState *app.State, folderRepository folder.FolderRepo, flashcardSetRepository flashcard_set.FlashcardSetRepo) {
	flashcardSets := []*domain.FlashcardSet{}

	folderPage := tview.NewPages()

	folderView := tview.NewGrid().
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

	folderFlashcardSetList := tview.NewTable().SetSelectable(true, false).
		SetSelectedFunc(func(row, _ int) {
			pos := row
			selectedFlashcardSet := flashcardSets[pos]
			selectedFlashcardSet.LastAccessed = time.Now()
			appState.SelectedFlashcardSet = selectedFlashcardSet

			flashcards, err := flashcardSetRepository.GetAllFlashcardsForSet(appState.Context, appState.SelectedFlashcardSet)
			if err != nil {
				return
			}

			appState.SelectedFlashcardSet.Flashcards = flashcards

			appState.Navigation.GoToView(app.VIEW_NAMES.FlashcardSetPreview)
		})

	addFlashcardSetButton := NewButton("Add Flashcard Sets").
		SetSelectedFunc(func() {
			folderPage.SwitchToPage("flashcard")
		})

	sortDropdown.
		SetSelectedFunc(func(_ string, option int) {
			sortFlashcardSetCollection(option, filteredFolderFlashcardSetList)
			folderFlashcardSetList.Clear()
			for i, flashcardSet := range filteredFolderFlashcardSetList {
				folderFlashcardSetList.SetCell(
					i,
					0,
					tview.NewTableCell(flashcardSet.String()).SetExpansion(1),
				)
			}
		})

	SetBorderFocusAndBlurCallbacks(folderFlashcardSetList.Box)

	folderFlashcardSetList.
		SetBorderColor(Focused).
		SetTitleColor(Focused).
		SetTitle("Flashcards").
		SetTitleAlign(tview.AlignLeft)

	searchFolderInputField := tview.NewInputField().
		SetChangedFunc(func(text string) {
			folderFlashcardSetList.Clear()
			filteredFolderFlashcardSetList = []domain.FlashcardSet{}

			for _, flashcardSet := range appState.SelectedFolder.FlashcardSets {
				if text == "" || strings.Contains(strings.ToLower(flashcardSet.Name), strings.ToLower(text)) {
					filteredFolderFlashcardSetList = append(filteredFolderFlashcardSetList, flashcardSet)
				}
			}

			for i, flashcardSet := range filteredFolderFlashcardSetList {
				folderFlashcardSetList.SetCell(
					i,
					0,
					tview.NewTableCell(flashcardSet.String()).SetExpansion(1),
				)
			}
		})

	SetBorderFocusAndBlurCallbacks(searchFolderInputField.Box)

	searchFolderInputField.
		SetFieldBackgroundColor(Background).
		SetTitle("Search folder").
		SetTitleAlign(tview.AlignLeft)

	addFlashcardSetView := tview.NewGrid().
		SetRows(3, -2, 3)

	flashcardSetList := tview.NewTable().SetSelectable(true, false)
	flashcardSetList.
		SetSelectedFunc(func(row, _ int) {
			pos := row
			selectedFlashcardSet := flashcardSets[pos]

			check := func(set domain.FlashcardSet) bool {
				return set.Id == selectedFlashcardSet.Id
			}

			if slices.ContainsFunc(appState.SelectedFolder.FlashcardSets, check) {
				appState.SelectedFolder.FlashcardSets = slices.DeleteFunc(appState.SelectedFolder.FlashcardSets, check)
				flashcardSetList.GetCell(row, 0).SetTextColor(tcell.ColorWhite)
			} else {
				appState.SelectedFolder.FlashcardSets = append(appState.SelectedFolder.FlashcardSets, *selectedFlashcardSet)
				flashcardSetList.GetCell(row, 0).SetTextColor(tcell.ColorGreen)
			}

		})

	SetBorderFocusAndBlurCallbacks(flashcardSetList.Box)

	flashcardSetList.
		SetTitle("Add Flashcard Sets").
		SetTitleAlign(tview.AlignLeft).
		SetTitleColor(Focused).
		SetBorderColor(Focused)

	flashcardSetListSearchInputField := tview.NewInputField().
		SetChangedFunc(func(text string) {
			flashcardSetList.Clear()
			filteredFlashcardSetList = []*domain.FlashcardSet{}

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
		})

	SetBorderFocusAndBlurCallbacks(flashcardSetListSearchInputField.Box)
	flashcardSetListSearchInputField.
		SetFieldBackgroundColor(Background).
		SetTitle("Search Flashcards").
		SetTitleAlign(tview.AlignLeft)

	closeFlashcardSetAddViewButton := NewButton("Close").
		SetSelectedFunc(func() {
			folderPage.SwitchToPage("main")
		})

	flashcardSetListSearchInputField.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyTab:
			appState.SetFocus(flashcardSetList)
			return nil
		}

		return event
	})

	flashcardSetList.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyBacktab:
			appState.SetFocus(flashcardSetListSearchInputField)
			return nil
		case tcell.KeyTab:
			appState.SetFocus(closeFlashcardSetAddViewButton)
			return nil
		}

		return event
	})

	closeFlashcardSetAddViewButton.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyBacktab:
			appState.SetFocus(flashcardSetList)
			return nil
		}

		return event
	})

	addFlashcardSetView.
		AddItem(flashcardSetListSearchInputField, 0, 0, 1, 1, 0, 0, false).
		AddItem(flashcardSetList, 1, 0, 1, 1, 0, 0, true).
		AddItem(closeFlashcardSetAddViewButton, 2, 0, 1, 1, 0, 0, false)

	folderFlashcardSetList.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
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
			appState.App.SetFocus(folderFlashcardSetList)
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
			appState.App.SetFocus(folderFlashcardSetList)
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
			appState.App.SetFocus(folderFlashcardSetList)
			return nil
		}

		return event
	})

	folderView.
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
		AddItem(folderFlashcardSetList, 1, 0, 1, 4, 0, 0, true).
		AddItem(addFlashcardSetButton, 2, 0, 1, 4, 0, 0, false)

	folderPage.SetChangedFunc(func() {
		if pg, _ := folderPage.GetFrontPage(); pg == "main" && appState.SelectedFolder != nil {
			searchFolderInputField.SetText("")

			folderFlashcardSetList.Clear()

			filteredFolderFlashcardSetList = slices.Clone(appState.SelectedFolder.FlashcardSets)

			for i, flashcardSet := range filteredFolderFlashcardSetList {
				folderFlashcardSetList.SetCell(
					i,
					0,
					tview.NewTableCell(flashcardSet.String()).SetExpansion(1),
				)
			}
		}
	})

	folderPage.
		AddPage("main", folderView, true, true).
		AddPage("flashcard", addFlashcardSetView, true, false)

	refresh := func() error {
		sets, err := flashcardSetRepository.List(appState.Context, 0, 50)
		if err != nil {
			return err
		}

		flashcardSets = sets

		folderFlashcardSetPointers, err := folderRepository.GetFlashcardSetsForFolder(appState.Context, appState.SelectedFolder)
		if err != nil {
			return err
		}

		folderFlashcardSets := make([]domain.FlashcardSet, 0, len(folderFlashcardSetPointers))
		for _, fc := range folderFlashcardSetPointers {
			folderFlashcardSets = append(folderFlashcardSets, *fc)
		}

		appState.SelectedFolder.FlashcardSets = folderFlashcardSets

		searchFolderInputField.SetText("")

		filteredFolderFlashcardSetList = slices.Clone(appState.SelectedFolder.FlashcardSets)
		filteredFlashcardSetList = slices.Clone(flashcardSets)

		option, _ := sortDropdown.GetCurrentOption()
		sortFlashcardSetCollection(option, filteredFolderFlashcardSetList)

		folderNameLabel.SetText(fmt.Sprintf("[ %s ]", appState.SelectedFolder.Name))

		folderFlashcardSetList.Clear()
		flashcardSetList.Clear()

		for i, flashcardSet := range filteredFolderFlashcardSetList {
			folderFlashcardSetList.SetCell(
				i,
				0,
				tview.NewTableCell(flashcardSet.String()).SetExpansion(1),
			)
		}

		for i, flashcardSet := range filteredFlashcardSetList {
			var col tcell.Color
			if slices.ContainsFunc(appState.SelectedFolder.FlashcardSets, func(set domain.FlashcardSet) bool {
				return set.Name == flashcardSet.Name &&
					set.Description == flashcardSet.Description &&
					slices.Equal(set.Flashcards, flashcardSet.Flashcards)
			}) {
				col = tcell.ColorGreen
			} else {
				col = tcell.ColorWhite
			}

			flashcardSetList.SetCell(
				i,
				0,
				&tview.TableCell{Text: flashcardSet.String(), BackgroundColor: Background, Color: col, Expansion: 1},
			)
		}

		return nil
	}

	exit := func() error {
		return folderRepository.Update(appState.Context, appState.SelectedFolder)
	}

	appState.Navigation.AddView(app.VIEW_NAMES.Folder, folderPage, false, refresh, exit)
}
