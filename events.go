package main

import (
	"strconv"

	"github.com/gdamore/tcell/v2"
)

type InputContext struct {
	SortMenu bool
}

func (ic *InputContext) Clear() {
	ic.SortMenu = false
}

func (b *Backend) ProcessScreenEventsWorker() {
	defer func() {
		if r := recover(); r != nil {
			b.Screen.Fini()
			panic(r)
		}
	}()

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

func (b *Backend) HandleEvents() {
	ech := b.DirListEventsChan
	ich := b.InputChan
	select {
	case <-ech:
		b.ActiveTab.RunListChangedFunc()
		b.Draw()
	case kev := <- ich:
		b.HandleKeyEvent(kev)
	}
}

func (b *Backend) HandleKeyEvent(kev *tcell.EventKey) {
	inputContext := b.InputContext
	t := b.ActiveTab
	switch kev.Key() {
	case tcell.KeyCtrlC, tcell.KeyCtrlD:
		b.ExitApp()
	case tcell.KeyEsc:
		t.InputCount = ""
		inputContext.Clear()
		return
	}

	inputKeyRune := kev.Rune()
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
	case 'e':
		t.BackendPointer.OpenEditor(t.ActiveDirList.GetSelectedItem().Path)
	case 'q', 'Q':
		b.ExitApp()
	case 's', 'S':
		inputContext.SortMenu = true
	}
	t.RunListChangedFunc()
	t.ActiveDirList.AdjustOffset()
}
