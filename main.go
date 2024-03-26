package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime/pprof"
)
var cpuprofile = flag.String("cpuprofile", "", "write cpu profile to file")
var memprofile = flag.String("memprofile", "", "write mem profile to file")
type AppOptions struct {
	WorkingDirectory *string
}

func main() {
	wd, _ := os.Getwd()
  wd, _ = filepath.Abs(wd)
	wdFlag := flag.String("wd", wd, "Initial working directory")
	flag.Parse()

	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			log.Fatal(err)
		}
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}


	if !PathExists(*wdFlag) {
		fmt.Fprintln(os.Stdout, "Error: Invalid path.")
		os.Exit(1)
	}
	options := AppOptions{wdFlag}

	RunApp(options)

	if *memprofile != "" {
		f, err := os.Create(*memprofile)
		if err != nil {
			log.Fatal(err)
		}
		pprof.WriteHeapProfile(f)
	}
}
