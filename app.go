package main

import (
	"time"
)

const (
	TickTime = time.Millisecond * 10
)

var ExitProgramFlag bool = false


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

	go b.ProcessScreenEventsWorker()

	for {
		// b.CheckSelectedFilePreview()
		b.Draw()
		b.HandleEvents()
		if ExitProgramFlag {
			return
		}
		time.Sleep(TickTime)
	}
}

func (b *Backend) ExitApp() {
	b.Screen.Fini()
	ExitProgramFlag = true
}
