package main

import (
	"log"
	"path/filepath"

	"github.com/gdamore/tcell/v2"
)

const (
	DirectionLeft = iota
	DirectionRight
	DirectionUp
	DirectionDown
)

const (
	InputChannelSize = 5
)

type AppBackend struct {
	Tabs         []Tab
	ActiveTab    *Tab
	DirListCache map[string]*DirList
	Screen       tcell.Screen
	InputChan    chan *tcell.EventKey
}

func InitAppBackend(startingPath string) *AppBackend {
	dlc := make(map[string]*DirList)
	tabs := make([]Tab, 0, 0)
	inputChan := make(chan *tcell.EventKey, InputChannelSize)
	s, _ := tcell.NewScreen()
	b := &AppBackend{tabs, nil, dlc, s, inputChan}

	ui := InitUI()
	dl := b.DirListCacheAdd(startingPath)
	t := &Tab{nil, ui, false, "", b}
	t.MakeActiveDirlist(dl)
	b.ActiveTab = t

	return b
}

// func NewAppBackend() *AppBackend {
// 	dlc := make(map[string]*DirList)
// 	inputChan := make(chan *tcell.EventKey, InputChannelSize)
// 	ui := InitUI()
// 	s, _ := tcell.NewScreen()
// 	b := AppBackend{nil, dlc, ui, false, "", s, inputChan}
// 	return &b
// }

// func (b *AppBackend) StartAppBackend(startingPath string) {
// 	dl := b.DirListCacheAdd(startingPath)
// 	b.MakeActiveDirlist(dl)
// }

func (b *AppBackend) DirListCacheAdd(path string) *DirList {
	dlc := b.DirListCache
	path, err := filepath.Abs(path)
	if err != nil {
		log.Fatalf("%+v", err)
	}
	if PathExists(path) {
		dl, err := NewDirList(path)
		if err != nil {
			log.Fatalf("%+v", err)
		}
		dlP := &dl
		dlc[path] = dlP

		return dlP
	}

	return nil
}

func (b *AppBackend) DirListCacheGetOtherwiseAdd(path string) *DirList {
	dlc := b.DirListCache
	dl, ok := dlc[path]

	if ok {
		return dl
	} else {
		return b.DirListCacheAdd(path)
	}
}

func (b *AppBackend) Draw() {
	s := b.Screen
	grid := b.ActiveTab.UI.Grid

	w, h := s.Size()
	s.Clear()
	grid.SetRect(0, 0, w, h)
	grid.Draw(b.Screen)
	s.Show()
}
