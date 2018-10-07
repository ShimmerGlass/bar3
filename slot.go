package main

import (
	"time"
)

type Slot interface {
	Out() <-chan string
	Signal()
}

type BasicSlot chan string

func (b BasicSlot) Out() <-chan string {
	return b
}

func (b BasicSlot) Signal() {
}

type CallbackSlot struct {
	out        chan string
	mainSignal chan struct{}
}

func NewCbSlot(exec func() string, signals ...chan struct{}) *CallbackSlot {
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

	return &CallbackSlot{
		out:        out,
		mainSignal: mainSignal,
	}
}

func (s *CallbackSlot) Out() <-chan string {
	return s.out
}

func (s *CallbackSlot) Signal() {
	s.mainSignal <- struct{}{}
}

func NewTimedSlot(interval time.Duration, exec func() string, signals ...chan struct{}) *CallbackSlot {
	tsignal := make(chan struct{})
	go func() {
		for range time.Tick(interval) {
			tsignal <- struct{}{}
		}
	}()

	return NewCbSlot(exec, append(signals, tsignal)...)
}

func Static(c string) *CallbackSlot {
	return NewCbSlot(func() string {
		return c
	})
}
