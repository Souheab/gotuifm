package main

import (
	"fmt"
	"os"
	"os/user"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func NewList(items []*FSItem) *tview.List {
	selectedStyle := tcell.StyleDefault.Foreground(tcell.ColorBlue).Reverse(true)
	list := tview.NewList().ShowSecondaryText(false).SetSelectedStyle(selectedStyle).SetHighlightFullLine(true)

	for _, item := range items {
		list.AddItem(item.Name, "", 0, nil)
	}

	f := BackendPointer.GetListChangedFunc()
	list.SetChangedFunc(f)

	return list
}

type UI struct {
	ListCache map[string]*tview.List
	Grid      *tview.Grid
	CurrentPath *tview.TextView
}

func (ui *UI) SetMainList(l *tview.List) {
	ui.Grid.AddItem(l, 1, 1, 3, 1, 0, 0, true)
}

func InitUI() UI {
	tview.Styles.PrimitiveBackgroundColor = tcell.ColorDefault

	hostname, _ := os.Hostname()
	user, err := user.Current()
	if err != nil {
		panic(err)
	}
	username := user.Username

	usernameHostnameTextBox := tview.NewTextView()
	currentPath := tview.NewTextView().SetText("test")
	header := tview.NewFlex()
	footer := tview.NewTextView()

	usernameHostnameTextBox.SetTextStyle(tcell.StyleDefault.Foreground(tcell.ColorGreen).Bold(true))
	usernameHostnameString := fmt.Sprintf("%v@%v:", username, hostname)
	usernameHostnameTextBox.SetText(usernameHostnameString)

	header.AddItem(usernameHostnameTextBox, len(usernameHostnameString), 1, false)
	header.AddItem(currentPath, 0, 1, false)

	grid := tview.NewGrid().
		SetRows(1, 0, 1).
		SetColumns(0, 0, 0).
		SetBorders(true).
		AddItem(header, 0, 0, 1, 3, 0, 0, false).
		AddItem(footer, 2, 0, 1, 3, 0, 0, false)

	listCache := make(map[string]*tview.List)

	return UI{listCache , grid, currentPath}
}
