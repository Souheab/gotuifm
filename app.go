package main

//import "os"
import (
	"os"
	"path/filepath"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

const (
	InputRateLimit = time.Millisecond * 100
)


func runApp() {
	cwd , _ := os.Getwd()
	cwd, _ = filepath.Abs(cwd)
	backend := InitAppBackend(cwd)
	ui := &backend.UI

	// Rate limit input
	lastInput := time.Now()

	app := tview.NewApplication()
	inputHandler := func(event *tcell.EventKey) *tcell.EventKey {
		now := time.Now()
		if now.Sub(lastInput) < InputRateLimit {
			return nil
		}

		switch event.Key() {
		case tcell.KeyCtrlC:
			app.Stop()
		}

		switch event.Rune() {
		case 'h', 'H':
			backend.Select(1, 0, DirectionLeft)
		case 'j', 'J':
			backend.Select(1, 0, DirectionDown)
		case 'k', 'K':
			backend.Select(1, 0, DirectionUp)
		case 'l', 'L':
			backend.Select(1, 0, DirectionRight)
		case 'q', 'Q':
			app.Stop()
		}
		return nil
	}



	if err := app.SetInputCapture(inputHandler).SetRoot(ui.Grid, true).SetFocus(ui.Grid).Run(); err != nil {
		panic(err)
	}
}


