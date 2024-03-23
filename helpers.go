package main

import (
	"fmt"
	"os"
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
