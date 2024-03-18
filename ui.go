package main

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// Global function that will take AppBackend clojure
// This function is instantiated in InitAppBackend function
var NewList func(items []*FSItem) *tview.List

type UI struct {
	ListCache map[string]*tview.List
	Grid      *tview.Grid
}

func (ui *UI) SetMainList(l *tview.List) {
	ui.Grid.AddItem(l, 1, 1, 3, 1, 0, 0, true)
}

func InitUI() UI {
	tview.Styles.PrimitiveBackgroundColor = tcell.ColorDefault
	// selectedStyle := tcell.StyleDefault.Reverse(true).Foreground(tcell.ColorBlue)

	// midList := tview.NewList().
	// 	AddItem("hello", "", 0, nil).
	// 	AddItem("World", "", 0, nil).
	// 	AddItem("!", "", 0, nil).
	// 	AddItem("qw[yellow::l]ertyi[-]opasdfghjkl;", "", 0, nil).
	// 	ShowSecondaryText(false).
	// 	SetHighlightFullLine(true).
	// 	SetSelectedStyle(selectedStyle)

	// leftList := tview.NewList().
	// 	ShowSecondaryText(false)
	// rightList := tview.NewList().
	// 	ShowSecondaryText(false)
	header := tview.NewTextView()
	footer := tview.NewTextView()

	grid := tview.NewGrid().
		SetRows(3, 0, 3).
		SetColumns(0, 0, 0).
		SetBorders(true).
		AddItem(header, 0, 0, 1, 3, 0, 0, false).
		AddItem(footer, 2, 0, 1, 3, 0, 0, false)

	listCache := make(map[string]*tview.List)

	return UI{listCache , grid}
}
