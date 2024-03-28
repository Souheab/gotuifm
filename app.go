package main

import (
	"time"
)

const (
	TickTime = time.Millisecond * 10
)

var ExitProgramFlag bool = false

type InputContext struct {
	SortMenu bool
}

func (ic *InputContext) Clear(){
	ic.SortMenu = false
}

func RunApp(options AppOptions) {
	b := InitAppBackend(*options.WorkingDirectory)
	s := b.Screen
	inputContext := &InputContext{false}

	defer func() {
		if r := recover(); r != nil {
			s.Fini()
			panic(r)
		}
	}()

	s.Init()

	go b.ProcessScreenEventsWorker()
	go b.ProcessDirListEventsWorker()

	for {
		b.Draw()
		b.HandleKeyEvent(inputContext)
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
