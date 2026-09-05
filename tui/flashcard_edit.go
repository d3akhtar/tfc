package tui

import (
	"strconv"

	"github.com/d3akhtar/tfc/app"
	"github.com/d3akhtar/tfc/db/flashcard_set"
	"github.com/d3akhtar/tfc/domain"
	"github.com/d3akhtar/tfc/utils"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type flashcardEditPrimitiveInfo struct {
	Layout   *tview.Flex
	Question *tview.TextArea
	Answer   *tview.TextArea
}

func InitFlashcardEditUi(appState *app.State, flashcardSetRepository flashcard_set.FlashcardSetRepo) {
	var selectedFlashcardSet *domain.FlashcardSet = nil

	var window *utils.SlidingWindow[domain.Flashcard]

	maxFlashcardsShownInPreviewFlashcardList := 3

	lastSelectedFlashcardPrimitive := 0
	activeFlashcardPrimitives := make([]flashcardEditPrimitiveInfo, maxFlashcardsShownInPreviewFlashcardList)

	flashcardList := tview.NewFlex().
		SetDirection(tview.FlexRow)

	flashcardList.
		SetTitle("Flashcards").
		SetTitleAlign(tview.AlignLeft)

	SetBorderFocusAndBlurCallbacks(flashcardList.Box)

	newFlashcardPrimitive := func(flashcard domain.Flashcard, pos int) tview.Primitive {
		layout := tview.NewFlex()

		layout.
			SetBorder(true).
			SetBorderPadding(1, 1, 1, 1).
			SetFocusFunc(func() {
				lastSelectedFlashcardPrimitive = pos
				layout.SetBorderColor(Focused).SetTitleColor(Focused)
			}).
			SetBlurFunc(func() {
				layout.SetBorderColor(FlashcardPrimitiveBorder).SetTitleColor(FlashcardPrimitiveBorder)
			}).
			SetTitle(strconv.Itoa((pos + 1))).
			SetTitleAlign(tview.AlignLeft).
			SetBorderColor(FlashcardPrimitiveBorder).
			SetTitleColor(FlashcardPrimitiveBorder)

		question := tview.NewTextArea().
			SetText(flashcard.Question, true)

		answer := tview.NewTextArea().
			SetText(flashcard.Answer, true)

		question.
			SetBorder(true).
			SetBorderPadding(1, 1, 1, 1).
			SetTitle("Question").
			SetTitleAlign(tview.AlignLeft)

		SetBorderFocusAndBlurCallbacks(question.Box)

		answer.
			SetBorder(true).
			SetBorderPadding(1, 1, 1, 1).
			SetTitle("Answer").
			SetTitleAlign(tview.AlignLeft)

		SetBorderFocusAndBlurCallbacks(answer.Box)

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
							activeFlashcardPrimitives[idx].Question.SetText(flashcard.Question, true)
							activeFlashcardPrimitives[idx].Answer.SetText(flashcard.Answer, true)
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
							activeFlashcardPrimitives[idx].Question.SetText(flashcard.Question, true)
							activeFlashcardPrimitives[idx].Answer.SetText(flashcard.Answer, true)
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
			case tcell.KeyEsc:
				appState.SetFocus(flashcardList)
				return nil
			}

			return event
		})

		question.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
			switch event.Key() {
			case tcell.KeyTab, tcell.KeyBacktab:
				appState.SetFocus(answer)
				return nil
			case tcell.KeyEsc:
				appState.SetFocus(layout)
			}

			return event
		})

		question.SetChangedFunc(func() {
			selectedFlashcardSet.Flashcards[lastSelectedFlashcardPrimitive+window.Start].Question = question.GetText()
		})

		answer.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
			switch event.Key() {
			case tcell.KeyTab, tcell.KeyBacktab:
				appState.SetFocus(question)
				return nil
			case tcell.KeyEsc:
				appState.SetFocus(layout)
			}

			return event
		})

		answer.SetChangedFunc(func() {
			selectedFlashcardSet.Flashcards[lastSelectedFlashcardPrimitive+window.Start].Answer = answer.GetText()
		})

		activeFlashcardPrimitives[pos] = flashcardEditPrimitiveInfo{
			Layout:   layout,
			Question: question,
			Answer:   answer,
		}

		return layout
	}

	flashcardEdit := tview.NewPages()

	edit := tview.NewGrid().
		SetRows(-1, -2, -10, -2)

	titleInput := tview.NewInputField().
		SetFieldBackgroundColor(Background).
		SetChangedFunc(func(text string) {
			selectedFlashcardSet.Name = text
		})

	titleInput.
		SetBorder(true).
		SetTitle("Title").
		SetTitleAlign(tview.AlignLeft)

	SetBorderFocusAndBlurCallbacks(titleInput.Box)
	titleInput.
		SetBorderColor(Focused).SetTitleColor(Focused)

	descriptionTextArea := tview.NewTextArea()
	descriptionTextArea.
		SetChangedFunc(func() {
			selectedFlashcardSet.Description = descriptionTextArea.GetText()
		})

	descriptionTextArea.
		SetBorder(true).
		SetTitle("Description").
		SetTitleAlign(tview.AlignLeft).
		SetBackgroundColor(Background)

	SetBorderFocusAndBlurCallbacks(descriptionTextArea.Box)

	SetBorderFocusAndBlurCallbacks(flashcardList.Box)

	addCardButton := NewButton("Add Card").
		SetSelectedFunc(func() {
			flashcardEdit.ShowPage("flashcard")
		})

	finishButton := NewButton("Finish").
		SetSelectedFunc(func() {
			if appState.SelectedFlashcardSet() == nil {
				err := flashcardSetRepository.Create(appState.Context, selectedFlashcardSet)
				if err != nil {
					return
				}

				appState.SetSelectedFlashcardSet(selectedFlashcardSet)
			}

			appState.Navigation.GoToView(app.VIEW_NAMES.FlashcardSetPreview)
		})

	actionButtonGroup := tview.NewFlex().
		AddItem(addCardButton, 0, 1, true).
		AddItem(nil, 2, 0, false).
		AddItem(finishButton, 0, 1, false)

	newFlashcardInputForm := tview.NewForm().
		SetButtonsAlign(tview.AlignCenter).
		SetFieldStyle(tcell.StyleDefault.Background(Background))

	newFlashcardInputForm.
		SetBorder(true)

	newFlashcardInputFormModal := NewModal(newFlashcardInputForm, 60, 18)

	questionInputTextArea := tview.NewTextArea()
	SetBorderFocusAndBlurCallbacks(questionInputTextArea.Box)
	questionInputTextArea.
		SetBorder(true).
		SetTitle("Question").
		SetTitleAlign(tview.AlignLeft)

	answerInputTextArea := tview.NewTextArea()
	SetBorderFocusAndBlurCallbacks(answerInputTextArea.Box)
	answerInputTextArea.
		SetBorder(true).
		SetTitle("Answer").
		SetTitleAlign(tview.AlignLeft)

	newFlashcardInputForm.
		AddFormItem(questionInputTextArea).
		AddFormItem(answerInputTextArea).
		AddButton("Save", func() {
			question := questionInputTextArea.GetText()
			answer := answerInputTextArea.GetText()
			selectedFlashcardSet.AddFlashcard(question, answer)

			flashcardList.Clear()

			maxFlashcardsShownInPreviewFlashcardList = min(3, len(selectedFlashcardSet.Flashcards))

			window = utils.NewSlidingWindow(0, maxFlashcardsShownInPreviewFlashcardList, selectedFlashcardSet.GetFlashcards())

			for i := window.Start; i <= window.End; i++ {
				flashcardList.AddItem(
					newFlashcardPrimitive(window.Collection[i], i),
					0,
					1,
					false,
				)
			}

			flashcardEdit.HidePage("flashcard")
		}).
		AddButton("Cancel", func() {
			flashcardEdit.HidePage("flashcard")
		})

	confirmDeleteButton := NewButton("Confirm")
	cancelDeleteButton := NewButton("Cancel").
		SetSelectedFunc(func() {
			flashcardEdit.HidePage("delete")
		})

	confirmDeleteDialog := tview.NewFrame(
		tview.NewFlex().
			AddItem(confirmDeleteButton, 0, 1, false).
			AddItem(cancelDeleteButton, 0, 1, true),
	).AddText("Are you sure you want to delete this card?", true, tview.AlignCenter, tcell.ColorWhite)

	titleInput.
		SetDoneFunc(func(_ tcell.Key) {
			appState.SetFocus(descriptionTextArea)
		})

	descriptionTextArea.
		SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
			switch event.Key() {
			case tcell.KeyTab:
				appState.SetFocus(flashcardList)
				return nil
			case tcell.KeyBacktab:
				appState.SetFocus(titleInput)
				return nil
			}

			return event
		})

	flashcardList.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if activeFlashcardPrimitives[lastSelectedFlashcardPrimitive].Layout.HasFocus() {
			return event
		}

		switch event.Key() {
		case tcell.KeyEnter:
			appState.SetFocus(activeFlashcardPrimitives[lastSelectedFlashcardPrimitive].Layout)
			return nil
		case tcell.KeyBacktab:
			appState.SetFocus(descriptionTextArea)
			return nil
		case tcell.KeyTab:
			appState.SetFocus(addCardButton)
			return nil
		}

		return event
	})

	addCardButton.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyBacktab:
			appState.SetFocus(flashcardList)
			return nil
		case tcell.KeyTab, tcell.KeyRight:
			appState.SetFocus(finishButton)
			return nil
		}

		return event
	})

	finishButton.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyBacktab, tcell.KeyLeft:
			appState.SetFocus(addCardButton)
			return nil
		}

		return event
	})

	edit.
		AddItem(titleInput, 0, 0, 1, 1, 0, 0, true).
		AddItem(descriptionTextArea, 1, 0, 1, 1, 0, 0, false).
		AddItem(flashcardList, 2, 0, 1, 1, 0, 0, false).
		AddItem(actionButtonGroup, 3, 0, 1, 1, 0, 0, false)

	flashcardEdit.AddPage("main", edit, true, true)
	flashcardEdit.AddPage("flashcard", newFlashcardInputFormModal, true, false)
	flashcardEdit.AddPage("delete", confirmDeleteDialog, true, false)

	refresh := func() error {
		flashcardList.Clear()

		if selectedFlashcardSet != appState.SelectedFlashcardSet() {
			selectedFlashcardSet = appState.SelectedFlashcardSet()
		}

		if selectedFlashcardSet == nil {
			selectedFlashcardSet = domain.NewFlashcardSet("Unnamed")
			finishButton.SetLabel("Create")
		} else {
			finishButton.SetLabel("Finish")
		}

		titleInput.SetText(selectedFlashcardSet.Name)
		descriptionTextArea.SetText(selectedFlashcardSet.Description, true)

		maxFlashcardsShownInPreviewFlashcardList = min(3, len(selectedFlashcardSet.Flashcards))

		window = utils.NewSlidingWindow(0, maxFlashcardsShownInPreviewFlashcardList, selectedFlashcardSet.GetFlashcards())

		for i := window.Start; i <= window.End; i++ {
			flashcardList.AddItem(
				newFlashcardPrimitive(window.Collection[i], i),
				0,
				1,
				false,
			)
		}

		return nil
	}

	exit := func() error {
		if appState.SelectedFlashcardSet() != nil {
			err := flashcardSetRepository.Update(appState.Context, selectedFlashcardSet)
			if err != nil {
				return err
			}
		}

		return nil
	}

	appState.Navigation.AddView(app.VIEW_NAMES.FlashcardEdit, flashcardEdit, false, refresh, exit)
}
