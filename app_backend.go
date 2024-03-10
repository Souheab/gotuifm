package main

import (
	"log"
	"path/filepath"

)

type AppBackend struct {
	DirListCache map[string]*DirList
	UI           UI
}

func InitAppBackend(startingPath string) AppBackend {
	dlc := make(map[string]*DirList)
	ui := InitUI()
  
	backend := AppBackend{dlc, ui}
	dl := DirListCacheAdd(dlc, startingPath)
	backend.UI.Grid.AddItem(dl.UI, 1, 1, 1, 1, 0, 0, true)
	
	return backend
}

func DirListCacheAdd(dlc map[string]*DirList, path string) *DirList{
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
