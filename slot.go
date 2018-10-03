package main

import (
	"time"
)

type Slot interface {
	Out() <-chan string
	Signal()
}

type BasicSlot struct {
	out        chan string
	mainSignal chan struct{}
}

func NewBasicSlot(exec func() string, signals ...chan struct{}) *BasicSlot {
	mainSignal := make(chan struct{})
	signals = append(signals, mainSignal)

	out := make(chan string)

	go func() {
		agg := make(chan struct{})
		for _, s := range signals {
			go func(s chan struct{}) {
				for range s {
					agg <- struct{}{}
				}
			}(s)
		}
		for range agg {
			out <- exec()
		}
	}()

	return &BasicSlot{
		out:        out,
		mainSignal: mainSignal,
	}
}

func (s *BasicSlot) Out() <-chan string {
	return s.out
}

func (s *BasicSlot) Signal() {
	s.mainSignal <- struct{}{}
}

func NewTimedSlot(interval time.Duration, exec func() string, signals ...chan struct{}) *BasicSlot {
	tsignal := make(chan struct{})
	go func() {
		for range time.Tick(interval) {
			tsignal <- struct{}{}
		}
	}()

	return NewBasicSlot(exec, append(signals, tsignal)...)
}

func Static(c string) *BasicSlot {
	return NewBasicSlot(func() string {
		return c
	})
}
