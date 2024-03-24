package main

import (
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/gdamore/tcell/v2"
)

const (
	TickTime = time.Millisecond * 10
)

func RunApp() {
	cwd, _ := os.Getwd()
	cwd, _ = filepath.Abs(cwd)
	b := NewAppBackend()
	b.StartAppBackend(cwd)
	s := b.Screen
	s.Init()
	go b.ProcessEventWorker()

	for {
		b.Draw()
		b.HandleKeyEvent()
		time.Sleep(TickTime)
	}
}

func (b *AppBackend) ProcessEventWorker() {
	ch := b.InputChan

	for {
		ev := b.Screen.PollEvent()

		switch ev := ev.(type) {
		case *tcell.EventKey:
			select {
			case ch <- ev:
			default:
			}
		}
	}
}

func (b *AppBackend) HandleKeyEvent() {
	ch := b.InputChan
	select {
	case ev := <-ch:
		switch ev.Key() {
		case tcell.KeyCtrlC, tcell.KeyCtrlD:
			b.ExitApp()
		case tcell.KeyEsc:
			b.InputCount = ""
			return
		}

		inputKeyRune := ev.Rune()
		if inputKeyRune >= '0' && inputKeyRune <= '9' {
			b.AddToInputCount(inputKeyRune)
			return
		}

		inputTimes := 1
		if b.InputCount != "" {
			inputTimes, _ = strconv.Atoi(b.InputCount)
			b.ClearInputCount()
		}

		switch inputKeyRune {
		case 'h', 'H':
			b.Select(1, 0, DirectionLeft)
		case 'j', 'J':
			b.Select(inputTimes, 0, DirectionDown)
		case 'k', 'K':
			b.Select(inputTimes, 0, DirectionUp)
		case 'l', 'L':
			b.Select(1, 0, DirectionRight)
		case '.':
			b.ToggleDotfilesVisibility()
		case 'q', 'Q':
			b.ExitApp()
		}
		b.RunListChangedFunc()
	default:
		return
	}

}

func (b *AppBackend) ExitApp() {
	b.Screen.Fini()
	os.Exit(0)
}
