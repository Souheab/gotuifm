package main

import (
	"fmt"
	"log"
	"path/filepath"

	"github.com/rivo/tview"
)

type Tab struct {
	ActiveDirList  *DirList
	UI             *UI
	DotfilesFlag   bool
	InputCount     string
	BackendPointer *Backend
}

func (t *Tab) Select(n int, initialIndex int, direction int) {
	if n < 0 {
		return
	}

	if n == 0 {
		n = 1
	}

	acDl := t.ActiveDirList

	switch direction {
	case DirectionLeft:
		path := acDl.Path
		for range n {
			path = filepath.Dir(path)
		}
		dl := t.BackendPointer.DirListCacheGetOtherwiseAdd(path)
		t.MakeActiveDirlist(dl)
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
			dl := t.BackendPointer.DirListCacheGetOtherwiseAdd(fsItem.Path)
			t.MakeActiveDirlist(dl)
		}
	}

	t.UpdateFooter()
}

func (t *Tab) MakeActiveDirlist(dl *DirList) {
	t.UI.Grid.RemoveItem(t.ActiveDirList)
	t.ActiveDirList = dl
	t.UI.Grid.AddItem(dl, 1, 1, 1, 1, 0, 0, false)
	dl.SetDotfilesVisibility(t.DotfilesFlag)
	t.RunListChangedFunc()
	t.UI.CurrentPath.SetText(dl.Path)
}

func (t *Tab) SetDotfilesVisibility(visibility bool) {
	if t.DotfilesFlag == visibility {
		return
	}

	t.DotfilesFlag = visibility
	t.ActiveDirList.SetDotfilesVisibility(visibility)
}

func (t *Tab) ToggleDotfilesVisibility() {
	t.SetDotfilesVisibility(!t.DotfilesFlag)
}

func (t *Tab) RunListChangedFunc() {
	activeDl := t.ActiveDirList
	fsItem := activeDl.GetItemAtIndex(t.ActiveDirList.GetCurrentItem())
	if fsItem == nil {
		log.Fatalf("error in listChangedFunc")
	} else {
		if fsItem.Metadata.Type == Folder {
			if fsItem.Metadata.Readable {
				dl := t.BackendPointer.DirListCacheGetOtherwiseAdd(fsItem.Path)
				dl.SetDotfilesVisibility(t.DotfilesFlag)
				if len(dl.FilteredItems) == 0 {
					t.UI.Grid.AddItem(EmptyDirTextBox, 1, 2, 1, 1, 0, 0, false)
				} else {
					t.UI.Grid.AddItem(dl, 1, 2, 1, 1, 0, 0, false)
				}
			} else {
				t.UI.Grid.AddItem(PermissionDeniedTextBox, 1, 2, 1, 1, 0, 0, false)
			}
		} else {
			textBox := tview.NewTextView().SetLabel(fmt.Sprintf("File: %v, \npath: %v", fsItem.Name, fsItem.Path))
			t.UI.Grid.AddItem(textBox, 1, 2, 1, 1, 0, 0, false)
		}

		if activeDl.Path == "/" {
			t.UI.Grid.AddItem(tview.NewBox(), 1, 0, 1, 1, 0, 0, false)
		} else {
			parentPath := filepath.Dir(activeDl.Path)
			dl := t.BackendPointer.DirListCacheGetOtherwiseAdd(parentPath)
			dl.SetDotfilesVisibility(t.DotfilesFlag)
			t.UI.Grid.AddItem(dl, 1, 0, 1, 1, 0, 0, false)
		}
	}
}

func (t *Tab) UpdateFooter() {
	fsItem := t.ActiveDirList.GetSelectedItem()
	timeString := fsItem.Metadata.LastModified.Format("02-01-2006 15:04:05")
	permsString := fsItem.Metadata.PermsString
	fileSizeString := GetFileSizeHumanReadableString(fsItem.Metadata.FileSize)
	userString := fsItem.Metadata.UserString
	groupString := fsItem.Metadata.GroupString
	footerString := fmt.Sprintf("%s %s %s %s %s", permsString, userString, groupString, timeString, fileSizeString)
	t.UI.Footer.SetText(footerString)
}

func (t *Tab) AddToInputCount(inputKeyRune rune) {
	t.InputCount = fmt.Sprintf("%s%c", t.InputCount, inputKeyRune)
	t.UI.Footer.SetText(t.InputCount)
}

func (t *Tab) ClearInputCount() {
	t.InputCount = ""
	t.UI.Footer.SetText("")
}
