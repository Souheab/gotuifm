package main

import (
	"fmt"
	"log"
	"os/user"
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
	Tabs              []Tab
	ActiveTab         *Tab
	DirListCache      map[string]*DirList
	Screen            tcell.Screen
	InputChan         chan *tcell.EventKey
	DirListEventsChan chan *string
	GroupNameCache    map[uint32]string
}

func InitAppBackend(startingPath string) *Backend {
	dlc := make(map[string]*DirList)
	gnc := make(map[uint32]string)
	tabs := make([]Tab, 0, 0)
	inputChan := make(chan *tcell.EventKey, InputChannelSize)
	dirListEventsChan := make(chan *string)
	s, _ := tcell.NewScreen()
	b := &Backend{tabs, nil, dlc, s, inputChan, dirListEventsChan, gnc}

	ui := InitUI()
	b.DirListCacheAddNonConcurrent(startingPath)
	dl := b.DirListCache[startingPath]
	t := &Tab{nil, ui, false, "", b, DefaultSort}
	t.MakeActiveDirlist(dl)
	b.ActiveTab = t

	return b
}

func (b *Backend) DirListCacheAdd(path string) {
	dlc := b.DirListCache
	_ , ok := dlc[path]
	if ok {
		return
	}

	path, err := filepath.Abs(path)
	if err != nil {
		log.Fatalf("%+v", err)
	}
	if PathExists(path) {
		dl, err := b.NewDirList(path)
		if err != nil {
			log.Fatalf("%+v", err)
		}
		dlc[path] = dl
	}
	b.DirListEventsChan <- &path
}


func (b *Backend) DirListCacheAddNonConcurrent(path string) {
	dlc := b.DirListCache
	_ , ok := dlc[path]
	if ok {
		return
	}

	path, err := filepath.Abs(path)
	if err != nil {
		log.Fatalf("%+v", err)
	}
	if PathExists(path) {
		dl, err := b.NewDirList(path)
		if err != nil {
			log.Fatalf("%+v", err)
		}
		dlc[path] = dl
	}
}

func (b *Backend) GetGroupName(gid uint32) string {
	gName, ok := b.GroupNameCache[gid]
	if !ok {
		grp, _ := user.LookupGroupId(fmt.Sprintf("%d", gid))
		gName = grp.Name
		b.GroupNameCache[gid] = gName
	}

	return gName
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
