package main

import "os"

func PathExists(s string) bool {
	_, err := os.Stat(s)
	return err == nil
}

func PathReadable(path string) bool {
	_, err := os.ReadDir(path)
	return err == nil
}
