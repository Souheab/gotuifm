package main

//import "os"
import (
	"github.com/rivo/tview"
)


func runApp() {
	// dirListCache := make(map[string]DirList)
	// cwd , _ := os.Getwd()
	// DirListCacheAdd(dirListCache, cwd)
	// backend := AppBackend{dirListCache}
	ui := InitUI()

	if err := tview.NewApplication().SetRoot(ui.Grid, true).SetFocus(ui.Grid).Run(); err != nil {
		panic(err)
	}
}


