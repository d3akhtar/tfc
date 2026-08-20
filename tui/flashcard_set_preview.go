package tui

import (
	"fmt"
	"strconv"

	"github.com/d3akhtar/tfc/app"
	"github.com/d3akhtar/tfc/db"
	"github.com/d3akhtar/tfc/utils"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type flashcardPrimitiveInfo struct {
	Layout   *tview.Flex
	Question *tview.TextView
	Answer   *tview.TextView
}

func InitFlashcardSetPreview(appState *app.State) {
	var window *utils.SlidingWindow[db.Flashcard]

	maxFlashcardsShownInPreviewFlashcardList := 4

	lastSelectedFlashcardPrimitive := 0
	activeFlashcardPrimitives := make([]flashcardPrimitiveInfo, maxFlashcardsShownInPreviewFlashcardList)

	flashcardList := tview.NewFlex().
		SetDirection(tview.FlexRow)

	newFlashcardPrimitive := func(flashcard db.Flashcard, pos int) tview.Primitive {
		layout := tview.NewFlex()

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
				if pos == 0 {
					if window.CanRetreat() {
						idx := 0
						window.Retreat()
						for i := window.Start; i <= window.End; i++ {
							flashcard := window.Collection[i]
							activeFlashcardPrimitives[idx].Layout.SetTitle(strconv.Itoa(i + 1))
							activeFlashcardPrimitives[idx].Question.SetText(flashcard.Question)
							activeFlashcardPrimitives[idx].Answer.SetText(flashcard.Answer)
							idx++
						}
					}
				} else {
					appState.SetFocus(activeFlashcardPrimitives[pos-1].Layout)
				}

				return nil
			case tcell.KeyDown:
				if pos == maxFlashcardsShownInPreviewFlashcardList-1 {
					if window.CanAdvance() {
						idx := 0
						window.Advance()
						for i := window.Start; i <= window.End; i++ {
							flashcard := window.Collection[i]
							activeFlashcardPrimitives[idx].Layout.SetTitle(strconv.Itoa(i + 1))
							activeFlashcardPrimitives[idx].Question.SetText(flashcard.Question)
							activeFlashcardPrimitives[idx].Answer.SetText(flashcard.Answer)
							idx++
						}
					}
				} else {
					appState.SetFocus(activeFlashcardPrimitives[pos+1].Layout)
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

		activeFlashcardPrimitives[pos] = flashcardPrimitiveInfo{
			Layout:   layout,
			Question: question,
			Answer:   answer,
		}

		return layout
	}

	flashcardSetPreview := tview.NewPages()

	preview := tview.NewGrid().
		SetRows(-1, -30, -3)

	flashcardSetNameLabel := tview.NewTextView().
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

	studyButton := tview.NewButton("Study").
		SetSelectedFunc(func() {
			if appState.SelectedFlashcardSet.TrackProgress {
				appState.Navigation.GoToView(app.VIEW_NAMES.QuizProgressTrack)
			} else {
				appState.Navigation.GoToView(app.VIEW_NAMES.QuizNormal)
			}
		})

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

	shuffleCheckbox := tview.NewCheckbox().
		SetLabel("Shuffle").
		SetChangedFunc(func(checked bool) {
			appState.SelectedFlashcardSet.SetShuffle(checked)
			appState.SelectedFlashcardSet.Quiz = nil

			window = utils.NewSlidingWindow(0, maxFlashcardsShownInPreviewFlashcardList, appState.SelectedFlashcardSet.GetFlashcards())

			for i := window.Start; i <= window.End; i++ {
				activeFlashcardPrimitives[i].Layout.
					SetTitle(strconv.Itoa(i + 1))

				activeFlashcardPrimitives[i].Question.SetText(window.Collection[i].Question)
				activeFlashcardPrimitives[i].Answer.SetText(window.Collection[i].Answer)
			}

			lastSelectedFlashcardPrimitive = 0
		}).
		SetCheckedString("✔")

	frontDropdown := tview.NewDropDown().
		SetLabel("Front").
		SetOptions([]string{"Question", "Answer"}, func(_ string, option int) {
			appState.SelectedFlashcardSet.Front = db.FlashcardFront(option)
		})

	settings := tview.NewForm().
		SetButtonsAlign(tview.AlignCenter).
		AddFormItem(shuffleCheckbox).
		AddFormItem(frontDropdown).
		AddButton("Reset Progress", func() {
			appState.SelectedFlashcardSet.ResetQuizProgress()
			flashcardSetPreview.HidePage("settings")
			appState.SetFocus(settingsButton)
		}).
		AddButton("Close", func() {
			flashcardSetPreview.HidePage("settings")
			appState.SetFocus(settingsButton)
		})

	settings.
		SetBorder(true).
		SetBorderPadding(1, 1, 3, 3)

	settingsModal := NewModal(settings, 40, 10)

	preview.
		AddItem(flashcardSetNameLabel, 0, 0, 1, 1, 0, 0, false).
		AddItem(flashcardList, 1, 0, 1, 1, 0, 0, true).
		AddItem(buttonGroup, 2, 0, 1, 1, 0, 0, false)

	flashcardSetPreview.AddPage("main", preview, true, true)
	flashcardSetPreview.AddPage("settings", settingsModal, true, false)

	settingsButton.SetSelectedFunc(func() {
		flashcardSetPreview.ShowPage("settings")
	})

	flashcardList.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if activeFlashcardPrimitives[lastSelectedFlashcardPrimitive].Layout.HasFocus() {
			return event
		}

		switch event.Key() {
		case tcell.KeyEnter:
			appState.SetFocus(activeFlashcardPrimitives[lastSelectedFlashcardPrimitive].Layout)
			return nil
		case tcell.KeyTab:
			appState.SetFocus(trackProgressButton)
		}

		return event
	})

	trackProgressButton.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyTab, tcell.KeyRight:
			appState.SetFocus(studyButton)
			return nil
		case tcell.KeyBacktab, tcell.KeyLeft:
			appState.SetFocus(flashcardList)
			return nil
		}

		return event
	})

	studyButton.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyTab, tcell.KeyRight:
			appState.SetFocus(editButton)
			return nil
		case tcell.KeyBacktab, tcell.KeyLeft:
			appState.SetFocus(trackProgressButton)
			return nil
		}

		return event
	})

	editButton.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyTab, tcell.KeyRight:
			appState.SetFocus(settingsButton)
			return nil
		case tcell.KeyBacktab, tcell.KeyLeft:
			appState.SetFocus(studyButton)
			return nil
		}

		return event
	})

	settingsButton.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyBacktab, tcell.KeyLeft:
			appState.SetFocus(editButton)
			return nil
		}

		return event
	})

	appState.Navigation.AddView(app.VIEW_NAMES.FlashcardSetPreview, flashcardSetPreview, false, func() {
		maxFlashcardsShownInPreviewFlashcardList = min(4, len(appState.SelectedFlashcardSet.Flashcards))

		window = utils.NewSlidingWindow(0, maxFlashcardsShownInPreviewFlashcardList, appState.SelectedFlashcardSet.GetFlashcards())

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

		for i := window.Start; i <= window.End; i++ {
			flashcardList.AddItem(
				newFlashcardPrimitive(window.Collection[i], i),
				0,
				1,
				false,
			)
		}

		shuffleCheckbox.SetChecked(appState.SelectedFlashcardSet.Shuffled())
		frontDropdown.SetCurrentOption(int(appState.SelectedFlashcardSet.Front))
	})
}
