package main

import (
	"fmt"
	"log"
	"path/filepath"

	"github.com/rivo/tview"
)

type AppBackend struct {
	ActiveDirList *DirList
	DirListCache  map[string]*DirList
	UI            UI
}

func InitAppBackend(startingPath string) *AppBackend {
	dlc := make(map[string]*DirList)
	ui := InitUI()
	dl := DirListCacheAdd(dlc, startingPath)

	backend := AppBackend{dl, dlc, ui}
	backend.UI.Grid.AddItem(dl.UI, 1, 1, 1, 1, 0, 0, true)

	listChangedFunc := func(index int, mainText string, secondaryText string, shortcut rune) {
		bP := &backend
		activeDl := bP.ActiveDirList
		fsItem := activeDl.GetItemAtIndex(index)
		if (fsItem == nil) {
			log.Fatalf("error in listChangedFunc")
		} else {
			if fsItem.Metadata.Type == Folder {
				dl := bP.DirListCacheGetOtherwiseAdd( fsItem.Path)
				// textBox := tview.NewTextView().SetLabel(fmt.Sprintf("Folder: %v, \n path: %v", fsItem.Name, fsItem.Path))
				bP.UI.Grid.AddItem( dl.UI, 1, 2 , 1, 1 , 0, 0, false)
			} else {
				textBox := tview.NewTextView().SetLabel(fmt.Sprintf("File: %v, \npath: %v", fsItem.Name, fsItem.Path))
				bP.UI.Grid.AddItem( textBox, 1, 2 , 1, 1 , 0, 0, false)
			}

			parentPath := filepath.Dir(filepath.Dir(fsItem.Path))
			dl := bP.DirListCacheGetOtherwiseAdd(parentPath)
			bP.UI.Grid.AddItem( dl.UI, 1, 0 , 1, 1 , 0, 0, false)
		}
	}

	dl.UI.SetChangedFunc(listChangedFunc)


	return &backend
}

func (b *AppBackend) DirListCacheGetOtherwiseAdd(path string) *DirList {
	return DirListCacheGetOtherwiseAdd(b.DirListCache, path)
}

func DirListCacheAdd(dlc map[string]*DirList, path string) *DirList {
	path, err := filepath.Abs(path)
	if err != nil {
		log.Fatalf("%+v", err)
	}
	if PathExists(path) {
		dl, err := NewDirList(path)
		if err != nil {
			log.Fatalf("%+v", err)
		}
		dlc[path] = &dl

		return &dl
	}

	return nil
}

func DirListCacheGetOtherwiseAdd(dlc map[string]*DirList, path string) *DirList {
	dl, ok := dlc[path]

	if ok {
		return dl
	} else {
		return DirListCacheAdd(dlc, path)
	}

}
