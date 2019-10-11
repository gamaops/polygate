package main

import (
	"time"
)

type ResetableTimerStatus uint8

const (
	RTCancel  ResetableTimerStatus = 0
	RTTimeout ResetableTimerStatus = 1
)

type ResetableTimer struct {
	resetCh  chan bool
	duration time.Duration
	Status   chan ResetableTimerStatus
}

func (r *ResetableTimer) Reset() {
	r.resetCh <- true
}

func (r *ResetableTimer) Start() {
	select {
	case <-time.After(r.duration):
		r.Status <- RTTimeout
	case <-r.resetCh:
		go r.Start()
	}
}

func NewResetableTimer(duration time.Duration, minCancellations int32) *ResetableTimer {
	r := &ResetableTimer{
		duration: duration,
		resetCh:  make(chan bool, 1),
		Status:   make(chan ResetableTimerStatus, 1),
	}

	go r.Start()

	return r
}
