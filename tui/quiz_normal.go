package tui

import (
	"fmt"
	"slices"

	"github.com/d3akhtar/tfc/app"
	"github.com/d3akhtar/tfc/domain"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func InitQuizNormalUi(appState *app.State) {
	var quizFlashcards []domain.Flashcard

	currentFlashcardIndex := 0
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
		SetBorderColor(Focused).
		SetBorderPadding(4, 4, 4, 4)

	contentView.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyLeft:
			if currentFlashcardIndex > 0 {
				showAnswer = false || appState.SelectedFlashcardSet.Front == domain.Answer
				currentFlashcardIndex--
			}

		case tcell.KeyRight:
			if currentFlashcardIndex < len(quizFlashcards)-1 {
				showAnswer = false || appState.SelectedFlashcardSet.Front == domain.Answer
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

	controls := tview.NewFlex()

	previousButton := NewButton("←").
		SetSelectedFunc(func() {
			if currentFlashcardIndex > 0 {
				showAnswer = false || appState.SelectedFlashcardSet.Front == domain.Answer
				currentFlashcardIndex--
			}

			contentView.SetText(quizFlashcards[currentFlashcardIndex].Question)
			contentViewContainer.SetTitle(fmt.Sprintf("%d / %d", currentFlashcardIndex+1, len(quizFlashcards)))

			appState.SetFocus(contentView)
		})

	nextButton := NewButton("→").
		SetSelectedFunc(func() {
			if currentFlashcardIndex < len(quizFlashcards)-1 {
				showAnswer = false || appState.SelectedFlashcardSet.Front == domain.Answer
				currentFlashcardIndex++
			}

			contentView.SetText(quizFlashcards[currentFlashcardIndex].Question)
			contentViewContainer.SetTitle(fmt.Sprintf("%d / %d", currentFlashcardIndex+1, len(quizFlashcards)))

			appState.SetFocus(contentView)
		})

	controls.
		AddItem(previousButton, 0, 1, false).
		AddItem(nil, 5, 0, false).
		AddItem(nextButton, 0, 1, false)

	quiz.
		AddItem(contentViewContainer, 0, 0, 1, 1, 0, 0, true).
		AddItem(controls, 1, 0, 1, 1, 0, 0, false)

	appState.Navigation.AddView(app.VIEW_NAMES.QuizNormal, quiz, false, func() {
		if !slices.Equal(quizFlashcards, appState.SelectedFlashcardSet.GetFlashcards()) {
			quizFlashcards = appState.SelectedFlashcardSet.GetFlashcards()
			currentFlashcardIndex = 0
			showAnswer = false || appState.SelectedFlashcardSet.Front == domain.Answer
		}

		if showAnswer {
			contentView.SetText(quizFlashcards[0].Answer)
		} else {
			contentView.SetText(quizFlashcards[0].Question)
		}

		contentViewContainer.SetTitle(fmt.Sprintf("%d / %d", currentFlashcardIndex+1, len(quizFlashcards)))
	})
}
