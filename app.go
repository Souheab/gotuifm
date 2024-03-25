package main

import (
	"os"
	"time"
)

const (
	TickTime = time.Millisecond * 10
)

func RunApp(options AppOptions) {
	b := InitAppBackend(*options.WorkingDirectory)
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

func (b *Backend) ExitApp() {
	b.Screen.Fini()
	os.Exit(0)
}
