package tui

import (
	"fmt"

	"github.com/d3akhtar/tfc/app"
	"github.com/d3akhtar/tfc/db"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

var (
	mainPageName      = "main"
	newFolderPageName = "folder"
)

var formNewFolderName = ""

func InitHomeUi(app *tview.Application, pages *tview.Pages, appState *app.State) tview.Primitive {
	home := tview.NewPages()

	recentSetsStudies := tview.NewTable().
		SetSelectable(true, false)

	recentSetsStudies.
		SetFocusFunc(func() {
			recentSetsStudies.SetBorderColor(tcell.ColorGreen)
			recentSetsStudies.SetTitleColor(tcell.ColorGreen)
		}).
		SetBlurFunc(func() {
			recentSetsStudies.SetBorderColor(tcell.ColorWhite)
			recentSetsStudies.SetTitleColor(tcell.ColorWhite)
		}).
		SetBorder(true).
		SetTitle("Recents Sets Studies").
		SetTitleAlign(tview.AlignLeft).
		SetBorderColor(tcell.ColorGreen).
		SetTitleColor(tcell.ColorGreen).
		SetBorderPadding(1, 1, 1, 1)

	for i := range 30 {
		recentSetsStudies.SetCell(i, 0, tview.NewTableCell(fmt.Sprintf("○ Set %d", i)).SetExpansion(1))
	}

	folders := tview.NewTable().
		SetBorders(false).
		SetSelectable(true, true)

	folders.SetSelectedFunc(func(row, column int) {
		pos := (row * 3) + column
		appState.SetSelectedFolder(appState.Folders[pos])
		pages.ShowPage(VIEW_NAMES.Folder)
		pages.HidePage(VIEW_NAMES.Library)
	})

	for i, loadedFolder := range appState.Folders {
		row := i / 3
		column := i % 3
		folders.SetCell(row, column, tview.NewTableCell(fmt.Sprintf("○ %s", loadedFolder.Name)).SetExpansion(1))
	}

	folders.
		SetBorder(true).
		SetTitle("Folders").
		SetTitleAlign(tview.AlignLeft).
		SetFocusFunc(func() {
			folders.SetBorderColor(tcell.ColorGreen)
			folders.SetTitleColor(tcell.ColorGreen)
		}).
		SetBlurFunc(func() {
			folders.SetBorderColor(tcell.ColorWhite)
			folders.SetTitleColor(tcell.ColorWhite)
		}).
		SetBorderPadding(1, 1, 2, 0)

	createFlashcardButton := tview.NewButton("Flashcard")
	createFolderButton := tview.NewButton("Folder")

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
		app.SetFocus(createFlashcardButton)

		appState.Folders = append(appState.Folders, db.Folder{
			Name:        formNewFolderName,
			Collections: []db.Collection{},
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
			app.SetFocus(createFlashcardButton)
		}).
		SetButtonsAlign(tview.AlignCenter).
		SetSubmitFunc(onFolderFormSubmit)

	newFolderForm.
		SetBorder(true).
		SetTitle("Enter new folder name").
		SetTitleAlign(tview.AlignLeft)

	newFolderFormLayout := NewModal(newFolderForm, 100, 7)

	createFlashcardButton.SetSelectedFunc(func() {
		pages.ShowPage(VIEW_NAMES.FlashcardEdit)
		pages.HidePage(VIEW_NAMES.Home)
	})

	createFolderButton.SetSelectedFunc(func() {
		home.ShowPage(newFolderPageName)
		app.SetFocus(newFolderForm)
	})

	create := tview.NewGrid().
		SetRows(-1, -1).
		AddItem(NewPaddedFrameXY(2, 1).SetPrimitive(createFlashcardButton), 0, 0, 1, 1, 0, 0, true).
		AddItem(NewPaddedFrameXY(2, 1).SetPrimitive(createFolderButton), 1, 0, 1, 1, 0, 0, false)

	create.
		SetBorder(true).
		SetTitle("Create").
		SetTitleAlign(tview.AlignLeft).
		SetFocusFunc(func() {
			create.SetBorderColor(tcell.ColorGreen)
			create.SetTitleColor(tcell.ColorGreen)
		}).
		SetBlurFunc(func() {
			create.SetBorderColor(tcell.ColorWhite)
			create.SetTitleColor(tcell.ColorWhite)
		})

	goToLibraryButton := tview.NewButton("Go To Library")
	goToLibrary := NewPaddedFrameAllSides(4).SetPrimitive(goToLibraryButton)

	goToLibrary.
		SetBorder(true).
		SetFocusFunc(func() {
			goToLibrary.SetBorderColor(tcell.ColorGreen)
			goToLibrary.SetTitleColor(tcell.ColorGreen)
		}).
		SetBlurFunc(func() {
			goToLibrary.SetBorderColor(tcell.ColorWhite)
			goToLibrary.SetTitleColor(tcell.ColorWhite)
		})

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
			app.SetFocus(folders)
		}

		return event
	})

	folders.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyTab:
			app.SetFocus(create)
			return nil
		case tcell.KeyBacktab:
			app.SetFocus(recentSetsStudies)
			return nil
		}

		return event
	})

	create.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyTab:
			app.SetFocus(goToLibraryButton)
		case tcell.KeyBacktab:
			app.SetFocus(folders)
		case tcell.KeyUp:
			app.SetFocus(createFlashcardButton)
		case tcell.KeyDown:
			app.SetFocus(createFolderButton)
		}

		return event
	})

	goToLibraryButton.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyTab:
			app.SetFocus(recentSetsStudies)
		case tcell.KeyBacktab:
			app.SetFocus(create)
		}

		return event
	})

	goToLibraryButton.SetSelectedFunc(func() {
		pages.SwitchToPage(VIEW_NAMES.Library)
	})

	home.AddPage(mainPageName, main, true, true)
	home.AddPage(newFolderPageName, newFolderFormLayout, true, false)

	return home
}
