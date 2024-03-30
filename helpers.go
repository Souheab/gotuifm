package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/gabriel-vasile/mimetype"
	"github.com/gdamore/tcell/v2"
)

func PathExists(s string) bool {
	_, err := os.Stat(s)
	return err == nil
}

func PathReadable(path string) bool {
	_, err := os.ReadDir(path)
	return err == nil
}

func GetFileSizeHumanReadableString(fileSize int64) string {
	if fileSize < 1024 {
		return fmt.Sprintf("%d B", fileSize)
	}
	units := []string{"B", "KB", "MB", "GB", "TB", "PB", "EB"}
	var i int
	value := float64(fileSize)

	for value > 1024 {
		value /= 1024
		i++
	}
	return fmt.Sprintf("%.1f %s", value, units[i])
}

func ClearArea(s tcell.Screen, x, y, width, height int) {
	for i := range height {
		for j := range width {
			s.SetContent(x+j, y+i, ' ', nil, tcell.StyleDefault)
		}
	}
}

// Opens with xdg-open, assumes user has xdg-utils package installed
// Maybe create a seperate package for this?
func (b *Backend) Open(path string) {
	mime, _ := mimetype.DetectFile(path)
	str := mime.String()
	parts := strings.Split(str, "/")
	str = parts[0]

	if str == "text" {
		b.OpenEditor(path)
		return
	}

	cmd := exec.Command("xdg-open", path)
	cmd.Start()
	cmd.Process.Release()
}

func (b *Backend) OpenEditor(path string) {
	defaultEditor := os.Getenv("EDITOR")
	cmd := exec.Command(defaultEditor, path)
	b.Screen.Suspend()
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Run()
	b.Screen.Resume()
}
