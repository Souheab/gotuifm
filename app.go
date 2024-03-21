package main

import (
	"os"
	"path/filepath"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"golang.org/x/time/rate"
)

// Single Source of Truth for this program
// Accesible Globally
var BackendPointer *AppBackend

func runApp() {
	cwd, _ := os.Getwd()
	cwd, _ = filepath.Abs(cwd)
	BackendPointer = CreateAppBackend()
	BackendPointer.StartAppBackend(cwd)
	ui := &BackendPointer.UI

	// Rate limit input
	limiter := rate.NewLimiter(rate.Limit(20), 5)
	app := tview.NewApplication()
	inputHandler := func(event *tcell.EventKey) *tcell.EventKey {
		if !limiter.Allow() {
			return nil
		}

		switch event.Key() {
		case tcell.KeyCtrlC, tcell.KeyCtrlD:
			app.Stop()
		}

		switch event.Rune() {
		case 'h', 'H':
			BackendPointer.Select(1, 0, DirectionLeft)
		case 'j', 'J':
			BackendPointer.Select(1, 0, DirectionDown)
		case 'k', 'K':
			BackendPointer.Select(1, 0, DirectionUp)
		case 'l', 'L':
			BackendPointer.Select(1, 0, DirectionRight)
		case '.':
			BackendPointer.ToggleDotfilesVisibility()
		case 'q', 'Q':
			app.Stop()
		}
		return nil
	}

	if err := app.SetInputCapture(inputHandler).SetRoot(ui.Grid, true).SetFocus(ui.Grid).Run(); err != nil {
		panic(err)
	}
}
