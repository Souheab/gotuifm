package main

import (
	"fmt"
	"os"
	"os/user"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

var PermissionDeniedTextBox *tview.TextView
var EmptyDirTextBox *tview.TextView
var EmptyBox *tview.Box
var LoadingTextBox *tview.TextView

type UI struct {
	ListCache   map[string]*tview.List
	Grid        *tview.Grid
	LeftPane    tview.Primitive
	RightPane   tview.Primitive
	CurrentPath *tview.TextView
	Footer      *tview.TextView
}

func InitUI() *UI {
	tview.Styles.PrimitiveBackgroundColor = tcell.ColorDefault

	PermissionDeniedTextBox = tview.NewTextView().SetLabel("[white:red]Permission Denied")
	EmptyDirTextBox = tview.NewTextView().SetLabel("[white::r]empty")
	EmptyBox = tview.NewBox()
	LoadingTextBox = tview.NewTextView().SetLabel("[white::r]loading...")

	hostname, _ := os.Hostname()
	user, err := user.Current()
	if err != nil {
		panic(err)
	}
	username := user.Username

	usernameHostnameTextBox := tview.NewTextView()
	currentPath := tview.NewTextView().SetText("test")
	header := tview.NewFlex()
	footer := tview.NewTextView().SetTextStyle(tcell.StyleDefault.Foreground(tcell.ColorYellow))

	usernameHostnameTextBox.SetTextStyle(tcell.StyleDefault.Foreground(tcell.ColorGreen).Bold(true))
	usernameHostnameString := fmt.Sprintf("%v@%v:", username, hostname)
	usernameHostnameTextBox.SetText(usernameHostnameString)

	header.AddItem(usernameHostnameTextBox, len(usernameHostnameString), 1, false)
	header.AddItem(currentPath, 0, 1, false)

	grid := tview.NewGrid().
		SetRows(1, 0, 1).
		SetColumns(0, -2, -3).
		SetBorders(true).
		AddItem(header, 0, 0, 1, 3, 0, 0, false).
		AddItem(footer, 2, 0, 1, 3, 0, 0, false)

	listCache := make(map[string]*tview.List)

	return &UI{listCache, grid,  nil, nil,  currentPath, footer}
}
