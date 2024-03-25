package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
)

type AppOptions struct {
	WorkingDirectory *string
}

func main() {
	wd, _ := os.Getwd()
  wd, _ = filepath.Abs(wd)
	wdFlag := flag.String("wd", wd, "Initial working directory")
	flag.Parse()

	if !PathExists(*wdFlag) {
		fmt.Fprintln(os.Stdout, "Error: Invalid path.")
		os.Exit(1)
	}
	options := AppOptions{wdFlag}

	RunApp(options)
}
