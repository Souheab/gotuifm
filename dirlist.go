package main

import (
	"os"
	"path/filepath"

	"github.com/gabriel-vasile/mimetype"
	"github.com/rivo/tview"
)

const (
	Folder = iota
	File
)

type DirList struct {
	*tview.List
	FSItems       []*FSItem
	FilteredItems []*FSItem
	Path          string
	DotfilesFlag  bool
}

func (dl *DirList) GetItemAtIndex(index int) *FSItem {
	fsItems := dl.FilteredItems

	if index >= len(fsItems) {
		return nil
	} else {
		return fsItems[index]
	}
}

func (dl *DirList) SetFilter(f func(fsItem *FSItem) bool) {
	filteredItems := make([]*FSItem, 0, 0)

	for _, fsItem := range dl.FSItems {
		if f(fsItem) {
			filteredItems = append(filteredItems, fsItem)
		}
	}

	dl.FilteredItems = filteredItems
	dl.List = NewList(filteredItems)
}

func (dl *DirList) RemoveFilter() {
	dl.FilteredItems = dl.FSItems
	dl.List = NewList(dl.FSItems)
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
	Type     int
	Readable bool
	MimeType *mimetype.MIME
}

func NewDirList(path string) (DirList, error) {

	fsDirEntry, err := os.ReadDir(path)
	if err != nil {
		return DirList{nil, nil, nil, path, true}, err
	}

	files := make([]*FSItem, 0, 0)
	folders := make([]*FSItem, 0, 0)

	for _, fsEntry := range fsDirEntry {
		name := fsEntry.Name()
		fsItemPath, _ := filepath.Abs(filepath.Join(path, name))
		mime, _ := mimetype.DetectFile(fsItemPath)

		if fsEntry.IsDir() {
			metadata := FSItemMetadata{Folder, PathReadable(fsItemPath), mime}
			folder := FSItem{fsItemPath, fsEntry.Name(), metadata}
			folders = append(folders, &folder)
		} else {
			metadata := FSItemMetadata{File, true, mime}
			file := FSItem{fsItemPath, fsEntry.Name(), metadata}
			files = append(files, &file)
		}
	}

	fsItems := append(folders, files...)

	list := NewList(fsItems)

	return DirList{list, fsItems, fsItems, path, true}, nil
}
