package main

import (
	"fmt"
	"log"
	"path/filepath"

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
	InputCount    string
}

func CreateAppBackend() *AppBackend {
	dlc := make(map[string]*DirList)
	ui := InitUI()
	b := AppBackend{nil, dlc, ui, false, ""}
	return &b
}

func (b *AppBackend) StartAppBackend(startingPath string) {
	dl := b.DirListCacheAdd(startingPath)
	b.MakeActiveDirlist(dl)
}

func (b *AppBackend) Select(n int, initialIndex int, direction int) {
	if n < 0 {
		return
	}

	if n == 0 {
		n = 1
	}

	acDl := b.ActiveDirList

	switch direction {
	case DirectionLeft:
		path := acDl.Path
		for range n {
			path = filepath.Dir(path)
		}
		dl := b.DirListCacheGetOtherwiseAdd(path)
		b.MakeActiveDirlist(dl)
	case DirectionUp, DirectionDown:
		currentIndex := acDl.GetCurrentItem()
		targetIndex := 0
		if direction == DirectionUp {
			targetIndex = currentIndex - n
			if targetIndex < 0 {
				targetIndex = 0
			}
		} else {
			targetIndex = currentIndex + n
		}
		acDl.SetCurrentItem(targetIndex)
	case DirectionRight:
		fsItem := acDl.GetItemAtIndex(acDl.GetCurrentItem())
		if fsItem.Metadata.Type == Folder && fsItem.Metadata.Readable {
			dl := b.DirListCacheGetOtherwiseAdd(fsItem.Path)
			b.MakeActiveDirlist(dl)
		}
	}
}

func (b *AppBackend) SetDotfilesVisibility(visibility bool) {
	if b.DotfilesFlag == visibility {
		return
	}

	b.DotfilesFlag = visibility
	b.ActiveDirList.SetDotfilesVisibility(visibility)
}

func (b *AppBackend) ToggleDotfilesVisibility() {
	b.SetDotfilesVisibility(!b.DotfilesFlag)
}

func (b *AppBackend) MakeActiveDirlist(dl *DirList) {
	b.UI.Grid.RemoveItem(b.ActiveDirList)
	b.ActiveDirList = dl
	b.UI.Grid.AddItem(dl, 1, 1, 1, 1, 0, 0, false)
	dl.SetDotfilesVisibility(b.DotfilesFlag)
	b.RunListChangedFunc()
	b.UI.CurrentPath.SetText(dl.Path)
}

func (b *AppBackend) GetListChangedFunc() func(index int, mainText string, secondaryText string, shortcut rune) {
	return func(index int, mainText string, secondaryText string, shortcut rune) {
		activeDl := b.ActiveDirList
		fsItem := activeDl.GetItemAtIndex(index)
		if fsItem == nil {
			log.Fatalf("error in listChangedFunc")
		} else {
			if fsItem.Metadata.Type == Folder {
				if fsItem.Metadata.Readable {
					dl := b.DirListCacheGetOtherwiseAdd(fsItem.Path)
					dl.SetDotfilesVisibility(b.DotfilesFlag)
					if len(dl.FilteredItems) == 0 {
						b.UI.Grid.AddItem(EmptyDirTextBox, 1, 2, 1, 1, 0, 0, false)
					} else {
						b.UI.Grid.AddItem(dl, 1, 2, 1, 1, 0, 0, false)
					}
				} else {
					b.UI.Grid.AddItem(PermissionDeniedTextBox, 1, 2, 1, 1, 0, 0, false)
				}
			} else {
				textBox := tview.NewTextView().SetLabel(fmt.Sprintf("File: %v, \npath: %v", fsItem.Name, fsItem.Path))
				b.UI.Grid.AddItem(textBox, 1, 2, 1, 1, 0, 0, false)
			}

			if activeDl.Path == "/" {
				b.UI.Grid.AddItem(tview.NewBox(), 1, 0, 1, 1, 0, 0, false)
			} else {
				parentPath := filepath.Dir(activeDl.Path)
				dl := b.DirListCacheGetOtherwiseAdd(parentPath)
				dl.SetDotfilesVisibility(b.DotfilesFlag)
				b.UI.Grid.AddItem(dl, 1, 0, 1, 1, 0, 0, false)
			}
		}
	}
}

// Usually it should run automatically but sometimes we will need to run it automatically
func (b *AppBackend) RunListChangedFunc() {
	f := b.GetListChangedFunc()
	f(b.ActiveDirList.GetCurrentItem(), "", "", 0)
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

func (b *AppBackend) AddToInputCount(inputKeyRune rune) {
	b.InputCount = fmt.Sprintf("%s%c", b.InputCount, inputKeyRune)
	b.UI.InputCount.SetText(b.InputCount)
}

func (b *AppBackend) ClearInputCount() {
	b.InputCount = ""
	b.UI.InputCount.SetText("")
}
