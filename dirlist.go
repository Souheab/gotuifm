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
	SortingCriteria   int
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
	var selectedItem *FSItem = nil

	if dl.selectedItemIndex >= 0 && dl.selectedItemIndex < len(dl.FilteredItems) {
		selectedItem = dl.FilteredItems[dl.selectedItemIndex]
	}

	selectedItemFiltered := false
	for _, fsItem := range dl.FSItems {
		if f(fsItem) {
			filteredItems = append(filteredItems, fsItem)
			if fsItem == selectedItem {
				dl.selectedItemIndex = len(filteredItems) - 1
				selectedItemFiltered = true
			}
		}
	}

	filteredItemsLength := len(filteredItems)
	if !selectedItemFiltered && dl.selectedItemIndex >= filteredItemsLength {
		dl.selectedItemIndex = filteredItemsLength - 1
	}

	dl.FilteredItems = filteredItems
	// TODO: deal with this offset thing
	dl.AdjustOffset()
}

func (dl *DirList) RemoveFilter() {
	// dl.FilteredItems = dl.FSItems
	dl.SetFilter(func(_ *FSItem) bool {
		return true
	})
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
	Icon          rune
	Dotfile       bool
}

func (b *Backend) NewDirList(path string, selectedItemPath *string) (*DirList, error) {

	fsDirEntry, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}

	fsItems := make([]*FSItem, 0, 0)
	selectedItemIndex := 0

	for i, fsEntry := range fsDirEntry {
		name := fsEntry.Name()
		fsItemPath, _ := filepath.Abs(filepath.Join(path, name))
		fileInfo, _ := fsEntry.Info()
		lastModified := fileInfo.ModTime()
		permsString := fileInfo.Mode().String()
		fileSize := fileInfo.Size()
		stat := fileInfo.Sys().(*syscall.Stat_t)
		uid := stat.Uid
		gid := stat.Gid
		fileExtension := ""
		readable := true
		fsItemType := File
		icon := FolderIcon

		ownerUser, _ := user.LookupId(fmt.Sprintf("%d", uid))
		groupName := b.GetGroupName(gid)

		isDotfile := name[0] == '.'

		var metadata FSItemMetadata
		if fsEntry.IsDir() {
			readable = PathReadable(fsItemPath)
			fsItemType = Folder
		} else {
			fileExtension = filepath.Ext(fsItemPath)
			icon = GetFileIcon(fileExtension)
		}

		metadata = FSItemMetadata{fsItemType, readable, fileExtension, lastModified, permsString, fileSize, ownerUser.Username, groupName, icon, isDotfile}

		if selectedItemPath != nil && *selectedItemPath == fsItemPath {
			selectedItemIndex = i
		}

		fsItem := &FSItem{fsItemPath, name, metadata}
		fsItems = append(fsItems, fsItem)
	}

	filteredItems := append([]*FSItem(nil), fsItems...)
	SortByCriteria(fsItems, DefaultSort)

	return &DirList{tview.NewBox(), fsItems, filteredItems, path, true, selectedItemIndex, 0, DefaultSort}, nil
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
	// ClearArea(screen, x, y, width, bottomLimit)

	if (dl.selectedItemIndex + dl.itemOffset) > bottomLimit {
	}

	for i, item := range dl.FilteredItems {
		if i < dl.itemOffset {
			continue
		}

		if y >= bottomLimit {
			break
		}

		printStyle := textStyle
		if i == dl.selectedItemIndex {
			printStyle = selectedStyle
		}

		if item.Metadata.Type == Folder {
			printStyle = printStyle.Foreground(tcell.ColorBlue)
		}

		PrintWithStyle(screen, item.Metadata.Icon, item.Name, x, y, width, printStyle)

		y++
	}
}

// TODO: Get a better algorithm to adjust offset
func (dl *DirList) AdjustOffset() {
	_, _, _, height := dl.GetInnerRect()
	if height == 0 {
		return
	}
	if dl.selectedItemIndex < height {
		dl.itemOffset = 0
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

// This function is hackily made to support wide characters (nerd font)
// Try to figure out a way to make this deal with wide runes:
// Tcell docs have some mention of them so look for that
func PrintWithStyle(screen tcell.Screen, icon rune, text string, x, y, maxWidth int, style tcell.Style) {
	screen.SetContent(x, y, ' ', nil, style)
	screen.SetContent(x+1, y, icon, nil, style)
	screen.SetContent(x+2, y, ' ', nil, style)

	x = x + 3
	var i int
	var ru rune
	for i, ru = range text {
		if i >= maxWidth-2 {
			break
		}

		screen.SetContent(x+i, y, ru, nil, style)
	}

	for i <= maxWidth-4 {
		screen.SetContent(x+1+i, y, ' ', nil, style)
		i++
	}
}
