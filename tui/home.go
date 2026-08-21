package tui

import (
	"fmt"

	"github.com/d3akhtar/tfc/app"
	"github.com/d3akhtar/tfc/db"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func InitHomeUi(appState *app.State) {
	home := tview.NewPages()

	recentSetsStudies := tview.NewTable().
		SetSelectable(true, false).
		SetSelectedFunc(func(row, _ int) {
			pos := row
			appState.SelectedFlashcardSet = &appState.RecentlyStudied[pos]
			appState.Navigation.GoToView(app.VIEW_NAMES.FlashcardSetPreview)
		})

	for i, flashcardSet := range appState.RecentlyStudied {
		recentSetsStudies.SetCell(i, 0, tview.NewTableCell(flashcardSet.String()).SetExpansion(1).SetTextColor(Text))
	}

	recentSetsStudies.
		SetFocusFunc(func() {
			recentSetsStudies.SetBorderColor(Focused)
			recentSetsStudies.SetTitleColor(Focused)
		}).
		SetBlurFunc(func() {
			recentSetsStudies.SetBorderColor(BoxBorder)
			recentSetsStudies.SetTitleColor(BoxBorder)
		}).
		SetBorder(true).
		SetTitle("Recent Sets Studies").
		SetTitleAlign(tview.AlignLeft).
		SetBorderColor(Focused).
		SetTitleColor(Focused).
		SetBorderPadding(1, 1, 1, 1).
		SetBackgroundColor(Background)

	folders := tview.NewTable().
		SetBorders(false).
		SetSelectable(true, true)

	folders.SetSelectedFunc(func(row, column int) {
		pos := (row * 3) + column
		appState.SelectedFolder = &appState.Folders[pos]
		appState.Navigation.GoToView(app.VIEW_NAMES.Folder)
	})

	for i, loadedFolder := range appState.Folders {
		row := i / 3
		column := i % 3
		folders.SetCell(row, column, tview.NewTableCell(fmt.Sprintf("○ %s", loadedFolder.Name)).SetExpansion(1).SetTextColor(Text))
	}

	folders.
		SetBorder(true).
		SetTitle("Folders").
		SetTitleAlign(tview.AlignLeft).
		SetFocusFunc(func() {
			folders.SetBorderColor(Focused)
			folders.SetTitleColor(Focused)
		}).
		SetBlurFunc(func() {
			folders.SetBorderColor(BoxBorder)
			folders.SetTitleColor(BoxBorder)
		}).
		SetBorderPadding(1, 1, 2, 0).
		SetBackgroundColor(Background).
		SetBorderColor(BoxBorder).
		SetTitleColor(BoxBorder)

	createFlashcardSetButton := NewButton("Flashcard")

	createFolderButton := NewButton("Folder")

	newFolderNameInputField := tview.NewInputField().
		SetLabel("Name: ").
		SetPlaceholder("Enter new folder name...").
		SetAcceptanceFunc(func(textToCheck string, lastChar rune) bool {
			return len(textToCheck) > 0
		}).SetChangedFunc(func(text string) {
		formNewFolderName = text
	})

	onFolderFormSubmit := func() {
		home.HidePage(newFolderPageName)
		appState.App.SetFocus(createFlashcardSetButton)

		appState.Folders = append(appState.Folders, db.Folder{
			Name:          formNewFolderName,
			FlashcardSets: []db.FlashcardSet{},
		})

		newFolderNameInputField.SetText("")

		folders.Clear()
		for i, loadedFolder := range appState.Folders {
			row := i / 3
			column := i % 3
			folders.SetCell(row, column, tview.NewTableCell(fmt.Sprintf("○ %s", loadedFolder.Name)).SetExpansion(1))
		}
	}

	newFolderForm := tview.NewForm().
		AddFormItem(newFolderNameInputField).
		AddButton("Save", onFolderFormSubmit).
		AddButton("Quit", func() {
			home.HidePage(newFolderPageName)
			appState.App.SetFocus(createFlashcardSetButton)
		}).
		SetButtonsAlign(tview.AlignCenter).
		SetSubmitFunc(onFolderFormSubmit)

	newFolderForm.
		SetBorder(true).
		SetTitle("Enter new folder name").
		SetTitleAlign(tview.AlignLeft)

	newFolderFormLayout := NewModal(newFolderForm, 100, 7)

	createFlashcardSetButton.SetSelectedFunc(func() {
		appState.SelectedFlashcardSet = nil
		appState.Navigation.GoToView(app.VIEW_NAMES.FlashcardEdit)
	})

	createFolderButton.SetSelectedFunc(func() {
		home.ShowPage(newFolderPageName)
		appState.App.SetFocus(newFolderForm)
	})

	createFlashcardSetFrame := NewPaddedFrameXY(2, 1).SetPrimitive(createFlashcardSetButton)
	createFlashcardSetFrame.SetBackgroundColor(Background)

	createFolderFrame := NewPaddedFrameXY(2, 1).SetPrimitive(createFolderButton)
	createFolderFrame.SetBackgroundColor(Background)

	create := tview.NewGrid().
		SetRows(-1, -1).
		AddItem(createFlashcardSetFrame, 0, 0, 1, 1, 0, 0, true).
		AddItem(createFolderFrame, 1, 0, 1, 1, 0, 0, false)

	create.
		SetBorder(true).
		SetTitle("Create").
		SetTitleAlign(tview.AlignLeft).
		SetFocusFunc(func() {
			create.SetBorderColor(Focused)
			create.SetTitleColor(Focused)
		}).
		SetBlurFunc(func() {
			create.SetBorderColor(BoxBorder)
			create.SetTitleColor(BoxBorder)
		}).
		SetBackgroundColor(Background).
		SetBorderColor(BoxBorder).
		SetTitleColor(BoxBorder)

	goToLibraryButton := NewButton("Go To Library")

	goToLibrary := NewPaddedFrameAllSides(4).SetPrimitive(goToLibraryButton)

	goToLibrary.
		SetBorder(true).
		SetFocusFunc(func() {
			goToLibrary.SetBorderColor(Focused)
			goToLibrary.SetTitleColor(Focused)
		}).
		SetBlurFunc(func() {
			goToLibrary.SetBorderColor(BoxBorder)
			goToLibrary.SetTitleColor(BoxBorder)
		}).
		SetBackgroundColor(Background).
		SetBorderColor(BoxBorder).
		SetTitleColor(BoxBorder)

	main := tview.NewGrid().
		SetRows(-1, -1, -1).
		SetColumns(-1, -1).
		AddItem(recentSetsStudies, 0, 0, 1, 2, 0, 0, true).
		AddItem(folders, 1, 0, 1, 1, 0, 0, false).
		AddItem(create, 1, 1, 1, 1, 0, 0, false).
		AddItem(goToLibrary, 2, 0, 1, 2, 0, 0, false)

	recentSetsStudies.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyTab:
			appState.App.SetFocus(folders)
		}

		return event
	})

	folders.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyTab:
			appState.App.SetFocus(create)
			return nil
		case tcell.KeyBacktab:
			appState.App.SetFocus(recentSetsStudies)
			return nil
		}

		return event
	})

	create.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyTab:
			appState.App.SetFocus(goToLibraryButton)
		case tcell.KeyBacktab:
			appState.App.SetFocus(folders)
		case tcell.KeyUp:
			appState.App.SetFocus(createFlashcardSetButton)
		case tcell.KeyDown:
			appState.App.SetFocus(createFolderButton)
		}

		return event
	})

	goToLibraryButton.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyTab:
			appState.App.SetFocus(recentSetsStudies)
		case tcell.KeyBacktab:
			appState.App.SetFocus(create)
		}

		return event
	})

	goToLibraryButton.SetSelectedFunc(func() {
		appState.Navigation.GoToView(app.VIEW_NAMES.Library)
	})

	home.AddPage(mainPageName, main, true, true)
	home.AddPage(newFolderPageName, newFolderFormLayout, true, false)

	appState.Navigation.AddView(app.VIEW_NAMES.Home, home, true, func() {
		recentSetsStudies.Clear()

		for i, flashcardSet := range appState.RecentlyStudied {
			recentSetsStudies.SetCell(i, 0, tview.NewTableCell(flashcardSet.String()).SetExpansion(1))
		}
	})
}
