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
	inputCapture := func(event *tcell.EventKey) *tcell.EventKey {
		now := time.Now()
		if now.Sub(lastInput) < InputRateLimit {
			return nil
		}

		return event
	}


	if err := tview.NewApplication().SetInputCapture(inputCapture).SetRoot(ui.Grid, true).SetFocus(ui.Grid).Run(); err != nil {
		panic(err)
	}
}


