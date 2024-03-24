package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"

	"github.com/gdamore/tcell/v2"
)

func RunApp() {
	cwd, _ := os.Getwd()
	cwd, _ = filepath.Abs(cwd)
	b := NewAppBackend()
	b.StartAppBackend(cwd)
	s := b.Screen
	s.Init()

	// App Loop
	for {
		ev := s.PollEvent()
		b.HandleEvent(ev)
		b.Draw()
	}
}

func (b *AppBackend) HandleEvent(ev tcell.Event) {
	switch ev := ev.(type) {
	case *tcell.EventKey:
		b.HandleKeyEvent(ev)
	case *tcell.EventResize:
		w, h := ev.Size()
		fmt.Printf("Resized to %dx%d", w, h)
	}
}

func (b *AppBackend) HandleKeyEvent(ev *tcell.EventKey) {
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
}

func (b *AppBackend) ExitApp() {
	b.Screen.Fini()
	os.Exit(0)
}
