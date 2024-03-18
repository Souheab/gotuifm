package main

import (
	"fmt"
	"log"
	"path/filepath"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

const (
	DirectionLeft = iota
	DirectionRight
	DirectionUp
	DirectionDown
)

type AppBackend struct {
	ActiveDirList *DirList
	DirListCache  map[string]*DirList
	UI            UI
	DotfilesFlag  bool
}

func InitAppBackend(startingPath string) *AppBackend {
	dlc := make(map[string]*DirList)
	ui := InitUI()

	b := AppBackend{nil, dlc, ui, true}
	bP := &b

	NewList = func(items []*FSItem) *tview.List {
		selectedStyle := tcell.StyleDefault.Foreground(tcell.ColorBlue).Reverse(true)
		list := tview.NewList().ShowSecondaryText(false).SetSelectedStyle(selectedStyle).SetHighlightFullLine(true)

		for _, item := range items {
			list.AddItem( item.Name, "", 0, nil)
		}

		f := bP.GetListChangedFunc()
		list.SetChangedFunc(f)

		return list
	}

	dl := bP.DirListCacheAdd(startingPath)
	bP.UI.Grid.AddItem(dl, 1, 1, 1, 1, 0, 0, false)
	bP.ActiveDirList = dl


	return bP
}

func (b *AppBackend) Select(n int, initialIndex int, direction int) {
	if n <= 0 {
		return
	}

	acDl := b.ActiveDirList

	switch direction {
	case DirectionLeft:
		path := acDl.Path
		for range n {
			path = filepath.Dir(path)
		}
		dl := b.DirListCacheGetOtherwiseAdd(path)
		b.MakeActiveDirlist(dl, initialIndex)
	case DirectionUp, DirectionDown:
		currentIndex := acDl.GetCurrentItem()
		targetIndex := 0
		if direction == DirectionUp {
			targetIndex = currentIndex - n
		} else {
			targetIndex = currentIndex + n
		}
		targetIndex = targetIndex % b.ActiveDirList.GetItemCount()
		acDl.SetCurrentItem(targetIndex)
	case DirectionRight:
		fsItem := acDl.GetItemAtIndex(acDl.GetCurrentItem())
		if fsItem.Metadata.Type == Folder {
			dl := b.DirListCacheGetOtherwiseAdd(fsItem.Path)
			b.MakeActiveDirlist(dl, 0)
		}
	}
}

func (b *AppBackend) SetDotfilesVisibility (visibility bool) {
	if b.DotfilesFlag == visibility {
		return
	}

	b.DotfilesFlag = visibility
	b.ActiveDirList.SetDotfilesVisibility(visibility)
}

func (b *AppBackend) ToggleDotfilesVisibility() {
	b.SetDotfilesVisibility(!b.DotfilesFlag)
}

func (b *AppBackend) MakeActiveDirlist(dl *DirList, initialIndex int) {
	dl.List.SetCurrentItem(initialIndex)
	b.UI.Grid.RemoveItem(b.ActiveDirList)
	b.ActiveDirList = dl
	b.UI.Grid.AddItem(dl, 1, 1, 1, 1, 0, 0, false)
	dl.SetDotfilesVisibility(b.DotfilesFlag)
}

func (b *AppBackend) GetListChangedFunc() func(index int, mainText string, secondaryText string, shortcut rune) {
	return func(index int, mainText string, secondaryText string, shortcut rune) {
		activeDl := b.ActiveDirList
		fsItem := activeDl.GetItemAtIndex(index)
		if fsItem == nil {
			log.Fatalf("error in listChangedFunc")
		} else {
			if fsItem.Metadata.Type == Folder {
				dl := b.DirListCacheGetOtherwiseAdd(fsItem.Path)
				dl.SetDotfilesVisibility(b.DotfilesFlag)
				b.UI.Grid.AddItem(dl, 1, 2, 1, 1, 0, 0, false)
			} else {
				textBox := tview.NewTextView().SetLabel(fmt.Sprintf("File: %v, \npath: %v", fsItem.Name, fsItem.Path))
				b.UI.Grid.AddItem(textBox, 1, 2, 1, 1, 0, 0, false)
			}

			parentPath := filepath.Dir(filepath.Dir(fsItem.Path))
			dl := b.DirListCacheGetOtherwiseAdd(parentPath)
			dl.SetDotfilesVisibility(b.DotfilesFlag)
			b.UI.Grid.AddItem(dl, 1, 0, 1, 1, 0, 0, false)
		}
	}
}

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
