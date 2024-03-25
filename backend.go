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

type Backend struct {
	Tabs         []Tab
	ActiveTab    *Tab
	DirListCache map[string]*DirList
	Screen       tcell.Screen
	InputChan    chan *tcell.EventKey
}

func InitAppBackend(startingPath string) *Backend {
	dlc := make(map[string]*DirList)
	tabs := make([]Tab, 0, 0)
	inputChan := make(chan *tcell.EventKey, InputChannelSize)
	s, _ := tcell.NewScreen()
	b := &Backend{tabs, nil, dlc, s, inputChan}

	ui := InitUI()
	dl := b.DirListCacheAdd(startingPath)
	t := &Tab{nil, ui, false, "", b}
	t.MakeActiveDirlist(dl)
	b.ActiveTab = t

	return b
}

func (b *Backend) DirListCacheAdd(path string) *DirList {
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
		dlc[path] = dl

		return dl
	}

	return nil
}

func (b *Backend) DirListCacheGetOtherwiseAdd(path string) *DirList {
	dlc := b.DirListCache
	dl, ok := dlc[path]

	if ok {
		return dl
	} else {
		return b.DirListCacheAdd(path)
	}
}

func (b *Backend) Draw() {
	s := b.Screen
	grid := b.ActiveTab.UI.Grid

	w, h := s.Size()
	s.Clear()
	grid.SetRect(0, 0, w, h)
	grid.Draw(b.Screen)
	s.Show()
}
