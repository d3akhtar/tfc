package tui

import (
	"fmt"
	"slices"

	"github.com/d3akhtar/tfc/app"
	"github.com/d3akhtar/tfc/db"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func InitQuizUi(appState *app.State) {
	var quizFlashcards []db.Flashcard

	currentFlashcardIndex := 0
	showAnswer := false

	quiz := tview.NewGrid().
		SetRows(-20, -2)

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

	contentView.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyLeft:
			if currentFlashcardIndex > 0 {
				showAnswer = false
				currentFlashcardIndex--
			}

		case tcell.KeyRight:
			if currentFlashcardIndex < len(quizFlashcards)-1 {
				showAnswer = false
				currentFlashcardIndex++
			}

		case tcell.KeyRune:
			if event.Rune() == ' ' {
				showAnswer = !showAnswer
			}
		}

		var contents string
		if showAnswer {
			contents = quizFlashcards[currentFlashcardIndex].Answer
		} else {
			contents = quizFlashcards[currentFlashcardIndex].Question
		}

		contentView.SetText(contents)
		contentViewContainer.SetTitle(fmt.Sprintf("%d / %d", currentFlashcardIndex+1, len(quizFlashcards)))

		return event
	})

	quiz.
		AddItem(contentViewContainer, 0, 0, 1, 1, 0, 0, true).
		AddItem(tview.NewBox().SetBorder(true), 1, 0, 1, 1, 0, 0, false)

	appState.Navigation.AddView(app.VIEW_NAMES.Quiz, quiz, false, func() {
		if !slices.Equal(quizFlashcards, appState.SelectedFlashcardSet.Flashcards) {
			quizFlashcards = appState.SelectedFlashcardSet.GetFlashcards()
			currentFlashcardIndex = 0
			showAnswer = false
		}

		contentView.SetText(quizFlashcards[0].Question)
		contentViewContainer.SetTitle(fmt.Sprintf("%d / %d", currentFlashcardIndex+1, len(quizFlashcards)))
	})
}
