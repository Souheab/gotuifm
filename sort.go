package main

import (
	"sort"
	"strings"
)

const (
	DefaultSort = iota
	TimeSort
	TimeSortReverse
)

const (
	Less = iota
	Equal
	More
)

func SortByCriteria(fsItems []*FSItem, sortingCriteria int) {
	switch sortingCriteria {
	case DefaultSort:
		sort.Sort(FSItemsDefaultSort(fsItems))
	case TimeSort:
		sort.Sort(FSItemsTimeSort(fsItems))
	case TimeSortReverse:
		sort.Sort(FSItemsTimeSortReverse(fsItems))
	}
}

type FSItemsDefaultSort []*FSItem
type FSItemsTimeSort []*FSItem
type FSItemsTimeSortReverse []*FSItem

func (fsi FSItemsDefaultSort) Less(i, j int) bool {
	fsI := fsi[i]
	fsJ := fsi[j]

	if (fsI.Metadata.Dotfile == fsJ.Metadata.Dotfile) && (fsI.Metadata.Type == fsJ.Metadata.Type) {
		iName := strings.ToLower(fsI.Name)
		jName := strings.ToLower(fsJ.Name)

		return iName < jName
	}

	if fsI.Metadata.Dotfile && fsI.Metadata.Type == File {
		return false
	}

	if fsI.Metadata.Dotfile && fsI.Metadata.Type == Folder {
		return true
	}

	if fsI.Metadata.Type == File {
		return !(fsJ.Metadata.Type == Folder)
	}

	if fsI.Metadata.Type == Folder {
		return fsJ.Metadata.Type == File
	}

	return false
}

func (fsi FSItemsTimeSort) Less(i, j int) bool {
	fsI := fsi[i]
	fsJ := fsi[j]

	switch SortingTypeSplitter(fsI, fsJ) {
	case Less:
		return true
	case More:
		return false
	}

	return fsI.Metadata.LastModified.Before(fsJ.Metadata.LastModified)
}

func (fsi FSItemsTimeSortReverse) Less(i, j int) bool {
	fsI := fsi[i]
	fsJ := fsi[j]

	switch SortingTypeSplitter(fsI, fsJ) {
	case Less:
		return true
	case More:
		return false
	}

	return fsJ.Metadata.LastModified.Before(fsI.Metadata.LastModified)
}

func SortingTypeSplitter(fsI *FSItem, fsJ *FSItem) int {
	if (fsI.Metadata.Dotfile == fsJ.Metadata.Dotfile) && (fsI.Metadata.Type == fsJ.Metadata.Type) {
		return Equal
	}

	if fsI.Metadata.Dotfile && fsI.Metadata.Type == File {
		return More
	}

	if fsI.Metadata.Dotfile && fsI.Metadata.Type == Folder {
		return Less
	}

	if fsI.Metadata.Type == File {
		if !(fsJ.Metadata.Type == Folder) {
			return Less
		} else {
			return More
		}
	}

	if fsI.Metadata.Type == Folder {
		if fsJ.Metadata.Type == File {
			return Less
		} else {
			return More
		}
	}

	return Less
}

// Len and Swap

func (fsi FSItemsDefaultSort) Len() int {
	return len(fsi)
}

func (fsi FSItemsTimeSort) Len() int {
	return len(fsi)
}

func (fsi FSItemsTimeSortReverse) Len() int {
	return len(fsi)
}

func (fsi FSItemsDefaultSort) Swap(i, j int) {
	fsI := fsi[i]
	fsi[i] = fsi[j]
	fsi[j] = fsI
}

func (fsi FSItemsTimeSort) Swap(i, j int) {
	fsI := fsi[i]
	fsi[i] = fsi[j]
	fsi[j] = fsI
}

func (fsi FSItemsTimeSortReverse) Swap(i, j int) {
	fsI := fsi[i]
	fsi[i] = fsi[j]
	fsi[j] = fsI
}
