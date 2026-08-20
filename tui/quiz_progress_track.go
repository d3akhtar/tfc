package tui

import (
	"fmt"
	"strconv"

	"github.com/d3akhtar/tfc/app"
	"github.com/d3akhtar/tfc/db"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func InitQuizProgressTrackUi(appState *app.State) {
	var flashcardQuiz *db.Quiz = nil
	var currentFlashcard db.Flashcard
	showAnswer := false

	quiz := tview.NewGrid().
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
		SetBorderColor(tcell.ColorGreen).
		SetBorderPadding(4, 4, 4, 4)

	controls := tview.NewFlex()

	unknownCount := NewVerticallyAlignedTextView("0").
		SetScrollable(false).
		SetTextColor(tcell.ColorRed)

	unknownCount.
		SetBorder(true).
		SetBorderColor(tcell.ColorRed)

	unknownButton := tview.NewButton("✖").
		SetStyle(tcell.StyleDefault.
			Background(tcell.ColorRed).
			Foreground(tcell.ColorBlack))

	dontKnowButtonFrame := tview.NewFrame(unknownButton).
		AddText("(←)", false, tview.AlignCenter, tcell.ColorWhite)

	dontKnowButtonFrame.
		SetBorder(true).
		SetBorderColor(tcell.ColorRed)

	undoButton := tview.NewButton("↺").
		SetStyle(tcell.StyleDefault.
			Background(tcell.NewHexColor(0x00fff2)).
			Foreground(tcell.ColorBlack))

	undoButtonFrame := tview.NewFrame(undoButton).
		AddText("(DEL)", false, tview.AlignCenter, tcell.ColorWhite)

	undoButtonFrame.
		SetBorder(true).
		SetBorderColor(tcell.NewHexColor(0x00fff2))

	knowButton := tview.NewButton("✔").
		SetStyle(tcell.StyleDefault.
			Background(tcell.ColorGreen).
			Foreground(tcell.ColorBlack))

	knowButtonFrame := tview.NewFrame(knowButton).
		AddText("(→)", false, tview.AlignCenter, tcell.ColorWhite)

	knowButtonFrame.
		SetBorder(true).
		SetBorderColor(tcell.ColorGreen)

	knowCount := NewVerticallyAlignedTextView("0").
		SetScrollable(false).
		SetTextColor(tcell.ColorGreen)

	knowCount.
		SetBorder(true).
		SetBorderColor(tcell.ColorGreen)

	contentView.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyLeft:
			if flashcardQuiz.AreCardsLeft() {
				showAnswer = false
				currentFlashcard = flashcardQuiz.NextCard(false)
			}

		case tcell.KeyRight:
			if flashcardQuiz.AreCardsLeft() {
				showAnswer = false
				currentFlashcard = flashcardQuiz.NextCard(true)
			}

		case tcell.KeyDelete, tcell.KeyBackspace, tcell.KeyDEL:
			if flashcardQuiz.CanUndo() {
				showAnswer = false
				currentFlashcard = flashcardQuiz.Undo()
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
		contentViewContainer.SetTitle(fmt.Sprintf("%d / %d", flashcardQuiz.CurrentlySelected+1, len(flashcardQuiz.Flashcards)))

		known, unknown := flashcardQuiz.GetKnownAndUnknownCount()
		knowCount.SetText(strconv.Itoa(known))
		unknownCount.SetText(strconv.Itoa(unknown))

		return event
	})

	controls.
		AddItem(unknownCount, 0, 1, false).
		AddItem(nil, 25, 0, false).
		AddItem(dontKnowButtonFrame, 0, 2, false).
		AddItem(undoButtonFrame, 0, 2, false).
		AddItem(knowButtonFrame, 0, 2, false).
		AddItem(nil, 25, 0, false).
		AddItem(knowCount, 0, 1, false)

	quiz.
		AddItem(contentViewContainer, 0, 0, 1, 1, 0, 0, true).
		AddItem(controls, 1, 0, 1, 1, 0, 0, false)

	appState.Navigation.AddView(app.VIEW_NAMES.QuizProgressTrack, quiz, false, func() {
		if flashcardQuiz == nil || flashcardQuiz != appState.SelectedFlashcardSet.Quiz {
			if appState.SelectedFlashcardSet.Quiz == nil {
				appState.SelectedFlashcardSet.StartQuiz()
			}

			flashcardQuiz = appState.SelectedFlashcardSet.Quiz
			showAnswer = false
		}

		contentView.SetText(flashcardQuiz.CurrentlySelectedCard().Question)
		contentViewContainer.SetTitle(fmt.Sprintf("%d / %d", flashcardQuiz.CurrentlySelected+1, len(flashcardQuiz.Flashcards)))

		known, unknown := flashcardQuiz.GetKnownAndUnknownCount()
		knowCount.SetText(strconv.Itoa(known))
		unknownCount.SetText(strconv.Itoa(unknown))
	})
}
