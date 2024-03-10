package main

//import "os"
import (
	"os"
	"path/filepath"
	"github.com/rivo/tview"
)


func runApp() {
	cwd , _ := os.Getwd()
	cwd, _ = filepath.Abs(cwd)
	backend := InitAppBackend(cwd)
	ui := &backend.UI


	if err := tview.NewApplication().SetRoot(ui.Grid, true).SetFocus(ui.Grid).Run(); err != nil {
		panic(err)
	}
}


