package main

import "os"

func PathExists(s string) bool {
	_, err := os.Stat(s)
	return err == nil
}
