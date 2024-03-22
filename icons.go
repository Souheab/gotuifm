package main

import (
	"strings"

	"github.com/gabriel-vasile/mimetype"
)

const (
	FolderIcon = ''
	FileIcon   = ''
	ImageIcon  = ''
	GoLangIcon = ''
)

func GetMimeTypeIcon(mime *mimetype.MIME) rune {
	mimeString := mime.String()
	parts := strings.Split(mimeString, "/")
	first_part := parts[0]

	switch first_part {
	case "image":
		return ImageIcon
	}

	return FileIcon
}
