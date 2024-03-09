package main

import (
	"os"
	"path/filepath"
)

type DirList struct {
	Folders []Folder
	Files   []File
}

type Folder struct {
	Path string
	Name string
}

type File struct {
	Path string
	Name string
}

func NewDirList(path string) (DirList, error) {
	
	fsItems, err := os.ReadDir(path)
	if err != nil {
		return DirList{nil, nil}, err
	}

	files := make([]File, 0, 0)
	folders := make([]Folder, 0, 0)

	for _, fsItem := range fsItems {
		name := fsItem.Name()
		path := filepath.Join(path, name)

		if fsItem.IsDir() {
			folder := Folder{path, fsItem.Name()}
			folders = append(folders, folder)
		} else {
			file := File{path, fsItem.Name()}
			files = append(files, file)
		}
	}

	return DirList{folders, files}, nil
}


