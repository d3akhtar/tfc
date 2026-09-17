package tui

import (
	"strings"
	"time"

	"github.com/d3akhtar/tfc/app"
	"github.com/d3akhtar/tfc/db/flashcard_set"
	"github.com/d3akhtar/tfc/db/folder"
	"github.com/d3akhtar/tfc/domain"
	"github.com/d3akhtar/tfc/importing"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/wizzymore/tinyfiledialogs"
)

func InitHomeUi(appState *app.State, flashcardSetRepository flashcard_set.FlashcardSetRepo, folderRepository folder.FolderRepo) {
	recentlyStudied := []*domain.FlashcardSet{}
	folders := []*domain.Folder{}

	home := tview.NewPages()

	recentSetsStudies := tview.NewTable().
		SetSelectable(true, false).
		SetSelectedFunc(func(row, _ int) {
			pos := row
			appState.SetSelectedFlashcardSet(recentlyStudied[pos])

			setFlashcards, err := flashcardSetRepository.GetAllFlashcardsForSet(appState.Context, appState.SelectedFlashcardSet())
			if err != nil {
				return
			}

			appState.SelectedFlashcardSet().Flashcards = setFlashcards

			appState.Navigation.GoToView(app.VIEW_NAMES.FlashcardSetPreview)
		})

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

	foldersTable := tview.NewTable().
		SetBorders(false).
		SetSelectable(true, true)

	foldersTable.SetSelectedFunc(func(row, _ int) {
		pos := row
		appState.SetSelectedFolder(folders[pos])
		appState.Navigation.GoToView(app.VIEW_NAMES.Folder)
	})

	foldersTable.
		SetBorder(true).
		SetTitle("Folders").
		SetTitleAlign(tview.AlignLeft).
		SetFocusFunc(func() {
			foldersTable.SetBorderColor(Focused)
			foldersTable.SetTitleColor(Focused)
		}).
		SetBlurFunc(func() {
			foldersTable.SetBorderColor(BoxBorder)
			foldersTable.SetTitleColor(BoxBorder)
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

		newFolder := &domain.Folder{
			Name:          formNewFolderName,
			FlashcardSets: []domain.FlashcardSet{},
			LastAccessed:  time.Now(),
		}

		err := folderRepository.Create(appState.Context, newFolder)
		if err != nil {
			return
		}

		folders = append(folders, newFolder)

		newFolderNameInputField.SetText("")
		row := len(folders) - 1
		foldersTable.SetCell(row, 0, tview.NewTableCell(newFolder.String()).SetExpansion(1))
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
		appState.SetSelectedFlashcardSet(nil)
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

	generalActionButtons := tview.NewGrid().
		SetRows(1, -1, 1).
		SetColumns(-1, 1, -1)

	SetBorderFocusAndBlurCallbacks(generalActionButtons.Box)

	goToLibraryButton := NewButton("Go To Library")
	importButton := NewButton("Import")

	generalActionButtons.
		AddItem(goToLibraryButton, 1, 0, 1, 1, 0, 0, true).
		AddItem(importButton, 1, 2, 1, 1, 0, 0, true)

	main := tview.NewGrid().
		SetRows(-2, -2, -1).
		SetColumns(-1, -1).
		AddItem(recentSetsStudies, 0, 0, 1, 2, 0, 0, true).
		AddItem(foldersTable, 1, 0, 1, 1, 0, 0, false).
		AddItem(create, 1, 1, 1, 1, 0, 0, false).
		AddItem(generalActionButtons, 2, 0, 1, 2, 0, 0, false)

	recentSetsStudies.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyTab:
			appState.App.SetFocus(foldersTable)
		}

		return event
	})

	foldersTable.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
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
			appState.App.SetFocus(generalActionButtons)
		case tcell.KeyBacktab:
			appState.App.SetFocus(foldersTable)
		case tcell.KeyUp:
			appState.App.SetFocus(createFlashcardSetButton)
		case tcell.KeyDown:
			appState.App.SetFocus(createFolderButton)
		}

		return event
	})

	generalActionButtons.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyTab:
			appState.App.SetFocus(recentSetsStudies)
		case tcell.KeyBacktab:
			appState.App.SetFocus(create)
		case tcell.KeyLeft:
			appState.App.SetFocus(goToLibraryButton)
		case tcell.KeyRight:
			appState.App.SetFocus(importButton)
		}

		return event
	})

	goToLibraryButton.SetSelectedFunc(func() {
		appState.Navigation.GoToView(app.VIEW_NAMES.Library)
	})

	home.AddPage(mainPageName, main, true, true)
	home.AddPage(newFolderPageName, newFolderFormLayout, true, false)

	refresh := func() error {
		list, err := flashcardSetRepository.List(appState.Context, 0, 50)
		if err != nil {
			return err
		}

		recentlyStudied = list

		folders, err = folderRepository.List(appState.Context, 0, 50)
		if err != nil {
			return err
		}

		recentSetsStudies.Clear()

		for i, flashcardSet := range recentlyStudied {
			recentSetsStudies.SetCell(i, 0, tview.NewTableCell(flashcardSet.String()).SetExpansion(1))
		}

		foldersTable.Clear()

		for i, loadedFolder := range folders {
			foldersTable.SetCell(i, 0, tview.NewTableCell(loadedFolder.String()).SetExpansion(1).SetTextColor(Text))
		}

		recentSetsStudies.Select(0, 0)
		foldersTable.Select(0, 0)

		return nil
	}

	importButton.SetSelectedFunc(func() {
		path, ok := tinyfiledialogs.OpenFileDialog("Select flashcard set or folder", "", []string{"*.tfcfs", "*.tfcf"}, "tfc files (*.tfcfs, *.tfcf)", false)
		if !ok {
			return
		}

		if strings.HasSuffix(path, ".tfcfs") {
			fc, err := importing.ImportFlashcardSet(path)
			if err != nil {
				return
			}

			fc.LastAccessed = time.Now()

			flashcardSetRepository.Create(appState.Context, fc)
		} else if strings.HasSuffix(path, ".tfcf") {
			f, err := importing.ImportFolder(path)
			if err != nil {
				return
			}

			for i, fc := range f.FlashcardSets {
				f.FlashcardSets[i].LastAccessed = time.Now()
				err = flashcardSetRepository.Create(appState.Context, &fc)
				if err != nil {
					return
				}

				f.FlashcardSets[i].Id = fc.Id
			}

			f.LastAccessed = time.Now()

			folderRepository.Create(appState.Context, f)
		}

		refresh()
	})

	appState.Navigation.AddView(app.VIEW_NAMES.Home, home, true, refresh, nil)
}
