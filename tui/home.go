package tui

import (
	"fmt"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

var (
	mainPageName      = "main"
	newFolderPageName = "folder"
)

func InitHomeUi(app *tview.Application, pages *tview.Pages) tview.Primitive {
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

	for column := range 3 {
		for row := range 40 {
			folders.SetCell(row, column, tview.NewTableCell(fmt.Sprintf("○ Folder %d,%d", row, column)).SetExpansion(1))
		}
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

	newFolderForm := tview.NewForm().
		AddInputField("Name", "", 0, nil, nil).
		AddButton("Save", func() {

		}).
		AddButton("Quit", func() {
			home.HidePage(newFolderPageName)
			app.SetFocus(createFlashcardButton)
		}).
		SetButtonsAlign(tview.AlignCenter).
		SetSubmitFunc(func() {
			home.HidePage(newFolderPageName)
			app.SetFocus(createFlashcardButton)
		})

	newFolderForm.
		SetBorder(true).
		SetTitle("Enter new folder name").
		SetTitleAlign(tview.AlignLeft)

	newFolderFormLayout := NewModal(newFolderForm, 100, 7)

	createFlashcardButton.SetSelectedFunc(func() {
		app.SetFocus(createFolderButton)
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
		pages.SwitchToPage(PAGE_NAMES.Library)
	})

	home.AddPage(mainPageName, main, true, true)
	home.AddPage(newFolderPageName, newFolderFormLayout, true, false)

	return home
}
