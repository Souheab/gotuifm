package main

import (
	"fmt"
	"log"
	"path/filepath"

	"github.com/rivo/tview"
)

type Tab struct {
	ActiveDirList   *DirList
	UI              *UI
	DotfilesFlag    bool
	InputCount      string
	BackendPointer  *Backend
	SortingCriteria int
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
		dl, ok := t.BackendPointer.DirListCache[path]
		if ok {
			t.MakeActiveDirlist(dl)
		}
	case DirectionUp, DirectionDown:
		currentIndex := acDl.selectedItemIndex
		targetIndex := 0
		if direction == DirectionUp {
			targetIndex = currentIndex - n
			if targetIndex < 0 {
				targetIndex = 0
			}
		} else {
			targetIndex = currentIndex + n
			if targetIndex >= len(acDl.FilteredItems) {
				targetIndex = len(acDl.FilteredItems) - 1
			}
		}
		acDl.selectedItemIndex = targetIndex
	case DirectionRight:
		fsItem := acDl.GetItemAtIndex(acDl.selectedItemIndex)
		if fsItem.Metadata.Type == Folder && fsItem.Metadata.Readable {
			dl, ok := t.BackendPointer.DirListCache[fsItem.Path]
			if len(dl.FilteredItems) != 0 && ok {
				t.MakeActiveDirlist(dl)
			}
		}
	}

	t.UpdateFooter()
}

func (t *Tab) MakeActiveDirlist(dl *DirList) {
	t.UI.Grid.RemoveItem(t.ActiveDirList)
	// Remove dl just in case it is already present somewhere in the grid to avoid making duplicates
	t.UI.Grid.RemoveItem(dl)
	t.ActiveDirList = dl
	t.EnsureCorrectSorting(dl)
	t.UI.Grid.AddItem(dl, 1, 1, 1, 1, 0, 0, false)
	dl.SetDotfilesVisibility(t.DotfilesFlag)
	t.RunListChangedFunc()
	t.UI.CurrentPath.SetText(dl.Path)
}

func (t *Tab) EnsureCorrectSorting(dl *DirList) {
	if dl.SortingCriteria == t.SortingCriteria || dl == nil{
		return
	}

	SortByCriteria(dl.FSItems, t.SortingCriteria)
	SortByCriteria(dl.FilteredItems, t.SortingCriteria)
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
	fsItem := activeDl.GetItemAtIndex(t.ActiveDirList.selectedItemIndex)
	var l tview.Primitive = nil
	var r tview.Primitive = nil
	if fsItem == nil {
		log.Fatalf("error in listChangedFunc")
	} else {
		if fsItem.Metadata.Type == Folder {
			if fsItem.Metadata.Readable {
				dl, ok := t.BackendPointer.DirListCache[fsItem.Path]

				if ok {
					dl.SetDotfilesVisibility(t.DotfilesFlag)
					if len(dl.FilteredItems) == 0 {
						r = EmptyDirTextBox
					} else {
						r = dl
					}
				} else {
					go t.BackendPointer.DirListCacheAdd(fsItem.Path)
					r = LoadingTextBox
				}
			} else {
				r = PermissionDeniedTextBox
			}
		} else {
			textBox := tview.NewTextView().SetLabel(fmt.Sprintf("File: %v, \npath: %v", fsItem.Name, fsItem.Path))
			r = textBox
		}

		if activeDl.Path == "/" {
			l = EmptyBox
		} else {
			parentPath := filepath.Dir(activeDl.Path)
			dl, ok := t.BackendPointer.DirListCache[parentPath]
			if ok {
				dl.SetDotfilesVisibility(t.DotfilesFlag)
				l = dl
			} else {
				go t.BackendPointer.DirListCacheAdd(parentPath)
				l = LoadingTextBox
			}
		}
	}

	d, ok := l.(*DirList)
	if ok {
		t.EnsureCorrectSorting(d)
	}
	d, ok = r.(*DirList)
	if ok {
		t.EnsureCorrectSorting(d)
	}
	t.UpdateGrid(l, r)
}

func (t *Tab) UpdateGrid(l, r tview.Primitive) {
	if t.UI.LeftPane != nil && l != nil && t.UI.LeftPane != t.ActiveDirList {
		t.UI.Grid.RemoveItem(t.UI.LeftPane)
	}

	if t.UI.RightPane != nil && r != nil && t.UI.RightPane != t.ActiveDirList {
		t.UI.Grid.RemoveItem(t.UI.RightPane)
	}

	if l != nil {
		t.UI.LeftPane = l
		t.UI.Grid.AddItem(t.UI.LeftPane, 1, 0, 1, 1, 0, 0, false)
	}

	if r != nil {
		t.UI.RightPane = r
		t.UI.Grid.AddItem(t.UI.RightPane, 1, 2, 1, 1, 0, 0, false)
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
