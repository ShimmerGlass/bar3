package main

import (
	"time"
)

type Slot interface {
	Out() <-chan []Part
	Signal()
}

type BasicSlot chan []Part

func (b BasicSlot) Out() <-chan []Part {
	return b
}

func (b BasicSlot) Signal() {
}

type CallbackSlot struct {
	out        chan []Part
	mainSignal chan struct{}
}

func NewCbSlot(exec func() []Part, signals ...chan struct{}) *CallbackSlot {
	mainSignal := make(chan struct{})
	signals = append(signals, mainSignal)

	out := make(chan []Part)

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

func (s *CallbackSlot) Out() <-chan []Part {
	return s.out
}

func (s *CallbackSlot) Signal() {
	s.mainSignal <- struct{}{}
}

func NewTimedSlot(interval time.Duration, exec func() []Part, signals ...chan struct{}) *CallbackSlot {
	tsignal := make(chan struct{})
	go func() {
		for range time.Tick(interval) {
			tsignal <- struct{}{}
		}
	}()

	return NewCbSlot(exec, append(signals, tsignal)...)
}

func Static(c []Part) *CallbackSlot {
	return NewCbSlot(func() []Part {
		return c
	})
}
