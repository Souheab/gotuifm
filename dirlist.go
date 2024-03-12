package main

import (
	"github.com/rivo/tview"
	"os"
	"path/filepath"
)

const (
	Folder = iota
	File
)

type DirList struct {
	*tview.List
	FSItems []FSItem
	Path    string
}

func (dl *DirList) GetItemAtIndex(index int) *FSItem {
	if index >= len(dl.FSItems) {
		return nil
	} else {
		return &dl.FSItems[index]
	}
}

type FSItem struct {
	Path     string
	Name     string
	Metadata FSItemMetadata
}

type FSItemMetadata struct {
	Type int
}

func NewDirList(path string) (DirList, error) {

	fsDirEntry, err := os.ReadDir(path)
	if err != nil {
		return DirList{nil, nil, path}, err
	}

	files := make([]FSItem, 0, 0)
	folders := make([]FSItem, 0, 0)

	for _, fsEntry := range fsDirEntry {
		name := fsEntry.Name()
		fsItemPath, _ := filepath.Abs(filepath.Join(path, name))

		if fsEntry.IsDir() {
			metadata := FSItemMetadata{Folder}
			folder := FSItem{fsItemPath, fsEntry.Name(), metadata}
			folders = append(folders, folder)
		} else {
			metadata := FSItemMetadata{File}
			file := FSItem{fsItemPath, fsEntry.Name(), metadata}
			files = append(files, file)
		}
	}

	fsItems := append(folders, files...)

	list := NewList(folders, files)

	return DirList{list, fsItems, path}, nil
}
