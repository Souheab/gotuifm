package main

import (
	"strconv"

	"github.com/gdamore/tcell/v2"
)

func (b *Backend) ProcessScreenEventsWorker() {
	ch := b.InputChan

	for {
		ev := b.Screen.PollEvent()

		switch ev := ev.(type) {
		case *tcell.EventKey:
			select {
			case ch <- ev:
			default:
			}
		case *tcell.EventResize:
			b.Draw()
		}
	}
}

func (b *Backend) ProcessDirListEventsWorker() {
	ch := b.DirListEventsChan

	for {
		<- ch
		b.ActiveTab.RunListChangedFunc()
		b.Draw()
	}
}

func (b *Backend) HandleKeyEvent(inputContext *InputContext) {
	ch := b.InputChan
	select {
	case ev := <-ch:
		t := b.ActiveTab
		switch ev.Key() {
		case tcell.KeyCtrlC, tcell.KeyCtrlD:
			b.ExitApp()
		case tcell.KeyEsc:
			t.InputCount = ""
			inputContext.Clear()
			return
		}

		inputKeyRune := ev.Rune()
		if inputContext.SortMenu {
			switch inputKeyRune {
			case 'd', 'D':
				b.ActiveTab.SetSortingCriteria(DefaultSort)
				b.ActiveTab.EchoFooter("Sorting criteria set to default")
			case 'T':
				b.ActiveTab.SetSortingCriteria(TimeSort)
				b.ActiveTab.EchoFooter("Sorting criteria set to sort by time")
			case 't':
				b.ActiveTab.SetSortingCriteria(TimeSortReverse)
				b.ActiveTab.EchoFooter("Sorting criteria set to sort by reverse time")
			}
			b.ActiveTab.EnsureCorrectSorting(t.ActiveDirList)
			inputContext.Clear()
		}

		if inputKeyRune >= '0' && inputKeyRune <= '9' {
			t.AddToInputCount(inputKeyRune)
			return
		}

		inputTimes := 1
		if t.InputCount != "" {
			inputTimes, _ = strconv.Atoi(t.InputCount)
			t.ClearInputCount()
		}

		switch inputKeyRune {
		case 'h', 'H':
			t.Select(1, 0, DirectionLeft)
		case 'j', 'J':
			t.Select(inputTimes, 0, DirectionDown)
		case 'k', 'K':
			t.Select(inputTimes, 0, DirectionUp)
		case 'l', 'L':
			t.Select(1, 0, DirectionRight)
		case '.':
			t.ToggleDotfilesVisibility()
		case 'q', 'Q':
			b.ExitApp()
		case 's', 'S':
			inputContext.SortMenu = true
		}
		t.RunListChangedFunc()
		t.ActiveDirList.AdjustOffset()
	}
}
