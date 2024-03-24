package main

import (
	"os"
	"path/filepath"
	"time"
)

const (
	TickTime = time.Millisecond * 10
)

func RunApp() {
	cwd, _ := os.Getwd()
	cwd, _ = filepath.Abs(cwd)
	b := InitAppBackend(cwd)
	s := b.Screen

	defer func() {
		if r := recover(); r != nil {
			s.Fini()
			panic(r)
		}
	}()

	s.Init()

	go b.ProcessEventsWorker()

	for {
		b.Draw()
		b.HandleKeyEvent()
		time.Sleep(TickTime)
	}
}

func (b *AppBackend) ExitApp() {
	b.Screen.Fini()
	os.Exit(0)
}
