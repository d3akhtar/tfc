package tui

import (
	"fmt"
	"strconv"

	"github.com/d3akhtar/tfc/app"
	"github.com/d3akhtar/tfc/db"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func InitFlashcardSetPreview(appState *app.State) {
	maxFlashcardsShownInPreviewFlashcardList := 4

	lastSelectedFlashcardPrimitive := 0
	activeFlashcardPrimitives := make([]tview.Primitive, maxFlashcardsShownInPreviewFlashcardList)

	flashcardList := tview.NewFlex().
		SetDirection(tview.FlexRow)

	newFlashcardPrimitive := func(flashcard db.Flashcard, pos int) tview.Primitive {
		layout := tview.NewFlex()

		activeFlashcardPrimitives[pos] = layout

		layout.
			SetBorder(true).
			SetBorderPadding(1, 1, 1, 1).
			SetFocusFunc(func() {
				lastSelectedFlashcardPrimitive = pos
				layout.SetBorderColor(tcell.ColorGreen)
				layout.SetTitleColor(tcell.ColorGreen)
			}).
			SetBlurFunc(func() {
				layout.SetBorderColor(tcell.ColorWhite)
				layout.SetTitleColor(tcell.ColorWhite)
			}).
			SetTitle(strconv.Itoa((pos + 1))).
			SetTitleAlign(tview.AlignLeft)

		question := tview.NewTextView().
			SetText(flashcard.Question)

		answer := tview.NewTextView().
			SetText(flashcard.Answer)

		question.
			SetBorder(true).
			SetBorderPadding(1, 1, 1, 1).
			SetFocusFunc(func() {
				question.SetBorderColor(tcell.ColorGreen)
				question.SetTitleColor(tcell.ColorGreen)
			}).
			SetBlurFunc(func() {
				question.SetBorderColor(tcell.ColorWhite)
				question.SetTitleColor(tcell.ColorWhite)
			})

		answer.
			SetBorder(true).
			SetBorderPadding(1, 1, 1, 1).
			SetFocusFunc(func() {
				answer.SetBorderColor(tcell.ColorGreen)
				answer.SetTitleColor(tcell.ColorGreen)
			}).
			SetBlurFunc(func() {
				answer.SetBorderColor(tcell.ColorWhite)
				answer.SetTitleColor(tcell.ColorWhite)
			})

		layout.AddItem(question, 0, 1, false)
		layout.AddItem(answer, 0, 1, false)

		layout.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
			if question.HasFocus() || answer.HasFocus() {
				return event
			}

			switch event.Key() {
			case tcell.KeyUp:
				if pos > 0 {
					appState.SetFocus(activeFlashcardPrimitives[pos-1])
				}

				return nil
			case tcell.KeyDown:
				if pos < maxFlashcardsShownInPreviewFlashcardList-1 {
					appState.SetFocus(activeFlashcardPrimitives[pos+1])
				}

				return nil
			case tcell.KeyEnter:
				appState.SetFocus(question)
				return nil
			case tcell.KeyRune:
				if event.Rune() == 'q' {
					appState.SetFocus(flashcardList)
				}

				return nil
			}

			return event
		})

		question.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
			switch event.Key() {
			case tcell.KeyLeft, tcell.KeyRight:
				appState.SetFocus(answer)
				return nil
			case tcell.KeyRune:
				if event.Rune() == 'q' {
					appState.SetFocus(layout)
				}
			}

			return event
		})

		answer.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
			switch event.Key() {
			case tcell.KeyLeft, tcell.KeyRight:
				appState.SetFocus(question)
				return nil
			case tcell.KeyRune:
				if event.Rune() == 'q' {
					appState.SetFocus(layout)
				}
			}

			return event
		})

		return layout
	}

	flashcardSetPreview := tview.NewPages()

	preview := tview.NewGrid().
		SetRows(-1, -30, -3)

	flashcardSetNameLabel := tview.NewTextView().
		SetText(fmt.Sprintf("[ %s ]", appState.SelectedFlashcardSet.Name)).
		SetScrollable(false).
		SetSize(2, 30)

	flashcardList.
		SetTitle("Flashcards").
		SetBorder(true).
		SetFocusFunc(func() {
			flashcardList.SetBorderColor(tcell.ColorGreen)
			flashcardList.SetTitleColor(tcell.ColorGreen)
		}).
		SetBlurFunc(func() {
			flashcardList.SetBorderColor(tcell.ColorWhite)
			flashcardList.SetTitleColor(tcell.ColorWhite)
		}).
		SetBorderColor(tcell.ColorGreen).
		SetTitleColor(tcell.ColorGreen)

	trackProgressButton := tview.NewButton("Track Progress")

	trackProgressButton.
		SetSelectedFunc(func() {
			appState.SelectedFlashcardSet.TrackProgress = !appState.SelectedFlashcardSet.TrackProgress
			var col tcell.Color
			var trackProgressButtonTextPrefix string
			if appState.SelectedFlashcardSet.TrackProgress {
				col = tcell.ColorGreen
				trackProgressButtonTextPrefix = "✔"
			} else {
				col = tcell.ColorRed
				trackProgressButtonTextPrefix = "✖"
			}

			trackProgressButton.SetStyle(tcell.StyleDefault.
				Background(col).
				Foreground(tcell.ColorBlack)).
				SetLabel(trackProgressButtonTextPrefix + " Track Progress")
		})

	studyButton := tview.NewButton("Study")

	editButton := tview.NewButton("Edit")

	settingsButton := tview.NewButton("Settings")

	buttonGroup := tview.NewFlex().
		SetDirection(tview.FlexColumn).
		AddItem(trackProgressButton, 0, 1, false).
		AddItem(nil, 2, 0, false).
		AddItem(studyButton, 0, 3, false).
		AddItem(nil, 2, 1, false).
		AddItem(editButton, 0, 3, false).
		AddItem(nil, 2, 1, false).
		AddItem(settingsButton, 0, 1, false)

	settings := tview.NewBox().
		SetTitle("Settings")

	settingsModal := NewModal(settings, 20, 50)

	preview.
		AddItem(flashcardSetNameLabel, 0, 0, 1, 1, 0, 0, false).
		AddItem(flashcardList, 1, 0, 1, 1, 0, 0, true).
		AddItem(buttonGroup, 2, 0, 1, 1, 0, 0, false)

	flashcardSetPreview.AddPage("main", preview, true, true)
	flashcardSetPreview.AddPage("settings", settingsModal, true, false)

	flashcardList.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if activeFlashcardPrimitives[lastSelectedFlashcardPrimitive].HasFocus() {
			return event
		}

		switch event.Key() {
		case tcell.KeyEnter:
			appState.SetFocus(activeFlashcardPrimitives[lastSelectedFlashcardPrimitive])
			return nil
		case tcell.KeyTab:
			appState.SetFocus(trackProgressButton)
		}

		return event
	})

	trackProgressButton.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyTab:
			appState.SetFocus(studyButton)
			return nil
		case tcell.KeyBacktab:
			appState.SetFocus(flashcardList)
			return nil
		}

		return event
	})

	studyButton.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyTab:
			appState.SetFocus(editButton)
			return nil
		case tcell.KeyBacktab:
			appState.SetFocus(trackProgressButton)
			return nil
		}

		return event
	})

	editButton.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyTab:
			appState.SetFocus(settingsButton)
			return nil
		case tcell.KeyBacktab:
			appState.SetFocus(studyButton)
			return nil
		}

		return event
	})

	settingsButton.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyBacktab:
			appState.SetFocus(editButton)
			return nil
		}

		return event
	})

	appState.Navigation.AddView(app.VIEW_NAMES.FlashcardSetPreview, flashcardSetPreview, false, func() {
		maxFlashcardsShownInPreviewFlashcardList = min(4, len(appState.SelectedFlashcardSet.Flashcards))

		var col tcell.Color
		var trackProgressButtonTextPrefix string
		if appState.SelectedFlashcardSet.TrackProgress {
			col = tcell.ColorGreen
			trackProgressButtonTextPrefix = "✔"
		} else {
			col = tcell.ColorRed
			trackProgressButtonTextPrefix = "✖"
		}

		trackProgressButton.SetStyle(tcell.StyleDefault.
			Background(col).
			Foreground(tcell.ColorBlack)).
			SetLabel(trackProgressButtonTextPrefix + " Track Progress")

		flashcardList.Clear()

		flashcardSetNameLabel.SetText(fmt.Sprintf("[ %s ]", appState.SelectedFlashcardSet.Name))

		for i, flashcard := range appState.SelectedFlashcardSet.Flashcards {
			flashcardList.AddItem(
				newFlashcardPrimitive(flashcard, i),
				0,
				1,
				false,
			)
		}
	})
}
