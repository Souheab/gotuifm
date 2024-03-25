package main

import (
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"syscall"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

const (
	Folder = iota
	File
)

type DirList struct {
	*tview.Box
	FSItems           []*FSItem
	FilteredItems     []*FSItem
	Path              string
	DotfilesFlag      bool
	selectedItemIndex int
	itemOffset        int
}

func (dl *DirList) GetItemAtIndex(index int) *FSItem {
	fsItems := dl.FilteredItems

	if index >= len(fsItems) {
		return nil
	} else {
		return fsItems[index]
	}
}

func (dl *DirList) GetSelectedItem() *FSItem {
	return dl.GetItemAtIndex(dl.selectedItemIndex)
}

func (dl *DirList) SetFilter(f func(fsItem *FSItem) bool) {
	filteredItems := make([]*FSItem, 0, 0)

	for _, fsItem := range dl.FSItems {
		if f(fsItem) {
			filteredItems = append(filteredItems, fsItem)
		}
	}

	dl.FilteredItems = filteredItems
}

func (dl *DirList) RemoveFilter() {
	dl.FilteredItems = dl.FSItems
}

func (dl *DirList) SetDotfilesVisibility(visible bool) {
	if dl.DotfilesFlag == visible {
		return
	}

	if visible {
		dl.DotfilesFlag = true
		dl.RemoveFilter()
	} else {
		f := func(fsItem *FSItem) bool {
			return fsItem.Name[0] != '.'
		}

		dl.SetFilter(f)
		dl.DotfilesFlag = false
	}
}

type FSItem struct {
	Path     string
	Name     string
	Metadata FSItemMetadata
}

type FSItemMetadata struct {
	Type          int
	Readable      bool
	FileExtension string
	LastModified  time.Time
	PermsString   string
	FileSize      int64
	UserString    string
	GroupString   string
}

func NewDirList(path string) (*DirList, error) {

	fsDirEntry, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}

	files := make([]*FSItem, 0, 0)
	folders := make([]*FSItem, 0, 0)

	for _, fsEntry := range fsDirEntry {
		name := fsEntry.Name()
		fsItemPath, _ := filepath.Abs(filepath.Join(path, name))
		fileInfo, _ := fsEntry.Info()
		lastModified := fileInfo.ModTime()
		permsString := fileInfo.Mode().String()
		fileSize := fileInfo.Size()
		stat := fileInfo.Sys().(*syscall.Stat_t)
		uid := stat.Uid
		gid := stat.Gid

		ownerUser, _ := user.LookupId(fmt.Sprintf("%d", uid))
		group, _ := user.LookupGroupId(fmt.Sprintf("%d", gid))

		if fsEntry.IsDir() {
			metadata := FSItemMetadata{Folder, PathReadable(fsItemPath), "", lastModified, permsString, fileSize, ownerUser.Username, group.Name}
			folder := FSItem{fsItemPath, fsEntry.Name(), metadata}
			folders = append(folders, &folder)
		} else {
			fileExtension := filepath.Ext(fsItemPath)
			metadata := FSItemMetadata{File, true, fileExtension, lastModified, permsString, fileSize, ownerUser.Username, group.Name}
			file := FSItem{fsItemPath, fsEntry.Name(), metadata}
			files = append(files, &file)
		}
	}

	fsItems := append(folders, files...)

	return &DirList{tview.NewBox(), fsItems, fsItems, path, true, 0, 0}, nil
}

func (dl *DirList) Draw(screen tcell.Screen) {
	textStyle := tcell.StyleDefault
	selectedStyle := textStyle.Reverse(true)

	dl.Box.DrawForSubclass(screen, dl)

	x, y, width, height := dl.GetInnerRect()
	bottomLimit := y + height
	_, totalHeight := screen.Size()
	if bottomLimit > totalHeight {
		bottomLimit = totalHeight
	}

	if (dl.selectedItemIndex + dl.itemOffset) > bottomLimit {
	}

	for i, item := range dl.FilteredItems {
		if i < dl.itemOffset {
			continue
		}

		itemString := item.Name
		if y >= bottomLimit {
			break
		}

		printStyle := textStyle
		if i == dl.selectedItemIndex {
			printStyle = selectedStyle
		}

		PrintWithStyle(screen, itemString, x, y, width, printStyle)

		y++
	}
}

func (dl *DirList) AdjustOffset() {
	_, _, _, height := dl.GetInnerRect()
	if height == 0 {
		return
	}
	if dl.selectedItemIndex < dl.itemOffset {
		dl.itemOffset = dl.selectedItemIndex
	} else {
		if dl.selectedItemIndex-dl.itemOffset >= height {
			dl.itemOffset = dl.selectedItemIndex + 1 - height
		}
	}
}
