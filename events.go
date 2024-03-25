package main

import (
	"strconv"

	"github.com/gdamore/tcell/v2"
)

func (b *Backend) ProcessEventsWorker() {
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

func (b *Backend) HandleKeyEvent() {
	ch := b.InputChan
	select {
	case ev := <-ch:
		t := b.ActiveTab
		switch ev.Key() {
		case tcell.KeyCtrlC, tcell.KeyCtrlD:
			b.ExitApp()
		case tcell.KeyEsc:
			t.InputCount = ""
			return
		}

		inputKeyRune := ev.Rune()
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
		}
		t.RunListChangedFunc()
		t.ActiveDirList.AdjustOffset()
	}
}
