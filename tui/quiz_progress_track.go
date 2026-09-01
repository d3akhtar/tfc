package tui

import (
	"fmt"
	"strconv"

	"github.com/d3akhtar/tfc/app"
	"github.com/d3akhtar/tfc/db/quiz"
	"github.com/d3akhtar/tfc/domain"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func InitQuizProgressTrackUi(appState *app.State, quizRepository quiz.QuizRepo) {
	finishedMessageFormat := "You finished the quiz with %d/%d cards learned! What do you want to do next?"

	quizProgressTrack := tview.NewPages()

	var quiz *domain.Quiz = nil
	var currentFlashcard domain.Flashcard
	showAnswer := false

	quizView := tview.NewGrid().
		SetRows(44, -1)

	contentView := tview.NewTextView().
		SetTextAlign(tview.AlignCenter)

	contentViewContainer := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(nil, 0, 1, false).
		AddItem(contentView, 0, 1, true).
		AddItem(nil, 0, 1, false)

	contentViewContainer.
		SetBorder(true).
		SetBorderColor(Focused).
		SetBorderPadding(4, 4, 4, 4)

	controls := tview.NewFlex()

	unknownCount := NewVerticallyAlignedTextView("0").
		SetScrollable(false).
		SetTextColor(tcell.ColorRed)

	unknownCount.
		SetBorder(true).
		SetBorderColor(tcell.ColorRed)

	unknownButton := NewButton("✖").
		SetStyle(tcell.StyleDefault.
			Background(tcell.ColorRed).
			Foreground(tcell.ColorBlack))

	unknownButtonFrame := tview.NewFrame(unknownButton).
		AddText("(←)", false, tview.AlignCenter, tcell.ColorWhite)

	unknownButtonFrame.
		SetBorder(true).
		SetBorderColor(tcell.ColorRed)

	undoButton := NewButton("↺").
		SetStyle(tcell.StyleDefault.
			Background(tcell.NewHexColor(0x00fff2)).
			Foreground(tcell.ColorBlack))

	undoButtonFrame := tview.NewFrame(undoButton).
		AddText("(DEL)", false, tview.AlignCenter, tcell.ColorWhite)

	undoButtonFrame.
		SetBorder(true).
		SetBorderColor(tcell.NewHexColor(0x00fff2))

	knowButton := NewButton("✔").
		SetStyle(tcell.StyleDefault.
			Background(Focused).
			Foreground(tcell.ColorBlack))

	knowButtonFrame := tview.NewFrame(knowButton).
		AddText("(→)", false, tview.AlignCenter, tcell.ColorWhite)

	knowButtonFrame.
		SetBorder(true).
		SetBorderColor(Focused)

	knowCount := NewVerticallyAlignedTextView("0").
		SetScrollable(false).
		SetTextColor(Focused)

	knowCount.
		SetBorder(true).
		SetBorderColor(Focused)

	quizFinishedMessage := tview.NewTextView().
		SetTextAlign(tview.AlignCenter)

	studyRemainingCardsButton := NewButton("Study Remaining Cards").
		SetSelectedFunc(func() {
			appState.SelectedFlashcardSet.StartQuiz()
			quiz = appState.SelectedFlashcardSet.Quiz
			showAnswer = false || appState.SelectedFlashcardSet.Front == domain.Answer

			contentView.SetText(quiz.CurrentlySelectedCard().Question)
			contentViewContainer.SetTitle(fmt.Sprintf("%d / %d", 1, len(quiz.Flashcards)))

			knowCount.SetText("0")
			unknownCount.SetText("0")

			quizProgressTrack.SwitchToPage("main")

			currentFlashcard = quiz.CurrentlySelectedCard()
		})

	resetProgressButton := NewButton("Reset Progress").
		SetSelectedFunc(func() {
			appState.SelectedFlashcardSet.ResetQuizProgress()
			quiz = appState.SelectedFlashcardSet.Quiz
			showAnswer = false || appState.SelectedFlashcardSet.Front == domain.Answer

			contentView.SetText(quiz.CurrentlySelectedCard().Question)
			contentViewContainer.SetTitle(fmt.Sprintf("%d / %d", 1, len(quiz.Flashcards)))

			knowCount.SetText("0")
			unknownCount.SetText("0")

			quizProgressTrack.SwitchToPage("main")

			currentFlashcard = quiz.CurrentlySelectedCard()
		})

	goHomeButton := NewButton("Home").
		SetSelectedFunc(func() {
			appState.Navigation.GoToView(app.VIEW_NAMES.Home)
		})

	contentView.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyLeft:
			if quiz.AreCardsLeft() {
				quiz.GoToNextCard(false)
				showAnswer = false || appState.SelectedFlashcardSet.Front == domain.Answer

				if quiz.Finished() {
					nKnown, _ := quiz.GetKnownAndUnknownCount()
					quizFinishedMessage.SetText(fmt.Sprintf(finishedMessageFormat, nKnown, len(quiz.Flashcards)))
					quizProgressTrack.SwitchToPage("finished")
					return event
				}

				currentFlashcard = quiz.CurrentlySelectedCard()
			}

		case tcell.KeyRight:
			if quiz.AreCardsLeft() {
				quiz.GoToNextCard(true)
				showAnswer = false || appState.SelectedFlashcardSet.Front == domain.Answer

				if quiz.Finished() {
					nKnown, _ := quiz.GetKnownAndUnknownCount()
					quizFinishedMessage.SetText(fmt.Sprintf(finishedMessageFormat, nKnown, len(quiz.Flashcards)))
					quizProgressTrack.SwitchToPage("finished")

					if nKnown == len(quiz.Flashcards) {
						studyRemainingCardsButton.SetDisabled(true)
						appState.SetFocus(resetProgressButton)
					}

					return event
				}

				currentFlashcard = quiz.CurrentlySelectedCard()
			}

		case tcell.KeyDelete, tcell.KeyBackspace, tcell.KeyDEL:
			if quiz.CanUndo() {
				showAnswer = false || appState.SelectedFlashcardSet.Front == domain.Answer
				currentFlashcard = quiz.Undo()
			}

		case tcell.KeyRune:
			if event.Rune() == ' ' {
				showAnswer = !showAnswer
			}
		}

		var contents string
		if showAnswer {
			contents = currentFlashcard.Answer
		} else {
			contents = currentFlashcard.Question
		}

		contentView.SetText(contents)
		contentViewContainer.SetTitle(fmt.Sprintf("%d / %d", quiz.CurrentlySelectedIndex+1, len(quiz.Flashcards)))

		known, unknown := quiz.GetKnownAndUnknownCount()
		knowCount.SetText(strconv.Itoa(known))
		unknownCount.SetText(strconv.Itoa(unknown))

		return event
	})

	unknownButton.SetSelectedFunc(func() {
		if quiz.AreCardsLeft() {
			quiz.GoToNextCard(false)
			showAnswer = false || appState.SelectedFlashcardSet.Front == domain.Answer

			if quiz.Finished() {
				nKnown, _ := quiz.GetKnownAndUnknownCount()
				quizFinishedMessage.SetText(fmt.Sprintf(finishedMessageFormat, nKnown, len(quiz.Flashcards)))
				quizProgressTrack.SwitchToPage("finished")
				return
			}

			currentFlashcard = quiz.CurrentlySelectedCard()
		}

		var contents string
		if showAnswer {
			contents = currentFlashcard.Answer
		} else {
			contents = currentFlashcard.Question
		}

		contentView.SetText(contents)
		contentViewContainer.SetTitle(fmt.Sprintf("%d / %d", quiz.CurrentlySelectedIndex+1, len(quiz.Flashcards)))

		known, unknown := quiz.GetKnownAndUnknownCount()
		knowCount.SetText(strconv.Itoa(known))
		unknownCount.SetText(strconv.Itoa(unknown))
	})

	undoButton.SetSelectedFunc(func() {
		if quiz.CanUndo() {
			showAnswer = false || appState.SelectedFlashcardSet.Front == domain.Answer
			currentFlashcard = quiz.Undo()
		}

		var contents string
		if showAnswer {
			contents = currentFlashcard.Answer
		} else {
			contents = currentFlashcard.Question
		}

		contentView.SetText(contents)
		contentViewContainer.SetTitle(fmt.Sprintf("%d / %d", quiz.CurrentlySelectedIndex+1, len(quiz.Flashcards)))

		known, unknown := quiz.GetKnownAndUnknownCount()
		knowCount.SetText(strconv.Itoa(known))
		unknownCount.SetText(strconv.Itoa(unknown))
	})

	knowButton.SetSelectedFunc(func() {
		if quiz.AreCardsLeft() {
			quiz.GoToNextCard(true)
			showAnswer = false || appState.SelectedFlashcardSet.Front == domain.Answer

			if quiz.Finished() {
				nKnown, _ := quiz.GetKnownAndUnknownCount()
				quizFinishedMessage.SetText(fmt.Sprintf(finishedMessageFormat, nKnown, len(quiz.Flashcards)))
				quizProgressTrack.SwitchToPage("finished")

				if nKnown == len(quiz.Flashcards) {
					studyRemainingCardsButton.SetDisabled(true)
					appState.SetFocus(resetProgressButton)
				}

				return
			}

			currentFlashcard = quiz.CurrentlySelectedCard()
		}

		var contents string
		if showAnswer {
			contents = currentFlashcard.Answer
		} else {
			contents = currentFlashcard.Question
		}

		contentView.SetText(contents)
		contentViewContainer.SetTitle(fmt.Sprintf("%d / %d", quiz.CurrentlySelectedIndex+1, len(quiz.Flashcards)))

		known, unknown := quiz.GetKnownAndUnknownCount()
		knowCount.SetText(strconv.Itoa(known))
		unknownCount.SetText(strconv.Itoa(unknown))
	})

	controls.
		AddItem(unknownCount, 0, 1, false).
		AddItem(nil, 25, 0, false).
		AddItem(unknownButtonFrame, 0, 2, false).
		AddItem(undoButtonFrame, 0, 2, false).
		AddItem(knowButtonFrame, 0, 2, false).
		AddItem(nil, 25, 0, false).
		AddItem(knowCount, 0, 1, false)

	quizView.
		AddItem(contentViewContainer, 0, 0, 1, 1, 0, 0, true).
		AddItem(controls, 1, 0, 1, 1, 0, 0, false)

	studyRemainingCardsButton.
		SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
			switch event.Key() {
			case tcell.KeyTab, tcell.KeyDown:
				appState.SetFocus(resetProgressButton)
				return nil
			}

			return event
		})

	resetProgressButton.
		SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
			switch event.Key() {
			case tcell.KeyBacktab, tcell.KeyUp:
				appState.SetFocus(studyRemainingCardsButton)
				return nil
			case tcell.KeyTab, tcell.KeyDown:
				appState.SetFocus(goHomeButton)
				return nil
			}

			return event
		})

	goHomeButton.
		SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
			switch event.Key() {
			case tcell.KeyBacktab, tcell.KeyUp:
				appState.SetFocus(resetProgressButton)
				return nil
			}

			return event
		})

	finishedQuizOptions := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(quizFinishedMessage, 5, 0, false).
		AddItem(studyRemainingCardsButton, 0, 1, true).
		AddItem(nil, 5, 0, false).
		AddItem(resetProgressButton, 0, 1, false).
		AddItem(nil, 5, 0, false).
		AddItem(goHomeButton, 0, 1, false)

	finishedQuizOptions.
		SetBorderPadding(5, 5, 5, 5)

	quizProgressTrack.
		AddPage("main", quizView, true, true).
		AddPage("finished", finishedQuizOptions, true, false)

	quizProgressTrack.
		SetChangedFunc(func() {
			if pg, _ := quizProgressTrack.GetFrontPage(); pg == "main" {
				studyRemainingCardsButton.SetDisabled(false)
			}
		})

	appState.Navigation.AddView(app.VIEW_NAMES.QuizProgressTrack, quizProgressTrack, false, func() {
		studyRemainingCardsButton.SetDisabled(false)

		if quiz == nil || quiz != appState.SelectedFlashcardSet.Quiz {
			if appState.SelectedFlashcardSet.Quiz == nil {
				appState.SelectedFlashcardSet.StartQuiz()
			}

			quiz = appState.SelectedFlashcardSet.Quiz
			showAnswer = false || appState.SelectedFlashcardSet.Front == domain.Answer
		}

		if quiz.Finished() {
			nKnown, _ := quiz.GetKnownAndUnknownCount()
			quizFinishedMessage.SetText(fmt.Sprintf(finishedMessageFormat, nKnown, len(quiz.Flashcards)))
			quizProgressTrack.SwitchToPage("finished")

			appState.SetFocus(studyRemainingCardsButton)

			if nKnown == len(quiz.Flashcards) {
				studyRemainingCardsButton.SetDisabled(true)
				appState.SetFocus(resetProgressButton)
			}

			return
		}

		quizProgressTrack.SwitchToPage("main")

		if showAnswer {
			contentView.SetText(quiz.CurrentlySelectedCard().Answer)
		} else {
			contentView.SetText(quiz.CurrentlySelectedCard().Question)
		}

		currentFlashcard = quiz.CurrentlySelectedCard()

		contentViewContainer.SetTitle(fmt.Sprintf("%d / %d", quiz.CurrentlySelectedIndex+1, len(quiz.Flashcards)))

		known, unknown := quiz.GetKnownAndUnknownCount()
		knowCount.SetText(strconv.Itoa(known))
		unknownCount.SetText(strconv.Itoa(unknown))
	})
}
